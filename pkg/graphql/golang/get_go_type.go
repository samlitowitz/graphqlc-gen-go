package golang

import (
	"github.com/dave/jennifer/jen"
	cfgPkg "github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang/config"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
)

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
