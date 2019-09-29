package resthelper

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type ErrorResponse struct {
	Code    int
	Message string
}


func SendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json, err := toJSON(ErrorResponse{code, message})
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
