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

	"github.com/ISSuh/mystream-media_preprocessor/internal/media"
	"github.com/ISSuh/mystream-media_preprocessor/internal/rtmp"
	"github.com/ISSuh/mystream-media_preprocessor/internal/segment"
	"github.com/ISSuh/mystream-media_preprocessor/internal/transport"

	log "github.com/sirupsen/logrus"
)

type Session struct {
	sessionId int
	streamKey string

	sessionHandler Handler
	transporter    transport.Transporter
	context        *rtmp.Context

	stopSignal  chan struct{}
	stopRunning sync.Once

	muxer           *media.TsMuxer
	streamSegmgment *segment.StreamSegments
}

func NewSession(sessionHandler Handler, transporter transport.Transporter) *Session {
	session := &Session{
		sessionId:       -1,
		streamKey:       "",
		sessionHandler:  sessionHandler,
		transporter:     transporter,
		context:         rtmp.NewContext(),
		stopSignal:      make(chan struct{}),
		muxer:           media.NewTSMuxer(),
		streamSegmgment: nil,
	}

	session.context.RegistHandler(session, transporter)
	return session
}

func (s *Session) Run() {
	for {
		select {
		case <-s.stopSignal:
			log.Info("[Session][run][", s.sessionId, "] terminate session")
			return
		default:
			err := s.passStream()
			if err != nil {
				if err == io.EOF {
					log.Info("[Session][run][", s.sessionId, "] end of stream")
					s.sessionHandler.streamEnd(s)
				} else {
					log.Error("[Session][run][", s.sessionId, "] stream read error. ", err)
					s.sessionHandler.streamError(s)
				}
			}
		}
	}
}

func (s *Session) setSessionId(id int) {
	s.sessionId = id
}

func (s *Session) registStreamSegment(streamSegmgment *segment.StreamSegments) {
	s.streamSegmgment = streamSegmgment
}

func (s *Session) passStream() error {
	data, err := s.transporter.Read()
	if err != nil {
		return err
	}

	err = s.context.InputStream(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) stop() {
	s.stopRunning.Do(
		func() {
			log.Info("[Session][run][", s.sessionId, "] stop session")
			s.transporter.Close()
			close(s.stopSignal)
		})
}

func (s *Session) OnPrePare(appName, streamPath string) error {
	log.Info("[Session][OnPrePare] appName : ", appName, " streamPath : ", streamPath)
	s.streamKey = streamPath
	return s.sessionHandler.checkValidStream(s, appName, streamPath)
}

func (s *Session) OnPublish() {
	log.Info("[Session][OnPublish][", s.sessionId, "]")
	err := s.sessionHandler.streamStart(s)
	if err != nil {
		s.sessionHandler.streamEnd(s)
	}
}

func (s *Session) OnError() {
	log.Warn("[Session][OnError][", s.sessionId, "]")
	s.sessionHandler.streamError(s)
}

func (s *Session) OnVideoFrame(frame *media.VideoFrame) {
	log.Trace("[Session][OnVideoFrame][", s.sessionId, "]")

	buffer, err := s.muxer.MuxingVideo(frame)
	if err != nil {
		log.Warn("[Session][OnVideoFrame][", s.sessionId, "] video muxing fail. ", err)
		return
	}

	isIDRFraem := media.CheckIsIDRFrame(frame)
	err = s.streamSegmgment.WriteVideo(buffer, frame.Timestamp(), isIDRFraem)
	if err != nil {
		log.Warn("[Session][OnVideoFrame][", s.sessionId, "] segment write fail. ", err)
		return
	}
}

func (s *Session) OnAudioFrame(frame *media.AudioFrame) {
	log.Trace("[Session][OnAudioFrame][", s.sessionId, "]")

	buffer, err := s.muxer.MuxingAudio(frame)
	if err != nil {
		log.Warn("[Session][OnAudioFrame][", s.sessionId, "] audio muxing fail. ", err)
		return
	}

	err = s.streamSegmgment.WriteAudio(buffer, frame.Timestamp())
	if err != nil {
		log.Warn("[Session][OnAudioFrame][", s.sessionId, "] segment write fail. ", err)
		return
	}
}
