package api

import (
	"html/template"
	"net/http"

	"github.com/Strum355/log"
)

type Resp struct {
	Code    int                    `json:"code"`
	Type    template.HTML          `json:"type"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

const (
	successType = `<span style='color: limegreen'>Success</span>`
	errorType   = `<span style='color: red'>Error</span>`
)

func AccessErrorResponse(status int, message string, respErr error) *Resp {
	if respErr != nil {
		log.WithError(respErr).Error(message)
	} else {
		log.Error(message)
	}
	return &Resp{Code: status, Type: errorType, Message: message}
}

func AccessSuccessResponse(message string, uid string, gid string) *Resp {
	fields := map[string]interface{}{"user": uid, "guild": gid}
	log.WithFields(fields).Info(message)
	return &Resp{Code: http.StatusOK, Type: successType, Message: message, Fields: fields}
}
