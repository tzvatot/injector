package injector

import (
	"fmt"
	"testing"
)

type A struct {
	num int
	s   string
	B   *B
}

type B struct {
	names []string
}


func TestInject(t *testing.T) {
	injector := NewEngine()
	b := &B{names: []string{"a", "b", "c"}}
  a := &A{num: 1, s:   "s"}

	if err := injector.Register(b, a); err != nil {
		t.Fatalf("failed to register: %v", err)
	}
  if err := injector.Inject(); err != nil {
  	t.Fatalf("failed to inject: %v", err)
	}
  fmt.Printf("%v", a)
}

func TestGetFields(t *testing.T) {
	injector := NewEngine()
  fields, err := injector.getFields(&A{})
  if err != nil {
		t.Fatalf("failed to get fields: %v", err)
	}
  if len(fields) != 3 {
  	t.Errorf("expected 3 fields")
	}
  for _, field := range fields {
		fieldName := field.Name
		fieldType := field.Type
		fmt.Printf("Name: %s, Type: %s\n", fieldName, fieldType)
	}
}
