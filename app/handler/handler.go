package handler

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"mbsms-api/app/model"
	"mbsms-api/app/service"
	"mbsms-api/app/util"
)

// PostMessage send provided text through a SMS gateway
func PostMessage(SMSSender service.Sender, w rest.ResponseWriter, r *rest.Request) {
	payload := &model.BaseMessage{}
	err := r.DecodeJsonPayload(payload)
	if err != nil {
		rest.Error(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	if payload.Recipient == 0 {
		rest.Error(w, "MISSING_ARG_RECIPIENT", http.StatusBadRequest)
		return
	}

	if payload.Originator == "" {
		rest.Error(w, "MISSING_ARG_ORIGINATOR", http.StatusBadRequest)
		return
	}
	if len(payload.Originator) > util.OriginatorMaxLen {
		rest.Error(w, "INVALID_ARG_ORIGINATOR", http.StatusBadRequest)
		return
	}

	if payload.Message == "" {
		rest.Error(w, "MISSING_ARG_MESSAGE", http.StatusBadRequest)
		return
	}

	respChan := SMSSender.Send(payload)
	response := <-respChan

	w.WriteJson(&response)
}
