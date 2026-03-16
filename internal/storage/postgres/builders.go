package postgres

import (
	"strings"

	"github.com/htrandev/gophkeeper/internal/domain"
)

func buildAddQuery(kind domain.PayloadKind) string {
	var s strings.Builder
	s.Grow(12 + 7 + 1 + 14)

	s.WriteString("INSERT INTO ") // 12

	switch kind {
	case domain.PayloadLogPass:
		s.Grow(53 + 14)
		s.WriteString("log_pass_data(id, description, login, password, user_id) ") // 53
		s.WriteString("VALUES(")                                                   // 7
		s.WriteString("$1, $2, $3, $4")                                            // 14
		s.WriteString(")")                                                         // 1
	case domain.PayloadText:
		s.Grow(41 + 10)
		s.WriteString("text_data(id, description, content, user_id) ") // 41
		s.WriteString("VALUES(")
		s.WriteString("$1, $2, $3") // 10
		s.WriteString(")")
	case domain.PayloadFile:
		s.Grow(41 + 14)
		s.WriteString("file_data(id, description, name, content, user_id) ") // 41
		s.WriteString("VALUES(")
		s.WriteString("$1, $2, $3, $4") // 14
		s.WriteString(")")
	case domain.PayloadBankCard:
		s.Grow(53 + 14)
		s.WriteString("bank_card_data(id, description, holder, number, user_id) ") // 53
		s.WriteString("VALUES(")
		s.WriteString("$1, $2, $3, $4") // 14
		s.WriteString(")")
	}

	s.WriteString(" RETURNING id;") // 13
	return s.String()
}

func buildAddArgs(d domain.AddRequest) []any {
	args := make([]any, 0, 4)
	switch d.Kind {
	case domain.PayloadLogPass:
		args = append(args, d.Description, d.LogPass.Login, d.LogPass.Password, d.UserID)
	case domain.PayloadText:
		args = append(args, d.Description, d.Text.Text, d.UserID)
	case domain.PayloadFile:
		args = append(args, d.Description, d.File.Name, d.File.Content, d.UserID)
	case domain.PayloadBankCard:
		args = append(args, d.Description, d.BankCard.Holder, d.BankCard.Number, d.UserID)
	}
	return args
}

func buildOwnerQuery(kind domain.PayloadKind) string {
	var s strings.Builder

	s.Grow(20 + 13)
	s.WriteString("SELECT user_id FROM ") // 20

	switch kind {
	case domain.PayloadLogPass:
		s.WriteString("log_pass_data ")
	case domain.PayloadText:
		s.WriteString("text_data ")
	case domain.PayloadFile:
		s.WriteString("file_data ")
	case domain.PayloadBankCard:
		s.WriteString("bank_card_data ")
	}

	s.WriteString("WHERE id = $1") // 13

	return s.String()
}

func buildDeleteQuery(kind domain.PayloadKind) string {
	var s strings.Builder

	s.Grow(12 + 13)
	s.WriteString("DELETE from ") // 12

	switch kind {
	case domain.PayloadLogPass:
		s.WriteString("log_pass_data ")
	case domain.PayloadText:
		s.WriteString("text_data ")
	case domain.PayloadFile:
		s.WriteString("file_data ")
	case domain.PayloadBankCard:
		s.WriteString("bank_card_data ")
	}

	s.WriteString("WHERE id = $1") // 13

	return s.String()
}

func buildGetQuery(kind domain.PayloadKind) string {
	var s strings.Builder

	s.Grow(7 + 13)
	s.WriteString("SELECT ") // 12

	switch kind {
	case domain.PayloadLogPass:
		s.WriteString("id, description, login, password FROM log_pass_data ")
	case domain.PayloadText:
		s.WriteString("id, description, content FROM text_data ")
	case domain.PayloadFile:
		s.WriteString("id, description, name, content FROM file_data ")
	case domain.PayloadBankCard:
		s.WriteString("id, description, holder, number FROM bank_card_data ")
	}

	s.WriteString("WHERE id = $1") // 13

	return s.String()
}
