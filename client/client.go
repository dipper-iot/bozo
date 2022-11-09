package client

import (
	"github.com/dipper-iot/bozo/client/selector"
	"github.com/dipper-iot/bozo/logger"
	"github.com/dipper-iot/bozo/registry"
	"github.com/dipper-iot/bozo/service"
	"google.golang.org/grpc"
	"strings"
)

func GetClient(options *service.Options, nameAndVersion string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	arrayData := strings.Split(nameAndVersion, "@")
	o := []registry.GetOption{}
	if len(arrayData) > 1 {
		o = append(o, registry.GetVersion(arrayData[1]))
	}
	nodes, err := options.Registry.GetService(arrayData[0], o...)
	if err != nil {
		return nil, err
	}

	next := selector.Random(nodes)

	node, err := next()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if opts == nil {
		opts = make([]grpc.DialOption, 0)
	}
	opts = append(opts, options.GrpcClientOptions...)

	conn, err = grpc.Dial(node.Id, opts...)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return
}
