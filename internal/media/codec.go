package media

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
