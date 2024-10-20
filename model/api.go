package model

import "github.com/LasramR/sclng-backend-test-lasramR/util"

type ApiResponse[T any] struct {
	Count            int                            `json:"count"`
	Content          T                              `json:"content"`
	IncompleteResult bool                           `json:"incomplete_result"`
	Page             int                            `json:"page,omitempty"`
	Previous         util.NullableJsonField[string] `json:"previous,omitempty"`
	Next             string                         `json:"next,omitempty"`
}

type ApiError[T any] struct {
	Status int `json:"status"`
	Reason T   `json:"reasons"`
}
