package service

import (
	"log"
	"net"
)

type SocketServiceConfiguration struct {
	SocketPath string
}

type SocketServiceRequest struct {
	Payload string
}

type SocketServiceResponse struct {
	rawPayload string
}

func (h *SocketServiceResponse) GetPayload() string {
	return h.rawPayload
}

func (h *SocketServiceRequest) GetRequestType() SERVICE_REQUEST_TYPE {
	//TODO
	return SERVICE_REQUEST_ALL_FLAGS
}

func (h *SocketServiceRequest) GenerateServiceResponse(body string) IServiceResponse {
	return &SocketServiceResponse{
		rawPayload: body,
	}
}

type SocketService struct {
	SocketServiceConfiguration *SocketServiceConfiguration
}

func (h *SocketService) Serve(handlerFunc func(IServiceRequest) IServiceResponse) error {

	l, err := net.Listen("unix", h.SocketServiceConfiguration.SocketPath)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		go func() {
			for {
				buf := make([]byte, 512)
				nr, err := fd.Read(buf)
				if err != nil {
					return
				}

				data := buf[0:nr]

				response := handlerFunc(&SocketServiceRequest{
					Payload: string(data),
				})

				_, err = fd.Write([]byte(response.GetPayload()))
				if err != nil {
					log.Fatal("Write: ", err)
				}
			}
		}()
	}
}
