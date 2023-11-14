package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"

	"github.com/yapingcat/gomedia/go-codec"
	"github.com/yapingcat/gomedia/go-mpeg2"
	"github.com/yapingcat/gomedia/go-rtmp"
)

var port = flag.String("port", "1935", "rtmp server listen port")

type MediaCenter map[string]*MediaProducer

var center MediaCenter
var mtx sync.Mutex

func init() {
	center = make(map[string]*MediaProducer)
}

func (c *MediaCenter) register(name string, p *MediaProducer) {
	mtx.Lock()
	defer mtx.Unlock()
	(*c)[name] = p
}

func (c *MediaCenter) unRegister(name string) {
	mtx.Lock()
	defer mtx.Unlock()
	delete(*c, name)
}

func (c *MediaCenter) find(name string) *MediaProducer {
	mtx.Lock()
	defer mtx.Unlock()
	if p, found := (*c)[name]; found {
		return p
	} else {
		return nil
	}
}

type MediaFrame struct {
	cid   codec.CodecID
	frame []byte
	pts   uint32
	dts   uint32
}

func (f *MediaFrame) clone() *MediaFrame {
	tmp := &MediaFrame{
		cid: f.cid,
		pts: f.pts,
		dts: f.dts,
	}
	tmp.frame = make([]byte, len(f.frame))
	copy(tmp.frame, f.frame)
	return tmp
}

type MediaProducer struct {
	name      string
	session   *MediaSession
	mtx       sync.Mutex
	consumers []*MediaSession
	quit      chan struct{}
	die       sync.Once
	muxer     *mpeg2.TSMuxer
	vpid      uint16
	apid      uint16

	tsfilename string
	tsfile     *os.File
}

func newMediaProducer(name string, sess *MediaSession) *MediaProducer {
	mediaProducer := &MediaProducer{
		name:      name,
		session:   sess,
		consumers: make([]*MediaSession, 0, 10),
		quit:      make(chan struct{}),
		muxer:     mpeg2.NewTSMuxer(),
	}

	mediaProducer.muxer.OnPacket = func(pkg []byte) {
		// mpeg.ShowPacketHexdump(pkg)
		_, err := mediaProducer.tsfile.Write(pkg)
		if err != nil {
			fmt.Println("write packet -- err : ", err)
		}

		// fmt.Println("write packet -- ", len(pkg), " / ", n)

		// os.Stdin.Read(make([]byte, 1))
		mediaProducer.tsfile.Sync()
	}

	mediaProducer.vpid = mediaProducer.muxer.AddStream(mpeg2.TS_STREAM_H264)
	mediaProducer.apid = mediaProducer.muxer.AddStream(mpeg2.TS_STREAM_AAC)
	mediaProducer.tsfilename = "./test.ts"

	tsfile, err := os.OpenFile(mediaProducer.tsfilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// defer tsfile.Close()

	mediaProducer.tsfile = tsfile

	fmt.Println("producer.vpid ", mediaProducer.vpid)

	return mediaProducer
}

func (producer *MediaProducer) start() {
	go producer.dispatch()
}

func (producer *MediaProducer) stop() {
	producer.die.Do(func() {
		close(producer.quit)
		center.unRegister(producer.name)
	})
}

func (producer *MediaProducer) dispatch() {
	defer func() {
		fmt.Println("quit dispatch")
		producer.stop()

		producer.tsfile.Close()
	}()
	for {
		select {
		case frame := <-producer.session.C:
			fmt.Println("media frame ", codec.CodecString(frame.cid))
			if frame == nil {
				continue
			}

			if frame.cid == codec.CODECID_VIDEO_H264 {
				codec.SplitFrameWithStartCode(frame.frame, func(nalu []byte) bool {
					// fmt.Println("wtite nalu")
					if codec.H264NaluType(nalu) <= codec.H264_NAL_I_SLICE {
						frame.pts += 40
						frame.dts += 40
						fmt.Println(frame.dts)
					}
					// fmt.Println(producer.muxer.Write(producer.vpid, nalu, uint64(frame.pts), uint64(frame.dts)))
					producer.muxer.Write(producer.vpid, nalu, uint64(frame.pts), uint64(frame.dts))
					return true
				})
			}

			producer.mtx.Lock()
			tmp := make([]*MediaSession, len(producer.consumers))
			copy(tmp, producer.consumers)
			producer.mtx.Unlock()
			for _, c := range tmp {
				if c.ready() {
					tmp := frame.clone()
					c.play(tmp)
				}
			}
		case <-producer.session.quit:
			return
		case <-producer.quit:
			return
		}
	}
}

func (producer *MediaProducer) addConsumer(consumer *MediaSession) {
	producer.mtx.Lock()
	defer producer.mtx.Unlock()
	producer.consumers = append(producer.consumers, consumer)
}

func (producer *MediaProducer) removeConsumer(id string) {
	producer.mtx.Lock()
	defer producer.mtx.Unlock()
	for i, consume := range producer.consumers {
		if consume.id == id {
			producer.consumers = append(producer.consumers[i:], producer.consumers[i+1:]...)
		}
	}
}

type MediaSession struct {
	handle    *rtmp.RtmpServerHandle
	conn      net.Conn
	lists     []*MediaFrame
	mtx       sync.Mutex
	id        string
	isReady   bool
	frameCome chan struct{}
	quit      chan struct{}
	source    *MediaProducer
	die       sync.Once
	C         chan *MediaFrame
}

func newMediaSession(conn net.Conn) *MediaSession {
	id := fmt.Sprintf("%d", rand.Uint64())
	return &MediaSession{
		id:        id,
		conn:      conn,
		handle:    rtmp.NewRtmpServerHandle(),
		quit:      make(chan struct{}),
		frameCome: make(chan struct{}, 1),
		C:         make(chan *MediaFrame, 30),
	}
}

func (sess *MediaSession) init() {

	sess.handle.OnPlay(func(app, streamName string, start, duration float64, reset bool) rtmp.StatusCode {
		fmt.Println("OnPlay - ", app, " / ", streamName)

		if source := center.find(streamName); source == nil {
			return rtmp.NETSTREAM_PLAY_NOTFOUND
		}
		return rtmp.NETSTREAM_PLAY_START
	})

	sess.handle.OnPublish(func(app, streamName string) rtmp.StatusCode {
		fmt.Println("OnPublish - ", streamName)

		return rtmp.NETSTREAM_PUBLISH_START
	})

	sess.handle.SetOutput(func(b []byte) error {
		_, err := sess.conn.Write(b)
		return err
	})

	sess.handle.OnStateChange(func(newState rtmp.RtmpState) {
		if newState == rtmp.STATE_RTMP_PLAY_START {
			fmt.Println("play start")
			name := sess.handle.GetStreamName()
			source := center.find(name)
			sess.source = source
			if source != nil {
				source.addConsumer(sess)
				fmt.Println("ready to play")
				sess.isReady = true
				go sess.sendToClient()
			}
		} else if newState == rtmp.STATE_RTMP_PUBLISH_START {
			fmt.Println("publish start")
			sess.handle.OnFrame(func(cid codec.CodecID, pts, dts uint32, frame []byte) {
				f := &MediaFrame{
					cid:   cid,
					frame: frame, //make([]byte, len(frame)),
					pts:   pts,
					dts:   dts,
				}
				//copy(f.frame, frame)
				sess.C <- f
			})
			name := sess.handle.GetStreamName()
			p := newMediaProducer(name, sess)
			go p.dispatch()
			center.register(name, p)
		}
	})
}

func (sess *MediaSession) start() {
	fmt.Println("MediaSession.start")

	defer sess.stop()
	for {
		buf := make([]byte, 65536)
		n, err := sess.conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = sess.handle.Input(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (sess *MediaSession) stop() {
	fmt.Println("MediaSession.stop")

	sess.die.Do(func() {
		close(sess.quit)
		if sess.source != nil {
			sess.source.removeConsumer(sess.id)
			sess.source = nil
		}
		sess.conn.Close()
	})
}

func (sess *MediaSession) ready() bool {
	fmt.Println("MediaSession.ready")

	return sess.isReady
}

func (sess *MediaSession) play(frame *MediaFrame) {
	fmt.Println("MediaSession.play")

	sess.mtx.Lock()
	sess.lists = append(sess.lists, frame)
	sess.mtx.Unlock()
	select {
	case sess.frameCome <- struct{}{}:
	default:
	}
}

func (sess *MediaSession) sendToClient() {
	fmt.Println("MediaSession.sendToClient")

	firstVideo := true
	for {
		select {
		case <-sess.frameCome:
			sess.mtx.Lock()
			frames := sess.lists
			sess.lists = nil
			sess.mtx.Unlock()
			for _, frame := range frames {
				if firstVideo { //wait for I frame
					if frame.cid == codec.CODECID_VIDEO_H264 && codec.IsH264IDRFrame(frame.frame) {
						firstVideo = false
					} else if frame.cid == codec.CODECID_VIDEO_H265 && codec.IsH265IDRFrame(frame.frame) {
						firstVideo = false
					} else {
						continue
					}
				}
				err := sess.handle.WriteFrame(frame.cid, frame.frame, frame.pts, frame.dts)
				if err != nil {
					sess.stop()
					return
				}
			}
		case <-sess.quit:
			return
		}
	}
}

func startRtmpServer() {
	listen, _ := net.Listen("tcp4", "localhost:1935")
	for {
		conn, _ := listen.Accept()
		sess := newMediaSession(conn)
		sess.init()
		go sess.start()
	}
}

func main() {
	flag.Parse()
	go startRtmpServer()
	select {}
}
