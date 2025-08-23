package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/skip2/go-qrcode"
)

var ctx = context.Background()

// Redis client
var rdb *redis.Client

// BuyTicket handles the creation of a new QR code for a ticket
func BuyTicket(c *gin.Context) {
	// Simulated ticket ID
	ticketID := "ticket:123"

	// Check if the QR exists and its TTL
	ttl, err := rdb.TTL(ctx, ticketID).Result()
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking QR TTL"})
		return
	}

	var qrCode string
	if ttl <= 0 {
		// Generate a new QR code if expired or does not exist
		qrCode = fmt.Sprintf("QR-%d", time.Now().UnixNano())
		err := rdb.Set(ctx, ticketID, qrCode, 1*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
			return
		}
		log.Println("Generated new QR:", qrCode)
	} else {
		// Retrieve existing QR code
		qrCode, err = rdb.Get(ctx, ticketID).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching QR code"})
			return
		}
		log.Printf("Reusing existing QR: %s, TTL: %d seconds\n", qrCode, int64(ttl.Seconds()))
	}

	// Generate the QR image
	var png []byte
	png, err = qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR image"})
		return
	}

	// Send QR image as base64 (for frontend rendering)
	c.JSON(http.StatusOK, gin.H{
		"ticket_id": ticketID,
		"qr_code":   qrCode,
		"ttl":       int64(ttl.Seconds()),
		"qr_image":  fmt.Sprintf("data:image/png;base64,%s", png),
	})
}

// ValidateQR handles the validation of a QR code
func ValidateQR(c *gin.Context) {
	var request struct {
		TicketID string `json:"ticket_id"`
		QRCode   string `json:"qr_code"`
	}

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Fetch QR code from Redis
	storedQRCode, err := rdb.Get(ctx, request.TicketID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "QR code expired or invalid"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching QR code"})
		return
	}

	// Validate QR code
	if storedQRCode != request.QRCode {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid QR code"})
		return
	}

	// If valid, return success
	c.JSON(http.StatusOK, gin.H{"message": "QR code validated successfully"})
}
