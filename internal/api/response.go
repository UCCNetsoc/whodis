package api

import (
	"net/http"

	"github.com/Strum355/log"
)

type Resp struct {
	Code    int                    `json:"code"`
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

func AccessErrorResponse(status int, message string, respErr error) *Resp {
	if respErr != nil {
		log.WithError(respErr).Error(message)
	} else {
		log.Error(message)
	}
	return &Resp{Code: status, Type: "Error", Message: message}
}

func AccessSuccessResponse(message string, uid string, gid string) *Resp {
	fields := map[string]interface{}{"user": uid, "guild": gid}
	log.WithFields(fields).Info(message)
	return &Resp{Code: http.StatusOK, Type: "Success", Message: message, Fields: fields}
}
