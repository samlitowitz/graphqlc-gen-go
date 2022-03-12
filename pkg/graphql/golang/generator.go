package golang

import (
	"github.com/dave/jennifer/jen"
	"github.com/samlitowitz/graphqlc-gen-echo/pkg/graphqlc/echo"
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
		Name:    g.config.OutputPrefix + "_types.go",
		Content: g.String(),
	})

	// Go repository interface definitions
	g.Reset()
	f = jen.NewFile(g.config.Package)
	g.writeRepositoryInterfaces(f)
	g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
		Name:    g.config.OutputPrefix + "_repositories.go",
		Content: g.String(),
	})

	// GraphQL definitions (types and schema)
	// Using https://github.com/graphql-go/graphql
	//g.Reset()
	//g.writeFileHeader()
	//g.writeImports([]string{"github.com/graphql-go/graphql"})
	//g.Response.File = append(g.Response.File, &graphqlc.CodeGeneratorResponse_File{
	//	Name:    g.config.OutputPrefix + "_graphql.go",
	//	Content: g.String(),
	//})
}

func (g *Generator) writeRepositoryInterfaces(f *jen.File) {
	if g.genFiles == nil {
		g.genFiles = buildGenFilesMap(g.Request.FileToGenerate)
	}
	repoFnsByType := make(map[string]map[string]*jen.Statement)
	for _, fd := range g.Request.GraphqlFile {
		if _, ok := g.genFiles[fd.Name]; !ok {
			continue
		}
		for _, objDef := range fd.GetObjects() {
			fields := objDef.GetFields()
			for _, field := range fields {
				fieldArgs := field.GetArguments()
				if len(fieldArgs) == 0 {
					continue
				}
				fieldTypeName := extractNamedType(field.GetType())
				if _, ok := repoFnsByType[fieldTypeName]; !ok {
					repoFnsByType[fieldTypeName] = make(map[string]*jen.Statement)
				}
				fieldName := field.GetName()
				if _, ok := repoFnsByType[fieldTypeName][fieldName]; !ok {
					argStmts := make(jen.Statement, 0, len(fieldArgs))
					for _, fieldArg := range fieldArgs {
						argStmts = append(
							argStmts,
							jen.Id(fieldArg.GetName()).Add(getGoType(fieldArg.GetType(), g.config.ScalarMap, true)),
						)
					}

					stmt := jen.Id(strings.Title(fieldName)).Params(argStmts...).Add(getGoType(field.GetType(), g.config.ScalarMap, true))
					repoFnsByType[fieldTypeName][fieldName] = stmt
				}
			}
		}
	}

	for typ, repoFns := range repoFnsByType {
		fnStmts := make(jen.Statement, 0, len(repoFns))
		for _, fnStmt := range repoFns {
			fnStmts = append(fnStmts, fnStmt)
		}
		f.Type().Id(typ + "Repository").Interface(fnStmts...)
	}
	f.Render(g)
}

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
