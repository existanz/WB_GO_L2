/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
	Фасад предоставляет простой интерфейс к сложной подсистеме. Эдакая инкапсуляция на уровне архитектуры.
	+ простой интерфейс
	+ сокрытие сложности подсистемы
	- возможна потеря гибкости из-за того, что фасад станет большим неповоротливым объектом с размытой ответственностью,
	так называемый "божественный объект"

	На практике фасад встречается практически в любой программе сложнее приветмира.
	К примеру, любая библиотека предоставляет внешнее API, для которого нужен фасад.
*/

package pattern

import (
	"fmt"
	"math/rand"
	"time"
)

type FacadeService interface {
	GetSomethingById(id int) string
	GetSomethingByName(name string) string
}

type service struct {
	db    *database
	cashe *cashe
}

func NewFacade(db *database, cashe *cashe) FacadeService {
	return &service{}
}

func (s *service) GetSomethingById(id int) string {
	some, err := s.cashe.getSomethingById(id)
	if err == nil {
		return some
	}

	some, err = s.db.getSomethingById(id)
	if err == nil {
		s.cashe.set(id, some)
		return some
	}

	return ""
}

func (s *service) GetSomethingByName(name string) string {
	id, err := s.db.getSomethingIdByName(name)
	if err == nil {
		return s.GetSomethingById(id)
	}
	return ""
}

type database struct {
	//
}

func (d *database) getSomethingById(id int) (string, error) {
	fmt.Println("try to get from database")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	if rand.Intn(100) < 50 {
		return "", fmt.Errorf("database error")
	}
	return fmt.Sprintf("Something with id: %d", id), nil
}

func (d *database) getSomethingIdByName(name string) (int, error) {
	fmt.Println("try to get from database")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	if rand.Intn(100) < 50 {
		return 0, fmt.Errorf("database error")
	}
	id := rand.Intn(100)
	fmt.Printf("Something with name: %s has id: %d", name, id)
	return id, nil
}

type cashe struct {
	//
}

func (c *cashe) getSomethingById(id int) (string, error) {
	fmt.Println("try to get from cashe")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	if rand.Intn(100) < 50 {
		return "", fmt.Errorf("cashe error")
	}
	return fmt.Sprintf("Something with id: %d", id), nil
}

func (c *cashe) set(id int, value string) error {
	fmt.Println("try to save to cashe")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	fmt.Printf("Value %s saved to cashe with id %d \n", value, id)
	return nil
}
