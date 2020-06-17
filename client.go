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
	"google.golang.org/grpc/credentials"
)

//RunClient -
func RunClient(addr string) error {

	isClose := false
	conn, err := net.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		panic(err)

	}
	defer conn.Close()
	log.Printf("Connect to ROUTE-SERVICE %s success", addr)

	muxCli, err := yamux.Client(conn, yamux.DefaultConfig())
	if err != nil {
		panic(err)

	}

	svConn, err := muxCli.Open()
	if err != nil {
		panic(err)

	}

	defer svConn.Close()

	svMux, err := yamux.Server(svConn, yamux.DefaultConfig())
	if err != nil {
		panic(err)

	}
	defer svMux.Close()

	go func() {

		s := Server{}

		creds, err := credentials.NewServerTLSFromFile(sslCrtFile, sslKeyFile)
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer(grpc.Creds(creds))

		packet.RegisterSayHiServer(grpcServer, &s)
		log.Println("client handler server")
		defer grpcServer.Stop()
		err = grpcServer.Serve(svMux)
		if err != nil {
			//panic(err)
		}
		isClose = true
	}()

	time.Sleep(time.Second * 2)

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	cliConn, err := grpc.Dial(":7777", grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithDialer(func(target string, timeout time.Duration) (net.Conn, error) {
			return muxCli.Open()
		}),
	)

	if err != nil {
		panic(err)
	}

	c := packet.NewSayHiClient(cliConn)

	for isClose == false {
		//loop send to server
		response, err := c.Say(context.Background(), &packet.PacketData{Id: rand.Int31(), Msg: "send from client"})

		if err != nil {
			log.Printf("err: %s", err)
			return nil
		}

		fmt.Println(response)

		time.Sleep(time.Second * 2)
	}

	return nil
}
