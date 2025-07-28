package dto

type LoginResult struct {
	User  PublicUser
	Token string
	Exp   int
	Err   error
}
