package sport

import (
	"github.com/mctofu/computrainer"
	sportgrpc "github.com/mctofu/sport-grpc"
)

type computrainerPublisher struct {
	deviceID    string
	publishChan chan *sportgrpc.SportData
	lastPower   uint16
	lastSpeed   float64
}

func (c *computrainerPublisher) publish(msg *computrainer.Message) {
	switch msg.Type {
	case computrainer.DataSpeed:
		speed := float64(msg.Value) * .01 * .9 // m/s * 2.237 for m/s to mph
		if speed != c.lastSpeed {
			c.lastSpeed = speed
			c.sendPerformance(sportgrpc.PerformanceType_PERFORMANCE_TYPE_SPEED, float64(speed))
		}
	case computrainer.DataPower:
		power := msg.Value
		if power != c.lastPower {
			c.lastPower = power
			c.sendPerformance(sportgrpc.PerformanceType_PERFORMANCE_TYPE_POWER, float64(power))
		}
	case computrainer.DataRRC:
		rrc := float64(msg.Value) / 256
		c.sendPerformance(sportgrpc.PerformanceType_PERFORMANCE_TYPE_CALIBRATION, rrc)
	}

	if msg.Buttons != computrainer.ButtonNone {
		controlData := &sportgrpc.ControlData{}

		if msg.Buttons&computrainer.ButtonPlus > 0 {
			controlData.Pressed = append(controlData.Pressed, sportgrpc.Button_BUTTON_PLUS)
		}
		if msg.Buttons&computrainer.ButtonMinus > 0 {
			controlData.Pressed = append(controlData.Pressed, sportgrpc.Button_BUTTON_MINUS)
		}
		if msg.Buttons&computrainer.ButtonReset > 0 {
			controlData.Pressed = append(controlData.Pressed, sportgrpc.Button_BUTTON_RESET)
		}
		if msg.Buttons&computrainer.ButtonF1 > 0 {
			controlData.Pressed = append(controlData.Pressed, sportgrpc.Button_BUTTON_F1)
		}
		if msg.Buttons&computrainer.ButtonF2 > 0 {
			controlData.Pressed = append(controlData.Pressed, sportgrpc.Button_BUTTON_F2)
		}
		if msg.Buttons&computrainer.ButtonF3 > 0 {
			controlData.Pressed = append(controlData.Pressed, sportgrpc.Button_BUTTON_F3)
		}

		c.sendControl(controlData)
	}
}

func (c *computrainerPublisher) sendPerformance(t sportgrpc.PerformanceType, val float64) {
	select {
	case c.publishChan <- &sportgrpc.SportData{
		DeviceId: c.deviceID,
		PerformanceData: &sportgrpc.PerformanceData{
			Type:  t,
			Value: val,
		},
	}:
	default:
	}
}

func (c *computrainerPublisher) sendControl(cd *sportgrpc.ControlData) {
	select {
	case c.publishChan <- &sportgrpc.SportData{
		DeviceId:    c.deviceID,
		ControlData: cd,
	}:
	default:
	}
}
