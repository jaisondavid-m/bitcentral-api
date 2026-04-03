package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"server/config"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	DB *sql.DB
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		DB: config.DB,
	}
}
func (h *AdminHandler) GetUsers(c *gin.Context) {
	client, _ := config.FirebaseApp.Auth(context.Background())

	iter := client.Users(context.Background(), "")

	var users []models.User

	for {
		u, err := iter.Next()
		if err != nil {
			break
		}

		users = append(users, models.User{
			UID:            u.UID,
			Email:          u.Email,
			DisplayName:    u.DisplayName,
			PhotoURL:       u.PhotoURL,
			CreationTime:   utils.TsToString(u.UserMetadata.CreationTimestamp),
			LastSignInTime: utils.TsToString(u.UserMetadata.LastLogInTimestamp),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
}

// DELETE USER
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	uid := c.Param("uid")

	client, _ := config.FirebaseApp.Auth(context.Background())
	client.DeleteUser(context.Background(), uid)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) UpdatePSToken(c *gin.Context) {
	var body models.PSToken

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false,"message":"invalid payload"})
		return
	}

	query := `
	INSERT INTO ps_tokens (token_key, token, updated_at, updated_by)
	VALUES ('ps_token', ?, ?, ?)
	ON DUPLICATE KEY UPDATE token=?, updated_at=?`

	now := time.Now()

	_, err := h.DB.Exec(query, body.Token, now, "admin", body.Token, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false , "err":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token updated",
	})
}