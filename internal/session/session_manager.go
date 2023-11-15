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

	"github.com/ISSuh/my-stream-media/internal/transport"
)

type SessionManager struct {
	sessions map[int]*Session
	rand     *rand.Rand
}

func NewSessionManager() *SessionManager {
	seed := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(seed)

	return &SessionManager{
		sessions: make(map[int]*Session),
		rand:     rand,
	}
}

func (sm *SessionManager) CreateNewSession(transporter transport.Transport) *Session {
	sessionId := sm.generateUniqSessionId()
	session := NewSession(transporter)

	sm.sessions[sessionId] = session
	return session
}

func (sm *SessionManager) TerminateSession() {
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
