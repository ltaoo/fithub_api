package web

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Register serves the two SPAs after API routing. static_dir contains admin/ and h5/.
func Register(router *gin.Engine, static_dir string) {
	router.NoRoute(func(ctx *gin.Context) {
		request_path := ctx.Request.URL.Path
		if static_dir == "" || (ctx.Request.Method != http.MethodGet && ctx.Request.Method != http.MethodHead) || request_path == "/api" || strings.HasPrefix(request_path, "/api/") || request_path == "/health" {
			ctx.Status(http.StatusNotFound)
			return
		}
		if request_path == "/admin" {
			target := "/admin/"
			if ctx.Request.URL.RawQuery != "" {
				target += "?" + ctx.Request.URL.RawQuery
			}
			ctx.Redirect(http.StatusPermanentRedirect, target)
			return
		}
		app_name, relative_path := "h5", strings.TrimPrefix(request_path, "/")
		if strings.HasPrefix(request_path, "/admin/") {
			app_name, relative_path = "admin", strings.TrimPrefix(request_path, "/admin/")
		}
		// Never serve hidden files or use SPA fallback for traversal, asset misses or directories.
		for _, segment := range strings.Split(relative_path, "/") {
			if strings.HasPrefix(segment, ".") || strings.Contains(segment, "\\") {
				ctx.Status(http.StatusNotFound)
				return
			}
		}
		app_root, err := filepath.Abs(filepath.Join(static_dir, app_name))
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		filename := filepath.Join(app_root, filepath.FromSlash(relative_path))
		real_filename, err := filepath.EvalSymlinks(filename)
		if err == nil {
			relative, rel_err := filepath.Rel(app_root, real_filename)
			if rel_err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				ctx.Status(http.StatusNotFound)
				return
			}
			file_info, stat_err := os.Stat(real_filename)
			if stat_err == nil && file_info.Mode().IsRegular() {
				ctx.Header("X-Content-Type-Options", "nosniff")
				if path.Base(relative_path) == "index.html" {
					ctx.Header("Cache-Control", "no-cache")
				}
				ctx.File(real_filename)
				return
			}
			if relative_path != "" && file_info != nil && file_info.IsDir() {
				ctx.Status(http.StatusNotFound)
				return
			}
		}
		accept := ctx.GetHeader("Accept")
		if path.Ext(relative_path) != "" || strings.HasPrefix(relative_path, "assets/") || (accept != "" && !strings.Contains(accept, "text/html") && !strings.Contains(accept, "*/*")) {
			ctx.Status(http.StatusNotFound)
			return
		}
		index_file := filepath.Join(app_root, "index.html")
		if file_info, stat_err := os.Stat(index_file); stat_err != nil || !file_info.Mode().IsRegular() {
			ctx.String(http.StatusServiceUnavailable, "Frontend assets are not built. Run pnpm build:web.")
			return
		}
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.File(index_file)
	})
}
