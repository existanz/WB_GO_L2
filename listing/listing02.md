Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
2
1

```
defer выполняется после return, но до возвращения управления вызывающей функции
в случае с именованными возвращаемыми переменными, он успевает изменять их значения
поэтому в первом случае мы получаем 2, а во втором defer не успевает изменить 
переменную и мы получаем 1