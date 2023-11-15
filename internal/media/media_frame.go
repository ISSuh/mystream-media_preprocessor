package media

type MediaType int

const (
	MEDIA_NONE MediaType = iota
	MEDIA_VIDEO
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

func (frame *MediaFrame) Dts() uint64 {
	return frame.timestamp.Dts
}

func (frame *MediaFrame) Pts() uint64 {
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

func NewVideoFrame(codec VideoCodec, timestamp Timestamp, data []byte) *VideoFrame {
	return &VideoFrame{
		MediaFrame: MediaFrame{
			mediaType: MEDIA_VIDEO,
			data:      data,
			timestamp: timestamp,
		},
		codec: codec,
	}
}

func (frame *VideoFrame) Codec() VideoCodec {
	return frame.codec
}

type AudioFrame struct {
	MediaFrame
	codec AudioCodec
}

func NewAudioFrame(codec AudioCodec, timestamp Timestamp, data []byte) *AudioFrame {
	return &AudioFrame{
		MediaFrame: MediaFrame{
			mediaType: MEDIA_AUDIO,
			data:      data,
			timestamp: timestamp,
		},
		codec: CODEC_AUDIO_AAC,
	}
}

func (frame *AudioFrame) Codec() AudioCodec {
	return frame.codec
}
