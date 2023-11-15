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

package configure

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ServerConfigure struct {
	RtmpPort               string `yaml:"rtmpPort"`
	BroadcastServerAddress string `yaml:"broadcastServerAddress"`
	PacketSize             int    `yaml:"packetSize"`
}

type MediaConfigure struct {
	RecordPath string `yaml:"recordPath"`
}

type Configure struct {
	Server ServerConfigure `yaml:"server"`
	Media  MediaConfigure  `yaml:"media"`
}

func LoadConfigure(filePath string) (*Configure, error) {
	if len(filePath) <= 0 {
		return nil, errors.New("Invalid option file path")
	}

	var buffer []byte
	var err error
	if buffer, err = loadFile(filePath); err != nil {
		return nil, err
	}

	configure := &Configure{}
	if err = yaml.Unmarshal(buffer, configure); err != nil {
		return nil, err
	}

	return configure, nil
}

func loadFile(path string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return buf, nil
}