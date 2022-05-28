package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpServiceConfiguration struct {
	Port int32
}

type HttpServiceRequest struct {
	Payload string
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		response := handlerFunc(&HttpServiceRequest{
			Payload: string(body),
		})

		w.Write([]byte(response.GetPayload()))
	})
	http.ListenAndServe(fmt.Sprintf(":%d", h.HttpServiceConfiguration.Port), nil)

	return nil
}
