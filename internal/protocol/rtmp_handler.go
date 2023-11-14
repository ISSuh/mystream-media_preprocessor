package protocol

import (
	"github.com/ISSuh/my-stream-media/internal/media"
)

type RtmpHandler interface {
	OnPrepare(appName, streamName string) error

	OnFrame(mediaType media.MediaType, codec int, timestamp media.Timestamp, data []byte)

	OnPublish()

	OnStop()
}
