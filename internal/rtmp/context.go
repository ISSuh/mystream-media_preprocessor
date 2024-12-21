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

package rtmp

import (
	"github.com/ISSuh/mystream-media_preprocessor/internal/media"
	"github.com/ISSuh/mystream-media_preprocessor/internal/transport"
	log "github.com/sirupsen/logrus"

	"github.com/yapingcat/gomedia/go-codec"
	"github.com/yapingcat/gomedia/go-rtmp"
)

type Context struct {
	handler     ServerHandler
	transporter transport.Transporter

	internalHandler *rtmp.RtmpServerHandle
}

func NewContext() *Context {
	return &Context{
		handler:         nil,
		transporter:     nil,
		internalHandler: rtmp.NewRtmpServerHandle(),
	}
}

func (c *Context) RegistHandler(handler ServerHandler, transporter transport.Transporter) {
	c.handler = handler
	c.transporter = transporter

	c.internalHandler.SetOutput(
		func(data []byte) error {
			return c.transporter.Write(data)
		})

	c.internalHandler.OnPlay(
		func(_, _ string, _, _ float64, _ bool) rtmp.StatusCode {
			log.Warn("[RtmpContext][OnPlay] not support")
			return rtmp.NETSTREAM_CONNECT_REJECTED
		})

	c.internalHandler.OnPublish(
		func(appName, streamPath string) rtmp.StatusCode {
			err := c.handler.OnPrePare(appName, streamPath)
			if err != nil {
				return rtmp.NETCONNECT_CONNECT_REJECTED
			}
			return rtmp.NETSTREAM_PUBLISH_START
		})

	c.internalHandler.OnStateChange(
		func(newState rtmp.RtmpState) {
			switch newState {
			case rtmp.STATE_RTMP_PUBLISH_START:
				c.handler.OnPublish()
			case rtmp.STATE_RTMP_PUBLISH_FAILED:
				c.handler.OnError()
			}
		})

	c.internalHandler.OnFrame(
		func(cid codec.CodecID, pts, dts uint32, frame []byte) {
			mediaType, codec := media.ConvertCodec(cid)
			timestamp := media.Timestamp{Pts: uint64(pts), Dts: uint64(dts)}

			switch mediaType {
			case media.MEDIA_VIDEO:
				c.handler.OnVideoFrame(
					media.NewVideoFrame(media.VideoCodec(codec), timestamp, frame, false))
			case media.MEDIA_AUDIO:
				c.handler.OnAudioFrame(
					media.NewAudioFrame(media.AudioCodec(codec), timestamp, frame))
			}
		})
}

func (c *Context) InputStream(data []byte) error {
	return c.internalHandler.Input(data)
}
