package golang

import (
	"fmt"
	"github.com/samlitowitz/graphqlc-gen-echo/pkg/graphqlc/echo"
	cfgPkg "github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang/config"
	"github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang/repository"
	typPkg "github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang/type"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
)

type Generator struct {
	*echo.Generator

	config   *Config
	genFiles map[string]bool
	typeMap  map[string]typPkg.Definition
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
	g.writeTypeDefinitions()
	g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
		Name:    g.config.Output + "_types.go",
		Content: g.String(),
	})

	// Go repository interface definitions
	g.Reset()
	g.writeFileHeader()
	g.writeRepositoryInterfaces()
	g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
		Name:    g.config.Output + "_types.go",
		Content: g.String(),
	})

	// GraphQL definitions (types and schema)
	// Using https://github.com/graphql-go/graphql
	g.Reset()
	g.writeFileHeader()
	g.writeImports([]string{"github.com/graphql-go/graphql"})
	g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
		Name:    g.config.Output + "_graphql.go",
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
	_, err = g.WriteString(")\n")
	if err != nil {
		g.Error(err)
	}
}

func (g *Generator) writeRepositoryInterfaces() {
	if g.genFiles == nil {
		g.genFiles = buildGenFilesMap(g.Request.FileToGenerate)
	}

	repoDefs := make(map[string]map[string]*repository.Function)

	for _, fd := range g.Request.GraphqlFile {
		if _, ok := g.genFiles[fd.Name]; !ok {
			continue
		}
		for _, objDef := range fd.Objects {
			if _, ok := g.typeMap[objDef.GetName()]; !ok {
				g.Error(fmt.Errorf("%s missing type definition", objDef.GetName()))
			}

			for _, fieldDef := range objDef.Fields {
				if len(fieldDef.Arguments) == 0 {
					continue
				}
				repoFnDefs, ok := repoDefs[fieldDef.Type.String()]
				if !ok {
					repoDef := repository.NewDefinition(objDef.G)
					repoDefs[objDef.GetName()] = repoDef
				}
				repoFnDef := &repository.Function{FieldDefinitionDescriptorProto: fieldDef}
				if _, ok := repoFnDefs[string(repoFnDef.Hash())]; ok {
					continue
				}
				repoFnDefs[string(repoFnDef.Hash())] = repoFnDef
			}
		}
	}

	for typName, repoFnDefs := range repoDefs {

	}
}

func (g *Generator) writeTypeDefinitions() {
	if g.genFiles == nil {
		g.genFiles = buildGenFilesMap(g.Request.FileToGenerate)
	}

	g.typeMap = buildBaseTypeMap()

	for _, fd := range g.Request.GraphqlFile {
		if _, ok := g.genFiles[fd.Name]; !ok {
			continue
		}

		importPackages := make([]string, 0)

		for _, scalarDef := range fd.Scalars {
			var customType *cfgPkg.ScalarType
			if _, ok := g.config.ScalarMap[scalarDef.GetName()]; ok {
				customType = &cfgPkg.ScalarType{}
				*customType = g.config.ScalarMap[scalarDef.GetName()]
				importPackages = append(importPackages, customType.Package)
			}
			if typDef, ok := g.typeMap[scalarDef.GetName()]; !ok || typDef == nil {
				def := &typPkg.GoScalarDefinition{
					ScalarTypeDefinitionDescriptorProto: scalarDef,
					CustomType:                          customType,
				}
				g.typeMap[scalarDef.GetName()] = def
			}
		}
		for _, enumDef := range fd.Enums {
			if typDef, ok := g.typeMap[enumDef.GetName()]; !ok || typDef == nil {
				def := &typPkg.GoEnumDefinition{
					EnumTypeDefinitionDescriptorProto: enumDef,
				}
				g.typeMap[enumDef.GetName()] = def
			}
		}
		for _, ifaceDef := range fd.Interfaces {
			if typDef, ok := g.typeMap[ifaceDef.GetName()]; !ok || typDef == nil {
				def := &typPkg.GoInterfaceDefinition{
					InterfaceTypeDefinitionDescriptorProto: ifaceDef,
				}
				g.typeMap[ifaceDef.GetName()] = def
			}
		}
		for _, iObjDef := range fd.InputObjects {
			if typDef, ok := g.typeMap[iObjDef.GetName()]; !ok || typDef == nil {
				def := &typPkg.GoInputObjectDefinition{
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
				def := &typPkg.GoObjectDefinition{
					ObjectTypeDefinitionDescriptorProto: objDef,
				}
				g.typeMap[objDef.GetName()] = def
			}
		}
		for _, unionDef := range fd.Unions {
			if typDef, ok := g.typeMap[unionDef.GetName()]; !ok || typDef == nil {
				def := &typPkg.GoUnionDefinition{
					UnionTypeDefinitionDescriptorProto: unionDef,
				}
				g.typeMap[unionDef.GetName()] = def

			}
		}

		for name, typDef := range g.typeMap {
			switch def := typDef.(type) {
			case *typPkg.GoEnumDefinition:
				g.typeMap[name] = def

			case *typPkg.GoInterfaceDefinition:
				def.TypeMap = g.typeMap
				g.typeMap[name] = def

			case *typPkg.GoInputObjectDefinition:
				def.TypeMap = g.typeMap
				g.typeMap[name] = def

			case *typPkg.GoObjectDefinition:
				def.TypeMap = g.typeMap
				g.typeMap[name] = def
			}
		}

		g.writeImports(importPackages)
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

func buildBaseTypeMap() map[string]typPkg.Definition {
	return map[string]typPkg.Definition{
		"Int": &typPkg.GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "string",
			},
		},
		"Float": &typPkg.GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "float64",
			},
		},
		"String": &typPkg.GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "string",
			},
		},
		"Boolean": &typPkg.GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "bool",
			},
		},
		"ID": &typPkg.GoScalarDefinition{
			ScalarTypeDefinitionDescriptorProto: &graphqlc.ScalarTypeDefinitionDescriptorProto{
				Name: "string",
			},
		},
	}
}
