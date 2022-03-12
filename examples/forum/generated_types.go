package main

import "time"

type UserRole int

const (
	UserRole_ADMINISTRATOR UserRole = 0
	UserRole_MODERATOR     UserRole = 1
	UserRole_MEMBER        UserRole = 2
)

type Content interface {
	isContent()
}
type Node interface {
	GetId() string
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
type User struct {
	Posts *PostConnection
	Role  UserRole
	Id    string
}

func (o User) GetId() string {
	return o.Id
}

type Topic struct {
	Name      string
	Author    User
	CreatedAt time.Time
	Id        string
}

func (o Topic) GetId() string {
	return o.Id
}

type Post struct {
	Author    User
	CreatedAt time.Time
	Replies   *ReplyConnection
	Title     string
	Message   string
	Id        string
}

func (o Post) GetId() string {
	return o.Id
}
func (o Post) isContent() {}

type Reply struct {
	Author    User
	CreatedAt time.Time
	Replies   *ReplyConnection
	Message   string
	Id        string
}

func (o Reply) GetId() string {
	return o.Id
}
func (o Reply) isContent() {}

type PageInfo struct {
	HasPreviousPage bool
	HasNextPage     bool
}
type UserEdge struct {
	Node   *User
	Cursor string
}
type UserConnection struct {
	Edge     []*UserEdge
	PageInfo PageInfo
}
type TopicEdge struct {
	Node   *Topic
	Cursor string
}
type TopicConnection struct {
	Edge     []*TopicEdge
	PageInfo PageInfo
}
type PostEdge struct {
	Node   *Post
	Cursor string
}
type PostConnection struct {
	Edge     []*PostEdge
	PageInfo PageInfo
}
type ReplyEdge struct {
	Node   *Reply
	Cursor string
}
type ReplyConnection struct {
	Edge     []*ReplyEdge
	PageInfo PageInfo
}
type CreatePostOutput struct {
	Post *Post
}
type DeletePostOutput struct {
	Post *Post
}
type UpdatePostOutput struct {
	Post *Post
}
