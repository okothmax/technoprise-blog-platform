package models

import (
	"math"
	"regexp"
	"strings"
	"unicode"
)

// GenerateSlug creates a URL-friendly slug from a title
func GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")
	
	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")
	
	// Limit length to 100 characters
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.Trim(slug, "-")
	}
	
	return slug
}

// CalculateReadingTime estimates reading time based on content length
// Average reading speed: 200 words per minute
func CalculateReadingTime(content string) int {
	if content == "" {
		return 0
	}
	
	// Count words (simple word count by splitting on whitespace)
	words := strings.Fields(stripHTMLTags(content))
	wordCount := len(words)
	
	// Calculate reading time (minimum 1 minute)
	readingTime := int(math.Ceil(float64(wordCount) / 200.0))
	if readingTime < 1 {
		readingTime = 1
	}
	
	return readingTime
}

// stripHTMLTags removes HTML tags from content for word counting
func stripHTMLTags(content string) string {
	// Simple HTML tag removal regex
	reg := regexp.MustCompile(`<[^>]*>`)
	return reg.ReplaceAllString(content, " ")
}

// SanitizeString removes or replaces potentially harmful characters
func SanitizeString(input string) string {
	// Remove control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			continue
		}
		result.WriteRune(r)
	}
	return strings.TrimSpace(result.String())
}

// truncateText truncates text to specified length with ellipsis
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	
	// Find the last space before maxLength to avoid cutting words
	truncated := text[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = text[:lastSpace]
	}
	
	return truncated + "..."
}

// GenerateExcerpt creates an excerpt from content if not provided
func GenerateExcerpt(content string, maxLength int) string {
	if maxLength == 0 {
		maxLength = 300
	}
	
	// Strip HTML tags and clean up
	cleaned := stripHTMLTags(content)
	cleaned = strings.ReplaceAll(cleaned, "\n", " ")
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)
	
	return truncateText(cleaned, maxLength)
}
