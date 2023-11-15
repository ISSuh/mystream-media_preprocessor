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
		context:        nil,
		stopSignal:     make(chan struct{}),
	}

	session.context.RegistHandler(session, transporter)
	return session
}

func (session *Session) run() {
	for {
		select {
		case <-session.stopSignal:
			break
		default:
			err := session.passStream()
			if err != nil {
				return
			}
		}

	}
}

func (session *Session) passStream() error {
	data, err := session.transporter.Read()
	if err != nil {
		if err == io.EOF {
			log.Error("[Session][run][", session.sessionId, "] end of stream")
		} else {
			log.Error("[Session][run][", session.sessionId, "] stream read error. ", err)
		}
		return err
	}

	err = session.context.InputStream(data)
	if err != nil {
		log.Error("[Session][run][", session.sessionId, "] stream input error. ", err)
		return err
	}
	return nil
}

func (session *Session) stop() {
	session.stopRunning.Do(
		func() {
			close(session.stopSignal)
		})
}

func (session *Session) OnPrePare(appName, streamPath string) error {
	log.Info("[Session][OnPublish][", session.sessionId, "] appName : ", appName, " streamPath : ", streamPath)
	return session.sessionHandler.checkValidStream(session.sessionId, appName, streamPath)
}

func (session *Session) OnPublish() {
	log.Info("[Session][OnPlay][", session.sessionId, "]")
	session.sessionHandler.streamStart(session.sessionId)
}

func (session *Session) OnError() {
	log.Info("[Session][OnError][", session.sessionId, "]")
	session.sessionHandler.streamError(session.sessionId)
}

func (session *Session) OnVideoFrame(frame *media.VideoFrame) {
	log.Error("[Session][", session.sessionId, "] OnVideoFrame")

}

func (session *Session) OnAudioFrame(frame *media.AudioFrame) {
	log.Error("[Session][", session.sessionId, "] OnAudioFrame")

}
