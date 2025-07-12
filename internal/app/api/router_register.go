package api

import "github.com/gin-gonic/gin"

type RouteGroup interface {
	GET(relativePath string, handlers ...gin.HandlerFunc)
	POST(relativePath string, handlers ...gin.HandlerFunc)
	PUT(relativePath string, handlers ...gin.HandlerFunc)
	PATCH(relativePath string, handlers ...gin.HandlerFunc)
	DELETE(relativePath string, handlers ...gin.HandlerFunc)
	Group(relativePath string, handlers ...gin.HandlerFunc) RouteGroup
}

type GinRouteGroup struct {
	group *gin.RouterGroup
}

func newGinRouteGroup(g *gin.RouterGroup) *GinRouteGroup {
	return &GinRouteGroup{group: g}
}

func (g *GinRouteGroup) GET(relativePath string, handlers ...gin.HandlerFunc) {
	g.group.GET(relativePath, handlers...)
}

func (g *GinRouteGroup) POST(relativePath string, handlers ...gin.HandlerFunc) {
	g.group.POST(relativePath, handlers...)
}

func (g *GinRouteGroup) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	g.group.PUT(relativePath, handlers...)
}

func (g *GinRouteGroup) PATCH(relativePath string, handlers ...gin.HandlerFunc) {
	g.group.PATCH(relativePath, handlers...)
}

func (g *GinRouteGroup) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	g.group.DELETE(relativePath, handlers...)
}

func (g *GinRouteGroup) Group(relativePath string, handlers ...gin.HandlerFunc) RouteGroup {
	return newGinRouteGroup(g.group.Group(relativePath, handlers...))
}

type GinRouterRegister interface {
	WithGroup(path string, handlers ...gin.HandlerFunc) RouteGroup
	NoGroup(handlers ...gin.HandlerFunc) RouteGroup
}

type GinRouterRegisterImpl struct {
	r gin.IRouter
}

var _ GinRouterRegister = (*GinRouterRegisterImpl)(nil)

func NewGinRouterRegisterImpl(r gin.IRouter) *GinRouterRegisterImpl {
	return &GinRouterRegisterImpl{r: r}
}

func (g *GinRouterRegisterImpl) WithGroup(path string, handlers ...gin.HandlerFunc) RouteGroup {
	grp := g.r.Group(path, handlers...)
	return newGinRouteGroup(grp)
}

func (g *GinRouterRegisterImpl) NoGroup(handlers ...gin.HandlerFunc) RouteGroup {
	g.r.Use(handlers...)
	rootGroup := g.r.Group("")
	return newGinRouteGroup(rootGroup)
}

type GinController interface {
	Register(r GinRouterRegister) error
}
