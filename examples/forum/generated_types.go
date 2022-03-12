package main

type UserRole int

const (
	UserRole_ADMINISTRATOR UserRole = 0
	UserRole_MODERATOR     UserRole = 1
	UserRole_MEMBER        UserRole = 2
)

type Node interface {
	Id() string
}
type CreatePostInput struct {
	AuthorId string
	Title    string
	Message  string
}
type DeletePostInput struct {
	Id string
}
type UpdatePostInput struct {
	Message string
	Id      string
}
