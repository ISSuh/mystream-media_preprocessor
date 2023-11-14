package service

import (
	"fmt"
	"net"

	"github.com/ISSuh/my-stream-media/internal/session"
)

const (
	NETWORK_TCP_V4     = "tcp4"
	NETWORK_DEFAULT_IP = "0.0.0.0"
)

type Service struct {
	sessionManager *session.SessionManager
}

func (service *Service) Run(port string) error {
	address := NETWORK_DEFAULT_IP + ":" + port
	listen, err := net.Listen(NETWORK_TCP_V4, address)
	if err != nil {
		return err
	}

	for {
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("[Service.Run] connection error. ", err)
			continue
		}

		connection.Close()
		// session := sessionManager.newSession()
		// session.start()
	}
}
