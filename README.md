# graphqlc-gen-go
[![Go Report Card](https://goreportcard.com/badge/github.com/samlitowitz/graphqlc-gen-go)](https://goreportcard.com/report/github.com/samlitowitz/graphqlc-gen-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/samlitowitz/graphqlc-gen-go.svg)](https://pkg.go.dev/github.com/samlitowitz/graphqlc-gen-go)

This is a code generator designed to work with [graphqlc](https://github.com/samlitowitz/graphqlc).

Generate Go definitions from a GraphQL schema.
See the [examples/](https://github.com/samlitowitz/graphqlc-gen-relayify/tree/master/examples) directory for more... examples.

# Installation
Install [graphqlc](https://github.com/samlitowitz/graphqlc).

`go get -u github.com/samlitowitz/graphqlc-gen-go/cmd/graphqlc-gen-go`

# Usage

## Parameters
  * config, required, name of configuration file as defined above

`graphqlc --go_out=config=golang.json:. schema.graphql`
  