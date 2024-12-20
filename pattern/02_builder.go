/*
Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
	Строитель предоставляет возможность гибкой настройки создаваемого объекта.
	+ гибкость: можно создавать разные варианты объектов
	+ понятность: можно наглядно прочесть какой объект будет создан
	- усложнение: для сложных объектов может потребоваться много методов

	На практике "строитель" встречается довольно часто, например есть целый класс
	библиотек генерирующих SQL запросы, которые так и называются sqlbuilder'ы, в
	которых мы пошагово собираем запрос.
*/

package pattern

import (
	"io"
	"os"
)

type LogBuilder interface {
	SetLogLevel(logLevel LoggerLevel)
	SetWriter(w io.Writer)
	SetLogType(logType LogType)
	IncludeSourceLine()
}

type LoggerLevel int

const (
	DebugLevel LoggerLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

type LogType int

const (
	JSONLogType LogType = iota
	TextLogType
)

type logger struct {
	logLevel          LoggerLevel
	logType           LogType
	writer            io.Writer
	includeSourceLine bool
}

func (l *logger) SetLogLevel(logLevel LoggerLevel) {
	l.logLevel = logLevel
}

func (l *logger) SetWriter(w io.Writer) {
	l.writer = w
}

func (l *logger) SetLogType(logType LogType) {
	l.logType = logType
}

func (l *logger) IncludeSourceLine() {
	l.includeSourceLine = true
}

func NewLogger() LogBuilder {
	return &logger{}
}

func NewJSONLogger(ll LoggerLevel) LogBuilder {
	l := NewLogger()
	l.SetLogLevel(ll)
	l.SetWriter(io.Discard)
	l.SetLogType(JSONLogType)
	l.IncludeSourceLine()
	return l
}

func NewTextLogger(ll LoggerLevel) LogBuilder {
	l := NewLogger()
	l.SetLogLevel(ll)
	l.SetWriter(os.Stdout)
	l.SetLogType(TextLogType)
	return l
}
