package verify_otp

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"net/http"
	"os"
	"qvickly/env"
	"time"
)

// rate limit using mobile number

// VerifySMSOTP godoc
// @Summary Verify OTP and Authenticate User
// @Description Verify the OTP code sent via SMS or WhatsApp, create user account if new, and return JWT access token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body Request true "OTP verification request containing phone number and code"
// @Router /api/auth/verify-otp [post]
func VerifySMSOTP(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	client := twilio.NewRestClient()

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(req.Phone)
	params.SetCode(req.Code)

	resp, err := client.VerifyV2.CreateVerificationCheck(env.TwilioOTPVerifierServiceSid, params)
	if err != nil || *resp.Status != "approved" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Create or find user
	var user User
	err = db.Get(&user, "SELECT * FROM user_profile.users WHERE phone_number = $1", req.Phone)
	if errors.Is(err, sql.ErrNoRows) {
		// Create new user
		user = User{ID: uuid.New(), Phone: req.Phone, CreatedAt: time.Now()}
		_, err = db.NamedExec(`INSERT INTO user_profile.users (id, phone_number, created_at) VALUES (:id, :phone, :created_at)`, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"user":         gin.H{"id": user.ID, "phone": user.Phone},
	})
}

func generateAccessToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
