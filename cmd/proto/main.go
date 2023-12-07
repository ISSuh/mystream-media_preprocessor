package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ISSuh/mystream-media_preprocessor/message"
	"google.golang.org/grpc"
)

func main() {
	serverAddr := "127.0.0.1:50051"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	c := message.NewGreeterClient(conn)
	// msg, err := c.SayHello(context.Background(), &message.HelloRequest{Name: "test"})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// fmt.Println("msg : ", msg)

	res := []*message.HelloRequest{
		&message.HelloRequest{Name: "test1"},
		&message.HelloRequest{Name: "test2"},
		&message.HelloRequest{Name: "test3"},
	}

	stream, err := c.SayHello(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stream.CloseAndRecv()

	for _, r := range res {
		stream.Send(r)
	}

	time.Sleep(5 * time.Second)
}
