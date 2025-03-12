package expiration

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/marianozunino/drop/config"
)

// ExpirationManager handles the file expiration process
type ExpirationManager struct {
	Config     config.Config
	configPath string
	stopChan   chan struct{}
}

// FileMetadata stores information about uploaded files
type FileMetadata struct {
	Token        string    `json:"token"`
	OriginalName string    `json:"original_name,omitempty"`
	UploadDate   time.Time `json:"upload_date"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type,omitempty"`
}

// NewExpirationManager creates a new expiration manager
func NewExpirationManager(configPath string) (*ExpirationManager, error) {
	manager := &ExpirationManager{
		Config:     config.DefaultConfig,
		configPath: configPath,
		stopChan:   make(chan struct{}),
	}

	// Try to load config file, if it exists
	if err := manager.LoadConfig(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load expiration config: %v", err)
	}

	spew.Dump(manager.Config)

	return manager, nil
}

// LoadConfig loads the configuration from a file
func (m *ExpirationManager) LoadConfig() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	var config config.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("invalid config file format: %v", err)
	}

	// Validate config values
	if config.MinAge <= 0 {
		return fmt.Errorf("min_age_days must be greater than 0")
	}
	if config.MaxAge <= config.MinAge {
		return fmt.Errorf("max_age_days must be greater than min_age_days")
	}
	if config.MaxSize <= 0 {
		return fmt.Errorf("max_size_mib must be greater than 0")
	}
	if config.CheckInterval <= 0 {
		return fmt.Errorf("check_interval_min must be greater than 0")
	}

	m.Config = config
	return nil
}

// SaveConfig saves the current configuration to a file
func (m *ExpirationManager) SaveConfig() error {
	data, err := json.MarshalIndent(m.Config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0o644)
}

// Start begins the expiration checking process
func (m *ExpirationManager) Start() {
	go func() {
		m.cleanupExpiredFiles() // Run once immediately

		ticker := time.NewTicker(time.Duration(m.Config.CheckInterval) * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.cleanupExpiredFiles()
			case <-m.stopChan:
				log.Println("Expiration manager stopped")
				return
			}
		}
	}()
	log.Printf("Expiration manager started, checking every %d minutes", m.Config.CheckInterval)
}

// Stop halts the expiration checking process
func (m *ExpirationManager) Stop() {
	close(m.stopChan)
}

// calculateRetention determines how long a file should be kept based on its size
func (m *ExpirationManager) calculateRetention(fileSize float64) time.Duration {
	// Convert file size to MiB
	fileSizeMiB := fileSize / (1024 * 1024)

	// If file is smaller than max size, use max age (longer retention)
	if fileSizeMiB <= m.Config.MaxSize {
		return time.Duration(m.Config.MaxAge) * 24 * time.Hour
	}

	// Apply the formula:
	// retention = min_age + (min_age - max_age) * pow((file_size / max_size - 1), 3)
	// Note: This formula decreases retention as file size increases
	fileSizeRatio := fileSizeMiB/m.Config.MaxSize - 1
	ageDiff := float64(m.Config.MinAge - m.Config.MaxAge)  // This will be negative (min_age < max_age)
	additionalDays := ageDiff * math.Pow(fileSizeRatio, 3) // This will further reduce the retention

	// Calculate total days, which will be less than min_age for large files
	totalDays := float64(m.Config.MinAge) + additionalDays

	// Ensure we don't go below min_age (the minimum retention period)
	if totalDays < float64(m.Config.MinAge) {
		totalDays = float64(m.Config.MinAge)
	}

	return time.Duration(totalDays) * 24 * time.Hour
}

func (m *ExpirationManager) calculateRetentionFixed(fileSize float64) time.Duration {
	// Convert file size to MiB
	fileSizeMiB := fileSize / (1024 * 1024)

	// If file is smaller than max size, use min age
	if fileSizeMiB <= m.Config.MaxSize {
		return time.Duration(m.Config.MinAge) * 24 * time.Hour
	}

	// Apply the corrected formula:
	// retention = min_age + (max_age - min_age) * relationship factor
	// This makes the retention INCREASE with file size, not decrease
	fileSizeRatio := fileSizeMiB/m.Config.MaxSize - 1
	ageDiff := float64(m.Config.MaxAge - m.Config.MinAge) // Changed from (min_age - max_age)
	additionalDays := ageDiff * fileSizeRatio             // Linear relationship for simplicity

	// Clamp to ensure we don't exceed max_age
	totalDays := float64(m.Config.MinAge) + additionalDays
	if totalDays > float64(m.Config.MaxAge) {
		totalDays = float64(m.Config.MaxAge)
	}

	return time.Duration(totalDays) * 24 * time.Hour
}

// CheckMetadataExpiration checks if a file has expired based on its metadata
func (m *ExpirationManager) CheckMetadataExpiration(metadataPath string) (bool, error) {
	// Read metadata file
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return false, err
	}

	// Parse metadata
	var metadata FileMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return false, err
	}

	// Check if file has explicit expiration date
	if !metadata.ExpiresAt.IsZero() {
		return time.Now().After(metadata.ExpiresAt), nil
	}

	// If no expiration date is set in metadata, calculate based on retention policy
	if metadata.UploadDate.IsZero() {
		// If upload date is missing, we can't calculate expiration
		return false, nil
	}

	// Calculate retention based on file size and upload date
	retention := m.calculateRetention(float64(metadata.Size))
	expirationTime := metadata.UploadDate.Add(retention)

	return time.Now().After(expirationTime), nil
}

// cleanupExpiredFiles checks all files and removes those that have expired
func (m *ExpirationManager) cleanupExpiredFiles() {
	if !m.Config.Enabled {
		return
	}
	uploadPath := m.Config.UploadPath

	log.Println("Checking for expired files...")

	files, err := os.ReadDir(uploadPath)
	if err != nil {
		log.Printf("Error reading upload directory: %v", err)
		return
	}

	var removed, total int
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Skip metadata files, we'll handle them with their corresponding files
		if strings.HasSuffix(file.Name(), ".meta") {
			continue
		}

		total++
		filePath := filepath.Join(uploadPath, file.Name())
		metadataPath := filePath + ".meta"

		// Check if metadata exists
		if _, err := os.Stat(metadataPath); err == nil {
			// Use metadata for expiration check
			expired, err := m.CheckMetadataExpiration(metadataPath)
			if err != nil {
				log.Printf("Error checking metadata expiration for %s: %v", file.Name(), err)
				continue
			}

			if expired {
				log.Printf("Removing expired file (from metadata): %s", file.Name())
				if err := os.Remove(filePath); err != nil {
					log.Printf("Error removing expired file %s: %v", filePath, err)
				} else {
					// Also remove the metadata file
					if err := os.Remove(metadataPath); err != nil {
						log.Printf("Error removing metadata file %s: %v", metadataPath, err)
					}
					removed++
				}
			}
			continue
		}

		// Fallback to standard file expiration check
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Error getting file info for %s: %v", filePath, err)
			continue
		}

		// Calculate the retention period based on file size
		retention := m.calculateRetention(float64(fileInfo.Size()))
		expirationTime := fileInfo.ModTime().Add(retention)

		// If the file has expired, remove it
		if time.Now().After(expirationTime) {
			if err := os.Remove(filePath); err != nil {
				log.Printf("Error removing expired file %s: %v", filePath, err)
			} else {
				log.Printf("Removed expired file: %s (age: %v, size: %.2f MiB)",
					file.Name(),
					time.Since(fileInfo.ModTime()).Round(time.Hour),
					float64(fileInfo.Size())/(1024*1024))
				removed++
			}
		}
	}

	log.Printf("Expiration check complete. Removed %d of %d files", removed, total)
}

// GetExpirationDate calculates when a file will expire based on its size
func (m *ExpirationManager) GetExpirationDate(fileSize int64) time.Time {
	spew.Dump(fileSize)
	retention := m.calculateRetention(float64(fileSize))
	return time.Now().Add(retention)
}
