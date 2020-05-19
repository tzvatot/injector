package injector

import (
	"fmt"
	"strings"
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

type Incrementor interface {
	Inc(int) int
}

type MyImplementation struct {
}

func (m *MyImplementation) Inc(x int) int {
	return x + 1
}

type MyStruct struct {
	MyIncrementor Incrementor `inject:"injector.MyImplementation"`
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
		t.Fatalf("injection failed - property is nil")
	}
	fmt.Printf("%v", a)
}

func TestInjectInterfaceImplementation(t *testing.T) {
	injector := NewEngine()
	incrementor := &MyImplementation{}
	toInject := &MyStruct{}

	if err := injector.Register(incrementor, toInject); err != nil {
		t.Fatalf("failed to register: %v", err)
	}
	if err := injector.Inject(); err != nil {
		t.Fatalf("failed to inject: %v", err)
	}
	if toInject.MyIncrementor == nil {
		t.Fatalf("injection failed - property is nil")
	}
	actual := toInject.MyIncrementor.Inc(5)
	if actual != 6 {
		t.Fatalf("Unexpected result: %d", actual)
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
		t.Fatalf("injection failed - property is nil")
	}
	if bar.Foo == nil {
		t.Fatalf("injection failed - property is nil")
	}
	if bar.Foo.Bar != bar || foo.Bar.Foo != foo {
		t.Fatalf("injection failed - wrong reference")
	}
	fmt.Printf("%v", foo)
	fmt.Printf("%v", bar)
}

func TestDoubleRegister(t *testing.T) {
	injector := NewEngine()
	a := &A{}

	err := injector.Register(a, a)
	if err == nil {
		t.Fatalf("expected to fail")
	}
	if !strings.Contains(err.Error(), "already a registered") {
		t.Fatalf("wrong error message")
	}
}

type Collection struct {
	Elements []Element `inject:"injector.ElementImplementorA,injector.ElementImplementorB"`
}

type Element interface {
  Do()
}

type ElementImplementorA struct {}

func (impl *ElementImplementorA) Do() {
	fmt.Println("A")
}

type ElementImplementorB struct {}

func (impl *ElementImplementorB) Do() {
	fmt.Println("B")
}

func TestSliceInjection(t *testing.T) {
	injector := NewEngine()
	implA := &ElementImplementorA{}
	implB := &ElementImplementorB{}
	collection := &Collection{}
	if err := injector.Register(collection, implA, implB); err != nil {
		t.Fatalf("failed to register: %v", err)
	}
	if err := injector.Inject(); err != nil {
		t.Fatalf("failed to inject: %v", err)
	}
	if len(collection.Elements) < 2 {
		t.Fatalf("failed to inject multiple implementations")
	}
	if collection.Elements[0] == implA && collection.Elements[1] == implB {
		return // works
	}
	if collection.Elements[1] == implA && collection.Elements[0] == implB {
		return // works
	}
  t.Fatalf("wrong implementations injected")
}

func TestGetFields(t *testing.T) {
	injector := NewEngine()
	fields, err := injector.getFields(&A{})
	if err != nil {
		t.Fatalf("failed to get fields: %v", err)
	}
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields")
	}
	for _, field := range fields {
		fieldName := field.Name
		fieldType := field.Type
		fmt.Printf("Name: %s, Type: %s\n", fieldName, fieldType)
	}
}
