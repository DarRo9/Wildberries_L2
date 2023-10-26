/*Цепочка ответственности (Chain Of Responsibility)
Паттерн Chain Of Responsibility относится к поведенческим паттернам уровня объекта.

Паттерн Chain Of Responsibility позволяет избежать привязки объекта-отправителя запроса к объекту-получателю запроса, 
при этом давая шанс обработать этот запрос нескольким объектам. Получатели связываются в цепочку, и запрос передается по цепочке, 
пока не будет обработан каким-то объектом.

По сути это цепочка обработчиков, которые по очереди получают запрос, а затем решают, обрабатывать его или нет. 
Если запрос не обработан, то он передается дальше по цепочке. Если же он обработан, то паттерн сам решает передавать его дальше или нет. 
Если запрос не обработан ни одним обработчиком, то он просто теряется.

Требуется для реализации:

Базовый абстрактный класс Handler, описывающий интерфейс обработчиков в цепочки;
Класс ConcreteHandlerA, реализующий конкретный обработчик A;
Класс ConcreteHandlerB, реализующий конкретный обработчик B;
Класс ConcreteHandlerC, реализующий конкретный обработчик C;
Обратите внимание, что вместо хранения ссылок на всех кандидатов-получателей запроса, каждый отправитель хранит единственную ссылку на начало цепочки, 
а каждый получатель имеет единственную ссылку на своего преемника - последующий элемент в цепочке.

[!] В описании паттерна применяются общие понятия, такие как Класс, Объект, Абстрактный класс. Применимо к языку Go, это Пользовательский Тип, 
Значение этого Типа и Интерфейс. Также в языке Go за место общепринятого наследования используется агрегирование и встраивание.*/

// Package chain_of_responsibility - пример паттерна "Цепочка ответственности"
package chain_of_responsibility

// Receiver поддерживает интерфейс получателя.
type Receiver interface {
	SendRequest(message string) string
}

// ConcreteReceiverFirst реализует первого получателя.
type ConcreteReceiverFirst struct {
	next Receiver
}

// Answer реализация.
func (h *ConcreteReceiverFirst) Answer(message string) (result string) {
	if message == "First" {
		result = "First receiver is working"
	} else if h.next != nil {
		result = h.next.Answer(message)
	}
	return
}

// ConcreteReceiverSecond реализует второго получателя.
type ConcreteReceiverSecond struct {
	next Receiver
}

// Answer реализация.
func (h *ConcreteReceiverSecond) Answer(message string) (result string) {
	if message == "Second" {
		result = "Second receiver is working"
	} else if h.next != nil {
		result = h.next.Answer(message)
	}
	return
}

// ConcreteReceiverThird реализует тертьего получателя.
type ConcreteReceiverThird struct {
	next Receiver
}

// Answer реализация.
func (h *ConcreteReceiverThird) Answer(message string) (result string) {
	if message == "Third" {
		result = "Third receiver is working"
	} else if h.next != nil {
		result = h.next.Answer(message)
	}
	return
}