package golang

import "github.com/samlitowitz/graphqlc/pkg/graphqlc"

type GoScalarDefinition struct {
	*graphqlc.ScalarTypeDefinitionDescriptorProto
}

func (typDef *GoScalarDefinition) UnqualifiedName() string {
	return typDef.GetName()
}

func (typDef *GoScalarDefinition) Definition() string {
	return ""
}

func buildTypeGoTypeDefinition(typDef *graphqlc.TypeDescriptorProto, typeMap map[string]GoTypeDefinition) GoTypeDefinition {
	switch def := typDef.Type.(type) {
	case *graphqlc.TypeDescriptorProto_ListType:
		return &GoListDefinition{
			ListTypeDescriptorProto: def.ListType,
			typeMap:                 typeMap,
		}
	case *graphqlc.TypeDescriptorProto_NamedType:
		typ, ok := typeMap[def.NamedType.GetName()]
		if !ok {
			return nil
		}
		return &GoNamedDefinition{
			name: typ.UnqualifiedName(),
		}
	case *graphqlc.TypeDescriptorProto_NonNullType:
		return &GoNonNullDefinition{
			NonNullTypeDescriptorProto: def.NonNullType,
			typeMap:                    typeMap,
		}
	}
	return nil
}
