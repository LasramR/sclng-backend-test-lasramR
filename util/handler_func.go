package util

import "net/http"

type ScalingoHandlerFunc func(w http.ResponseWriter, r *http.Request, _ map[string]string) error
