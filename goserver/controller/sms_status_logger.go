package controller

import (
	"dt/logwrap"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func smsStatus(c *gin.Context) {
	type SmsStatus struct {
		State     string `json:"state"`
		Action    string `json:"action"`
		SmsID     string `json:"sms_id"`
		UserID    string `json:"user_id"`
		Price     string `json:"price"`
		ErrorCode string `json:"error_code"`
		UpdatedAt string `json:"updated_at"`
		Key       string `json:"key"`
		Parts     string `json:"parts"`
		Secret    string `json:"secret"`
	}

	status := SmsStatus{
		State:     c.Query("state"),
		Action:    c.Query("action"),
		SmsID:     c.Query("sms"),
		UserID:    c.Query("id"),
		Price:     c.Query("price"),
		ErrorCode: c.Query("error_code"),
		UpdatedAt: c.Query("updated_at"),
		Key:       c.Query("key"),
		Parts:     c.Query("parts"),
		Secret:    c.Query("secret"),
	}

	bytes, _ := json.Marshal(status)
	logwrap.Debug("sms status: %s", string(bytes))
	c.Status(http.StatusOK)
}
