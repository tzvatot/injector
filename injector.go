package injector

import (
	"fmt"
	"reflect"
	"strings"
)

type Injector interface {
	Register(beans ...interface{}) error
	Inject() error
}

type Engine struct {
	beans      []interface{}
	typeToBean map[string]interface{}
}

func NewEngine() *Engine {
	return &Engine{
		typeToBean: make(map[string]interface{}),
	}
}

func (inj *Engine) Register(beans ...interface{}) error {
	for _, bean := range beans {
		inj.beans = append(inj.beans, bean)
		val, err := inj.getValue(bean)
		if err != nil {
			return err
		}
		structType := val.Type()
		typeFullName := inj.getTypeFullName(structType)
		if inj.typeToBean[typeFullName] != nil {
			return fmt.Errorf("there's already a registered bean of type '%s'", typeFullName)
		}
		inj.typeToBean[typeFullName] = bean
	}
	return nil
}

func (inj *Engine) Inject() error {
	for _, bean := range inj.beans {
		if err := inj.injectBean(bean); err != nil {
			return err
		}
	}
	return nil
}

func (inj *Engine) injectBean(bean interface{}) error {
	val, err := inj.getValue(bean)
	if err != nil {
		return err
	}
	structType := val.Type()
	fields, err := inj.getFields(bean)
	if err != nil {
		return err
	}
	for _, field := range fields {
		fieldName := field.Name
		fieldType := field.Type
		if fieldType == structType {
			// not supporting self injection
			continue
		}
		childName := inj.getContainingTypeFullName(field)
		child := inj.typeToBean[childName]
		if child == nil {
			if childName, ok := field.Tag.Lookup("inject"); ok {
				child = inj.typeToBean[childName]
			}
		}
		if child != nil {
			f := val.FieldByName(fieldName)
			if f.CanAddr() && f.IsValid() {
				f.Set(reflect.ValueOf(child))
			}
		}
	}
	return nil
}

func (inj *Engine) getValue(bean interface{}) (val reflect.Value, err error) {
	val = reflect.ValueOf(bean) // could be any underlying type
	// if its a pointer, resolve its value
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	if val.Kind() != reflect.Struct {
		return val, fmt.Errorf("bean is not a struct type")
	}
	return val, nil
}

func (inj *Engine) getTypeFullName(refType reflect.Type) string {
	return fmt.Sprintf("%s", refType.String())
}

func (inj *Engine) getContainingTypeFullName(field reflect.StructField) string {
	typ := fmt.Sprintf("%s", field.Type)
	return strings.Replace(typ, "*", "", 1)
}

func (inj *Engine) getFields(bean interface{}) ([]reflect.StructField, error) {
	val, err := inj.getValue(bean)
	if err != nil {
		return nil, err
	}
	// now we grab our values as before (note: I assume table name should come from the struct type)
	structType := val.Type()
	fields := make([]reflect.StructField, 0)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fields = append(fields, field)
	}
	return fields, nil
}
