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

package discovery

import (
	"errors"
	"strconv"

	"github.com/ISSuh/mystream-media_preprocessor/internal/configure"
	"github.com/hudl/fargo"
)

const (
	BroadCastServiceName = "MYSTREAM-BROADCAST"
)

type DiscorveryClient struct {
	configure  fargo.Config
	connection fargo.EurekaConnection
}

func NewDiscorveryClient(discoveryConfigure *configure.DiscorveryConfigure) *DiscorveryClient {
	var configure fargo.Config
	configure.Eureka.ServiceUrls = discoveryConfigure.ServerUrls
	configure.Eureka.ConnectTimeoutSeconds = discoveryConfigure.ConnectionTimeout
	configure.Eureka.PollIntervalSeconds = discoveryConfigure.PollInterval
	configure.Eureka.Retries = discoveryConfigure.Retry

	client := &DiscorveryClient{
		configure:  configure,
		connection: fargo.NewConnFromConfig(configure),
	}
	return client
}

func (dc *DiscorveryClient) GetBroadcastServiceAddress() (string, error) {
	info, err := dc.connection.GetApp(BroadCastServiceName)
	if err != nil {
		return "", err
	}

	lastInstrance := info.Instances[len(info.Instances)-1]
	if (lastInstrance.Status != fargo.UP) &&
		(lastInstrance.Status != fargo.STARTING) {
		return "", errors.New("cat not connect discory server")
	}

	broadcastServiceAddress := lastInstrance.IPAddr + ":" + strconv.Itoa(lastInstrance.Port)
	return broadcastServiceAddress, nil
}
