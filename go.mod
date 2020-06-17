module github.com/mctofu/sport-grpc

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/go-multierror v1.1.0
	github.com/mctofu/computrainer v0.0.1
	go.bug.st/serial v1.1.0
	google.golang.org/grpc v1.29.1
)

replace github.com/mctofu/computrainer => ../computrainer/
