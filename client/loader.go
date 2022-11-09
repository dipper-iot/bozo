package client

import (
	"github.com/dipper-iot/bozo/cli"
	"github.com/dipper-iot/bozo/service"
	"google.golang.org/grpc"
)

type ClientLoader struct {
	conn       *grpc.ClientConn
	serverName string
}

func NewClientLoader(serverName string) *ClientLoader {
	return &ClientLoader{serverName: serverName}
}

func (c ClientLoader) Reference() []string {
	return []string{}
}

func (c ClientLoader) ReferenceWith(name string, object interface{}) {

}

func (c ClientLoader) Name() string {
	return "client"
}

func (c ClientLoader) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (c ClientLoader) Priority() int {
	return 1
}

func (c *ClientLoader) Start(o *service.Options, ci *cli.Context) error {
	var err error
	c.conn, err = GetClient(o, c.serverName)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientLoader) Client() *grpc.ClientConn {
	return c.conn
}

func (c *ClientLoader) Stop() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
