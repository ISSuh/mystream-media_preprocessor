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
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/mystream-media_preprocessor/internal/configure"
	"github.com/ISSuh/mystream-media_preprocessor/internal/segment"
	"github.com/ISSuh/mystream-media_preprocessor/internal/session/dto"
	"github.com/ISSuh/mystream-media_preprocessor/internal/transport"
)

const (
	HttpScheme            = "http://"
	StreamUrlPathPrefix   = "/api/broadcast/v1/streams/"
	StreamActiveUrlPath   = StreamUrlPathPrefix + "active"
	StreamDeactiveUrlPath = StreamUrlPathPrefix + "deactive"
)

type SessionManager struct {
	configure *configure.Configure
	sessions  map[int]*Session
	rand      *rand.Rand

	httpClient *http.Client

	segmentManager *segment.SegmentManager
}

func NewSessionManager(configure *configure.Configure) *SessionManager {
	seed := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(seed)

	sessionManager := &SessionManager{
		configure:      configure,
		sessions:       make(map[int]*Session),
		rand:           rand,
		httpClient:     nil,
		segmentManager: segment.NewSessionManager(configure.Segment),
	}

	sessionManager.httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: sessionManager.dialTimeout,
		},
	}

	return sessionManager
}

func (sm *SessionManager) dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(sm.configure.Server.RequestTimeout)*time.Millisecond)
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

func (sm *SessionManager) checkValidStream(sessionId int, appName, streamKey string) error {
	log.Info("[SessionManager][checkValidStream][", sessionId, "]")

	err := sm.requestValidateStreamKey(sessionId, streamKey)
	if err != nil {
		return err
	}

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

// request about streamKey is validated to mystream-broadcast service
func (sm *SessionManager) requestValidateStreamKey(sessionId int, streamKey string) error {
	streamActive := dto.NewStreamActive(streamKey)
	jsonStr, err := json.Marshal(streamActive)
	if err != nil {
		log.Error("[SessionManager][checkValidStream][", sessionId, "] cat not convert StreamActive to json. ", err)
		return err
	}

	url := HttpScheme + sm.configure.Server.BroadcastServerAddress + StreamActiveUrlPath
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error("[SessionManager][checkValidStream][", sessionId, "] cat not create http request. ", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := sm.httpClient.Do(req)
	if err != nil {
		log.Error("[SessionManager][checkValidStream][", sessionId, "] http response error. ", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("[SessionManager][checkValidStream][", sessionId, "] body parse error. ", err)
		return err
	}

	apiResponse := dto.ApiResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Error("[SessionManager][checkValidStream][", sessionId, "] body parse error. ", err)
		return err
	}

	log.Info("[SessionManager][checkValidStream][", sessionId, "] resp : ", apiResponse)

	if !apiResponse.Success {
		log.Error("[SessionManager][checkValidStream][", sessionId, "] validate fail from broadcast service.", apiResponse.Error.Message)
		return errors.New("validate fail from broadcast service. " + apiResponse.Error.Message)
	}

	return nil
}
