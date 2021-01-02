//go:generate graphqlc --go_out=config=golang.json:. schema.graphql
//go:generate go fmt .
package main
