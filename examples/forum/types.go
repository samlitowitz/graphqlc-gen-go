// DO NOT EDIT!!!
package main

import (
	"time"
)

type UserRole int64

const (
	UserRole_ADMINISTRATOR UserRole = iota
	UserRole_MODERATOR
	UserRole_MEMBER
)

type UserEdge struct {
	Node   User
	Cursor string
}

type PostConnection struct {
	Edge     []PostEdge
	PageInfo PageInfo
}

type ReplyEdge struct {
	Node   Reply
	Cursor string
}

type UpdatePostOutput struct {
	Post Post
}

type Reply struct {
	Auther    User
	CreatedAt time.Time
	Replies   ReplyConnection
	Message   string
	Id        string
}

func (rep *Reply) GetId() string {
	return rep.Id
}

func (rep *Reply) IsContent() {}

type PageInfo struct {
	HasPreviousPage bool
	HasNextPage     bool
}

type UserConnection struct {
	Edge     []UserEdge
	PageInfo PageInfo
}

type TopicEdge struct {
	Node   Topic
	Cursor string
}

type PostEdge struct {
	Node   Post
	Cursor string
}

type Node interface {
	GetId() string
}

type Post struct {
	Author    User
	CreatedAt time.Time
	Replies   ReplyConnection
	Title     string
	Message   string
	Id        string
}

func (pos *Post) GetId() string {
	return pos.Id
}

func (pos *Post) IsContent() {}

type ReplyConnection struct {
	Edge     []ReplyEdge
	PageInfo PageInfo
}

type TopicConnection struct {
	Edge     []TopicEdge
	PageInfo PageInfo
}

type CreatePostOutput struct {
	Post Post
}

type Content interface {
	IsContent()
}

type DeletePostInput struct {
	Id string
}

type UpdatePostInput struct {
	Message string
	Id      string
}

type Topic struct {
	Name      string
	Author    User
	CreatedAt time.Time
	Id        string
}

func (top *Topic) GetId() string {
	return top.Id
}

type DeletePostOutput struct {
	Post Post
}

type CreatePostInput struct {
	AuthorId string
	Title    string
	Message  string
}

type User struct {
	Posts PostConnection
	Role  UserRole
	Id    string
}

func (use *User) GetId() string {
	return use.Id
}
