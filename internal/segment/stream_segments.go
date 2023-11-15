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
	"fmt"
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
	lastTimestamp  media.Timestamp
	streamTime     float64

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
		lastTimestamp:  media.Timestamp{Pts: 0, Dts: 0},
		streamTime:     0,
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
	s.currentSegment.close()
}

func (s *StreamSegments) Write(data []byte, timestamp media.Timestamp) error {
	if s.needNewSegment(timestamp) {
		segment, err := s.createSegment()
		if err != nil {
			return nil
		}

		if s.currentSegment != nil {
			s.currentSegment.close()
			s.segments = append(s.segments, s.currentSegment)
		}

		s.currentSegment = segment
		s.lastTimestamp = media.Timestamp{Pts: 0, Dts: 0}
	}

	return s.currentSegment.write(data, timestamp)
}

func (s *StreamSegments) needNewSegment(timestamp media.Timestamp) bool {
	if s.currentSegment == nil {
		fmt.Println("[TEST][StreamSegments] needNewSegment")
		return true
	}

	s.streamTime += float64(timestamp.Pts) * 0.014
	fmt.Println("[TEST][StreamSegments]["+time.Now().String()+"]needNewSegment - ", timestamp.Pts, " / ", s.streamTime)
	if (timestamp.Pts - s.currentSegment.beginTime.Pts) > uint64(2*1000) {
		fmt.Println("[TEST][StreamSegments] needNewSegment - ",
			timestamp.Pts, " / ", s.currentSegment.beginTime.Pts, " // ",
			(timestamp.Pts - s.currentSegment.beginTime.Pts))
		return true
	}

	return false
}

func (s *StreamSegments) createSegment() (*Segment, error) {
	fmt.Println("[TEST][StreamSegments] createSegment")

	now := time.Now().Format("20060102150405")
	segmentFileName := s.streamBasePath + "/" + strconv.Itoa(s.idCounter) + "_" + now + ".ts"
	segment := NewSegment(s.idCounter, segmentFileName)

	if err := segment.open(); err != nil {
		return nil, err
	}

	s.idCounter++
	return segment, nil
}
