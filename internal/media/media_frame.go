package media

type MediaType int

const (
	MEDIA_VIDEO MediaType = iota
	MEDIA_AUDIO
)

type MediaFrame struct {
	mediaType MediaType
	data      []byte
	timestamp Timestamp
}

func (frame *MediaFrame) Data() []byte {
	return frame.data
}

func (frame *MediaFrame) Dts() uint32 {
	return frame.timestamp.Dts
}

func (frame *MediaFrame) Pts() uint32 {
	return frame.timestamp.Pts
}

func (frame *MediaFrame) MediaType() MediaType {
	return frame.mediaType
}

func (frame *MediaFrame) Codec() int {
	return 0
}

type VideoFrame struct {
	MediaFrame
	codec VideoCodec
}

func newVideoFrame() *VideoFrame {
	return &VideoFrame{
		MediaFrame: MediaFrame{
			mediaType: MEDIA_VIDEO,
			data:      make([]byte, 0),
			timestamp: Timestamp{
				Pts: 0,
				Dts: 0,
			},
		},
		codec: CODEC_VIDEO_H264,
	}
}

func (frame *VideoFrame) Codec() VideoCodec {
	return frame.codec
}

type AudioFrame struct {
	MediaFrame
	codec AudioCodec
}

func newAudioFrame() *AudioFrame {
	return &AudioFrame{
		MediaFrame: MediaFrame{
			mediaType: MEDIA_AUDIO,
			data:      make([]byte, 0),
			timestamp: Timestamp{
				Pts: 0,
				Dts: 0,
			},
		},
		codec: CODEC_AUDIO_AAC,
	}
}

func (frame *AudioFrame) Codec() AudioCodec {
	return frame.codec
}
