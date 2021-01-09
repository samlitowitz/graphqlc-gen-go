# GraphQL Cross-compile to Go

# Table of Contents
1. [Types](#types)
   1. [Scalars](#scalars)
   1. [Scalar Extentions](#scalar-extensions)
   1. [Objects](#objects)
   1. [Interfaces](#interfaces)
   1. [Unions](#unions)
   1. [Enums](#enums)
   1. [Input Objects](#input-objects)
   1. [List](#list)
   1. [Non-null](#non-null)
   1. [Directives](#directives)
1. [Repository Interfaces](#repository-interfaces)
1. [Go GraphQL Schema Definition](#go-graphql-schema-definition)
1. [References](#references)

# Types
## Scalars
Convert directly to Go basic types

| GraphQL Type | Go Type |
| --- | --- |
| Integer | int32 |
| Float | float32 |
| Boolean | bool |
| String | string |
| ID | string |

### Scalar Extensions
Map scalar extension type to Go type, e.g. map GraphQLs `scalar DateTime` to Gos `time.Time`.

## Objects
Convert GraphQL object type to Go struct. 
Convert GraphQL object fields to externally visible Go struct fields.

GraphQL
```graphql
type Person {
    name: String
    picture(size: Int): String
}
```

Go
```go
type Person struct {
    Name string
    Picture string
}
```

## Interfaces
Convert GraphQL interface type to Go interface type.
Add `is<TypeName>()` function to Go interface type.
GraphQL interfaces implemented by GraphQL objects need to be implemented by Go structs.

GraphQL
```graphql
interface Node {
    id: ID!
}

type Person implements Node {
    id: ID!
}
```

Go
```go
type Node interface {
    isNode()
}

type Person struct {}
func (p *Person) isNode() {}
```

## Unions
Convert GraphQL union type to Go interface type.
Add `is<TypeName>()` function to Go interface type.
GraphQL objects included a GraphQL union must have their corresponding Go structs implement the Go union interface.

GraphQL
```graphql
type Person {}
type Photo {}
union Result = Person | Photo
```

Go
```go
type Result interface {
    isResult()
}

type Person struct {}
func (p *Person) isResult() {}

type Photo struct {}
func (p *Photo) isResult() {}
```

## Enums
Convert GraphQL enum type to a Go type redefinition of int64.
Convert each GraphQL enum value to a Go const of the type redefinition prefixed with the type name.

GraphQL
```graphql
enum Direction {
    NORTH
    SOUTH
    EAST
    WEST
}
```

Go
```go
type Direction int64
const (
    Direction_NORTH Direction = iota
    Direction_SOUTH Direction
    Direction_EAST Direction
    Direction_WEST Direction
)
```
## Input Objects
GraphQL input objects converted to Go in the same manner as [GraphQL Objects](#objects).

## List
GraphQL list types are converted to Go arrays.

GraphQL
```graphql
type Person {
    photos: [String]!
}
```

Go
```go
type Person struct {
    Photos []string
}
```

## Non-null
GraphQL non-null types are ignored.

## Directives
Not currently supported

# Repository Interfaces

# Go GraphQL Schema Definition
Implement a Go mapping of the GraphQL schema using [graphql-go/graphql](https://github.com/graphql-go/graphql).

# References
1. [GraphQL Spec](https://spec.graphql.org/June2018/)
1. [graphql-go/graphql](https://github.com/graphql-go/graphql)