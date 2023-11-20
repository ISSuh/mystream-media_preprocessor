/*
MIT License

Copyright (c) 2023 ISSuh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package service

import (
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/mystream-media_preprocessor/internal/configure"
	"github.com/ISSuh/mystream-media_preprocessor/internal/session"
	"github.com/ISSuh/mystream-media_preprocessor/internal/transport"
)

const (
	NETWORK_TCP_V4     = "tcp4"
	NETWORK_DEFAULT_IP = "0.0.0.0"
)

type Service struct {
	configure      *configure.Configure
	sessionManager *session.SessionManager
}

func NewService(configure *configure.Configure) *Service {
	return &Service{
		configure:      configure,
		sessionManager: session.NewSessionManager(configure),
	}
}

func (service *Service) Run() error {
	log.Info("[Service][Run] service running")
	address := NETWORK_DEFAULT_IP + ":" + service.configure.Server.RtmpPort
	listen, err := net.Listen(NETWORK_TCP_V4, address)
	if err != nil {
		return err
	}

	for {
		connection, err := listen.Accept()
		if err != nil {
			log.Warn("[Service][Run] connection error. ", err)
			continue
		}

		socketTransport := transport.NewSocketTransport(connection, service.configure.Server.PacketSize)
		sessionId := service.sessionManager.CreateNewSession(socketTransport)
		service.sessionManager.RunSession(sessionId)
	}
}
