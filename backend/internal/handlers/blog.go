package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"technoprise-blog-backend/internal/models"
)

// BlogHandler handles blog-related HTTP requests
type BlogHandler struct {
	db *gorm.DB
}

// NewBlogHandler creates a new blog handler
func NewBlogHandler(db *gorm.DB) *BlogHandler {
	return &BlogHandler{db: db}
}

// GetBlogs handles GET /api/v1/blogs
// @Summary Get paginated list of blog posts
// @Description Retrieve blog posts with pagination and search functionality
// @Tags blogs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param featured query bool false "Filter by featured posts"
// @Param published query bool false "Filter by published posts" default(true)
// @Success 200 {object} models.BlogListResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /blogs [get]
func (h *BlogHandler) GetBlogs(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	featuredParam := c.Query("featured")
	publishedParam := c.DefaultQuery("published", "true")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Build query
	query := h.db.Model(&models.Blog{})

	// Filter by published status
	if published, err := strconv.ParseBool(publishedParam); err == nil {
		query = query.Where("published = ?", published)
	}

	// Filter by featured status
	if featuredParam != "" {
		if featured, err := strconv.ParseBool(featuredParam); err == nil {
			query = query.Where("featured = ?", featured)
		}
	}

	// Search functionality
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(excerpt) LIKE ? OR LOWER(tags) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count blogs",
		})
		return
	}

	// Calculate pagination
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Fetch blogs
	var blogs []models.Blog
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch blogs",
		})
		return
	}

	// Convert to response format
	blogResponses := make([]models.BlogResponse, len(blogs))
	for i, blog := range blogs {
		blogResponses[i] = blog.ToResponse(false) // Don't include full content in list
	}

	// Prepare response
	response := models.BlogListResponse{
		Blogs:      blogResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	// Set accessibility headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(limit))

	c.JSON(http.StatusOK, response)
}

// GetBlogBySlug handles GET /api/v1/blogs/:slug
// @Summary Get a single blog post by slug
// @Description Retrieve a blog post by its slug and increment view count
// @Tags blogs
// @Accept json
// @Produce json
// @Param slug path string true "Blog slug"
// @Success 200 {object} models.BlogResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /blogs/{slug} [get]
func (h *BlogHandler) GetBlogBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var blog models.Blog
	if err := h.db.Where("slug = ? AND published = ?", slug, true).First(&blog).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Blog post not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch blog post",
		})
		return
	}

	// Increment view count
	if err := h.db.Model(&blog).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		// Log error but don't fail the request
		// log.Printf("Failed to increment view count for blog %d: %v", blog.ID, err)
	}

	// Set SEO and accessibility headers
	c.Header("X-Meta-Title", blog.MetaTitle)
	c.Header("X-Meta-Description", blog.MetaDesc)
	c.Header("X-Reading-Time", strconv.Itoa(blog.ReadingTime))

	response := blog.ToResponse(true) // Include full content for single blog view
	c.JSON(http.StatusOK, response)
}

// CreateBlog handles POST /api/v1/blogs
// @Summary Create a new blog post
// @Description Create a new blog post with accessibility validation
// @Tags blogs
// @Accept json
// @Produce json
// @Param blog body models.CreateBlogRequest true "Blog data"
// @Success 201 {object} models.BlogResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /blogs [post]
func (h *BlogHandler) CreateBlog(c *gin.Context) {
	var req models.CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Generate slug if not provided
	slug := models.GenerateSlug(req.Title)
	
	// Check if slug already exists
	var existingBlog models.Blog
	if !h.db.Where("slug = ?", slug).First(&existingBlog).RecordNotFound() {
		// Append timestamp to make slug unique
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	// Generate excerpt if not provided
	excerpt := req.Excerpt
	if excerpt == "" {
		excerpt = models.GenerateExcerpt(req.Content, 300)
	}

	// Create blog post
	blog := models.Blog{
		Title:     models.SanitizeString(req.Title),
		Slug:      slug,
		Content:   models.SanitizeString(req.Content),
		Excerpt:   models.SanitizeString(excerpt),
		Author:    models.SanitizeString(req.Author),
		Published: req.Published,
		Featured:  req.Featured,
		Tags:      models.SanitizeString(req.Tags),
		MetaTitle: models.SanitizeString(req.MetaTitle),
		MetaDesc:  models.SanitizeString(req.MetaDesc),
	}

	if err := h.db.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create blog post",
		})
		return
	}

	response := blog.ToResponse(true)
	c.JSON(http.StatusCreated, response)
}

// UpdateBlog handles PUT /api/v1/blogs/:id
// @Summary Update a blog post
// @Description Update an existing blog post
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path int true "Blog ID"
// @Param blog body models.UpdateBlogRequest true "Updated blog data"
// @Success 200 {object} models.BlogResponse
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /blogs/{id} [put]
func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid blog ID",
		})
		return
	}

	var req models.UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	var blog models.Blog
	if err := h.db.First(&blog, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Blog post not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch blog post",
		})
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	
	if req.Title != nil {
		updates["title"] = models.SanitizeString(*req.Title)
		// Regenerate slug if title changed
		updates["slug"] = models.GenerateSlug(*req.Title)
	}
	if req.Content != nil {
		updates["content"] = models.SanitizeString(*req.Content)
	}
	if req.Excerpt != nil {
		updates["excerpt"] = models.SanitizeString(*req.Excerpt)
	}
	if req.Author != nil {
		updates["author"] = models.SanitizeString(*req.Author)
	}
	if req.Published != nil {
		updates["published"] = *req.Published
	}
	if req.Featured != nil {
		updates["featured"] = *req.Featured
	}
	if req.Tags != nil {
		updates["tags"] = models.SanitizeString(*req.Tags)
	}
	if req.MetaTitle != nil {
		updates["meta_title"] = models.SanitizeString(*req.MetaTitle)
	}
	if req.MetaDesc != nil {
		updates["meta_desc"] = models.SanitizeString(*req.MetaDesc)
	}

	if err := h.db.Model(&blog).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update blog post",
		})
		return
	}

	// Fetch updated blog
	if err := h.db.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch updated blog post",
		})
		return
	}

	response := blog.ToResponse(true)
	c.JSON(http.StatusOK, response)
}

// DeleteBlog handles DELETE /api/v1/blogs/:id
// @Summary Delete a blog post
// @Description Delete a blog post by ID
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path int true "Blog ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /blogs/{id} [delete]
func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid blog ID",
		})
		return
	}

	var blog models.Blog
	if err := h.db.First(&blog, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Blog post not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch blog post",
		})
		return
	}

	if err := h.db.Delete(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete blog post",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
