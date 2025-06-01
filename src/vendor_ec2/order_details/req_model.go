package order_details

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
}
