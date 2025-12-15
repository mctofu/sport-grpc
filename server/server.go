package server

import (
	"context"
	fmt "fmt"
	"log"

	sportgrpc "github.com/mctofu/sport-grpc"
	"github.com/mctofu/sport-grpc/sport"
)

var _ sportgrpc.ControllerServer = (*Server)(nil)

type Server struct {
	controller *sport.Controller
	doneChan   <-chan struct{}
}

func NewServer(ctlr *sport.Controller, doneChan <-chan struct{}) *Server {
	return &Server{
		controller: ctlr,
		doneChan:   doneChan,
	}
}

func (s *Server) ReadData(_ *sportgrpc.DataRequest, ds sportgrpc.Controller_ReadDataServer) error {
	ctx := ds.Context()

	log.Println("Client connected")

	// TODO: support multiple readers
	for {
		select {
		case <-ctx.Done():
			log.Println("Client done")
			return nil
		case <-s.doneChan:
			log.Println("Server done. Aborting ReadData.")
			return nil
		case data := <-s.controller.Reader():
			if err := ds.Send(data); err != nil {
				return fmt.Errorf("send: %v", err)
			}
		}
	}
}

func (s *Server) SetLoad(ctx context.Context, req *sportgrpc.LoadRequest) (*sportgrpc.LoadResponse, error) {
	if err := s.controller.SetLoad(ctx, req.DeviceId, req.TargetLoad); err != nil {
		return nil, err
	}
	return &sportgrpc.LoadResponse{}, nil
}

func (s *Server) Recalibrate(ctx context.Context, req *sportgrpc.RecalibrateRequest) (*sportgrpc.RecalibrateResponse, error) {
	if err := s.controller.Recalibrate(ctx, req.DeviceId); err != nil {
		return nil, err
	}
	return &sportgrpc.RecalibrateResponse{}, nil
}
