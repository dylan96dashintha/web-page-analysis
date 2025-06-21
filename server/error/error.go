package error

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Msg struct {
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
	Code             int64  `json:"code"`
}

type ErrorMsg struct {
	Error Msg `json:"error"`
}

func BadRequestError(developerMessage string, w http.ResponseWriter) {

	errorMessage := Msg{
		Message:          "bad request",
		DeveloperMessage: developerMessage,
		Code:             400,
	}

	errMsg := ErrorMsg{Error: errorMessage}
	data, err := json.Marshal(errMsg)
	if err != nil {
		log.Errorf("ERROR in marshalling message, err: %+v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(data)
}
