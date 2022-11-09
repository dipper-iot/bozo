package service

import (
	"context"
	"github.com/dipper-iot/bozo/cli"
	"github.com/dipper-iot/bozo/registry"
	_ "github.com/joho/godotenv"
	"google.golang.org/grpc"
	"os"
	"reflect"
	"sync"
)

type Options struct {
	cliApp                             *cli.App
	Address                            string
	Version                            string
	Name                               string
	Id                                 string
	Server                             *grpc.Server
	Context                            context.Context
	Registry                           registry.Registry
	cancel                             context.CancelFunc
	listBeforeStart                    []Callback
	listAfterStart                     []Callback
	listBeforeStop                     []CallbackStop
	listAfterStop                      []CallbackStop
	one                                sync.Once
	grpcOptions                        []grpc.ServerOption
	GrpcClientOptions                  []grpc.DialOption
	loaders                            []ILoader
	data                               map[string]interface{}
	grpcOptionsStreamInterceptors      []grpc.StreamServerInterceptor
	grpcOptionsUnaryServerInterceptors []grpc.UnaryServerInterceptor
	commands                           map[string]*cli.Command
	flags                              []cli.Flag
	MetaData                           map[string]string
}

var defaultFlag = []cli.Flag{
	cli.StringFlag{
		Name:         "server_address",
		Aliases:      []string{"a"},
		Env:          []string{"SERVER_ADDRESS"},
		Usage:        "Server Address",
		DefaultValue: ":",
	},
}

func NewOptions() *Options {
	ctx, cancel := context.WithCancel(context.Background())
	id := RandStringRunes(8)
	return &Options{
		Id:                                 id,
		Registry:                           registry.DefaultRegistry,
		Context:                            ctx,
		cancel:                             cancel,
		listBeforeStart:                    []Callback{},
		listAfterStart:                     []Callback{},
		listBeforeStop:                     []CallbackStop{},
		listAfterStop:                      []CallbackStop{},
		grpcOptions:                        []grpc.ServerOption{},
		one:                                sync.Once{},
		grpcOptionsStreamInterceptors:      []grpc.StreamServerInterceptor{},
		grpcOptionsUnaryServerInterceptors: []grpc.UnaryServerInterceptor{},
		GrpcClientOptions:                  []grpc.DialOption{},

		MetaData: map[string]string{},
		loaders:  []ILoader{},
		data:     map[string]interface{}{},
		commands: map[string]*cli.Command{},
		flags:    defaultFlag,
	}
}

type Callback = func(o *Options, c *cli.Context) error
type CallbackStop = func(o *Options) error
type Option = func(o *Options)

func (o *Options) Set(name string, data interface{}) {
	o.data[name] = data
}

func (o *Options) Get(name string) (interface{}, bool) {
	data, success := o.data[name]
	return data, success
}

func (o Options) Cli() cli.App {

	return *o.cliApp
}

func (o *Options) Command(command *cli.Command) {
	if command == nil {
		return
	}
	if o.cliApp != nil {
		o.cliApp.Commands = append(o.cliApp.Commands, command)
		return
	}
	o.commands[command.Name] = command
}

func (o *Options) Flags(flags []cli.Flag) {
	if flags == nil || len(flags) < 1 {
		return
	}
	o.flags = append(o.flags, flags...)
}

func (o *Options) OptionServerStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) {
	if len(interceptors) > 0 {
		o.grpcOptionsStreamInterceptors = append(o.grpcOptionsStreamInterceptors, interceptors...)
	}
}

func (o *Options) OptionServerUnaryServerInterceptor(interceptors ...grpc.UnaryServerInterceptor) {
	if len(interceptors) > 0 {
		o.grpcOptionsUnaryServerInterceptors = append(o.grpcOptionsUnaryServerInterceptors, interceptors...)
	}
}

func (o *Options) BeforeStart(callback Callback) {
	o.listBeforeStart = append(o.listBeforeStart, callback)
}

func (o *Options) AfterStart(callback Callback) {
	o.listAfterStart = append(o.listAfterStart, callback)
}

func (o *Options) BeforeStop(callback CallbackStop) {
	o.listBeforeStop = append(o.listBeforeStop, callback)
}

func (o *Options) AfterStop(callback CallbackStop) {
	o.listAfterStop = append(o.listAfterStop, callback)
}

func (o *Options) AddOptionGrpc(options ...grpc.ServerOption) {
	o.grpcOptions = append(o.grpcOptions, options...)
}

func (o *Options) AddOptionGrpcClient(options ...grpc.DialOption) {
	if len(options) > 0 {
		o.GrpcClientOptions = append(o.GrpcClientOptions, options...)
	}
}

func Flag(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Flags(flags)
	}
}

func Command(command *cli.Command) Option {
	return func(o *Options) {
		o.Command(command)
	}
}

func Address(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

func SetMetaData(name string, value string) Option {
	return func(o *Options) {
		o.MetaData[name] = value
	}
}

func (o *Options) run(before cli.ActionBeforeFunc, action cli.ActionFunc, after cli.ActionFunc) error {
	name := "Cli Service"
	if o.Name != "" {
		name = o.Name
	}

	commands := make([]*cli.Command, 0)
	for _, command := range o.commands {
		commands = append(commands, command)
	}

	runLoaderFlag(o.loaders, o)

	o.cliApp = &cli.App{
		Name:        name,
		Description: "Cli Command service",
		Flags:       o.flags,
		Commands:    commands,
		Context:     o.Context,
		Before: func(c *cli.ContextBefore) error {

			if before != nil {
				err := before(c)
				if err != nil {
					return err
				}
			}

			err := runCallBack(o, o.listBeforeStart, c.Context)
			if err != nil {
				return err
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			if action != nil {
				err := action(c)
				if err != nil {
					return err
				}
			}
			return nil
		},
		After: func(c *cli.Context) error {

			if after != nil {
				err := after(c)
				if err != nil {
					return err
				}
			}

			err := runCallBack(o, o.listAfterStart, c)
			if err != nil {
				return err
			}

			return nil
		},
	}

	err := o.cliApp.Run(os.Args)
	if err != nil {
		return err
	}
	return nil
}

func ServerUnaryServerInterceptor(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		o.OptionServerUnaryServerInterceptor(interceptors...)
	}
}

func ServerStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(o *Options) {
		o.OptionServerStreamInterceptor(interceptors...)
	}
}

func ServerOption(options ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.AddOptionGrpc(options...)
	}
}

func ClientOption(options ...grpc.DialOption) Option {
	return func(o *Options) {
		o.AddOptionGrpcClient(options...)
	}
}

func Loader(loaders ...ILoader) Option {
	return func(o *Options) {
		for _, loader := range loaders {
			data := reflect.ValueOf(loader)
			if !data.IsNil() {
				o.loaders = append(o.loaders, loader)
			}
		}
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Version(name string) Option {
	return func(o *Options) {
		o.Version = name
	}
}

func Register(register registry.Registry) Option {
	return func(o *Options) {
		o.Registry = register
	}
}

func SetContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context, o.cancel = context.WithCancel(ctx)
	}
}

func BeforeStart(fn Callback) Option {
	return func(o *Options) {
		o.listBeforeStart = append(o.listBeforeStart, fn)
	}
}

func AfterStart(fn Callback) Option {
	return func(o *Options) {
		o.listAfterStart = append(o.listAfterStart, fn)
	}
}

func BeforeStop(fn CallbackStop) Option {
	return func(o *Options) {
		o.listBeforeStop = append(o.listBeforeStop, fn)
	}
}

func AfterStop(fn CallbackStop) Option {
	return func(o *Options) {
		o.listAfterStop = append(o.listAfterStop, fn)
	}
}

func runCallBack(o *Options, list []Callback, c *cli.Context) error {

	for _, callback := range list {
		if callback == nil {
			continue
		}
		err := callback(o, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func runCallBackStop(o *Options, list []CallbackStop) error {

	for _, callback := range list {
		if callback == nil {
			continue
		}
		err := callback(o)
		if err != nil {
			return err
		}
	}
	return nil
}
