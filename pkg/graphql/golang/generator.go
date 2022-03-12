package golang

import (
	"github.com/dave/jennifer/jen"
	"github.com/samlitowitz/graphqlc-gen-echo/pkg/graphqlc/echo"
	cfgPkg "github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang/config"
	typPkg "github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang/type"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
	"strings"
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
	//g.writeFileHeader()
	f := jen.NewFile(g.config.Package)
	g.writeTypeDefinitions(f)
	g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
		Name:    g.config.Output + "types.go",
		Content: g.String(),
	})

	// Go repository interface definitions
	//g.Reset()
	//g.writeFileHeader()
	//g.writeRepositoryInterfaces()
	//g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
	//	Name:    g.config.Output + "_types.go",
	//	Content: g.String(),
	//})

	// GraphQL definitions (types and schema)
	// Using https://github.com/graphql-go/graphql
	//g.Reset()
	//g.writeFileHeader()
	//g.writeImports([]string{"github.com/graphql-go/graphql"})
	//g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
	//	Name:    g.config.Output + "_graphql.go",
	//	Content: g.String(),
	//})
}

//func (g *Generator) writeFileHeader() {
//	_, err := g.WriteString("// DO NOT EDIT!!!\n")
//	if err != nil {
//		g.Error(err)
//	}
//
//	_, err = g.WriteString("package " + g.config.Package + "\n")
//	if err != nil {
//		g.Error(err)
//	}
//}
//
//func (g *Generator) writeImports(packages []string) {
//	_, err := g.WriteString("import (")
//	if err != nil {
//		g.Error(err)
//	}
//	for _, pkg := range packages {
//		_, err = g.WriteString("\t\"" + pkg + "\"")
//		if err != nil {
//			g.Error(err)
//		}
//	}
//	_, err = g.WriteString(")\n")
//	if err != nil {
//		g.Error(err)
//	}
//}

//func (g *Generator) writeRepositoryInterfaces() {
//	if g.genFiles == nil {
//		g.genFiles = buildGenFilesMap(g.Request.FileToGenerate)
//	}
//
//	repoDefs := make(map[string]map[string]*repository.Function)
//
//	for _, fd := range g.Request.GraphqlFile {
//		if _, ok := g.genFiles[fd.Name]; !ok {
//			continue
//		}
//		for _, objDef := range fd.Objects {
//			if _, ok := g.typeMap[objDef.GetName()]; !ok {
//				g.Error(fmt.Errorf("%s missing type definition", objDef.GetName()))
//			}
//
//			for _, fieldDef := range objDef.Fields {
//				if len(fieldDef.Arguments) == 0 {
//					continue
//				}
//				repoFnDefs, ok := repoDefs[fieldDef.Type.String()]
//				if !ok {
//					repoDef := repository.NewDefinition(objDef.G)
//					repoDefs[objDef.GetName()] = repoDef
//				}
//				repoFnDef := &repository.Function{FieldDefinitionDescriptorProto: fieldDef}
//				if _, ok := repoFnDefs[string(repoFnDef.Hash())]; ok {
//					continue
//				}
//				repoFnDefs[string(repoFnDef.Hash())] = repoFnDef
//			}
//		}
//	}
//
//	for typName, repoFnDefs := range repoDefs {
//
//	}
//}

func (g *Generator) writeTypeDefinitions(f *jen.File) {
	if g.genFiles == nil {
		g.genFiles = buildGenFilesMap(g.Request.FileToGenerate)
	}

	//g.typeMap = buildBaseTypeMap()

	for _, fd := range g.Request.GraphqlFile {
		if _, ok := g.genFiles[fd.Name]; !ok {
			continue
		}

		//importPackages := make([]string, 0)
		//
		//for _, scalarDef := range fd.Scalars {
		//	var customType *cfgPkg.ScalarType
		//	if _, ok := g.config.ScalarMap[scalarDef.GetName()]; ok {
		//		customType = &cfgPkg.ScalarType{}
		//		*customType = g.config.ScalarMap[scalarDef.GetName()]
		//		importPackages = append(importPackages, customType.Package)
		//	}
		//	if typDef, ok := g.typeMap[scalarDef.GetName()]; !ok || typDef == nil {
		//		def := &typPkg.GoScalarDefinition{
		//			ScalarTypeDefinitionDescriptorProto: scalarDef,
		//			CustomType:                          customType,
		//		}
		//		g.typeMap[scalarDef.GetName()] = def
		//	}
		//}
		for _, enumDef := range fd.GetEnums() {
			enumValues := enumDef.GetValues()
			f.Type().Id(enumDef.GetName()).Int()
			valDefs := make(jen.Statement, 0, len(enumValues))
			for i, val := range enumValues {
				valDefs = append(
					valDefs,
					jen.Id(enumDef.GetName()+"_"+val.GetValue()).Qual("", enumDef.GetName()).Op("=").Lit(i),
				)
			}
			f.Const().Defs(valDefs...)
		}
		for _, ifaceDef := range fd.GetInterfaces() {
			fields := ifaceDef.GetFields()
			fnDefs := make(jen.Statement, 0, len(fields))
			for _, field := range fields {
				typ := field.GetType()
				stmt := jen.Id(strings.Title(field.GetName())).Params().Add(getGoType(typ, g.config.ScalarMap, true))
				fnDefs = append(fnDefs, stmt)
			}
			f.Type().Id(ifaceDef.GetName()).Interface(fnDefs...)
		}
		for _, iObjDef := range fd.InputObjects {
			fields := iObjDef.GetFields()
			fieldDefs := make(jen.Statement, 0, len(fields))
			for _, field := range fields {
				typ := field.GetType()
				stmt := jen.Id(strings.Title(field.GetName())).Add(getGoType(typ, g.config.ScalarMap, true))
				fieldDefs = append(fieldDefs, stmt)
			}
			f.Type().Id(iObjDef.GetName()).Struct(fieldDefs...)
		}
		//for _, objDef := range fd.Objects {
		//	if objDef.GetName() == fd.Schema.Mutation.GetName() {
		//		continue
		//	}
		//	if objDef.GetName() == fd.Schema.Query.GetName() {
		//		continue
		//	}
		//	if typDef, ok := g.typeMap[objDef.GetName()]; !ok || typDef == nil {
		//		def := &typPkg.GoObjectDefinition{
		//			ObjectTypeDefinitionDescriptorProto: objDef,
		//		}
		//		g.typeMap[objDef.GetName()] = def
		//	}
		//}
		//for _, unionDef := range fd.Unions {
		//	if typDef, ok := g.typeMap[unionDef.GetName()]; !ok || typDef == nil {
		//		def := &typPkg.GoUnionDefinition{
		//			UnionTypeDefinitionDescriptorProto: unionDef,
		//		}
		//		g.typeMap[unionDef.GetName()] = def
		//
		//	}
		//}
		//
		//for name, typDef := range g.typeMap {
		//	switch def := typDef.(type) {
		//	case *typPkg.GoEnumDefinition:
		//		g.typeMap[name] = def
		//
		//	case *typPkg.GoInterfaceDefinition:
		//		def.TypeMap = g.typeMap
		//		g.typeMap[name] = def
		//
		//	case *typPkg.GoInputObjectDefinition:
		//		def.TypeMap = g.typeMap
		//		g.typeMap[name] = def
		//
		//	case *typPkg.GoObjectDefinition:
		//		def.TypeMap = g.typeMap
		//		g.typeMap[name] = def
		//	}
		//}
		//
		//g.writeImports(importPackages)
		//for _, typDef := range g.typeMap {
		//	g.WriteString(typDef.Definition() + "\n")
		//}
	}
	f.Render(g)
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

func getGoType(typDef *graphqlc.TypeDescriptorProto, customTypes map[string]cfgPkg.ScalarType, nullable bool) *jen.Statement {
	switch v := typDef.GetType().(type) {
	case *graphqlc.TypeDescriptorProto_ListType:
		return jen.Index().Add(getGoType(v.ListType.GetType(), customTypes, true))
	case *graphqlc.TypeDescriptorProto_NamedType:
		switch v.NamedType.GetName() {
		case "Int":
			return jen.Int()
		case "Float":
			return jen.Float64()
		case "String":
			return jen.String()
		case "Boolean":
			return jen.Bool()
		case "ID":
			return jen.String()
		default:
			// custom scalar
			if customTyp, ok := customTypes[v.NamedType.GetName()]; ok {
				if nullable {
					return jen.Op("*").Qual(customTyp.Package, customTyp.Type)
				}
				return jen.Qual(customTyp.Package, customTyp.Type)
			}
			if nullable {
				return jen.Op("*").Qual("", v.NamedType.GetName())
			}
			return jen.Qual("", v.NamedType.GetName())
		}
	case *graphqlc.TypeDescriptorProto_NonNullType:
		if listType := v.NonNullType.GetListType(); listType != nil {
			return getGoType(
				&graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_ListType{
						ListType: listType,
					},
				},
				customTypes,
				false,
			)
		}
		if namedType := v.NonNullType.GetNamedType(); namedType != nil {
			return getGoType(
				&graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NamedType{
						NamedType: namedType,
					},
				},
				customTypes,
				false,
			)
		}
	}
	return nil
}
