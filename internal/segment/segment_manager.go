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

package segment

import (
	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/mystream-media_preprocessor/internal/configure"
)

type SegmentManager struct {
	segmentConfigure configure.SegmentConfigure
	mediaConfigure   configure.MediaConfigure
	streams          map[int]*StreamSegments
}

func NewSessionManager(segmentConfigure configure.SegmentConfigure, mediaConfigure configure.MediaConfigure) *SegmentManager {
	return &SegmentManager{
		segmentConfigure: segmentConfigure,
		mediaConfigure:   mediaConfigure,
		streams:          make(map[int]*StreamSegments, 0),
	}
}

func (sm *SegmentManager) OpenStreamSegments(streamId int, uri string) (*StreamSegments, error) {
	log.Info("[SegmentManager][OpenStreamSegments][", streamId, "]")
	streamSegmentBasePath := sm.segmentConfigure.BasePath + uri

	streamSegments := NewStreamSegments(sm.mediaConfigure, streamSegmentBasePath)
	if err := streamSegments.Open(); err != nil {
		return nil, err
	}

	sm.streams[streamId] = streamSegments
	return streamSegments, nil
}

func (sm *SegmentManager) CloseStreamSegments(userId int) {
	log.Info("[SegmentManager][CloseStreamSegments][", userId, "]")

	streamSegments := sm.streams[userId]
	if streamSegments != nil {
		streamSegments.Close()
	}

	delete(sm.streams, userId)
}
