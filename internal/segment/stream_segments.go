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
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/ISSuh/my-stream-media/internal/media"
)

type StreamSegments struct {
	streamBasePath string

	segmentRange int

	currentSegment *Segment
	segments       []*Segment

	idCounter int
}

func NewStreamSegments(basePath string, segmentRange int) *StreamSegments {
	t := time.Now()
	streamBasePath := basePath + "/" + t.Format("20060102150405")

	return &StreamSegments{
		streamBasePath: streamBasePath,
		segmentRange:   segmentRange,
		currentSegment: nil,
		segments:       make([]*Segment, 0),
		idCounter:      0,
	}
}

func (s *StreamSegments) Open() error {
	if _, err := os.Stat(s.streamBasePath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(s.streamBasePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StreamSegments) Close() {
	if s.currentSegment != nil {
		s.currentSegment.close()
	}
}

func (s *StreamSegments) WriteVideo(data []byte, timeestamp media.Timestamp, isIDRFraem bool) error {
	if s.needNewSegment(isIDRFraem) {
		segment, err := s.createSegment()
		if err != nil {
			return err
		}

		if s.currentSegment != nil {
			s.currentSegment.close()
			s.segments = append(s.segments, s.currentSegment)
		}

		s.currentSegment = segment
	}

	return s.currentSegment.write(data, timeestamp)
}

func (s *StreamSegments) WriteAudio(data []byte, timeestamp media.Timestamp) error {
	return s.currentSegment.write(data, timeestamp)
}

func (s *StreamSegments) needNewSegment(isIDRFraem bool) bool {
	if s.currentSegment == nil {
		return true
	}

	if isIDRFraem {
		return true
	}

	return false
}

func (s *StreamSegments) createSegment() (*Segment, error) {
	now := time.Now().Format("20060102150405")
	segmentFileName := s.streamBasePath + "/" + strconv.Itoa(s.idCounter) + "_" + now + ".ts"
	segment := NewSegment(s.idCounter, segmentFileName)

	if err := segment.open(); err != nil {
		return nil, err
	}

	s.idCounter++
	return segment, nil
}
