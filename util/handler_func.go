package util

import "net/http"

// Handler type for Scalingo router.HandlerFunc
type ScalingoHandlerFunc func(w http.ResponseWriter, r *http.Request, _ map[string]string) error
