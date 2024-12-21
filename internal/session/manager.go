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
	"io"
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

type Manager struct {
	configure *configure.Configure
	sessions  map[int]*Session
	rand      *rand.Rand

	httpClient *http.Client

	segmentManager *segment.SegmentManager
}

func NewManager(configure *configure.Configure) *Manager {
	seed := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(seed)

	Manager := &Manager{
		configure:      configure,
		sessions:       make(map[int]*Session),
		rand:           rand,
		httpClient:     nil,
		segmentManager: segment.NewSessionManager(configure.Segment, configure.Media),
	}

	Manager.httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: Manager.dialTimeout,
		},
	}

	return Manager
}

func (sm *Manager) dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(sm.configure.Server.RequestTimeout)*time.Millisecond)
}

func (sm *Manager) CreateNewSession(transporter transport.Transporter) *Session {
	return NewSession(sm, transporter)
}

func (sm *Manager) TerminateAllSession() {
	for sessionId, session := range sm.sessions {
		session.stop()
		delete(sm.sessions, sessionId)
	}
}

func (sm *Manager) checkValidStream(session *Session, appName, streamKey string) error {
	log.Info("[Manager][checkValidStream]")
	streamStatus, err := sm.requestValidateStreamKey(streamKey)
	if err != nil {
		return err
	}

	if !streamStatus.Active || (len(streamStatus.Url) == 0) {
		return errors.New("invalide stream status")
	}

	if _, exist := sm.sessions[streamStatus.StreamId]; exist {
		return errors.New("alread exist session")
	}

	streamId := streamStatus.StreamId
	streamUrl := streamStatus.Url

	sm.addSession(streamId, session)

	streamSegments, err := sm.segmentManager.OpenStreamSegments(streamId, streamUrl)
	if err != nil {
		return err
	}

	session.registStreamSegment(streamSegments)
	return nil
}

func (sm *Manager) streamStart(session *Session) error {
	log.Info("[Manager][streamStart]")
	return nil
}

func (sm *Manager) streamEnd(session *Session) {
	log.Info("[Manager][streamEnd]")
	streamDeactive := dto.NewStreamActive(session.streamKey)
	sm.requestStreamStatus(streamDeactive, false)

	sm.segmentManager.CloseStreamSegments(session.sessionId)
	sm.stopSession(session)
}

func (sm *Manager) streamError(session *Session) {
	log.Info("[Manager][streamError]")
	streamDeactive := dto.NewStreamActive(session.streamKey)
	sm.requestStreamStatus(streamDeactive, false)

	sm.segmentManager.CloseStreamSegments(session.sessionId)
	sm.stopSession(session)
}

func (sm *Manager) stopSession(session *Session) {
	_, exist := sm.sessions[session.sessionId]
	if exist {
		delete(sm.sessions, session.sessionId)
	}

	session.stop()
}

func (sm *Manager) addSession(streamId int, session *Session) {
	sm.sessions[streamId] = session
	session.setSessionId(streamId)
}

// request about streamKey is validated to mystream-broadcast service
func (sm *Manager) requestValidateStreamKey(streamKey string) (*dto.StreamStatus, error) {
	streamActive := dto.NewStreamActive(streamKey)
	response, err := sm.requestStreamStatus(streamActive, true)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		log.Error("[Manager][requestValidateStreamKey] validate fail from broadcast service. ", response.Error.Message)
		return nil, errors.New("validate fail from broadcast service. " + response.Error.Message)
	}

	return &response.Result, nil
}

func (sm *Manager) requestStreamStatus(streamActive dto.StreamActive, active bool) (*dto.ApiResponse, error) {
	jsonStr, err := json.Marshal(streamActive)
	if err != nil {
		log.Error("[Manager][requestStreamStatus] cat not convert StreamActive to json. ", err)
		return nil, err
	}

	path := StreamActiveUrlPath
	if !active {
		path = StreamDeactiveUrlPath
	}

	response, err := sm.requestToBroadcastService(path, string(jsonStr))
	if err != nil {
		return nil, err
	}

	apiResponse := &dto.ApiResponse{}
	err = json.Unmarshal(response, apiResponse)
	if err != nil {
		log.Error("[Manager][requestStreamStatus] body parse error. ", err, " / ", string(response))
		return nil, err
	}

	log.Info("[Manager][requestStreamStatus] response : ", string(response))
	return apiResponse, nil
}

func (sm *Manager) requestToBroadcastService(uri string, requestBody string) ([]byte, error) {
	url := HttpScheme + sm.configure.Server.BroadcastServerAddress + uri
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Error("[Manager][requestToBroadcastService] cat not create http request. ", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := sm.httpClient.Do(req)
	if err != nil {
		log.Error("[Manager][requestToBroadcastService] http response error. ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("[Manager][requestToBroadcastService] body parse error. ", err)
		return nil, err
	}
	return bytes.Clone(body), nil
}
