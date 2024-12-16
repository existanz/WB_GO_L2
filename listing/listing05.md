Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error

```
Интерфейсы в go содержат в себе не только данные но и указатель на тип. Не смотря на то, что функция `test()` возвращает `nil` она указывает её тип `*customError`, таким образом `err` содержит указатель на тип данных и хотя данные равны `nil` сам `err != nil`