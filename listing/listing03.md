Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
<nil>
false

```
Чтобы понять это поведение, необходимо посмотреть, что из себя представляет интерфейс в go:
```go
// src/runtime/runtime2.go 205:
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
```
по сути это структура, которая содержит в себе указатель на тип (на самом деле `itab`, кроме типа интерфейса содержит ещё и ссылку на перечень методов, и hash для typa switch'а) и на сами данные. В момент объявления переменной err мы указываем её тип, а в качестве значения присваиваем nil. Поэтому 	`fmt.Println(err)` выводит `<nil>`, но `err != nil`  
Пустым интерфейсам не зачем хранить таблицу методов, поэтому для пустых интерфейсов используется другая структура.
```go
// src/runtime/runtime2.go 210:
type eface struct {
	_type *_type
	data  unsafe.Pointer
}
```
