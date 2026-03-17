package domain

type PayloadKind int

const (
	PayloadUnknown PayloadKind = iota
	PayloadText
	PayloadFile
	PayloadLogPass
	PayloadBankCard
)

func (k PayloadKind) String() string {
	switch k {
	case PayloadLogPass:
		return "LogPass"
	case PayloadText:
		return "Text"
	case PayloadFile:
		return "File"
	case PayloadBankCard:
		return "BankCard"
	default:
		return "Unknown"
	}
}

func Parse(kind string) PayloadKind {
	switch kind {
	case "LogPass":
		return PayloadLogPass
	case "Text":
		return PayloadText
	case "File":
		return PayloadFile
	case "BankCard":
		return PayloadBankCard
	default:
		return PayloadUnknown
	}
}
