Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
числа от 1 до 8 в рандомном порядке (возможно, что между ними проскочат нули) после будут нули

```
так происходит из-за того, что функция `merge` не проверяет на закрытие каналы и после их закрытия продолжает писать в выходной канал zero value. Сам `for range` по каналу идёт до тех пор пока канал не закроется. Чтобы решить проблему стоит исправить функцию `merge` добавив проверку на закрытие входящих каналов, а так же закрыть выходной канал после вычитки входящих. Например так:
```go
func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)  // закрываем выходной канал, когда горутина закончит свою работу
		for {
			select {
			case v, ok := <-a: // вроверяем закрыт ли читаемый канал a
				if !ok {
					a = nil
				} else {
					c <- v
				}
			case v, ok := <-b: // проверяем закрыт ли читаемый канал b
				if !ok {
					b = nil
				} else {
					c <- v
				}
			}

			if a == nil && b == nil {
				break
			}
		}
	}()
	return c
}
```
