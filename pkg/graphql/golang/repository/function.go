package repository

import (
	"crypto/sha256"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
	"strings"
)

type Function struct {
	*graphqlc.FieldDefinitionDescriptorProto
}

func (d *Function) Hash() string {
	list := make([]string, 0)
	for _, argDef := range d.Arguments {
		list = append(list, argDef.Name, argDef.Type.String())
	}
	hash := sha256.Sum256([]byte(strings.Join(list, ",")))
	return string(hash[:])
}
