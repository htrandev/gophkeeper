package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/htrandev/gophkeeper/internal/domain"
	pb "github.com/htrandev/gophkeeper/proto"
)

// GophkeeperServer определяет grpc-сервер для работы с секретными данными.
type GophkeeperServer struct {
	pb.UnimplementedGophkeeperServer
	opts ServerOptions
}

// NewGophkeeper возвращает новый экземпляр GophkeeperServer.
func NewGophkeeper(opts ...Option) *GophkeeperServer {
	s := &GophkeeperServer{}

	for _, opt := range opts {
		opt(&s.opts)
	}
	return s
}

// Add добавляет данные к пользователю.
func (s *GophkeeperServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// валидируем запрос
	var form AddForm
	if err := form.LoadAndValidate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate: %s", err.Error())
	}

	// получаем идентификатор пользователя
	uid, err := s.authorize(form.Req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get id from token: %s", err.Error())
	}
	form.Req.UserID = uid

	// добавляем данные к пользователю
	res, err := s.opts.keeper.Add(ctx, form.Req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server: add: %s", err.Error())
	}
	return &pb.AddResponse{Id: res}, nil
}

func (s *GophkeeperServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.EmptyResponse, error) {
	// валидируем запрос
	var form DeleteForm
	if err := form.LoadAndValidate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate: %s", err.Error())
	}

	// получаем идентификатор пользователя
	uid, err := s.authorize(form.Req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get id from token: %s", err.Error())
	}

	form.Req.UserID = uid

	if err := s.opts.keeper.Delete(ctx, form.Req); err != nil {
		// если пользователь пытается удалить чужие данные
		if errors.Is(err, domain.ErrPermissionDenied) {
			return nil, status.Errorf(codes.PermissionDenied, "server: delete: %s", domain.ErrPermissionDenied.Error())
		}
		if errors.Is(err, domain.ErrNotFound) {
			return &pb.EmptyResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "server: delete: %s", err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func (s *GophkeeperServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	// валидируем запрос
	var form GetForm
	if err := form.LoadAndValidate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate: %s", err.Error())
	}

	// получаем идентификатор пользователя
	uid, err := s.authorize(form.Req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get id from token: %s", err.Error())
	}

	form.Req.UserID = uid

	data, err := s.opts.keeper.Get(ctx, form.Req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server: get: %s", err.Error())
	}

	return &pb.GetResponse{Data: buildData(data)}, nil
}

func (s *GophkeeperServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	// валидируем запрос
	var form GetAllForm
	if err := form.LoadAndValidate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate: %s", err.Error())
	}

	// получаем идентификатор пользователя
	uid, err := s.authorize(form.Req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get id from token: %s", err.Error())
	}

	info, err := s.opts.keeper.GetAll(ctx, uid)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &pb.GetAllResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "server: getAll: %s", err.Error())
	}
	return &pb.GetAllResponse{Info: buildInfo(info)}, nil
}

func (s *GophkeeperServer) authorize(token string) (uuid.UUID, error) {
	id, err := s.opts.authorizer.GetIDFromToken(token)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("cant get id from token: %w", err)
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("cant parse user id [%s]: %w", id, err)
	}

	return uid, nil
}

func buildData(data domain.Data) *pb.Data {
	res := &pb.Data{
		Id:          data.ID,
		Description: data.Descriptoin,
	}
	switch data.Kind {
	case domain.PayloadLogPass:
		res.Payload = &pb.Data_LogPass{
			LogPass: &pb.LogPass{
				Login:    data.LogPass.Login,
				Password: data.LogPass.Password,
			},
		}
	case domain.PayloadText:
		res.Payload = &pb.Data_Text{
			Text: &pb.Text{Text: data.Text.Text},
		}
	case domain.PayloadFile:
		res.Payload = &pb.Data_File{
			File: &pb.File{
				Name:    data.File.Name,
				Content: data.File.Content,
			},
		}
	case domain.PayloadBankCard:
		res.Payload = &pb.Data_BankCard{
			BankCard: &pb.BankCard{
				Holder: data.BankCard.Holder,
				Number: data.BankCard.Number,
			},
		}
	}
	return res
}

func buildInfo(info []domain.Info) []*pb.Info {
	res := make([]*pb.Info, 0, len(info))
	for _, inf := range info {
		res = append(res, &pb.Info{
			Id:          inf.ID,
			Type:        pb.InfoType(inf.Kind),
			Description: inf.Description,
		})
	}
	return res
}
