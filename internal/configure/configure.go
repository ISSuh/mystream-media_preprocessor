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

type DiscorveryConfigure struct {
	ServerUrls        []string `yaml:"serverUrls"`
	ConnectionTimeout int      `yaml:"connectTimeout"`
	PollInterval      int      `yaml:"pollInterval"`
	Retry             int      `yaml:"retry"`
}

type ServerConfigure struct {
	RtmpPort               string              `yaml:"rtmpPort"`
	Discovery              DiscorveryConfigure `yaml:"discovery"`
	BroadcastServerAddress string              `yaml:"broadcastServerAddress"`
	PacketSize             int                 `yaml:"packetSize"`
	RequestTimeout         int                 `yaml:"requestTimeout"`
}

type MediaEncodingConfigure struct {
	Resolution  string `yaml:"resolution"`
	Frame       int    `yaml:"frame"`
	Keyint      int    `yaml:"keyint"`
	SegmentTime int    `yaml:"segmentTime"`
}

type MediaConfigure struct {
	Reserve  string                   `yaml:"reserve"`
	Encoding []MediaEncodingConfigure `yaml:"encoding"`
}

type SegmentConfigure struct {
	BasePath string `yaml:"basePath"`
	TsRange  int    `yaml:"tsRange"`
}

type Configure struct {
	Server  ServerConfigure  `yaml:"server"`
	Media   MediaConfigure   `yaml:"media"`
	Segment SegmentConfigure `yaml:"segment"`
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
