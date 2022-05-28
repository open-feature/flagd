package service

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpServiceConfiguration struct {
	Port int32
}

type HttpServiceRequest struct {
	// TODO
}

type HttpServiceResponse struct {
	rawPayload string
}

func (h *HttpServiceResponse) GetPayload() string {
	return h.rawPayload
}

func (h *HttpServiceRequest) GetRequestType() SERVICE_REQUEST_TYPE {
	//TODO
	return SERVICE_REQUEST_ALL_FLAGS
}

func (h *HttpServiceRequest) GenerateServiceResponse(body string) IServiceResponse {
	return &HttpServiceResponse{
		rawPayload: body,
	}
}

type HttpService struct {
	HttpServiceConfiguration *HttpServiceConfiguration
}

func (h *HttpService) Serve(handlerFunc func(IServiceRequest) IServiceResponse) error {
	if h.HttpServiceConfiguration == nil {
		return errors.New("http service configuration has not been initialised")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := handlerFunc(&HttpServiceRequest{})
		//TODO: Improve this significantly and add guards
		w.Write([]byte(response.GetPayload()))
	})
	http.ListenAndServe(fmt.Sprintf(":%d", h.HttpServiceConfiguration.Port), nil)

	return nil
}
