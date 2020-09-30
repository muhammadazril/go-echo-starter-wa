package gowa

import (
	"context"
)

// Repository for Whatsapp
type Repository interface {
	SendMessage(c context.Context, number string, message string) (string, error)
}
