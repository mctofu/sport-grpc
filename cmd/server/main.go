package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	sportgrpc "github.com/mctofu/sport-grpc"
	"github.com/mctofu/sport-grpc/server"
	"github.com/mctofu/sport-grpc/sport"
	"go.bug.st/serial"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
}

func run() error {
	ports, err := serial.GetPortsList()
	if err != nil {
		return fmt.Errorf("couldn't get port list: %v", err)
	}

	for _, port := range ports {
		log.Printf("Port: %v\n", port)
	}

	antSerialPort := flag.String("antserial", "", "Serial port for nRF24AP1")
	computrainerPort := flag.String("computrainer", "", "Serial port for CompuTrainer")

	flag.Parse()

	if computrainerPort == nil && antSerialPort == nil {
		return errors.New("specify antserial and/or computrainer ports")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctlr := sport.NewController()
	defer func() {
		if err := ctlr.Close(); err != nil {
			log.Printf("Controller cleanup error: %v", err)
		}
	}()
	if computrainerPort != nil {
		if err := ctlr.AddComputrainer(ctx, *computrainerPort); err != nil {
			return fmt.Errorf("addComputrainer: %v", err)
		}
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sportgrpc.RegisterControllerServer(s, server.NewServer(ctlr, ctx.Done()))

	go func() {
		defer signal.Stop(stop)
		<-stop
		log.Println("Stopping")
		cancel()
		s.GracefulStop()
		log.Println("Stopped")
	}()

	log.Println("Starting server")
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
