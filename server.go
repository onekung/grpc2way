package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"grpc2way/packet"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/hashicorp/yamux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

//RunServer -
func RunServer(addr string) {

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer lis.Close()
	log.Println("Server Listening on : ", addr)

	for lis.Addr().Network() != "" {

		incoming, err := lis.Accept()
		if err != nil {
			log.Printf("couldn't accept %s", err)
			continue
		}

		go handleClient(incoming)

	}

}

func handleClient(incoming net.Conn) {

	time.Sleep(time.Second)
	isClose := make(chan bool)
	muxSv, err := yamux.Server(incoming, yamux.DefaultConfig())
	if err != nil {
		panic(err)
	}
	defer muxSv.Close()

	go func() {
		//Server Handler
		s := Server{}
		creds, err := credentials.NewServerTLSFromFile(sslCrtFile, sslKeyFile)
		if err != nil {
			panic(err)
		}

		grpcServer := grpc.NewServer(grpc.Creds(creds))
		packet.RegisterSayHiServer(grpcServer, &s)
		defer grpcServer.Stop()
		err = grpcServer.Serve(muxSv)
		if err != nil {
			log.Println("client close")
		}
		isClose <- true
	}()

	svConn, err := muxSv.Accept()
	if err != nil {
		panic(err)
	}

	defer svConn.Close()

	muxCli, err := yamux.Client(svConn, yamux.DefaultConfig())
	if err != nil {
		panic(err)
	}
	defer muxCli.Close()

	stream, err := muxCli.Open()
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	config := &tls.Config{
		//change to false on use invalid ssl cert
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(":7777", grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithDialer(func(target string, timeout time.Duration) (net.Conn, error) {
			return stream, nil
		}),
	)

	if err != nil {
		panic(err)
	}
	defer conn.Close()
	log.Println("success connect back to client")
	go func() {
		//client Handler
		// loop send test
		for conn.GetState() <= connectivity.Ready {

			c := packet.NewSayHiClient(conn)
			response, err := c.Say(context.Background(), &packet.PacketData{Id: rand.Int31(), Msg: "send from server"})
			if err != nil {
				log.Printf("send err: %s", err)

				continue
			}

			fmt.Println(response)

			time.Sleep(time.Second * 1)

		}

		isClose <- true
	}()

	<-isClose
	log.Println("stop client handler thread")
}
