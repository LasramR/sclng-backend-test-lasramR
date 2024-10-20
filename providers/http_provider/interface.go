package http_provider

import (
	"net/http"
)

type HttpProvider interface {
	ReqUnmarshalledBody(req *http.Request, unMarshalledResBody any) error
}
