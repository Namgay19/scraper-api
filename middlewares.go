package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type customResponse struct {
	http.ResponseWriter
	statusCode int
}

func (res *customResponse) WriteHeader(code int) {
    res.statusCode = code
    res.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		checkError(err, w)
		defer logFile.Close()

		requestEntry := fmt.Sprintf("Host=%v Request=%v Method=%v RequestTime=%v", r.Host, r.RequestURI, r.Method, time.Now())
		logFile.Write([]byte(requestEntry))
		newResponse := customResponse{w, http.StatusOK}

		next.ServeHTTP(&newResponse, r)
		
		responseEntry := fmt.Sprintf("Response Status=%v \n", newResponse.statusCode)
		logFile.Write([]byte(responseEntry))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)	
	})
}
