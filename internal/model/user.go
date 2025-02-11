package model

type User struct {
	Login        string
	PasswordHash string // Строго 64 символа (SHA-256 в HEX)
	Balance      int64
}
