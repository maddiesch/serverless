package ginapp

import (
	"context"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/maddiesch/serverless/middleware"
	"github.com/maddiesch/serverless/router"
)

type ginRouter struct {
	router.Router

	engine *gin.Engine
}

func buildMiddlewares(fns []middleware.Handler) []gin.HandlerFunc {
	mid := make([]gin.HandlerFunc, len(fns))

	for i, fn := range fns {
		mid[i] = buildMiddleware(fn)
	}

	return mid
}

func buildMiddleware(fn middleware.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		aborter := make(chan error, 1)

		ctx := c.MustGet(ginSubContextKey).(context.Context)

		go func() {
			defer close(aborter)

			new := fn(ctx, c.Writer, c.Request, func(err error) {
				aborter <- err
				runtime.Goexit()
			})

			c.Set(ginSubContextKey, new)
		}()

		err := <-aborter

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.Next()
		}
	}
}

func buildHandler(fn router.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet(ginSubContextKey).(context.Context)

		if err := fn(ctx, c.Writer, c.Request); err != nil {
			c.Error(err)
		}

		c.Next()
	}
}

func (r *ginRouter) Handle(method string, path string, handler router.HandlerFunc) {
	r.engine.Handle(method, path, buildHandler(handler))
}

func (r *ginRouter) Group(path string, middleware ...middleware.Handler) router.Group {
	group := r.engine.Group(path, buildMiddlewares(middleware)...)

	return &ginGroup{
		router:  r,
		backing: group,
	}
}

func (r *ginRouter) Dispatch(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	r.engine.ServeHTTP(w, req)

	return nil
}

type ginGroup struct {
	router  *ginRouter
	parent  *ginGroup
	backing *gin.RouterGroup
}

func (r *ginGroup) Handle(method string, path string, handler router.HandlerFunc) {
	r.backing.Handle(method, path, buildHandler(handler))
}

func (r *ginGroup) Group(path string, middleware ...middleware.Handler) router.Group {
	group := r.backing.Group(path, buildMiddlewares(middleware)...)

	return &ginGroup{
		router:  r.router,
		parent:  r,
		backing: group,
	}
}

func (r *ginGroup) Parent() router.Group {
	return r.parent
}

func (r *ginGroup) Router() router.Router {
	return r.router
}

func (r *ginGroup) Middleware() []middleware.Handler {
	return []middleware.Handler{}
}
