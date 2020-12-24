package golang

type GoNamedDefinition struct {
	name string
}

func (typDef *GoNamedDefinition) UnqualifiedName() string {
	return typDef.name
}

func (typDef *GoNamedDefinition) Definition() string {
	return ""
}
