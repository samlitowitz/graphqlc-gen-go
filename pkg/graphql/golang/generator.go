package golang

import (
	"github.com/dave/jennifer/jen"
	"github.com/samlitowitz/graphqlc-gen-echo/pkg/graphqlc/echo"
	cfgPkg "github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang/config"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
	"strings"
)

type Generator struct {
	*echo.Generator

	config   *Config
	genFiles map[string]bool
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
	unionTypesMap := make(map[string]map[string]struct{})

	for _, fd := range g.Request.GraphqlFile {
		if _, ok := g.genFiles[fd.Name]; !ok {
			continue
		}

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
		for _, unionDef := range fd.Unions {
			membersTypes := unionDef.GetMemberTypes()

			for _, memberType := range membersTypes {
				if _, ok := unionTypesMap[memberType.GetName()]; !ok {
					unionTypesMap[memberType.GetName()] = make(map[string]struct{})
				}
				unionTypesMap[memberType.GetName()][unionDef.GetName()] = struct{}{}
			}

			f.Type().Id(unionDef.GetName()).Interface(
				jen.Id("is" + unionDef.GetName()).Params(),
			)
		}
		for _, ifaceDef := range fd.GetInterfaces() {
			fields := ifaceDef.GetFields()
			fnDefs := make(jen.Statement, 0, len(fields))
			for _, field := range fields {
				typ := field.GetType()
				stmt := jen.Id("Get" + strings.Title(field.GetName())).Params().Add(getGoType(typ, g.config.ScalarMap, true))
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
		for _, objDef := range fd.Objects {
			if objDef.GetName() == fd.Schema.Mutation.GetName() {
				continue
			}
			if objDef.GetName() == fd.Schema.Query.GetName() {
				continue
			}
			fields := objDef.GetFields()
			fieldDefs := make(jen.Statement, 0, len(fields))
			for _, field := range fields {
				typ := field.GetType()
				stmt := jen.Id(strings.Title(field.GetName())).Add(getGoType(typ, g.config.ScalarMap, true))
				fieldDefs = append(fieldDefs, stmt)
			}
			f.Type().Id(objDef.GetName()).Struct(fieldDefs...)

			ifaces := objDef.GetImplements()
			for _, iface := range ifaces {
				ifaceFields := iface.GetFields()
				for _, ifaceField := range ifaceFields {
					typ := ifaceField.GetType()
					f.Func().Params(
						jen.Id("o").Id(objDef.GetName()),
					).Id("Get" + strings.Title(ifaceField.GetName())).Params().Add(getGoType(typ, g.config.ScalarMap, true)).Block(
						jen.Return(
							jen.Id("o").Dot(strings.Title(ifaceField.GetName())),
						),
					)
				}
			}

			unions, ok := unionTypesMap[objDef.GetName()]
			if !ok {
				continue
			}

			for union := range unions {
				f.Func().Params(
					jen.Id("o").Id(objDef.GetName()),
				).Id("is" + union).Params().Block()
			}
		}
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

func getGoType(typDef *graphqlc.TypeDescriptorProto, customTypes map[string]cfgPkg.ScalarType, nullable bool) *jen.Statement {
	switch v := typDef.GetType().(type) {
	case *graphqlc.TypeDescriptorProto_ListType:
		return jen.Index().Add(getGoType(v.ListType.GetType(), customTypes, true))
	case *graphqlc.TypeDescriptorProto_NamedType:
		switch v.NamedType.GetName() {
		case "Int":
			if nullable {
				return jen.Op("*").Int()
			}
			return jen.Int()
		case "Float":
			if nullable {
				return jen.Op("*").Float64()
			}
			return jen.Float64()
		case "String":
			if nullable {
				return jen.Op("*").String()
			}
			return jen.String()
		case "Boolean":
			if nullable {
				return jen.Op("*").Bool()
			}
			return jen.Bool()
		case "ID":
			if nullable {
				return jen.Op("*").String()
			}
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
