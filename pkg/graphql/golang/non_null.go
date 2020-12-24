package golang

import "github.com/samlitowitz/graphqlc/pkg/graphqlc"

type GoNonNullDefinition struct {
	*graphqlc.NonNullTypeDescriptorProto

	typeMap map[string]GoTypeDefinition

	base GoTypeDefinition
}

func (typDef *GoNonNullDefinition) UnqualifiedName() string {
	if typDef.typeMap == nil {
		typDef.typeMap = make(map[string]GoTypeDefinition)
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
