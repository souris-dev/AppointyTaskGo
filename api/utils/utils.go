package utils

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

func GetHashed256(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return fmt.Sprintf("%x", hash)
}

// This function is a wrapper for checking the correct HTTP verb is used
func MakeCheckMethodHandler(method string, handlerFn func(writer http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, fmt.Sprintf("This endpoint only accepts %s requests!", method), http.StatusBadRequest)
			return
		} else {
			handlerFn(w, req)
		}
	}
}

// Function to add appropriate headers
func addHeadersUtil(headers map[string]string, w *http.ResponseWriter) {
	for key, val := range headers {
		(*w).Header().Set(key, val)
	}
}

// Function to add common headers
func AddCommonHeaders(w *http.ResponseWriter) {
	addHeadersUtil(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}, w)
}
