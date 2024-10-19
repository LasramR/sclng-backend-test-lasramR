package http_provider

type HttpProvider interface {
	GetJson(url string, payload any) error
}
