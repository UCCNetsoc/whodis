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

func AccessErrorResponse(status int, message string, respErr error) (code int, resp interface{}) {
	if respErr != nil {
		log.WithError(respErr).Error(message)
	} else {
		log.Error(message)
	}
	return status, &Resp{Code: status, Type: "Error", Message: message}
}

func AccessSuccessResponse(message string, uid string, gid string, rid string) (code int, resp interface{}) {
	fields := map[string]interface{}{"user": uid, "guild": gid, "role": rid}
	log.WithFields(fields).Info(message)
	return http.StatusOK, &Resp{Code: http.StatusOK, Type: "Success", Message: message, Fields: fields}
}
