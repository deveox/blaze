package encoder

import "strings"

type fields struct {
	short       bool
	fields      []string
	currentPath string
}

func (e *fields) Has(fieldName string, short bool) bool {
	if len(e.fields) == 0 {
		if e.short {
			return short
		}
		return true
	}

	if e.currentPath == "" {
		e.currentPath = fieldName
	} else {
		e.currentPath = e.currentPath + "." + fieldName
	}
	if e.short && short {
		return true
	}
	return e.has(e.currentPath)
}

func (e *fields) has(fieldName string) bool {
	for _, f := range e.fields {
		if f == fieldName {
			return true
		}
	}
	return false
}

func (e *fields) Init(fields []string, short bool) {
	e.fields = fields
	e.short = short
	for _, f := range e.fields {
		e.AddNested(f, false)
	}
}

func (e *fields) AddNested(field string, addNotNested bool) {
	if strings.Contains(field, ".") {
		parts := strings.Split(field, ".")
		for i := 0; i < len(parts)-1; i++ {
			field := strings.Join(parts[:i+1], ".")
			if !e.has(field) {
				e.fields = append(e.fields, field)
			}
		}
	}
	if addNotNested {
		e.fields = append(e.fields, field)
	}
}
