package _type

import (
	"fmt"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
	"strings"
)

type GoEnumDefinition struct {
	*graphqlc.EnumTypeDefinitionDescriptorProto
}

func (typDef *GoEnumDefinition) UnqualifiedName() string {
	return typDef.GetName()
}

func (typDef *GoEnumDefinition) Definition() string {
	if len(typDef.Values) == 0 {
		return ""
	}

	members := []string{
		fmt.Sprintf("%s_%s %s = iota", typDef.UnqualifiedName(), typDef.Values[0].GetValue(), typDef.UnqualifiedName()),
	}

	for _, value := range typDef.Values[1:] {
		members = append(members, fmt.Sprintf("%s_%s", typDef.UnqualifiedName(), value.GetValue()))
	}

	return fmt.Sprintf(
		"type %s int64\nconst(\n%s)\n",
		typDef.UnqualifiedName(),
		strings.Join(members, "\n"),
	)
}
