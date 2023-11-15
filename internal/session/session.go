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

package session

import (
	"io"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/my-stream-media/internal/media"
	"github.com/ISSuh/my-stream-media/internal/protocol"
	"github.com/ISSuh/my-stream-media/internal/transport"
)

type Session struct {
	sessionId int

	sessionHandler SessionHandler
	transporter    transport.Transport
	context        *protocol.RtmpContext

	stopSignal  chan struct{}
	stopRunning sync.Once
}

func NewSession(sessionId int, sessionHandler SessionHandler, transporter transport.Transport) *Session {
	session := &Session{
		sessionId:      sessionId,
		sessionHandler: sessionHandler,
		transporter:    transporter,
		context:        protocol.NewRtmpContext(),
		stopSignal:     make(chan struct{}),
	}

	session.context.RegistHandler(session, transporter)
	return session
}

func (session *Session) run() {
	for {
		select {
		case <-session.stopSignal:
			log.Info("[Session][run][", session.sessionId, "] terminate session")
			return
		default:
			err := session.passStream()
			if err != nil {
				if err == io.EOF {
					log.Info("[Session][run][", session.sessionId, "] end of stream")
					session.sessionHandler.streamEnd(session.sessionId)
				} else {
					log.Error("[Session][run][", session.sessionId, "] stream read error. ", err)
					session.sessionHandler.streamError(session.sessionId)
				}
			}
		}
	}
}

func (session *Session) passStream() error {
	data, err := session.transporter.Read()
	if err != nil {
		return err
	}

	err = session.context.InputStream(data)
	if err != nil {
		return err
	}
	return nil
}

func (session *Session) stop() {
	session.stopRunning.Do(
		func() {
			log.Info("[Session][run][", session.sessionId, "] stop session")
			session.transporter.Close()
			close(session.stopSignal)
		})
}

func (session *Session) OnPrePare(appName, streamPath string) error {
	log.Info("[Session][OnPrePare][", session.sessionId, "] appName : ", appName, " streamPath : ", streamPath)
	return session.sessionHandler.checkValidStream(session.sessionId, appName, streamPath)
}

func (session *Session) OnPublish() {
	log.Info("[Session][OnPublish][", session.sessionId, "]")
	err := session.sessionHandler.streamStart(session.sessionId)
	if err != nil {
		session.sessionHandler.streamEnd(session.sessionId)
	}
}

func (session *Session) OnError() {
	log.Warn("[Session][OnError][", session.sessionId, "]")
	session.sessionHandler.streamError(session.sessionId)
}

func (session *Session) OnVideoFrame(frame *media.VideoFrame) {
	log.Trace("[Session][OnVideoFrame][", session.sessionId, "]")

}

func (session *Session) OnAudioFrame(frame *media.AudioFrame) {
	log.Trace("[Session][OnAudioFrame][", session.sessionId, "]")

}
