package handlers

import (
	"belcamp/internal/domain/valueobject"
	"belcamp/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CRUDHandler[T any] struct {
	service *service.CRUDService[T]
	tmpl    string // Base template name for the entity
	BaseHandler
}

func NewCRUDHandler[T any](service *service.CRUDService[T], tmpl string) *CRUDHandler[T] {
	return &CRUDHandler[T]{
		service: service,
		tmpl:    tmpl,
	}
}

// RegisterRoute registers a route with the given method and handler
func (h *CRUDHandler[T]) RegisterRoute(r *gin.RouterGroup, path string, method string, handler func(c *gin.Context)) {
	switch method {
	case "GET":
		r.GET(path, handler)
	case "POST":
		r.POST(path, handler)
	case "PUT":
		r.PUT(path, handler)
	case "DELETE":
		r.DELETE(path, handler)
	}
}

func (h *CRUDHandler[T]) RegisterDefaultRoutes(r *gin.RouterGroup, path string) {
	group := r.Group(path)
	group.GET("", h.SmartTableList)
	group.GET("/:id", h.Get)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *CRUDHandler[T]) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	pagination := valueobject.NewPagination(page, pageSize)
	entities, pagination, err := h.service.List(c.Request.Context(), pagination)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": err.Error()})
		return
	}

	h.Render(c, h.tmpl+".index", gin.H{
		"entities":   entities,
		"pagination": pagination,
	}, h.tmpl+".table")
}

func (h *CRUDHandler[T]) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{"error": "Invalid ID"})
		return
	}

	entity, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error", gin.H{"error": err.Error()})
		return
	}

	h.Render(c, h.tmpl+".index", gin.H{"entity": entity}, h.tmpl+".show")
}

func (h *CRUDHandler[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBind(&entity); err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(c.Request.Context(), &entity); err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": err.Error()})
		return
	}

	h.Redirect(c, c.Request.URL.Path)
}

func (h *CRUDHandler[T]) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{"error": "Invalid ID"})
		return
	}

	existingEntity, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.", gin.H{"error": "Entity not found"})
		return
	}

	if err := c.ShouldBind(existingEntity); err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), existingEntity); err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": err.Error()})
		return
	}

	h.Redirect(c, c.Request.URL.Path)

	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, h.tmpl+"/detail_partial.html", gin.H{"entity": existingEntity})
		return
	}

	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
}

func (h *CRUDHandler[T]) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": err.Error()})
		return
	}

	h.Redirect(c, c.Request.URL.Path)
}
