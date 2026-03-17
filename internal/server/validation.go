package server

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/htrandev/gophkeeper/internal/domain"
	"github.com/htrandev/gophkeeper/pkg/luhn"
	pb "github.com/htrandev/gophkeeper/proto"
)

// AuthForm определяет форму запроса регистрации и авторизации пользователя.
type AuthForm struct {
	Req domain.AuthorizationRequest
}

// LoadAndValidate валидирует и наполняет запрос регистрации и авторизации пользователя.
func (a *AuthForm) LoadAndValidate(req *pb.AuthorizationRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Login, validation.Required),
		validation.Field(&req.Password, validation.Required),
	); err != nil {
		return err
	}

	a.Req.Login = req.GetLogin()
	a.Req.Password = req.GetPassword()

	return nil
}

// AuthForm определяет форму запроса для добавления информации пользователю.
type AddForm struct {
	Req domain.AddRequest
}

// LoadAndValidate валидирует и наполняет запрос для добавления информации пользователю.
func (a *AddForm) LoadAndValidate(req *pb.AddRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Token, validation.Required),
		validation.Field(&req.Payload, validation.Required),
	); err != nil {
		return err
	}

	switch p := req.GetPayload().(type) {
	case *pb.AddRequest_LogPass:
		a.Req.Kind = domain.PayloadLogPass
		logPass := domain.LogPass{
			Login:    p.LogPass.GetLogin(),
			Password: p.LogPass.GetPassword(),
		}
		a.Req.LogPass = &logPass
	case *pb.AddRequest_Text:
		a.Req.Kind = domain.PayloadText
		text := domain.Text{
			Text: p.Text.GetText(),
		}
		a.Req.Text = &text
	case *pb.AddRequest_File:
		a.Req.Kind = domain.PayloadFile
		file := domain.File{
			Name:    p.File.GetName(),
			Content: p.File.GetContent(),
		}
		a.Req.File = &file
	case *pb.AddRequest_BankCard:
		a.Req.Kind = domain.PayloadBankCard
		card := domain.BankCard{
			Holder: p.BankCard.GetHolder(),
			Number: p.BankCard.GetNumber(),
		}

		if !luhn.Check(card.Number) {
			return fmt.Errorf("invalid card number")
		}
		a.Req.BankCard = &card
	default:
		return fmt.Errorf("unknown payload type: %v", p)
	}

	a.Req.Token = req.GetToken()
	a.Req.Description = req.GetDescription()

	return nil
}

// DeleteForm определяет форму запроса для удалении информации пользователя.
type DeleteForm struct {
	Req domain.DeleteRequest
}

// LoadAndValidate валидирует и наполняет запрос для удаления информации пользователя.
func (f *DeleteForm) LoadAndValidate(req *pb.DeleteRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Token, validation.Required),
		validation.Field(&req.Id, validation.Required),
		validation.Field(&req.Type, validation.Required,
			validation.In(
				pb.InfoType_INFO_TYPE_TEXT,
				pb.InfoType_INFO_TYPE_FILE,
				pb.InfoType_INFO_TYPE_LOG_PASS,
				pb.InfoType_INFO_TYPE_BANK_CARD,
			),
		),
	); err != nil {
		return err
	}

	f.Req.Token = req.GetToken()
	f.Req.Kind = domain.PayloadKind(req.GetType())
	f.Req.DataID = req.GetId()

	return nil
}

// GetForm определяет форму запроса для получения информации пользователя.
type GetForm struct {
	Req domain.GetRequest
}

// LoadAndValidate валидирует и наполняет запрос для получения информации пользователя.
func (f *GetForm) LoadAndValidate(req *pb.GetRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Token, validation.Required),
		validation.Field(&req.Id, validation.Required),
		validation.Field(&req.Type, validation.Required,
			validation.In(
				pb.InfoType_INFO_TYPE_TEXT,
				pb.InfoType_INFO_TYPE_FILE,
				pb.InfoType_INFO_TYPE_LOG_PASS,
				pb.InfoType_INFO_TYPE_BANK_CARD,
			),
		),
	); err != nil {
		return err
	}

	f.Req.Token = req.GetToken()
	f.Req.Kind = domain.PayloadKind(req.GetType())
	f.Req.DataID = req.GetId()

	return nil
}

// GetForm определяет форму запроса для получения всей информации пользователя.
type GetAllForm struct {
	Req domain.GetAllRequest
}

// LoadAndValidate валидирует и наполняет запрос для получения все информации пользователя.
func (f *GetAllForm) LoadAndValidate(req *pb.GetAllRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Token, validation.Required),
	); err != nil {
		return nil
	}

	f.Req.Token = req.GetToken()
	return nil
}
 