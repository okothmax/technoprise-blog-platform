package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' fonts.googleapis.com; font-src 'self' fonts.gstatic.com; img-src 'self' data: https:; connect-src 'self'")
		
		// HSTS header for HTTPS
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		c.Next()
	})
}

// AccessibilityHeaders adds accessibility-related headers
func AccessibilityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Accessibility and performance headers
		c.Header("X-UA-Compatible", "IE=edge")
		c.Header("X-Accessibility-Compliant", "WCAG-2.2-AA")
		c.Header("X-Screen-Reader-Optimized", "true")
		c.Header("X-Keyboard-Navigation", "enabled")
		c.Header("X-High-Contrast-Support", "available")
		c.Header("X-Reduced-Motion-Support", "available")
		
		// Cache control for accessibility resources
		if c.Request.URL.Path == "/api/v1/accessibility" {
			c.Header("Cache-Control", "public, max-age=3600")
		}
		
		c.Next()
	})
}
