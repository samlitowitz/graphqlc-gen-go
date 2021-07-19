```go
Schema         *SchemaDescriptorProto                      `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
TypeExtensions []*TypeSystemExtensionDescriptorProto       `protobuf:"bytes,10,rep,name=type_extensions,json=typeExtensions,proto3" json:"type_extensions,omitempty"`
Directives     []*DirectiveDefinitionDescriptorProto       `protobuf:"bytes,9,rep,name=directives,proto3" json:"directives,omitempty"`
Scalars        []*ScalarTypeDefinitionDescriptorProto      `protobuf:"bytes,3,rep,name=scalars,proto3" json:"scalars,omitempty"`
Objects        []*ObjectTypeDefinitionDescriptorProto      `protobuf:"bytes,4,rep,name=objects,proto3" json:"objects,omitempty"`
Interfaces     []*InterfaceTypeDefinitionDescriptorProto   `protobuf:"bytes,5,rep,name=interfaces,proto3" json:"interfaces,omitempty"`
Unions         []*UnionTypeDefinitionDescriptorProto       `protobuf:"bytes,6,rep,name=unions,proto3" json:"unions,omitempty"`
Enums          []*EnumTypeDefinitionDescriptorProto        `protobuf:"bytes,7,rep,name=enums,proto3" json:"enums,omitempty"`
InputObjects   []*InputObjectTypeDefinitionDescriptorProto `protobuf:"bytes,8,rep,name=input_objects,json=inputObjects,proto3" json:"input_objects,omitempty"`
```

# Type
1. Map of GraphQL type name -> definition
## Definition
1. Pointer to AST file defined in
1. Definition
   1. Unqualified Name

## Usage (file scoped)
1. Import path
   1. Alias
