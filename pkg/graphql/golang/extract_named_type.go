package golang

import (
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
)

func extractNamedType(typDef *graphqlc.TypeDescriptorProto) string {
	switch v := typDef.GetType().(type) {
	case *graphqlc.TypeDescriptorProto_ListType:
		return extractNamedType(v.ListType.GetType())
	case *graphqlc.TypeDescriptorProto_NamedType:
		return v.NamedType.GetName()
	case *graphqlc.TypeDescriptorProto_NonNullType:
		if listType := v.NonNullType.GetListType(); listType != nil {
			return extractNamedType(
				&graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_ListType{
						ListType: listType,
					},
				},
			)
		}
		if namedType := v.NonNullType.GetNamedType(); namedType != nil {
			return extractNamedType(
				&graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NamedType{
						NamedType: namedType,
					},
				},
			)
		}
	}
	return ""
}
