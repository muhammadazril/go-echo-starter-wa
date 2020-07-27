package client

import (
	"context"
)

// Usecase : Use case for Merchant and MerchantUser
type Usecase interface {
	PostStreamRegisterClient(c context.Context, msg string) error
	// RegisterClient(c context.Context, cid string, key string) (*models.Client, error)
	// DeactivateClient(c context.Context, cid string) error
	// ActivateClient(c context.Context, cid string) error
}
