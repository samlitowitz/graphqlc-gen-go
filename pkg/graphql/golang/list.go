package golang

import (
	"fmt"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
)

type GoListDefinition struct {
	*graphqlc.ListTypeDescriptorProto

	typeMap map[string]GoTypeDefinition
	base    GoTypeDefinition
}

func (typDef *GoListDefinition) UnqualifiedName() string {
	if typDef.typeMap == nil {
		typDef.typeMap = make(map[string]GoTypeDefinition)
	}

	if typDef.base == nil {
		typDef.base = buildTypeGoTypeDefinition(typDef.Type, typDef.typeMap)
	}

	if typDef.base == nil {
		return ""
	}

	return fmt.Sprintf(
		"[]%s",
		typDef.base.UnqualifiedName(),
	)
}

func (typDef *GoListDefinition) Definition() string {
	return ""
}
