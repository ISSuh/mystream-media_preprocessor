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

package ffmpeg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ISSuh/mystream-media_preprocessor/internal/configure"
)

const (
	Command = "ffmpeg " +
		"-i pipe:0 " +
		"-c:v libx264 -x264opts keyint=120:no-scenecut -s 1920x1080 -r 60 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./temp/1080_60/out_1080_60_%03d.ts " +
		"-c:v libx264 -x264opts keyint=120:no-scenecut -s 1280x720 -r 60 -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./temp/720_60/out_720_60_%03d.ts " +
		"-c:v libx264 -x264opts keyint=60:no-scenecut -s 1280x720 -r 30  -profile:v main -preset veryfast -c:a aac -sws_flags bilinear -f segment -segment_time 2 ./temp/720_30/out_720_30_%03d.ts " +
		"-c:v libx264 -x264opts keyint=60:no-scenecut -s 852x480 -r 30  -profile:v main -preset veryfast -c:a aac -sws_flags bilinear  -f segment -segment_time 2 ./temp/480_30/out_480_30_%03d.ts "
)

type FFmpegWrapper struct {
	mediaConfigure configure.MediaConfigure
	basePath       string
	cmd            *exec.Cmd
	inputPipe      io.WriteCloser
}

func NewFFmpegWrapper(mediaConfigure configure.MediaConfigure, basePath string) *FFmpegWrapper {
	return &FFmpegWrapper{
		mediaConfigure: mediaConfigure,
		basePath:       basePath,
		cmd:            nil,
		inputPipe:      nil,
	}
}

func (w *FFmpegWrapper) Open() error {
	if err := w.createDirByResolution(w.basePath); err != nil {
		return err
	}

	args := strings.Fields(w.makeCommand())
	exec.Command(args[0], args[1:]...)
	w.cmd = exec.Command(args[0], args[1:]...)

	var err error
	w.inputPipe, err = w.cmd.StdinPipe()
	if err != nil {
		return err
	}

	return nil
}

func (w *FFmpegWrapper) Run() error {
	err := w.cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (w *FFmpegWrapper) Input(buffer []byte) error {
	if _, err := io.WriteString(w.inputPipe, string(buffer)); err != nil {
		return err
	}
	return nil
}

func (w *FFmpegWrapper) Stop() {
	w.cmd.Process.Kill()
}

func (s *FFmpegWrapper) makeCommand() string {
	command := "ffmpeg -i pipe:0 "

	fmt.Println("[TEST] config : ", s.mediaConfigure)

	for _, configure := range s.mediaConfigure.Encoding {
		videoSubCommand := fmt.Sprintf(
			"-c:v libx264 -x264opts keyint=%d:no-scenecut -s %s -r %d -profile:v main -preset veryfast ",
			configure.Keyint, configure.Resolution, configure.Frame)

		audoSubCommand := fmt.Sprintf("-c:a aac -sws_flags bilinear ")

		fileName := "%06d.ts"
		path := s.basePath + "/" + configure.Resolution + "_" + strconv.Itoa(configure.Frame) + "/" + fileName
		segmentSubCommand := fmt.Sprintf("-f segment -segment_time %d %s ", configure.SegmentTime, path)

		command += videoSubCommand + audoSubCommand + segmentSubCommand
	}

	fmt.Println("[TEST] command : ", command)

	return command
}

func (s *FFmpegWrapper) createDirByResolution(basePath string) error {
	for _, configure := range s.mediaConfigure.Encoding {
		path := s.basePath + "/" + configure.Resolution + "_" + strconv.Itoa(configure.Frame)
		fmt.Println("[TEST] path : ", path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}
