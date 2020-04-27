package injector

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultLogFormat       = "[%lvl%]: %time% - %msg%"
	defaultTimestampFormat = time.RFC3339
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

var log *logrus.Logger

func init() {
	setupLogger()
}

func setupLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	log = logrus.StandardLogger()
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
		log.Debugf("type full name: %s\n", typeFullName)
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
		log.Debugf("Name: %s, Type: %s\n", fieldName, fieldType)
		if fieldType == structType {
			log.Infof("not supporting self injection")
			continue
		}
		childName := inj.getContainingTypeFullName(field)
		log.Debugf("chile name: %s\n", childName)
		child := inj.typeToBean[childName]
		if child == nil {
			if childName, ok := field.Tag.Lookup("inject"); ok {
				log.Debugf("chile name: %s\n", childName)
				child = inj.typeToBean[childName]
			}
		}
		if child != nil {
			f := val.FieldByName(fieldName)
			if f.CanAddr() && f.IsValid() {
				f.Set(reflect.ValueOf(child))
				log.Debugf("Injected %s", child)
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
