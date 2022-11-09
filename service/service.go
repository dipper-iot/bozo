package service

import (
	"context"
	"fmt"
	"github.com/dipper-iot/bozo/cli"
	"github.com/dipper-iot/bozo/logger"
	"github.com/dipper-iot/bozo/registry"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Service struct {
	Options *Options
	listen  net.Listener
	init    bool
	ctx     context.Context
	cancel  context.CancelFunc
	run     chan error
}

func NewService(opts ...Option) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	srv := &Service{
		listen:  nil,
		Options: NewOptions(),
		init:    false,
		run:     make(chan error),
		ctx:     ctx,
		cancel:  cancel,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(srv.Options)
	}

	return srv
}

func (s *Service) Init(opts ...Option) error {

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(s.Options)
	}

	optionsServer := append(s.Options.grpcOptions,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(s.Options.grpcOptionsStreamInterceptors...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(s.Options.grpcOptionsUnaryServerInterceptors...)),
	)

	s.Options.Server = grpc.NewServer(optionsServer...)

	s.init = true
	return nil
}

func (s *Service) Run() error {

	if !s.init {
		s.Init()
	}

	err := s.Start()
	if err != nil {
		return err
	}

	return s.Stop()
}

func (s *Service) Start() error {
	var (
		err error
		run bool = false
	)

	signalChanel := make(chan os.Signal)
	signal.Notify(signalChanel, os.Interrupt)

	err = s.Options.run(
		// before
		func(c *cli.ContextBefore) error {

			sortLoaders(s.Options.loaders)

			err = runLoader(s.Options.loaders, s.Options, c.Context, true)
			if err != nil {
				return err
			}

			return nil
		},
		// action
		func(c *cli.Context) error {
			s.Options.Address = c.String("server_address")

			return nil
		},
		// after
		func(c *cli.Context) error {
			s.listen, err = net.Listen("tcp", s.Options.Address)
			if err != nil {
				logger.Errorf("Failed to listen: %v", err)
				return err
			}

			go func() {
				logger.Infof("Start Service %s at: %s", s.Options.Name, s.listen.Addr())
				if err := s.Options.Server.Serve(s.listen); err != nil {
					s.cancel()
					s.run <- err
					return
				}
				s.cancel()
			}()
			time.Sleep(100 * time.Millisecond)

			address := s.listen.Addr().String()
			if strings.Contains(address, "[::]:") {
				ips := strings.Split(address, "[::]:")
				address = fmt.Sprintf("127.0.0.1:%s", ips[len(ips)-1])
			}
			s.Options.MetaData["address"] = address

			service := &registry.Service{
				Name:     s.Options.Name,
				Version:  s.Options.Version,
				Metadata: s.Options.MetaData,
				Nodes: []*registry.Node{
					{
						Id:       s.Options.Id,
						Metadata: s.Options.MetaData,
						Address:  s.Options.MetaData["address"],
					},
				},
			}

			// registry
			err := s.Options.Registry.Register(service)
			if err != nil {
				logger.Error("Register: ", err)
			}

			run = true
			return nil
		})
	if err != nil {
		return err
	}

	if run {
		go func() {
			select {
			case <-signalChanel:
				{
					logger.Infof("Stop with Signal")
				}
			case <-s.ctx.Done():

			}

			service := &registry.Service{
				Name:     s.Options.Name,
				Version:  s.Options.Version,
				Metadata: s.Options.MetaData,
				Nodes: []*registry.Node{
					{
						Id:       s.Options.Id,
						Metadata: s.Options.MetaData,
						Address:  s.Options.MetaData["address"],
					},
				},
			}

			// registry
			err := s.Options.Registry.Deregister(service)
			if err != nil {
				logger.Error("Register: ", err)
			}

			s.run <- s.Stop()

		}()

		return <-s.run
	}
	return nil
}

func (s *Service) Server() *grpc.Server {
	return s.Options.Server
}

func (s *Service) Stop() error {

	if s.Options.cancel != nil {
		s.Options.cancel()
	}

	err := runCallBackStop(s.Options, s.Options.listBeforeStop)
	if err != nil {
		return err
	}

	s.Options.Server.Stop()
	if s.listen != nil {
		err = s.listen.Close()
		if err != nil {
			return err
		}
	}

	err = runCallBackStop(s.Options, s.Options.listAfterStop)
	if err != nil {
		return err
	}

	err = runLoader(s.Options.loaders, s.Options, nil, false)
	if err != nil {
		return err
	}

	return nil
}
