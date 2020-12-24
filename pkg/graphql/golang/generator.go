package golang

import (
	"github.com/samlitowitz/graphqlc-gen-echo/pkg/graphqlc/echo"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
)

type Generator struct {
	*echo.Generator

	config   *Config
	genFiles map[string]bool
	typeMap  map[string]GoTypeDefinition
}

func New() *Generator {
	g := new(Generator)
	g.Generator = echo.New()
	g.LogPrefix = "graphqlc-gen-go"
	return g
}

func (g *Generator) GenerateSchemaFiles() {}

func (g *Generator) CommandLineArguments(parameter string) {
	g.Generator.CommandLineArguments(parameter)

	for k, v := range g.Param {
		switch k {
		case "config":
			config, err := LoadConfig(v)
			if err != nil {
				g.Error(err)
			}
			g.config = config
		}
	}
	if g.config == nil {
		g.Fail("a configuration must be provided")
	}
}

func (g *Generator) GenerateAllFiles() {
	// Go type definitions
	g.Reset()
	g.writeFileHeader()
	g.writeGoTypeDefinitions()
	g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
		Name:    g.config.Output + "types.go",
		Content: g.String(),
	})

	// Go schema definition
	// Using https://github.com/graphql-go/graphql
	g.Reset()
	g.writeFileHeader()
	g.writeImports([]string{"github.com/graphql-go/graphql"})
	g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
		Name:    g.config.Output + "schema.go",
		Content: g.String(),
	})
}

func (g *Generator) writeFileHeader() {
	_, err := g.WriteString("// DO NOT EDIT!!!\n")
	if err != nil {
		g.Error(err)
	}

	_, err = g.WriteString("package " + g.config.Package + "\n")
	if err != nil {
		g.Error(err)
	}
}

func (g *Generator) writeImports(packages []string) {
	_, err := g.WriteString("import (")
	if err != nil {
		g.Error(err)
	}
	for _, pkg := range packages {
		_, err = g.WriteString("\t\"" + pkg + "\"")
		if err != nil {
			g.Error(err)
		}
	}
	_, err = g.WriteString(")")
	if err != nil {
		g.Error(err)
	}
}

func (g *Generator) writeGoTypeDefinitions() {
	if g.genFiles == nil {
		g.genFiles = buildGenFilesMap(g.Request.FileToGenerate)
	}

	g.typeMap = buildBaseTypeMap()

	for _, fd := range g.Request.GraphqlFile {
		if _, ok := g.genFiles[fd.Name]; !ok {
			continue
		}

		for _, enumDef := range fd.Enums {
			if typDef, ok := g.typeMap[enumDef.GetName()]; !ok || typDef == nil {
				def := &GoEnumDefinition{
					EnumTypeDefinitionDescriptorProto: enumDef,
				}
				g.typeMap[enumDef.GetName()] = def
			}
		}
		for _, ifaceDef := range fd.Interfaces {
			if typDef, ok := g.typeMap[ifaceDef.GetName()]; !ok || typDef == nil {
				def := &GoInterfaceDefinition{
					InterfaceTypeDefinitionDescriptorProto: ifaceDef,
				}
				g.typeMap[ifaceDef.GetName()] = def
			}
		}
		for _, iObjDef := range fd.InputObjects {
			if typDef, ok := g.typeMap[iObjDef.GetName()]; !ok || typDef == nil {
				def := &GoInputObjectDefinition{
					InputObjectTypeDefinitionDescriptorProto: iObjDef,
				}
				g.typeMap[iObjDef.GetName()] = def
			}
		}
		for _, objDef := range fd.Objects {
			if objDef.GetName() == fd.Schema.Mutation.GetName() {
				continue
			}
			if objDef.GetName() == fd.Schema.Query.GetName() {
				continue
			}
			if typDef, ok := g.typeMap[objDef.GetName()]; !ok || typDef == nil {
				def := &GoObjectDefinition{
					ObjectTypeDefinitionDescriptorProto: objDef,
				}
				g.typeMap[objDef.GetName()] = def
			}
		}
		for _, unionDef := range fd.Unions {
			if typDef, ok := g.typeMap[unionDef.GetName()]; !ok || typDef == nil {
				def := &GoUnionDefinition{
					UnionTypeDefinitionDescriptorProto: unionDef,
				}
				g.typeMap[unionDef.GetName()] = def

			}
		}

		for name, typDef := range g.typeMap {
			switch def := typDef.(type) {
			case *GoInterfaceDefinition:
				def.TypeMap = g.typeMap
				g.typeMap[name] = def

			case *GoInputObjectDefinition:
				def.TypeMap = g.typeMap
				g.typeMap[name] = def

			case *GoObjectDefinition:
				def.TypeMap = g.typeMap
				g.typeMap[name] = def

			}
		}

		for _, typDef := range g.typeMap {
			g.WriteString(typDef.Definition() + "\n")
		}
	}
}

func buildGenFilesMap(filesToGenerate []string) map[string]bool {
	genFiles := make(map[string]bool)
	for _, file := range filesToGenerate {
		genFiles[file] = true
	}
	return genFiles
}

func buildBaseTypeMap() map[string]GoTypeDefinition {
	return map[string]GoTypeDefinition{
		"Int": &GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "string",
			},
		},
		"Float": &GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "float64",
			},
		},
		"String": &GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "string",
			},
		},
		"Boolean": &GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "bool",
			},
		},
		"ID": &GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "string",
			},
		},
	}
}
