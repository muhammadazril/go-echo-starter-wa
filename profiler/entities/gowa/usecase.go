package gowa

import (
	"context"
)

// Usecase : Use case for Whatsapp
type Usecase interface {
	SendMessage(c context.Context, number string, message string) (string, error)
}
