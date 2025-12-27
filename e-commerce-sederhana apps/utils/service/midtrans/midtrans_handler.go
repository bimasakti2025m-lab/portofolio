package midtrans

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MidtransHandler struct {
	service MidtransService
}

func NewMidtransHandler(midtransService MidtransService) *MidtransHandler {
	return &MidtransHandler{
		service: midtransService,
	}
}

// HandleNotification adalah endpoint yang akan menerima notifikasi dari Midtrans
func (h *MidtransHandler) HandleNotification(c *gin.Context) {
	var notificationPayload map[string]interface{}

	// 1. Bind JSON body ke map
	if err := c.ShouldBindJSON(&notificationPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification payload"})
		return
	}

	// 2. Panggil service untuk memproses notifikasi
	err := h.service.HandleNotification(notificationPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Kirim respons 200 OK ke Midtrans untuk menandakan notifikasi berhasil diterima
	c.JSON(http.StatusOK, gin.H{"message": "Notification processed successfully"})
}
