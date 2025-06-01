package send_otp

import (
	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"net/http"
	"qvickly/env"
)

// SendWhatsappOTPController godoc
// @Summary Send OTP via WhatsApp
// @Description Send a one-time password (OTP) to the provided phone number via WhatsApp using Twilio Verify service
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body Request true "Phone number request"
// @Router /api/auth/otp/whatsapp [post]
func SendWhatsappOTPController(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone format"})
		return
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioAccountSid,
		Password: env.TwilioAuthToken,
	})

	params := &verify.CreateVerificationParams{}
	params.SetTo(req.Phone)

	params.SetChannel("whatsapp")
	_, err := client.VerifyV2.CreateVerification(env.TwilioOTPVerifierServiceSid, params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP via WhatsApp"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}

// rate limit using mobile number

// SendSmsOTPController godoc
// @Summary Send OTP via SMS
// @Description Send a one-time password (OTP) to the provided phone number via SMS using Twilio Verify service
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body Request true "Phone number request"
// @Router /api/auth/otp/sms [post]
func SendSmsOTPController(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone format"})
		return
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioAccountSid,
		Password: env.TwilioAuthToken,
	})

	params := &verify.CreateVerificationParams{}
	params.SetTo(req.Phone)

	params.SetChannel("sms")
	_, err := client.VerifyV2.CreateVerification(env.TwilioOTPVerifierServiceSid, params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP via SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}
