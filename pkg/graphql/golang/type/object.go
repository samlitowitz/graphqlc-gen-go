package _type

import (
	"fmt"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
	"strings"
)

type GoObjectDefinition struct {
	*graphqlc.ObjectTypeDefinitionDescriptorProto

	TypeMap map[string]Definition
}

func (typDef *GoObjectDefinition) UnqualifiedName() string {
	return typDef.GetName()
}

func (typDef *GoObjectDefinition) Definition() string {
	if typDef.TypeMap == nil {
		typDef.TypeMap = make(map[string]Definition)
	}

	fields := make([]string, 0)

	for _, fieldDef := range typDef.Fields {
		var fieldGoTypDef Definition = buildTypeGoTypeDefinition(fieldDef.Type, typDef.TypeMap)
		if fieldGoTypDef == nil {
			continue
		}
		fields = append(fields, fmt.Sprintf("%s %s", strings.Title(fieldDef.Name), fieldGoTypDef.UnqualifiedName()))
	}

	fnTyp := lcfirst(typDef.UnqualifiedName()[:3])
	functions := make([]string, 0)
	for _, ifaceDef := range typDef.Implements {
		typ, ok := typDef.TypeMap[ifaceDef.GetName()]
		if !ok {
			continue
		}
		ifaceGoDef, ok := typ.(*GoInterfaceDefinition)
		if !ok {
			continue
		}
		for _, def := range ifaceGoDef.functionDefinitions(typDef.TypeMap) {
			functions = append(
				functions,
				fmt.Sprintf(
					"func (%s *%s) %s {\nreturn %s.%s\n}\n",
					fnTyp,
					typDef.UnqualifiedName(),
					def.Definition(),
					fnTyp,
					def.UnqualifiedName(),
				),
			)
		}
	}

	for _, def := range typDef.TypeMap {
		unionDef, ok := def.(*GoUnionDefinition)
		if !ok {
			continue
		}
		if !unionDef.TypeIsMember(typDef.GetName()) {
			continue
		}
		functions = append(
			functions,
			fmt.Sprintf(
				"func (%s *%s) %s() {}\n",
				fnTyp,
				typDef.UnqualifiedName(),
				unionDef.MembershipFunctionName(),
			),
		)
	}

	fnSep := "";
	if len(functions) > 0 {
		fnSep = "\n\n"
	}

	return fmt.Sprintf(
		"type %s struct {\n%s\n}%s%s\n",
		typDef.UnqualifiedName(),
		strings.Join(fields, "\n"),
		fnSep,
		strings.Join(functions, "\n"),
	)
}
