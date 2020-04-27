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

type MyInterface interface {
	Foo(int) int
}

type MyImplementation struct {
}

func (m *MyImplementation) Foo(x int) int {
	return x + 1
}

type MyStruct struct {
	MyInterfaceField MyInterface `inject:"injector.MyImplementation"`
}

type Foo struct {
	Bar *Bar `inject:"injector.Bar"`
}

type Bar struct {
	Foo *Foo `inject:"injector.Foo"`
}

func TestInjectStructPointer(t *testing.T) {
	injector := NewEngine()
	b := &B{names: []string{"a", "b", "c"}}
	a := &A{num: 1, s: "s"}

	if err := injector.Register(b, a); err != nil {
		t.Fatalf("failed to register: %v", err)
	}
	if err := injector.Inject(); err != nil {
		t.Fatalf("failed to inject: %v", err)
	}
	if a.B == nil {
		t.Errorf("injection failed - property is nil")
	}
	fmt.Printf("%v", a)
}

func TestInjectInterfaceImplementation(t *testing.T) {
	injector := NewEngine()
	impl := &MyImplementation{}
	toInject := &MyStruct{}

	if err := injector.Register(impl, toInject); err != nil {
		t.Fatalf("failed to register: %v", err)
	}
	if err := injector.Inject(); err != nil {
		t.Fatalf("failed to inject: %v", err)
	}
	if toInject.MyInterfaceField == nil {
		t.Errorf("injection failed - property is nil")
	}
	actual := toInject.MyInterfaceField.Foo(5)
	if actual != 6 {
		t.Errorf("Unexpected result: %d", actual)
	}
	fmt.Printf("%v", toInject)
}


func TestCyclicInjection(t *testing.T) {
	injector := NewEngine()
	foo := &Foo{}
	bar := &Bar{}

	if err := injector.Register(foo, bar); err != nil {
		t.Fatalf("failed to register: %v", err)
	}
	if err := injector.Inject(); err != nil {
		t.Fatalf("failed to inject: %v", err)
	}
	if foo.Bar == nil {
		t.Errorf("injection failed - property is nil")
	}
	if bar.Foo == nil {
		t.Errorf("injection failed - property is nil")
	}
	if bar.Foo.Bar != bar || foo.Bar.Foo != foo {
		t.Errorf("injection failed - wrong reference")
	}
	fmt.Printf("%v", foo)
	fmt.Printf("%v", bar)
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
