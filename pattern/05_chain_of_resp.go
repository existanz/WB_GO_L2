/*
Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
	Цепочка вызовов позволяет передавать запросы по цепочке обработчиков,
	каждый из которых может обработать запрос или передать его дальше по цепочке
	+ гибкость: легко добавлять изменять обработчики
	+ снижение связности: обработчики не зависят друг от друга
	+ упрощение: позволяет избежать сложных условных конструкций
	+ SO из SOLID: реализует принцип единственной ответственности и принцип открытости/закрытости
	- отладка: может быть сложно отследить, какой обработчик обработал запрос
	- производительность: при длинной цепочке возможно снижение производительности

	Это довольно частый паттерн, используемый во многих системах.
	Например логирование: Обработчики могут быть использованы для логирования сообщений разного
	уровня (debug, info, warning, error).
  Или валидация данных: Валидация запросов может быть организована с помощью цепочки
	обработчиков, где каждый обработчик отвечает за свою часть валидации.
*/

package pattern

import (
	"fmt"
	"slices"
	"time"
)

type CV struct {
	name       string
	experience int
	birthDate  time.Time
	stack      []string
}

type Handler interface {
	SetNext(handler Handler)
	Handle(request CV)
}

type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *BaseHandler) Handle(request CV) {
	if h.next != nil {
		h.next.Handle(request)
	}
}

type HrHandler struct {
	BaseHandler
}

func (h *HrHandler) Handle(request CV) {
	if request.experience < 3 || !slices.Contains(request.stack, "Go") {
		fmt.Printf("[HR] Sorry, %s. Go and learn more Go\n", request.name)
		return
	}
	if request.birthDate.Before(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)) {
		fmt.Printf("[HR] Really sorry, %s. Too old for this shit \n", request.name)
		return
	}

	h.BaseHandler.Handle(request)
}

type TarologHandler struct {
	BaseHandler
}

func (h *TarologHandler) Handle(request CV) {
	y, m, d := request.birthDate.Date()
	if (y+int(m)+d)%23 != 0 {
		fmt.Printf("[Tarolog] The stars don't agree with hiring you, %s \n", request.name)
		return
	}

	h.BaseHandler.Handle(request)
}

type TeamLeadHandler struct {
	BaseHandler
}

func (h *TeamLeadHandler) Handle(request CV) {
	if !slices.Contains(request.stack, "like boss's jokes") {
		fmt.Println("[TeamLead] I don't like this person")
		return
	}
	fmt.Println("[TeamLead] Welcome aboard!")
	h.BaseHandler.Handle(request)
}

func mainCoR() {
	cvs := []CV{
		{
			name:       "John Doe",
			experience: 4,
			birthDate:  time.Date(1998, 12, 10, 0, 0, 0, 0, time.UTC),
			stack:      []string{"PHP", "Go", "Guitar", "Trombone", "like boss's jokes"},
		},
		{
			name:       "Jane Doe",
			experience: 7,
			birthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			stack:      []string{"Python", "С++", "Brainfuck"},
		},
		{
			name:       "Jin Di",
			experience: 3,
			birthDate:  time.Date(1990, 1, 10, 0, 0, 0, 0, time.UTC),
			stack:      []string{"Go", "Run", "like boss's jokes"},
		},
	}
	hr := HrHandler{}
	taro := TarologHandler{}
	lead := TeamLeadHandler{}
	hr.SetNext(&taro)
	taro.SetNext(&lead)
	for _, cv := range cvs {
		hr.Handle(cv)
	}
}
