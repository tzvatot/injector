# injector
Dependency Injection engine for Go

## Usage
Assume you have the following interface:
```go
package injector

type Incrementor interface {
	Inc(int) int
}
```
This interface have an implementation as well:
```go
package injector

type MyImplementation struct {
}

func (m *MyImplementation) Inc(x int) int {
	return x + 1
}
```
And there's a struct that would like to use that interface, and inject it with a certain implementation:

```go
package injector

type MyStruct struct {
	MyIncrementor Incrementor `inject:"injector.MyImplementation"`
}
```
The following example can invoke the injection:
```go
incrementor := &MyImplementation{}
toInject := &MyStruct{}
if err := injector.Register(incrementor, toInject); err != nil {
	panic(fmt.Sprintf("failed to register: %v", err))
}
if err := injector.Inject(); err != nil {
	panic(fmt.Sprintf("failed to inject: %v", err))
}

result := toInject.MyIncrementor.Inc(5)
fmt.Println(result) // 6
```

## Notes
- Injected field must be exported in order for assignment to work.
- Injection of struct pointer is also supported without tags.
