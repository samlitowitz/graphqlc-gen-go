package repository

type Definition struct {
	typName   string
	functions map[string]*Function
}

func NewDefinition(typName string) *Definition {
	return &Definition{
		typName:   typName,
		functions: make(map[string]*Function),
	}
}

func (d *Definition) TypeName() string {
	return d.typName
}

func (d *Definition) HasFunction(fnName string) bool {
	_, ok := d.functions[fnName]
	return ok
}
