package repository

import (
	"github.com/rimantoro/event_driven/profiler/entities/client"
)

type clientRepository struct {
	// DB         dbHelper.Database
	// Collection dbHelper.Collection
}

func NewClientRepository() client.StreamRepository {
	return &clientRepository{}
}
