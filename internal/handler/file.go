package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marianozunino/drop/internal/config"
)

// HandleFileAccess serves uploaded files with proper content type detection
func (h *Handler) HandleFileAccess(c echo.Context) error {
	filename := c.Param("filename")

	parts := strings.SplitN(filename, "/", 2)
	filename = parts[0]

	if strings.Contains(filename, "..") {
		return c.String(http.StatusBadRequest, "Invalid file path")
	}

	filePath := filepath.Join(config.DefaultUploadPath, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "File not found")
	}

	contentType := "application/octet-stream"

	metadataPath := filePath + ".meta"
	if _, err := os.Stat(metadataPath); err == nil {
		metadataBytes, err := os.ReadFile(metadataPath)
		if err == nil {
			var metadata FileMetadata
			if err := json.Unmarshal(metadataBytes, &metadata); err == nil && metadata.ContentType != "" {
				contentType = metadata.ContentType

				if !metadata.ExpiresAt.IsZero() {
					expiresMs := metadata.ExpiresAt.UnixNano() / int64(time.Millisecond)
					c.Response().Header().Set("X-Expires", fmt.Sprintf("%d", expiresMs))
				}
			}
		}
	}

	if contentType == "application/octet-stream" {
		ext := filepath.Ext(filename)
		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			contentType = mimeType
		} else {
			file, err := os.Open(filePath)
			if err == nil {
				defer file.Close()

				buffer := make([]byte, 512)

				_, err := file.Read(buffer)
				if err == nil {
					contentType = http.DetectContentType(buffer)
				}
			}
		}
	}

	c.Response().Header().Set("Content-Type", contentType)

	if strings.HasPrefix(contentType, "video/") ||
		strings.HasPrefix(contentType, "audio/") ||
		strings.HasPrefix(contentType, "image/") ||
		contentType == "application/pdf" ||
		strings.HasPrefix(contentType, "text/") {
		c.Response().Header().Set("Content-Disposition", "inline")
	}

	return c.File(filePath)
}

// HandleFileManagement manages file operations (delete, update expiration)
func (h *Handler) HandleFileManagement(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(32 << 20); err != nil { // feeling a bit smart ass today
		if err := c.Request().ParseForm(); err != nil {
			log.Printf("Info: Non-form request or parsing error: %v", err)
		}
	}

	filename := c.Param("filename")
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return c.String(http.StatusBadRequest, "Invalid file path")
	}

	filePath := filepath.Join(config.DefaultUploadPath, filename)
	metadataPath := filePath + ".meta"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "File not found")
	}

	token := c.FormValue("token")
	if token == "" {
		return c.String(http.StatusBadRequest, "Missing management token")
	}

	authenticated := false
	var metadata FileMetadata

	metadataBytes, err := os.ReadFile(metadataPath)
	if err == nil {
		if err := json.Unmarshal(metadataBytes, &metadata); err == nil {
			if metadata.Token == token {
				authenticated = true
			} else {
				return c.String(http.StatusUnauthorized, "Invalid management token")
			}
		} else {
			log.Printf("Warning: Failed to parse metadata for %s: %v", filename, err)
		}
	} else {
		log.Printf("Warning: Metadata file not found for %s: %v", filename, err)
		authenticated = true
	}

	if !authenticated {
		return c.String(http.StatusUnauthorized, "Authentication failed")
	}

	if _, deleteRequested := c.Request().Form["delete"]; deleteRequested {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Error: Failed to delete file %s: %v", filePath, err)
			return c.String(http.StatusInternalServerError, "Failed to delete file")
		}

		if err := os.Remove(metadataPath); err != nil {
			log.Printf("Warning: Failed to delete metadata %s: %v", metadataPath, err)
		}

		return c.String(http.StatusOK, "File deleted successfully")
	}

	if expiresStr := c.FormValue("expires"); expiresStr != "" {
		expirationDate, err := parseExpirationTime(expiresStr)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid expiration format: %v", err))
		}

		if metadataBytes != nil {
			metadata.ExpiresAt = expirationDate

			updatedMetadata, err := json.Marshal(metadata)
			if err != nil {
				log.Printf("Error: Failed to marshal metadata: %v", err)
				return c.String(http.StatusInternalServerError, "Failed to update expiration")
			}

			if err := os.WriteFile(metadataPath, updatedMetadata, 0o644); err != nil {
				log.Printf("Error: Failed to save metadata: %v", err)
				return c.String(http.StatusInternalServerError, "Failed to save expiration")
			}

			return c.String(http.StatusOK, "Expiration updated successfully")
		}

		newMetadata := FileMetadata{
			Token:      token,
			UploadDate: time.Now(),
			ExpiresAt:  expirationDate,
		}

		if fileInfo, err := os.Stat(filePath); err == nil {
			newMetadata.Size = fileInfo.Size()
		}

		newMetadataBytes, err := json.Marshal(newMetadata)
		if err != nil {
			log.Printf("Error: Failed to marshal new metadata: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to create expiration metadata")
		}

		if err := os.WriteFile(metadataPath, newMetadataBytes, 0o644); err != nil {
			log.Printf("Error: Failed to save new metadata: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to save expiration metadata")
		}

		return c.String(http.StatusOK, "Expiration created successfully")
	}

	return c.String(http.StatusBadRequest, "No valid operation specified. Use 'delete' or 'expires'.")
}

// HandleUpload processes file uploads
func (h *Handler) HandleUpload(c echo.Context) error {
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, config.MaxUploadSize)

	var fileContent []byte
	var fileName string
	var fileSize int64
	var contentType string

	file, header, err := c.Request().FormFile("file")
	if err == nil {
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to read file")
		}

		fileContent = content
		fileName = header.Filename
		fileSize = header.Size
		contentType = header.Header.Get("Content-Type")
	} else {
		url := c.FormValue("url")
		if url == "" {
			return c.String(http.StatusBadRequest, "No file or URL provided")
		}

		resp, err := http.Get(url)
		if err != nil {
			return c.String(http.StatusBadRequest, "Failed to download from URL")
		}
		defer resp.Body.Close()

		contentLength := resp.Header.Get("Content-Length")
		if contentLength == "" {
			return c.String(http.StatusBadRequest, "Remote server did not provide Content-Length")
		}

		length, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid Content-Length")
		}

		if length > config.MaxUploadSize {
			return c.String(http.StatusBadRequest, fmt.Sprintf("File too large (max %d bytes)", config.MaxUploadSize))
		}

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to read from URL")
		}

		urlPath := strings.Split(url, "/")
		if len(urlPath) > 0 {
			fileName = urlPath[len(urlPath)-1]
		} else {
			fileName = "download"
		}

		fileContent = content
		fileSize = int64(len(content))
		contentType = resp.Header.Get("Content-Type")
	}

	useSecretId := c.FormValue("secret") != ""

	var id string
	if useSecretId {
		id, err = generateID(8)
	} else {
		id, err = generateID(config.DefaultIDLength)
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Server error")
	}

	fileExt := filepath.Ext(fileName)

	filename := id
	if fileExt != "" {
		filename += fileExt
	}

	filePath := filepath.Join(config.DefaultUploadPath, filename)

	if err := os.WriteFile(filePath, fileContent, 0o644); err != nil {
		return c.String(http.StatusInternalServerError, "Server error")
	}

	var expirationDate time.Time

	expiresStr := c.FormValue("expires")
	if expiresStr != "" {
		expirationDate, err = parseExpirationTime(expiresStr)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid expiration format: %v", err))
		}
	} else if h.expManager != nil {
		expirationDate = h.expManager.GetExpirationDate(fileSize)
	}

	managementToken, err := generateID(16)
	if err != nil {
		log.Printf("Warning: Failed to generate management token: %v", err)
	}

	metadata := FileMetadata{
		Token:        managementToken,
		OriginalName: fileName,
		UploadDate:   time.Now(),
		Size:         fileSize,
		ContentType:  contentType,
	}

	if !expirationDate.IsZero() {
		metadata.ExpiresAt = expirationDate
	}

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("Warning: Failed to serialize file metadata: %v", err)
	} else {
		metadataPath := filePath + ".meta"
		if err := os.WriteFile(metadataPath, metadataBytes, 0o644); err != nil {
			log.Printf("Warning: Failed to write metadata file: %v", err)
		}
	}

	c.Response().Header().Set("X-Token", managementToken)

	fileURL := h.expManager.Config.BaseURL + filename

	if !expirationDate.IsZero() {
		expiresMs := expirationDate.UnixNano() / int64(time.Millisecond)
		c.Response().Header().Set("X-Expires", fmt.Sprintf("%d", expiresMs))
	}

	if strings.Contains(c.Request().Header.Get("Accept"), "application/json") {
		response := map[string]interface{}{
			"url":   fileURL,
			"size":  fileSize,
			"token": managementToken,
		}

		if !expirationDate.IsZero() {
			response["expires_at"] = expirationDate
			days := int(expirationDate.Sub(time.Now()).Hours() / 24)
			response["expires_in_days"] = days
		}

		return c.JSON(http.StatusOK, response)
	} else {
		c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
		return c.String(http.StatusOK, fileURL+"\n")
	}
}

// parseExpirationTime parses different expiration time formats
func parseExpirationTime(expiresStr string) (time.Time, error) {
	if hours, err := strconv.Atoi(expiresStr); err == nil {
		return time.Now().Add(time.Duration(hours) * time.Hour), nil
	}

	if ms, err := strconv.ParseInt(expiresStr, 10, 64); err == nil {
		return time.Unix(0, ms*int64(time.Millisecond)), nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, expiresStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized date/time format")
}

// generateID creates a random hex string of given length
func generateID(length int) (string, error) {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
