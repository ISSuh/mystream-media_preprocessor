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

import rtmpCodec "github.com/yapingcat/gomedia/go-codec"

type VideoCodec int
type AudioCodec int

const (
	CODEC_VIDEO_NONE VideoCodec = iota
	CODEC_VIDEO_H264
)

const (
	CODEC_AUDIO_NONE AudioCodec = iota
	CODEC_AUDIO_AAC
)

func ConvertCodec(codecId rtmpCodec.CodecID) (MediaType, int) {
	mediaType := MEDIA_NONE
	codec := 0

	switch codecId {
	case rtmpCodec.CODECID_VIDEO_H264:
		mediaType = MEDIA_VIDEO
		codec = int(CODEC_VIDEO_H264)
	case rtmpCodec.CODECID_VIDEO_H265:
		mediaType = MEDIA_VIDEO
		codec = int(CODEC_VIDEO_NONE)
	case rtmpCodec.CODECID_VIDEO_VP8:
		mediaType = MEDIA_VIDEO
		codec = int(CODEC_VIDEO_NONE)
	case rtmpCodec.CODECID_AUDIO_AAC:
		mediaType = MEDIA_AUDIO
		codec = int(CODEC_AUDIO_AAC)
	case rtmpCodec.CODECID_AUDIO_G711A:
		mediaType = MEDIA_AUDIO
		codec = int(CODEC_AUDIO_NONE)
	case rtmpCodec.CODECID_AUDIO_G711U:
		mediaType = MEDIA_AUDIO
		codec = int(CODEC_AUDIO_NONE)
	case rtmpCodec.CODECID_AUDIO_OPUS:
		mediaType = MEDIA_AUDIO
		codec = int(CODEC_AUDIO_NONE)
	case rtmpCodec.CODECID_AUDIO_MP3:
		mediaType = MEDIA_AUDIO
		codec = int(CODEC_AUDIO_NONE)
	case rtmpCodec.CODECID_UNRECOGNIZED:
		mediaType = MEDIA_NONE
		codec = int(CODEC_AUDIO_NONE)
	}

	return mediaType, codec
}
