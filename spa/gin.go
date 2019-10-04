package spa

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// InstallForGin creates a new Spa and installs the required routes for a Gin application
func InstallForGin(engine *gin.Engine, app Spa) error {
	if server := app.StaticFileServer(); server != nil {
		loc := server.Location
		if !strings.HasPrefix(loc, "/") {
			loc = "/" + loc
		}

		engine.Group(loc, ginStaticFileMiddleware(app)).StaticFS("", server.FileSystem)
	}
	engine.NoRoute(ginNoRouteHandler(app))
	return nil
}

func ginStaticFileMiddleware(app Spa) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "public,max-age=31536000,immutable")

		ctx.Next()
	}
}

func ginNoRouteHandler(app Spa) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch ctx.Request.Method {
		case "GET", "HEAD":
			if filepath.Ext(ctx.Request.URL.Path) == "" {
				var notModified bool

				eTag := app.ETag()

				if eTag != "" {
					ctx.Header("Cache-Control", "no-cache")
					ctx.Header("ETag", eTag)
					notModified = eTag == ctx.GetHeader("ETag")
				}

				if notModified {
					ctx.Status(http.StatusNotModified)
				} else {
					ctx.Data(http.StatusOK, "text/html", app.Index())
				}

				ctx.Abort()

				return
			}
		}

		ctx.Next()
	}
}
