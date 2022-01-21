package _type

import (
	"fmt"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
)

type GoUnionDefinition struct {
	*graphqlc.UnionTypeDefinitionDescriptorProto
}

func (typDef *GoUnionDefinition) UnqualifiedName() string {
	return typDef.GetName()
}

func (typDef *GoUnionDefinition) Definition() string {
	return fmt.Sprintf(
		"type %s interface {\nIs%s()\n}\n",
		typDef.GetName(),
		typDef.GetName(),
	)
}

func (typDef *GoUnionDefinition) MembershipFunctionName() string {
	return fmt.Sprintf("Is%s", typDef.GetName())
}

func (typDef *GoUnionDefinition) TypeIsMember(typ string) bool {
	for _, member := range typDef.MemberTypes {
		if member.GetName() == typ {
			return true
		}
	}
	return false
}
