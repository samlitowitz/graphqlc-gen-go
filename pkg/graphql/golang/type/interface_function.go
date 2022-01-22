package _type

import (
	"fmt"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
	"strings"
)

type GoInterfaceFunctionDefinition struct {
	name    string
	typ     *graphqlc.TypeDescriptorProto
	typeMap map[string]Definition
}

func (typDef *GoInterfaceFunctionDefinition) UnqualifiedName() string {
	return typDef.name
}

func (typDef *GoInterfaceFunctionDefinition) Definition() string {
	var fieldGoTypDef = buildTypeGoTypeDefinition(typDef.typ, typDef.typeMap)
	if fieldGoTypDef == nil {
		return ""
	}
	return fmt.Sprintf("Get%s() %s", typDef.name, fieldGoTypDef.UnqualifiedName())
}

type GoInputObjectDefinition struct {
	*graphqlc.InputObjectTypeDefinitionDescriptorProto

	TypeMap map[string]Definition
}

func (typDef *GoInputObjectDefinition) UnqualifiedName() string {
	return typDef.GetName()
}

func (typDef *GoInputObjectDefinition) Definition() string {
	if typDef.TypeMap == nil {
		typDef.TypeMap = make(map[string]Definition)
	}

	fields := make([]string, 0)

	for _, fieldDef := range typDef.Fields {
		var fieldGoTypDef = buildTypeGoTypeDefinition(fieldDef.Type, typDef.TypeMap)
		if fieldGoTypDef == nil {
			continue
		}
		fields = append(fields, fmt.Sprintf("%s %s", strings.Title(fieldDef.Name), fieldGoTypDef.UnqualifiedName()))
	}

	return fmt.Sprintf(
		"type %s struct {\n%s\n}\n",
		typDef.UnqualifiedName(),
		strings.Join(fields, "\n"),
	)
}
