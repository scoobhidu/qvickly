package send_otp

type Request struct {
	Phone string `json:"phone" binding:"required,e164" example:"+918010201921"` // E.164 format: +1234567890
}
