package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"

	"mbsms-api/app/model"
	"mbsms-api/app/service"
	"mbsms-api/app/util"
)

// PostMessage handles post message requests
func PostMessage(sms *service.SMSService) rest.HandlerFunc {
	fn := func(w rest.ResponseWriter, r *rest.Request) {
		payload := &model.BaseMessage{}
		err := r.DecodeJsonPayload(payload)
		if err != nil {
			log.Print(err.Error())
			rest.Error(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}

		if len(payload.Recipient) == 0 {
			rest.Error(w, "MISSING_ARG_RECIPIENT", http.StatusBadRequest)
			return
		}

		if _, err := strconv.ParseInt(payload.Recipient, 10, 64); err != nil {
			rest.Error(w, "INVALID_ARG_RECIPIENT", http.StatusBadRequest)
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

		// send text message through given SMS service
		respChan := sms.Send(payload)
		response := <-respChan

		w.WriteJson(&response)
	}

	return rest.HandlerFunc(fn)
}
