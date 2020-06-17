package sport

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/mctofu/computrainer"
	sportgrpc "github.com/mctofu/sport-grpc"
)

type Controller struct {
	trainers   map[string]*computrainer.Connection
	outputChan chan *sportgrpc.SportData
}

func NewController() *Controller {
	return &Controller{
		trainers:   make(map[string]*computrainer.Connection),
		outputChan: make(chan *sportgrpc.SportData, 100),
	}
}

func (c *Controller) Reader() <-chan *sportgrpc.SportData {
	return c.outputChan
}

func (c *Controller) AddComputrainer(ctx context.Context, port string) error {
	ct := &computrainer.Controller{}

	conn, err := ct.Start(port)
	if err != nil {
		return fmt.Errorf("computrainer start: %v\n", err)
	}

	deviceID := "default"

	c.trainers[deviceID] = conn

	publisher := computrainerPublisher{
		deviceID:    deviceID,
		publishChan: c.outputChan,
	}

	go func() {
		for msg := range conn.Messages {
			publisher.publish(&msg)
		}
	}()

	return nil
}

func (c *Controller) SetLoad(ctx context.Context, deviceID string, load int32) error {
	trainer, ok := c.trainers[deviceID]
	if !ok {
		return fmt.Errorf("trainer not found: %s", deviceID)
	}

	trainer.SetLoad(load)

	return nil
}

func (c *Controller) Close() error {
	var result error
	for k, v := range c.trainers {
		if err := v.Close(); err != nil {
			result = multierror.Append(result, fmt.Errorf("%s: %v", k, err))
		}
	}
	return result
}
