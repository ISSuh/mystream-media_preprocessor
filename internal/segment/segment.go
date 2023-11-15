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
	"os"

	"github.com/ISSuh/my-stream-media/internal/media"
)

type Segment struct {
	id       int
	filePath string
	file     *os.File

	beginTime media.Timestamp
	endTime   media.Timestamp
}

func NewSegment(id int, filePath string) *Segment {
	return &Segment{
		id:        id,
		filePath:  filePath,
		file:      &os.File{},
		beginTime: media.Timestamp{Pts: 0, Dts: 0},
		endTime:   media.Timestamp{Pts: 0, Dts: 0},
	}
}

func (s *Segment) Write(data []byte) error {
	if _, err := s.file.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *Segment) Close(timestamp media.Timestamp) {
	s.endTime = timestamp
	s.file.Sync()
	s.file.Close()
}

func (s *Segment) BeginTime() media.Timestamp {
	return s.beginTime
}
