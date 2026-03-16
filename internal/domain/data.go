package domain

// Data определяет формат данных.
type Data struct {
	ID          string
	Descriptoin string
	Kind        PayloadKind
	LogPass     *LogPass
	Text        *Text
	File        *File
	BankCard    *BankCard
}

// LogPass определяет пары логин/пароль.
type LogPass struct {
	Login    string
	Password string
}

// Text определяет произвольные текстовые данные.
type Text struct {
	Text string
}

// File произвольные бинарные данные.
type File struct {
	Name    string
	Content []byte
}

// BankCard орпделеяет данные банковских карт.
type BankCard struct {
	Holder string
	Number string
}

// Info орпделеяет формат информации о хранимых данных.
type Info struct {
	ID          string
	Description string
	Kind        PayloadKind
}
