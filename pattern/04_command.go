/*
Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Command_pattern
*/

/*
	Паттерн «комманда» превращает действия (запросы или операции) в объекты,
	которые можно выполнить и отменить.

	+ O из SOLID: реализует принцип открытости/закрытости.
	+ разделение ответственности: логика выполнения и инициирования команд разделены.
	+ гибкость: легко добавлять новые команды без изменения существующего кода.
	+ функциональность: позволяет реализовать простую отмену или повтор операций.
	+ функциональность: позволяет реализовать отложенный запуск операций.
	- сложность: приводит к увеличению количества структур и интерфейсов.

	В GUI интерфейсах: кнопки могут быть связаны с командами, которые выполняются при нажатии.
	Отложенные задачи: команды можно использовать для планирования задач.
	История изменений: в любых программах где нужно хранить и проходиться по истории изменений,
	например в текстовых редакторах.
*/

package pattern

import "fmt"

type Command interface {
	Execute()
	Undo()
}

type AddTextCommand struct {
	text string
	doc  *Document
}

func (c *AddTextCommand) Execute() {
	c.doc.Append(c.text)
}

func (c *AddTextCommand) Undo() {
	c.doc.Remove(len(c.text))
}

type RemoveTextCommand struct {
	text string
	doc  *Document
}

func (c *RemoveTextCommand) Execute() {
	c.doc.Remove(len(c.text))
}

func (c *RemoveTextCommand) Undo() {
	c.doc.Append(c.text)
}

type Document struct {
	text string
}

func (d *Document) Append(s string) {
	d.text += s
}

func (d *Document) Remove(n int) {
	if len(d.text) >= n {
		d.text = d.text[:len(d.text)-n]
	}
}

func (d *Document) GetText() string {
	return d.text
}

func mainCommand() {
	doc := &Document{}
	history := []Command{}

	addHello := &AddTextCommand{text: "Hello", doc: doc}
	addHello.Execute()
	history = append(history, addHello)
	fmt.Println("Text:", doc.GetText())

	removeO := &RemoveTextCommand{text: "o", doc: doc}
	removeO.Execute()
	history = append(history, removeO)
	fmt.Println("Text:", doc.GetText())

	lastCommand := history[len(history)-1]
	lastCommand.Undo()
	history = history[:len(history)-1]
	fmt.Println("Text after undo:", doc.GetText())

	addWorld := &AddTextCommand{text: " world!", doc: doc}
	addWorld.Execute()
	history = append(history, addWorld)
	fmt.Println("Text:", doc.GetText())

	lastCommand = history[len(history)-1]
	lastCommand.Execute()
	fmt.Println("Text after redo:", doc.GetText())

}
