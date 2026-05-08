package responseHelper

import (
	"encoding/json"
	"net/http"
	"sales-manager-back/pkg/domain/response"
)

func ResponseBuilder(status int, message string, data interface{}) ([]byte, error) {
	resp := response.Response{
		Message: message,
		Data:    data,
	}

	marshalResponse, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return marshalResponse, nil
}

func ResponseStatusChecker(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		return
	}
}

func WriteResponse(w http.ResponseWriter, status response.Status, data interface{}) {
	resp := response.Response{
		Message: status.StatusText,
		Data:    data,
	}

	marshalResponse, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ResponseStatusChecker(w, []byte("500: Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.StatusCode)
	ResponseStatusChecker(w, marshalResponse)
}
