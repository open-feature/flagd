package service

type IServiceConfiguration interface {
}

type SERVICE_REQUEST_TYPE int32

const (
	SERVICE_REQUEST_ALL_FLAGS = iota
)

type IServiceRequest interface {
	GetRequestType() SERVICE_REQUEST_TYPE
	GenerateServiceResponse(body string) IServiceResponse
}

type IServiceResponse interface {
	GetPayload() string
}

type IService interface {
	Serve(handlerFunc func(IServiceRequest) IServiceResponse) error
}
