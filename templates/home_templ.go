// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/marianozunino/drop/internal/config"
	"strconv"
)

func HomePage(config config.Config) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>MZ.DROP</title><style>\n\t\t\t\tpre {\n\t\t\t\t\twhite-space: pre;\n\t\t\t\t\tfont-family: monospace;\n\t\t\t\t\tline-height: 1.2;\n\t\t\t\t\toverflow-x: auto;\n\t\t\t\t}\n\t\t\t</style></head><body><h1>MZ.DROP</h1><p>Temporary file hoster, inspired by <a href=\"https://0x0.st/\">0x0.st</a> .</p><p><small>A personal service for quick file sharing.</small></p><div>min_age = ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(config.MinAge))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 30, Col: 43}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, " days<br>max_age = ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(config.MaxAge))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 32, Col: 43}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, " days<br>max_size = ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%.1f", config.MaxSize))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 34, Col: 52}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, " MiB<br>retention = min_age + (min_age - max_age) * pow((file_size / max_size - 1), 3)</div><pre>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = RetentionGraph(config).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</pre><details id=\"uploading\"><summary>Uploading files</summary><pre>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = ValidFields(config).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</pre><details><summary>cURL examples</summary><pre>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Examples(config).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</pre></details><p>It is possible to append a custom file name to any URL:<br><code>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(config.BaseURL)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 55, Col: 27}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "aaa.jpg/image.jpeg</code></p><p>File URLs are valid for at least ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(config.MinAge))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 57, Col: 69}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, " days and up to ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(config.MaxAge))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 57, Col: 116}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, " days (see above).</p><p>Expired files won't be removed immediately but within the next ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(config.CheckInterval))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 58, Col: 106}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, " minutes.</p><p>Maximum file size: 100MB</p></details> <details id=\"managing\"><summary>Managing your files</summary><pre>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = FileManagement().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "</pre><details><summary>cURL examples</summary><p>Delete a file immediately:</p><pre>curl -X POST -F'token=token_here' -F'delete=' ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 string
		templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(config.BaseURL)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 69, Col: 72}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "abc.txt</pre><p>Change the expiration date (see above):</p><pre>curl -X POST -F'token=token_here' -F'expires=3' ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(config.BaseURL)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/home.templ`, Line: 71, Col: 74}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "abc.txt</pre></details></details> <details><summary>Terms of Service</summary><p>This service is NOT a platform for:</p><ul><li>piracy</li><li>pornography and gore</li><li>extremist material of any kind</li><li>terrorist content</li><li>malware / botnet C&C</li><li>anything related to crypto currencies</li><li>backups</li><li>CI build artifacts</li><li>other automated mass uploads</li><li>doxxing, database dumps containing personal information</li><li>anything illegal</li></ul><p>Uploads found to be in violation of these rules will be removed, and the originating IP address may be blocked from further uploads.</p></details> <details><summary>Privacy Policy</summary><p>For the purpose of moderation, the following is stored with each uploaded file:</p><ul><li>IP address</li><li>User agent string</li></ul><p>This site generally does not log requests, but may enable logging if necessary for purposes such as threat mitigation.</p><p>No data is shared with third parties.</p></details><hr><p>Personal instance inspired by <a href=\"https://0x0.st/\">0x0.st</a>.</p><p>Hosted on mz.uy for personal use.</p></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func RetentionGraph(config config.Config) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var11 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var11 == nil {
			templ_7745c5c3_Var11 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.Raw(`
   days
   `+strconv.Itoa(config.MaxAge)+`  |  \
       |   \
       |    \
       |     \
       |      \
       |       \
       |        ..
       |          \
   `+strconv.Itoa((config.MinAge+config.MaxAge)/2)+`  | ----------..-------------------------------------------
       |             ..
       |               \
       |                ..
       |                  ...
       |                     ..
       |                       ...
       |                          ....
       |                              ......
    `+strconv.Itoa(config.MinAge)+`  |                                    ....................
         0                      `+fmt.Sprintf("%.1f", config.MaxSize/2)+`                      `+fmt.Sprintf("%.1f", config.MaxSize)+`
                                                             MiB
    `).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func ValidFields(config config.Config) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var12 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var12 == nil {
			templ_7745c5c3_Var12 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.Raw(`
Send HTTP POST requests to `+config.BaseURL+` with data encoded as multipart/form-data

Valid fields are:
  ┌─────────┬────────────┬──────────────────────────────────────────────────┐
  │ field   │ content    │ remarks                                          │
  ╞═════════╪════════════╪══════════════════════════════════════════════════╡
  │ file    │ data       │                                                  │
  ├─────────┼────────────┼──────────────────────────────────────────────────┤
  │ url     │ remote URL │ Mutually exclusive with "file".                  │
  │         │            │ Remote site must return Content-Length header.   │
  ├─────────┼────────────┼──────────────────────────────────────────────────┤
  │ secret  │ (ignored)  │ If present, a longer, hard-to-guess URL          │
  │         │            │ will be generated.                               │
  ├─────────┼────────────┼──────────────────────────────────────────────────┤
  │ expires │ time       │ Sets file expiration time. Accepts:              │
  │         │ format     │ - Hours as integer (e.g., 24)                    │
  │         │            │ - Milliseconds since epoch (e.g., 1681996320000) │
  │         │            │ - RFC3339 (e.g., 2006-01-02T15:04:05Z07:00)      │
  │         │            │ - ISO date (e.g., 2006-01-02)                    │
  │         │            │ - ISO datetime (e.g., 2006-01-02T15:04:05)       │
  │         │            │ - SQL datetime (e.g., 2006-01-02 15:04:05)       │
  └─────────┴────────────┴──────────────────────────────────────────────────┘

`).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func FileManagement() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var13 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var13 == nil {
			templ_7745c5c3_Var13 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.Raw(`
	Whenever a file that does not already exist or has expired is uploaded,
	the HTTP response header includes an X-Token field. You can use this
	to perform management operations on the file by sending POST requests
	to the file URL.
	When using cURL, you can add the -i option to view the response header.
	Valid fields are:
	┌─────────┬────────────┬──────────────────────────────────────────────────┐
	│ field   │ content    │ remarks                                          │
	╞═════════╪════════════╪══════════════════════════════════════════════════╡
	│ token   │ management │ Returned after upload in X-Token HTTP header     │
	│         │ token      │ field. Required.                                 │
	├─────────┼────────────┼──────────────────────────────────────────────────┤
	│ delete  │ (ignored)  │ Removes the file.                                │
	├─────────┼────────────┼──────────────────────────────────────────────────┤
	│ expires │ time       │ Sets file expiration time. Accepts:              │
	│         │ format     │ - Hours as integer (e.g., 24)                    │
	│         │            │ - Milliseconds since epoch (e.g., 1681996320000) │
	│         │            │ - RFC3339 (e.g., 2006-01-02T15:04:05Z07:00)      │
	│         │            │ - ISO date (e.g., 2006-01-02)                    │
	│         │            │ - ISO datetime (e.g., 2006-01-02T15:04:05)       │
	│         │            │ - SQL datetime (e.g., 2006-01-02 15:04:05)       │
	└─────────┴────────────┴──────────────────────────────────────────────────┘
`).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func Examples(config config.Config) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var14 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var14 == nil {
			templ_7745c5c3_Var14 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.Raw(`
    Uploading a file:
        curl -F'file=@yourfile.png' `+config.BaseURL+`

    Copy a file from a remote URL:
        curl -F'url=http://example.com/image.jpg' `+config.BaseURL+`

    Same, but with hard-to-guess URLs:
        curl -F'file=@yourfile.png' -Fsecret= `+config.BaseURL+`
        curl -F'url=http://example.com/image.jpg' -Fsecret= `+config.BaseURL+`

    Setting retention time in hours:
        curl -F'file=@yourfile.png' -Fexpires=24 `+config.BaseURL+`

    Setting expiration date with milliseconds since UNIX epoch:
        curl -F'file=@yourfile.png' -Fexpires=1681996320000 `+config.BaseURL+`

    Setting expiration date with RFC3339 format:
        curl -F'file=@yourfile.png' -Fexpires=2023-04-20T10:15:30Z `+config.BaseURL+`

    Setting expiration date with ISO date format:
        curl -F'file=@yourfile.png' -Fexpires=2023-04-20 `+config.BaseURL+`

    Setting expiration date with ISO datetime format:
        curl -F'file=@yourfile.png' -Fexpires=2023-04-20T10:15:30 `+config.BaseURL+`

    Setting expiration date with SQL datetime format:
        curl -F'file=@yourfile.png' -Fexpires="2023-04-20 10:15:30" `+config.BaseURL+`
	`).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
