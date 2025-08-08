package models

import (
	"strings"
	"time"
	"github.com/jinzhu/gorm"
)

// Blog represents a blog post with accessibility features
type Blog struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Title       string    `json:"title" gorm:"not null;size:255" validate:"required,min=1,max=255"`
	Slug        string    `json:"slug" gorm:"unique;not null;size:255" validate:"required,min=1,max=255"`
	Content     string    `json:"content" gorm:"type:text" validate:"required,min=10"`
	Excerpt     string    `json:"excerpt" gorm:"size:500" validate:"max=500"`
	Author      string    `json:"author" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	Published   bool      `json:"published" gorm:"default:false"`
	Featured    bool      `json:"featured" gorm:"default:false"`
	Tags        string    `json:"tags" gorm:"size:500"` // Comma-separated tags
	MetaTitle   string    `json:"meta_title" gorm:"size:60"` // SEO meta title
	MetaDesc    string    `json:"meta_description" gorm:"size:160"` // SEO meta description
	ReadingTime int       `json:"reading_time" gorm:"default:0"` // Estimated reading time in minutes
	ViewCount   int       `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// BlogResponse represents the API response structure
type BlogResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content,omitempty"` // Only included in single blog requests
	Excerpt     string     `json:"excerpt"`
	Author      string     `json:"author"`
	Published   bool       `json:"published"`
	Featured    bool       `json:"featured"`
	Tags        []string   `json:"tags"`
	MetaTitle   string     `json:"meta_title,omitempty"`
	MetaDesc    string     `json:"meta_description,omitempty"`
	ReadingTime int        `json:"reading_time"`
	ViewCount   int        `json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// BlogListResponse represents paginated blog list response
type BlogListResponse struct {
	Blogs      []BlogResponse `json:"blogs"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
	HasNext    bool           `json:"has_next"`
	HasPrev    bool           `json:"has_prev"`
}

// CreateBlogRequest represents the request structure for creating a blog
type CreateBlogRequest struct {
	Title     string `json:"title" validate:"required,min=1,max=255"`
	Content   string `json:"content" validate:"required,min=10"`
	Excerpt   string `json:"excerpt" validate:"max=500"`
	Author    string `json:"author" validate:"required,min=1,max=100"`
	Published bool   `json:"published"`
	Featured  bool   `json:"featured"`
	Tags      string `json:"tags"`
	MetaTitle string `json:"meta_title" validate:"max=60"`
	MetaDesc  string `json:"meta_description" validate:"max=160"`
}

// UpdateBlogRequest represents the request structure for updating a blog
type UpdateBlogRequest struct {
	Title     *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content   *string `json:"content,omitempty" validate:"omitempty,min=10"`
	Excerpt   *string `json:"excerpt,omitempty" validate:"omitempty,max=500"`
	Author    *string `json:"author,omitempty" validate:"omitempty,min=1,max=100"`
	Published *bool   `json:"published,omitempty"`
	Featured  *bool   `json:"featured,omitempty"`
	Tags      *string `json:"tags,omitempty"`
	MetaTitle *string `json:"meta_title,omitempty" validate:"omitempty,max=60"`
	MetaDesc  *string `json:"meta_description,omitempty" validate:"omitempty,max=160"`
}

// BeforeCreate hook to generate slug and calculate reading time
func (b *Blog) BeforeCreate(scope *gorm.Scope) error {
	if b.Slug == "" {
		b.Slug = GenerateSlug(b.Title)
	}
	b.ReadingTime = CalculateReadingTime(b.Content)
	if b.Published && b.PublishedAt == nil {
		now := time.Now()
		b.PublishedAt = &now
	}
	return nil
}

// BeforeUpdate hook to update reading time and published date
func (b *Blog) BeforeUpdate(scope *gorm.Scope) error {
	if scope.HasColumn("content") {
		b.ReadingTime = CalculateReadingTime(b.Content)
	}
	if b.Published && b.PublishedAt == nil {
		now := time.Now()
		b.PublishedAt = &now
	} else if !b.Published {
		b.PublishedAt = nil
	}
	return nil
}

// ToResponse converts Blog to BlogResponse
func (b *Blog) ToResponse(includeContent bool) BlogResponse {
	tags := []string{}
	if b.Tags != "" {
		tags = strings.Split(b.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	response := BlogResponse{
		ID:          b.ID,
		Title:       b.Title,
		Slug:        b.Slug,
		Excerpt:     b.Excerpt,
		Author:      b.Author,
		Published:   b.Published,
		Featured:    b.Featured,
		Tags:        tags,
		ReadingTime: b.ReadingTime,
		ViewCount:   b.ViewCount,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
		PublishedAt: b.PublishedAt,
	}

	if includeContent {
		response.Content = b.Content
		response.MetaTitle = b.MetaTitle
		response.MetaDesc = b.MetaDesc
	}

	return response
}
