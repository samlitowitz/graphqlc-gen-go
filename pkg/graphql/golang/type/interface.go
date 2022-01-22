package _type

import (
	"fmt"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
	"strings"
)

type GoInterfaceDefinition struct {
	*graphqlc.InterfaceTypeDefinitionDescriptorProto

	TypeMap map[string]Definition
}

func (typDef *GoInterfaceDefinition) UnqualifiedName() string {
	return typDef.GetName()
}

func (typDef *GoInterfaceDefinition) Definition() string {
	if typDef.TypeMap == nil {
		typDef.TypeMap = make(map[string]Definition)
	}
	fnDefs := typDef.functionDefinitions(typDef.TypeMap)
	fields := make([]string, 0)

	for _, fnDef := range fnDefs {
		fields = append(fields, fnDef.Definition())
	}

	return fmt.Sprintf(
		"type %s interface {\n%s\n}\n",
		typDef.UnqualifiedName(),
		strings.Join(fields, "\n"),
	)
}

func (typDef *GoInterfaceDefinition) functionDefinitions(typeMap map[string]Definition) []*GoInterfaceFunctionDefinition {
	fnDefs := make([]*GoInterfaceFunctionDefinition, 0)

	for _, fieldDef := range typDef.Fields {
		fnDefs = append(fnDefs, &GoInterfaceFunctionDefinition{
			name:    strings.Title(fieldDef.Name),
			typ:     fieldDef.Type,
			typeMap: typeMap,
		})
	}

	return fnDefs
}
