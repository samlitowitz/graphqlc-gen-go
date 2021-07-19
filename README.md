# graphqlc-gen-go
[![Go Report Card](https://goreportcard.com/badge/github.com/samlitowitz/graphqlc-gen-go)](https://goreportcard.com/report/github.com/samlitowitz/graphqlc-gen-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/samlitowitz/graphqlc-gen-go.svg)](https://pkg.go.dev/github.com/samlitowitz/graphqlc-gen-go)

# Table of Contents
1. Introduction (#introduction)
1. Installation(#installation)
1. Usage (#usage)
   1. Parameters (#parameters)
      1. Configuration File (#configuration-file)
         1. Package (#package)
         1. 

# Introduction
This is a code generator designed to work with [graphqlc](https://github.com/samlitowitz/graphqlc).

Generate Go definitions from a GraphQL schema.
See the [examples/](examples/) directory for more... examples.

# Installation
Install [graphqlc](https://github.com/samlitowitz/graphqlc).

`go get -u github.com/samlitowitz/graphqlc-gen-go/cmd/graphqlc-gen-go`

# Usage

## Parameters
  * config, required, name of configuration file as defined above

`graphqlc --go_out=config=golang.json:. schema.graphql`

### Configuration File
#### Package
Go package name to use for generated code

###


1. Scalar Type Map
1. Add ID
   1. Type {int, string}
1. Repository (CRUD)
   1. Interfaces
   1. Implementations
1. CRUD
   1. Functions
1. GraphQL
   1. Schema
      1. Mutation
         1. Add Create, Update, Delete operations
      1. Query
         1. Add Read operations
   1. `https://github.com/graphql-go/graphql` GraphQL Types
   1. Resolvers
   1. Relay
      1. Node
      1. Connection
  
# Reference
1. [Go AST Viewer](https://yuroyoro.github.io/goast-viewer)
