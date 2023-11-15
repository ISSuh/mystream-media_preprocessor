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
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/my-stream-media/internal/configure"
	"github.com/ISSuh/my-stream-media/internal/segment"
	"github.com/ISSuh/my-stream-media/internal/transport"
)

type SessionManager struct {
	configure *configure.Configure
	sessions  map[int]*Session
	rand      *rand.Rand

	segmentManager *segment.SegmentManager
}

func NewSessionManager(configure *configure.Configure) *SessionManager {
	seed := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(seed)

	return &SessionManager{
		configure:      configure,
		sessions:       make(map[int]*Session),
		rand:           rand,
		segmentManager: segment.NewSessionManager(configure.Segment),
	}
}

func (sm *SessionManager) CreateNewSession(transporter transport.Transport) int {
	sessionId := sm.generateUniqSessionId()
	session := NewSession(sessionId, sm, transporter)

	sm.sessions[sessionId] = session
	return sessionId
}

func (sm *SessionManager) RunSession(sessionId int) {
	session, exist := sm.sessions[sessionId]
	if !exist {
		log.Error("[SessionManager][RunSession] invalid session id. ", sessionId)
		return
	}

	go session.run()
}

func (sm *SessionManager) TerminateAllSession() {
	for sessionId, session := range sm.sessions {
		session.stop()
		delete(sm.sessions, sessionId)
	}
}

func (sm *SessionManager) checkValidStream(sessionId int, appName, streamPath string) error {
	log.Info("[SessionManager][checkValidStream][", sessionId, "]")

	streamSegments, err := sm.segmentManager.OpenStreamSegments(sessionId)
	if err != nil {
		return err
	}

	sm.sessions[sessionId].registStreamSegment(streamSegments)
	return nil
}

func (sm *SessionManager) streamStart(sessionId int) error {
	log.Info("[SessionManager][streamStart][", sessionId, "]")
	return nil
}

func (sm *SessionManager) streamEnd(sessionId int) {
	log.Info("[SessionManager][streamEnd][", sessionId, "]")
	sm.segmentManager.CloseStreamSegments(sessionId)
	sm.stopSession(sessionId)
}

func (sm *SessionManager) streamError(sessionId int) {
	log.Info("[SessionManager][streamError][", sessionId, "]")
	sm.segmentManager.CloseStreamSegments(sessionId)
	sm.stopSession(sessionId)
}

func (sm *SessionManager) stopSession(sessionId int) {
	session, exist := sm.sessions[sessionId]
	if !exist {
		log.Error("[SessionManager][stopSession] invalid session id. ", sessionId)
		return
	}

	session.stop()
	delete(sm.sessions, sessionId)
}

func (sm *SessionManager) generateUniqSessionId() int {
	isUniqe := false
	sessionId := 0
	for !isUniqe {
		sessionId = rand.Int()
		_, exsit := sm.sessions[sessionId]
		isUniqe = !exsit
	}
	return sessionId
}
