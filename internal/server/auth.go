package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/htrandev/gophkeeper/internal/domain"
	pb "github.com/htrandev/gophkeeper/proto"
)

// AuthServer определяет grpc-сервер для работы с пользователем.
type AuthServer struct {
	pb.UnimplementedAuthServer
	opts ServerOptions
}

// NewAuth возвращает новый экземпляр AuthServer.
func NewAuth(opts ...Option) *AuthServer {
	s := &AuthServer{}

	for _, opt := range opts {
		opt(&s.opts)
	}
	return s
}

// SignUp регистрирует нового пользователя.
func (s *AuthServer) SignUp(ctx context.Context, req *pb.AuthorizationRequest) (*pb.TokenResponse, error) {
	// валидируем запрос
	var form AuthForm
	if err := form.LoadAndValidate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate: %s", err.Error())
	}

	// регистрируем нового пользователя
	token, err := s.opts.auth.SignUp(ctx, form.Req)
	if err != nil {
		if errors.Is(err, domain.ErrNotUniqueLogin) {
			return nil, status.Errorf(codes.AlreadyExists, "server: signUp: %s", domain.ErrNotUniqueLogin.Error())
		}
		return nil, status.Errorf(codes.Internal, "server: signUp: %s", err.Error())
	}

	return &pb.TokenResponse{Token: token}, nil
}

// SignIn авторизовывает пользователя.
func (s *AuthServer) SignIn(ctx context.Context, req *pb.AuthorizationRequest) (*pb.TokenResponse, error) {
	// валидируем запрос
	var form AuthForm
	if err := form.LoadAndValidate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate: %s", err.Error())
	}

	// авторизовываем пользователя
	token, err := s.opts.auth.SignIn(ctx, form.Req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "server: signIn: %s", domain.ErrNotFound.Error())
		}
		return nil, status.Errorf(codes.Internal, "server: signIn: %s", err.Error())
	}

	return &pb.TokenResponse{Token: token}, nil
}
