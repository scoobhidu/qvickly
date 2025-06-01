package verify_otp

import (
	"github.com/google/uuid"
	"time"
)

type Request struct {
	Phone      string `json:"phone" binding:"required,e164"`
	Code       string `json:"code" binding:"required"`
	DeviceInfo struct {
		DeviceID   string `json:"device_id" binding:"required"`
		DeviceType string `json:"device_type"`
		AppVersion string `json:"app_version"`
		OSVersion  string `json:"os_version"`
	} `json:"device_info" binding:"required"`
}

type User struct {
	ID        uuid.UUID `db:"id"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
}
