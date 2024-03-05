// This file is auto-generated *DO NOT EDIT*

package router

import "github.com/gin-gonic/gin"

func (w *WrappedRouterGroup[T]) GET(
	relativePath string,
	handlers ...T,
) gin.IRoutes {
	return w.RouterGroup.GET(relativePath, w.wrapHandlers(handlers)...)
}

func (w *WrappedRouterGroup[T]) POST(
	relativePath string,
	handlers ...T,
) gin.IRoutes {
	return w.RouterGroup.POST(relativePath, w.wrapHandlers(handlers)...)
}

func (w *WrappedRouterGroup[T]) PUT(
	relativePath string,
	handlers ...T,
) gin.IRoutes {
	return w.RouterGroup.PUT(relativePath, w.wrapHandlers(handlers)...)
}

func (w *WrappedRouterGroup[T]) PATCH(
	relativePath string,
	handlers ...T,
) gin.IRoutes {
	return w.RouterGroup.PATCH(relativePath, w.wrapHandlers(handlers)...)
}

func (w *WrappedRouterGroup[T]) OPTIONS(
	relativePath string,
	handlers ...T,
) gin.IRoutes {
	return w.RouterGroup.OPTIONS(relativePath, w.wrapHandlers(handlers)...)
}

func (w *WrappedRouterGroup[T]) DELETE(
	relativePath string,
	handlers ...T,
) gin.IRoutes {
	return w.RouterGroup.DELETE(relativePath, w.wrapHandlers(handlers)...)
}
