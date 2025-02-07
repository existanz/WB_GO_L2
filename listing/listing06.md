Что выведет программа? Объяснить вывод программы. Рассказать про внутреннее устройство слайсов и что происходит при передачи их в качестве аргументов функции.

```go
package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
```

Ответ:
```
["3", "2", "3"]

```
Чтобы понять почему так получилось, посмотрим что такое слайс (динамический массив языка Go).
Слайс, как и почти всё в Go это структура, и она имеет следующие поля
```go
//src/runtime/slice.go 15:
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```
len - длина слайса, cap - капасити, array - указатель на базовый массив.  
При передаче в функцию аргументы в Go всегда копируются, то есть создаётся новый слайс, в который мы скопируем len, cap и самое главное array (указатель на базовый массив). Таким образом при оперировании с элементами нового слайса мы будем непосредственно оперировать c элементами массива, а значит и исходного слайса `s`.  
Но всё может измениться, когда мы добавляем элементы в слайс. При append'е в случае если len превышает cap происходит аллокация новой памяти под базовый массив. Таким образом слайс `i` внутри функции `modifySlice` перестаёт ссылаться на исходный массив. И значит все последующие изменения не коснутся слайса `s`.
