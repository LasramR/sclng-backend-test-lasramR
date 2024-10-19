package http_provider

import "context"

type HttpProvider interface {
	GetJson(ctx context.Context, url string, payload any) error
}
