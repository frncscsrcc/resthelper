package resthelper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type ErrorResponse struct {
	Error     bool
	ErrorCode int
	Message   string
}

func SendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json, err := toJSON(ErrorResponse{true, code, message})
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	fmt.Fprintf(w, json)
}

func SendResponse(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json, err := toJSON(object)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	fmt.Fprintf(w, json)
}

func toJSON(object interface{}) (string, error) {
	json, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func GetNewToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetSessionID(r *http.Request) string {
	// Search in the context
	sessionID, ok := r.Context().Value("sessionID").(string)
	if ok && sessionID != "" {
		return sessionID
	}
	// Search in URL
	sessionIDSlice, ok := r.URL.Query()["sessionID"]
	if ok != true {
		return ""
	}
	return sessionIDSlice[0]
	// Search in body
	// TODO
}

func GetToken(r *http.Request) string {
	// Search in the context
	token, ok := r.Context().Value("token").(string)
	if ok && token != "" {
		return token
	}
	// Search in URL
	tokenSlice, ok := r.URL.Query()["token"]
	if ok != true {
		return ""
	}
	return tokenSlice[0]
	// Search in body
	// TODO
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Context().Value("sessionID")
		log.Println(r.Method, "-", sessionID, "-", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func AddSessionID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := GetSessionID(r)
		ctx := context.WithValue(r.Context(), "sessionID", sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
