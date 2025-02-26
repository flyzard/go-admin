package handlers

import (
	"belcamp/internal/domain/entity"
	"belcamp/internal/service"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProductHandler extends the generic CRUD handler with product-specific functionality
type ProductHandler struct {
	*CRUDHandler[entity.Product]
	uploadDir string
}

// NewProductHandler creates a new product handler
func NewProductHandler(service *service.CRUDService[entity.Product], uploadDir string) *ProductHandler {
	return &ProductHandler{
		CRUDHandler: NewCRUDHandler(service, "products"),
		uploadDir:   uploadDir,
	}
}

// Override the Get method to use your custom template structure
func (h *ProductHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "new" {
		// Render the new product form
		h.Render(c, "products.edit", gin.H{
			"entity": &entity.Product{},
			"isNew":  true,
		}, "")
		return
	}

	// Use the parent Get method for existing products
	h.CRUDHandler.Get(c)
}

// Override the Update method to handle file uploads
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Get the existing product
	var existingProduct *entity.Product
	var err error

	if id != "new" {
		existingProduct, err = h.service.Get(c.Request.Context(), convertToUint(id))
		if err != nil {
			c.HTML(http.StatusNotFound, "error", gin.H{"error": "Product not found"})
			return
		}
	} else {
		existingProduct = &entity.Product{}
	}

	// Get the basic product data
	if err := c.ShouldBind(existingProduct); err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{"error": err.Error()})
		return
	}

	// Handle datasheet upload
	datasheetFilename, err := h.handleFileUpload(c, "datasheet", *existingProduct.Datasheet, []string{".pdf"})
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{"error": fmt.Sprintf("Error uploading datasheet: %v", err)})
		return
	}

	// Update the product with the new datasheet
	existingProduct.Datasheet = &datasheetFilename

	// Save the product
	var saveErr error
	if id == "new" {
		saveErr = h.service.Create(c.Request.Context(), existingProduct)
	} else {
		saveErr = h.service.Update(c.Request.Context(), existingProduct)
	}

	if saveErr != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": saveErr.Error()})
		return
	}

	// Redirect to the product list
	h.Redirect(c, "/products")
}

// handleFileUpload processes a file upload for the given field
func (h *ProductHandler) handleFileUpload(c *gin.Context, fieldName, existingFile string, allowedExts []string) (string, error) {
	// Check if there's a request to remove the existing file
	if c.PostForm(fmt.Sprintf("remove_%s", fieldName)) == "true" {
		// Remove the existing file if it exists
		if existingFile != "" {
			os.Remove(filepath.Join(h.uploadDir, existingFile))
		}
		return "", nil
	}

	// Check if a new file is being uploaded
	file, err := c.FormFile(fieldName)
	if err != nil {
		// If no new file and the request isn't to remove, keep the existing file
		if existingFile != "" && c.PostForm(fmt.Sprintf("existing_%s", fieldName)) != "" {
			return existingFile, nil
		}
		// No existing file and no new file
		return "", nil
	}

	// Validate file size (5MB max)
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("file size exceeds 5MB limit")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExt := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			validExt = true
			break
		}
	}

	if !validExt {
		return "", fmt.Errorf("invalid file type, only %s are allowed", strings.Join(allowedExts, ", "))
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s-%s%s", uuid.New().String(), sanitizeFilename(file.Filename), ext)
	filePath := filepath.Join(h.uploadDir, filename)

	// Ensure upload directory exists
	os.MkdirAll(h.uploadDir, 0755)

	// Save the file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	// Remove old file if it exists and a new one was uploaded
	if existingFile != "" {
		os.Remove(filepath.Join(h.uploadDir, existingFile))
	}

	return filename, nil
}

// sanitizeFilename removes potentially dangerous characters from a filename
func sanitizeFilename(filename string) string {
	// Remove path information
	filename = filepath.Base(filename)
	// Replace problematic characters
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.", r) {
			return r
		}
		return '-'
	}, filename)
}

// Helper function to convert string ID to uint
func convertToUint(id string) uint {
	var result uint
	fmt.Sscanf(id, "%d", &result)
	return result
}
