package main

type UserConnectionRepository interface {
	Users(first *int, after string, last *int, before string) *UserConnection
}
type TopicConnectionRepository interface {
	Topics(first *int, after string, last *int, before string) *TopicConnection
}
type PostConnectionRepository interface {
	NewestPosts(first *int, after string, last *int, before string) *PostConnection
	Posts(first *int, after string, last *int, before string) *PostConnection
}
type NodeRepository interface {
	Node(id string) Node
}
type ReplyConnectionRepository interface {
	Replies(first *int, after string, last *int, before string) *ReplyConnection
}
type CreatePostOutputRepository interface {
	CreatePost(input CreatePostInput) *CreatePostOutput
}
type DeletePostOutputRepository interface {
	DeletePost(input DeletePostInput) *DeletePostOutput
}
type UpdatePostOutputRepository interface {
	UpdatePost(input UpdatePostInput) *UpdatePostOutput
}
