package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"belcamp/internal/models"
	"belcamp/internal/service"

	"github.com/gin-gonic/gin"
)

// ProductHandler is a handler for products
type ProductHandler struct {
	BaseHandler
	productService service.ProductService
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(ps service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: ps,
	}
}

// List shows the product listing page
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Get products with pagination
	products, total, err := h.productService.ListProducts(page, pageSize, nil)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// If it's an HTMX request, return only the table
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "products/partials/table.html", gin.H{
			"products": products,
			"page":     page,
			"total":    total,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		})
		return
	}

	// Otherwise return the full page
	h.Render(c, "products/index.html", gin.H{
		"title":    "Products",
		"products": products,
		"page":     page,
		"total":    total,
		"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// Show displays a single product
func (h *ProductHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, err)
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.Render(c, "products/show.html", gin.H{
		"title":   product.Name,
		"product": product,
	})
}

// ShowCreateForm displays the product creation form
func (h *ProductHandler) ShowCreateForm(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "products/partials/form.html", gin.H{
			"action": "/products",
			"method": "POST",
		})
		return
	}

	h.Render(c, "products/create.html", gin.H{
		"title":  "Create Product",
		"action": "/products",
		"method": "POST",
	})
}

// Create handles product creation
func (h *ProductHandler) Create(c *gin.Context) {
	var product models.Product

	// Handle form data
	product.Name = c.PostForm("name")
	product.ShortDescription = c.PostForm("short_description")
	product.Description = c.PostForm("description")
	product.Status = c.PostForm("status") == "true"

	// Handle file uploads
	datasheet, err := c.FormFile("datasheet")
	if err == nil {
		// Handle datasheet upload
		filename := generateFilename(datasheet.Filename)
		if err := c.SaveUploadedFile(datasheet, "uploads/datasheets/"+filename); err != nil {
			h.handleError(c, err)
			return
		}
		product.Datasheet = filename
	}

	// Handle multiple photo uploads
	form, _ := c.MultipartForm()
	photos := form.File["photos"]
	if len(photos) > 0 {
		var photoNames []string
		for _, photo := range photos {
			filename := generateFilename(photo.Filename)
			if err := c.SaveUploadedFile(photo, "uploads/products/"+filename); err != nil {
				h.handleError(c, err)
				return
			}
			photoNames = append(photoNames, filename)
		}
		photosJSON, _ := json.Marshal(photoNames)
		product.Photos = string(photosJSON)
	}

	// Create the product
	if err := h.productService.CreateProduct(&product); err != nil {
		h.handleError(c, err)
		return
	}

	// Handle the response based on request type
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/products")
		return
	}

	c.Redirect(http.StatusFound, "/products")
}

// ShowEditForm displays the product edit form
func (h *ProductHandler) ShowEditForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, err)
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "products/partials/form.html", gin.H{
			"product": product,
			"action":  "/products/" + strconv.FormatUint(id, 10),
			"method":  "PUT",
		})
		return
	}

	h.Render(c, "products/edit.html", gin.H{
		"title":   "Edit " + product.Name,
		"product": product,
		"action":  "/products/" + strconv.FormatUint(id, 10),
		"method":  "PUT",
	})
}

// Update handles product updates
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, err)
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Update fields
	product.Name = c.PostForm("name")
	product.ShortDescription = c.PostForm("short_description")
	product.Description = c.PostForm("description")
	product.Status = c.PostForm("status") == "true"

	// Handle file updates
	datasheet, err := c.FormFile("datasheet")
	if err == nil {
		filename := generateFilename(datasheet.Filename)
		if err := c.SaveUploadedFile(datasheet, "uploads/datasheets/"+filename); err != nil {
			h.handleError(c, err)
			return
		}
		product.Datasheet = filename
	}

	// Update the product
	if err := h.productService.UpdateProduct(product); err != nil {
		h.handleError(c, err)
		return
	}

	// Handle the response based on request type
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/products")
		return
	}

	c.Redirect(http.StatusFound, "/products")
}

// Delete handles product deletion
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, err)
		return
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		h.handleError(c, err)
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		// Return success message or updated table
		c.HTML(http.StatusOK, "products/partials/table.html", gin.H{
			"message": "Product deleted successfully",
		})
		return
	}

	c.Redirect(http.StatusFound, "/products")
}

// Variants handles product variant management
func (h *ProductHandler) Variants(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, err)
		return
	}

	variants, err := h.productService.ListVariants(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "products/partials/variants.html", gin.H{
			"variants": variants,
		})
		return
	}

	h.Render(c, "products/variants.html", gin.H{
		"title":    "Product Variants",
		"variants": variants,
	})
}

// Search handles product search
func (h *ProductHandler) Search(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	products, total, err := h.productService.SearchProducts(query, page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "products/partials/table.html", gin.H{
			"products": products,
			"total":    total,
		})
		return
	}

	h.Render(c, "products/index.html", gin.H{
		"title":    "Search Results",
		"products": products,
		"total":    total,
		"query":    query,
	})
}

// handleError handles error responses
func (h *ProductHandler) handleError(c *gin.Context, err error) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusBadRequest, "shared/partials/error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	h.Render(c, "shared/error.html", gin.H{
		"title": "Error",
		"error": err.Error(),
	})
}

// Helper function to generate unique filenames
func generateFilename(original string) string {
	// Implement your filename generation logic here
	return original
}
