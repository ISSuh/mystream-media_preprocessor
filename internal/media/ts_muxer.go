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

package media

import (
	"errors"

	"github.com/yapingcat/gomedia/go-mpeg2"
)

type TsMuxer struct {
	videoStreamId uint16
	audioStreamId uint16
	context       *mpeg2.TSMuxer
	videoCodec    *CodecH264
}

func NewTSMuxer() *TsMuxer {
	tsMuxer := &TsMuxer{
		videoStreamId: 0,
		audioStreamId: 0,
		context:       mpeg2.NewTSMuxer(),
		videoCodec:    &CodecH264{},
	}

	tsMuxer.videoStreamId = tsMuxer.context.AddStream(mpeg2.TS_STREAM_H264)
	return tsMuxer
}

func (m *TsMuxer) MuxingVideo(frame *VideoFrame) ([]byte, error) {
	if (frame.mediaType != MEDIA_VIDEO) || (frame.codec != CODEC_VIDEO_H264) {
		return nil, errors.New("invalid frame")
	}

	buffer := make([]byte, 0)
	m.context.OnPacket = func(data []byte) {
		buffer = append(buffer, data...)
	}

	m.videoCodec.ProcessingFrame(frame)
	err := m.context.Write(m.videoStreamId, frame.Data(), frame.Pts(), frame.Dts())
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
