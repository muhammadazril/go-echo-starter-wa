package model

var (
	Topic string = "profiler.client.crud"
)

type Client struct {
	ID     string `json:"id" bson:"_id" validate:"required,unique"`
	Key    string `json:"key" bson:"key" validate:"required,min=32"`
	Status int16  `json:"status" bson:"status" validate:"oneof=0 1"`
}

type CrudPayload struct {
	ProcName string      `json:"procname" validate:"required,oneof=create update delete"`
	Payload  interface{} `json:"payload" validate:"required"`
}
