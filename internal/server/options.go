package server

import (
	"github.com/htrandev/gophkeeper/internal/contracts"
	"go.uber.org/zap"
)

type Option func(*ServerOptions)

type ServerOptions struct {
	logger     *zap.Logger
	auth       contracts.AuthService
	authorizer contracts.Authorizer
	keeper     contracts.KeeperService
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *ServerOptions) {
		opts.logger = logger
	}
}

func WithAuthorizerService(service contracts.Authorizer) Option {
	return func(opts *ServerOptions) {
		opts.authorizer = service
	}
}

func WithAuthService(service contracts.AuthService) Option {
	return func(opts *ServerOptions) {
		opts.auth = service
	}
}

func WithKeeperService(service contracts.KeeperService) Option {
	return func(opts *ServerOptions) {
		opts.keeper = service
	}
}
