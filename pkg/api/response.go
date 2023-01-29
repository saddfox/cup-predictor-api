package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status,omitempty"`
}

// writes response with given status code and payload
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
	fmt.Println("sent response ", data)
}

// send response with given error code and message
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
