package handler

import (
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"

	"mbsms-api/app/model"
	"mbsms-api/app/service"
)

// FakeSMSProvider for testing purpose
type FakeSMSProvider struct {
	response service.Response
}

func (f *FakeSMSProvider) Balance() {}
func (f *FakeSMSProvider) Send(bm *model.BaseMessage) <-chan service.Response {
	respChan := make(chan service.Response, 1)

	respChan <- f.response
	return respChan
}

func TestPostMessage(t *testing.T) {
	// define test cases for message post
	testCases := []struct {
		description        string
		smsProvider        service.SMSProvider
		url                string
		payload            interface{}
		expectedStatusCode int
		expectedBody       string
	}{
		{
			description: "succesfull post",
			smsProvider: &FakeSMSProvider{
				response: map[string]interface{}{
					"Status":         "Success",
					"TotalSentParts": 1,
				},
			},
			url: "/messages",
			payload: map[string]interface{}{
				"recipient":  "31612345678",
				"originator": "MessageBird",
				"message":    "This is a test message.",
			},
			expectedStatusCode: 200,
			expectedBody:       `{"Status":"Success","TotalSentParts":1}`,
		},
		{
			description:        "missing payload",
			smsProvider:        &FakeSMSProvider{},
			url:                "/messages",
			payload:            nil,
			expectedStatusCode: 500,
			expectedBody:       `{"Error":"INTERNAL_ERROR"}`,
		},
		{
			description: "missing argument recipient",
			smsProvider: &FakeSMSProvider{},
			url:         "/messages",
			payload: map[string]interface{}{
				"originator": "MessageBird",
				"message":    "This is a test message.",
			},
			expectedStatusCode: 400,
			expectedBody:       `{"Error":"MISSING_ARG_RECIPIENT"}`,
		},
		{
			description: "invalid argument recipient",
			smsProvider: &FakeSMSProvider{},
			url:         "/messages",
			payload: map[string]interface{}{
				"recipient":  "##############",
				"originator": "MessageBird",
				"message":    "This is a test message.",
			},
			expectedStatusCode: 400,
			expectedBody:       `{"Error":"INVALID_ARG_RECIPIENT"}`,
		},
		{
			description: "missing argument originator",
			smsProvider: &FakeSMSProvider{},
			url:         "/messages",
			payload: map[string]interface{}{
				"recipient": "31612345678",
				"message":   "This is a test message.",
			},
			expectedStatusCode: 400,
			expectedBody:       `{"Error":"MISSING_ARG_ORIGINATOR"}`,
		},
		{
			description: "missing argument message",
			smsProvider: &FakeSMSProvider{},
			url:         "/messages",
			payload: map[string]interface{}{
				"recipient":  "31612345678",
				"originator": "MessageBird",
			},
			expectedStatusCode: 400,
			expectedBody:       `{"Error":"MISSING_ARG_MESSAGE"}`,
		},
	}

	for _, tc := range testCases {
		sms := service.NewSMSService(tc.smsProvider)

		api := rest.NewApi()
		api.Use(rest.DefaultProdStack...)
		router, err := rest.MakeRouter(rest.Post("/messages", PostMessage(sms)))
		if err != nil {
			t.Fatal(err)
		}
		api.SetApp(router)

		t.Log("Test", tc.description)
		recorded := test.RunRequest(t, api.MakeHandler(),
			test.MakeSimpleRequest("POST", "http://127.0.0.1"+tc.url, tc.payload))

		recorded.ContentTypeIsJson()
		recorded.CodeIs(tc.expectedStatusCode)
		if tc.expectedBody != "" {
			recorded.BodyIs(tc.expectedBody)
		}
	}
}
