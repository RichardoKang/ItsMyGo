package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// Field 代表一个字段
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema 代表一个对象的结构化信息
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

// GetField 获取字段信息
func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 仅处理导出字段
		if ast.IsExported(p.Name) && !p.Anonymous {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}

			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest)) //获取dest对象的值
	var fieldValues []interface{}
	// 遍历schema的字段列表, 通过反射获取每个字段的值
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface()) // eg: User{Name: "Tom", Age: 18} -> []interface{}{"Tom", 18}
	}
	return fieldValues
}
