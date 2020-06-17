package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	"grpc2way/packet"
	"log"
	"math/rand"
	"os"
	"os/signal"
)

var (
	build   string
	version string
	runAddr = ":8080"
	myFlag  = ""
)

const (
	sslCrtFile = "ssl/self-signed-test.crt"
	sslKeyFile = "ssl/self-signed-test.key"
)

//Server -
type Server struct {
	packet.UnimplementedSayHiServer
}

func main() {

	expvar.NewString("build").Set(build)
	expvar.NewString("version").Set(version)
	log.Printf("Application initializing : build version %q - %q \n", build, version)

	runFlag := flag.String("mode", "", "run mode -mode=server or -mode=client")
	flag.Parse()

	if *runFlag == "" {
		log.Println("invalid run mode", *runFlag)
		return
	}

	fmt.Println("RUN MODE :", *runFlag)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	myFlag = *runFlag
	switch *runFlag {
	case "server":
		go RunServer(runAddr)
		break
	case "client":
		go RunClient(runAddr)
		break
	default:
		log.Println("invalid run mode", *runFlag)
		break
	}

	<-signalCh

	log.Println("stop")
}

//Say -
func (*Server) Say(ctx context.Context, req *packet.PacketData) (*packet.PacketData, error) {

	log.Println("Say Request : ", req)

	return &packet.PacketData{Id: rand.Int31(), Msg: myFlag + " reply from say"}, nil
}
