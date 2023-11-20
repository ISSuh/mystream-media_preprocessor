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

package protocol

import (
	"github.com/ISSuh/mystream-media_preprocessor/internal/media"
	"github.com/ISSuh/mystream-media_preprocessor/internal/transport"
	log "github.com/sirupsen/logrus"

	rtmpCodec "github.com/yapingcat/gomedia/go-codec"
	rtmp "github.com/yapingcat/gomedia/go-rtmp"
)

type RtmpContext struct {
	handler     RtmpHandler
	transporter transport.Transport

	internalHandler *rtmp.RtmpServerHandle
}

func NewRtmpContext() *RtmpContext {
	return &RtmpContext{
		handler:         nil,
		transporter:     nil,
		internalHandler: rtmp.NewRtmpServerHandle(),
	}
}

func (context *RtmpContext) RegistHandler(handler RtmpHandler, transporter transport.Transport) {
	context.handler = handler
	context.transporter = transporter

	context.internalHandler.SetOutput(
		func(data []byte) error {
			return context.transporter.Write(data)
		})

	context.internalHandler.OnPlay(
		func(_, _ string, _, _ float64, _ bool) rtmp.StatusCode {
			log.Warn("[RtmpContext][OnPlay] not support")
			return rtmp.NETSTREAM_CONNECT_REJECTED
		})

	context.internalHandler.OnPublish(
		func(appName, streamPath string) rtmp.StatusCode {
			err := context.handler.OnPrePare(appName, streamPath)
			if err != nil {
				return rtmp.NETCONNECT_CONNECT_REJECTED
			}
			return rtmp.NETSTREAM_PUBLISH_START
		})

	context.internalHandler.OnStateChange(
		func(newState rtmp.RtmpState) {
			switch newState {
			case rtmp.STATE_RTMP_PUBLISH_START:
				context.handler.OnPublish()
			case rtmp.STATE_RTMP_PUBLISH_FAILED:
				context.handler.OnError()
			}
		})

	context.internalHandler.OnFrame(
		func(cid rtmpCodec.CodecID, pts, dts uint32, frame []byte) {
			mediaType, codec := media.ConvertCodec(cid)
			timestamp := media.Timestamp{Pts: uint64(pts), Dts: uint64(dts)}

			switch mediaType {
			case media.MEDIA_VIDEO:
				context.handler.OnVideoFrame(
					media.NewVideoFrame(media.VideoCodec(codec), timestamp, frame, false))
			case media.MEDIA_AUDIO:
				context.handler.OnAudioFrame(
					media.NewAudioFrame(media.AudioCodec(codec), timestamp, frame))
			}
		})
}

func (contex *RtmpContext) InputStream(data []byte) error {
	return contex.internalHandler.Input(data)
}
