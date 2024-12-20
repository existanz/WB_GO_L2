/*
Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
	Паттерн "фабричный метод" позволяет создавать объекты через интерфейс, оставляя решение о том,
	какую конкретную структуру создавать наследникам.

	+ O из SOLID: Реализует принцип открытости/закрытости.
	+ гибкость: Упрощает добавление новых продуктов в программу.
	- бойлерплейт: Может привести к созданию больших параллельных иерархий классов,
	так как для каждого класса продукта надо создать свой подкласс создателя

	Фреймворки: Многие веб-фреймворки используют фабричный метод для создания компонентов
	(например, создание форм, кнопок и т.д.).
	Можно использовать фабричный метод для создания соединений с различными типами баз
	данных (MySQL, PostgreSQL и т.д.).
*/

package pattern

import "fmt"

type Pizza interface {
	Prepare()
	Cook()
	Cut()
	Box()
}

type CheesePizza struct {
	size int
	fill string
}

func (p *CheesePizza) Prepare() {
	fmt.Println("Preparing Cheese Pizza")
}

func (p *CheesePizza) Cook() {
	fmt.Println("Cooking Cheese Pizza")
}

func (p *CheesePizza) Cut() {
	fmt.Println("Cutting Cheese Pizza")
}

func (p *CheesePizza) Box() {
	fmt.Println("Boxing Cheese Pizza")
}

func NewCheesePizza() Pizza {
	return &CheesePizza{
		size: 45,
		fill: "mozarella",
	}
}

type PepperoniPizza struct {
	size int
	fill string
}

func (p *PepperoniPizza) Prepare() {
	fmt.Println("Preparing Pepperoni Pizza")
}

func (p *PepperoniPizza) Cook() {
	fmt.Println("Cooking Pepperoni Pizza")
}

func (p *PepperoniPizza) Cut() {
	fmt.Println("Cutting Pepperoni Pizza")
}

func (p *PepperoniPizza) Box() {
	fmt.Println("Boxing Pepperoni Pizza")
}

func NewPepperoniPizza() Pizza {
	return &PepperoniPizza{
		size: 30,
		fill: "sousage",
	}
}

type PizzaStore interface {
	OrderPizza() Pizza
}

type CheezePizzaStore struct {
	PizzaStore
}

func (c *CheezePizzaStore) OrderPizza() Pizza {
	pizza := NewCheesePizza()
	pizza.Prepare()
	pizza.Cook()
	pizza.Cut()
	pizza.Box()
	return pizza
}

type PepperoniPizzaStore struct {
	PizzaStore
}

func (p *PepperoniPizzaStore) OrderPizza() Pizza {
	pizza := NewPepperoniPizza()
	pizza.Prepare()
	pizza.Cook()
	pizza.Cut()
	pizza.Box()
	return pizza
}

func mainFM() {
	stores := []PizzaStore{
		&CheezePizzaStore{},
		&PepperoniPizzaStore{},
	}
	for _, store := range stores {
		store.OrderPizza()
	}
}
