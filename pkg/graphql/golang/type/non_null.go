package _type

import (
	"github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
)

type GoNonNullDefinition struct {
	*graphqlc.NonNullTypeDescriptorProto

	typeMap map[string]golang.GoTypeDefinition

	base golang.GoTypeDefinition
}

func (typDef *GoNonNullDefinition) UnqualifiedName() string {
	if typDef.typeMap == nil {
		typDef.typeMap = make(map[string]golang.GoTypeDefinition)
	}
	if typDef.base == nil {
		switch def := typDef.Type.(type) {
		case *graphqlc.NonNullTypeDescriptorProto_ListType:
			typDef.base = &GoListDefinition{
				ListTypeDescriptorProto: def.ListType,
				typeMap:                 typDef.typeMap,
			}
		case *graphqlc.NonNullTypeDescriptorProto_NamedType:
			typ, ok := typDef.typeMap[def.NamedType.GetName()]
			if !ok {
				break
			}
			typDef.base = &GoNamedDefinition{
				name: typ.UnqualifiedName(),
			}
		}
	}

	if typDef.base == nil {
		return ""
	}
	return typDef.base.UnqualifiedName()
}

func (typDef *GoNonNullDefinition) Definition() string {
	return ""
}
