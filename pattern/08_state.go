/*
Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/State_pattern
*/

/*
Паттерн «состояние» позволяет объекту изменять свое поведение в зависимости от
его состояния. Вместо использования больших условных операторов, поведение делегируется
отдельным объектам, представляющим каждое возможное состояние.

+ читаемость: заменяем множество условных операторов на цепочку состояний
+ гибкость: можно легко добавлять новые состояния
- сложность: приводит к увеличению количества структур и интерфейсов
- избыточность: если состояний не много, или их смена происходит нечасто, проще использовать условные операторы

Можно использовать при обработке заказов. Состояния заказа: новый, обработан, отправлен, доставлен.
Или в авторизации. Состояния: авторизован, не авторизован, заблокирован.
*/

package pattern

import (
	"fmt"
	"time"
)

type TrafficLightState interface {
	ChangeState() TrafficLightState
	Show()
}

type TrafficLight struct {
	curState TrafficLightState
}

func NewTrafficLight() *TrafficLight {
	return &TrafficLight{curState: &RedState{}}
}

func (tl *TrafficLight) Start() {
	tl.curState.Show()
	tl.curState = tl.curState.ChangeState()
}

type RedState struct {
	tl *TrafficLight
}

func (r *RedState) ChangeState() TrafficLightState {
	return &GreenState{tl: r.tl}
}

func (r *RedState) Show() {
	fmt.Println("\033[1m\033[31mRed\033[0m")
	time.Sleep(3 * time.Second)
}

type GreenState struct {
	tl *TrafficLight
}

func (g *GreenState) ChangeState() TrafficLightState {
	return &YellowState{tl: g.tl}
}

func (g *GreenState) Show() {
	fmt.Println("\033[1m\033[32mGreen\033[0m")
	time.Sleep(3 * time.Second)
}

type YellowState struct {
	tl *TrafficLight
}

func (y *YellowState) ChangeState() TrafficLightState {
	return &RedState{tl: y.tl}
}

func (y *YellowState) Show() {
	fmt.Println("\033[1m\033[33mYellow\033[0m")
	time.Sleep(1 * time.Second)
}

func mainState() {
	tl := NewTrafficLight()
	for range 10 {
		tl.Start()
	}
}
