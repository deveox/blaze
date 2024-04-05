package encoder

type fields struct {
	short       bool
	fields      []string
	currentPath string
	enabled     bool
}

func (e *fields) Has(fieldName string, short bool) bool {
	if !e.enabled {
		return true
	}

	if len(e.fields) > 0 {
		if e.currentPath == "" {
			e.currentPath = fieldName
		} else {
			e.currentPath = e.currentPath + "." + fieldName
		}
		if e.short && short {
			return true
		}
		for _, f := range e.fields {
			if f == e.currentPath {
				return true
			}
		}
		return false
	}
	return short
}

func (e *fields) Init(fields []string, short bool) {
	e.fields = fields
	e.short = short
	e.enabled = true
}
