/*Команда (Command)
Паттерн Command относится к поведенческим паттернам уровня объекта.

Паттерн Command позволяет представить запрос в виде объекта. Из этого следует, что команда - это объект. 
Такие запросы, например, можно ставить в очередь, отменять или возобновлять.

В этом паттерне мы оперируем следующими понятиями: Command - запрос в виде объекта на выполнение; 
Receiver - объект-получатель запроса, который будет обрабатывать нашу команду; Invoker - объект-инициатор запроса.

Паттерн Command отделяет объект, инициирующий операцию, от объекта, который знает, как ее выполнить. 
Единственное, что должен знать инициатор, это как отправить команду.

Требуется для реализации:

Базовый абстрактный класс Command описывающий интерфейс команды;
Класс ConcreteCommand, реализующий команду;
Класс Invoker, реализующий инициатора, записывающий команду и провоцирующий её выполнение;
Класс Receiver, реализующий получателя и имеющий набор действий, которые команда можем запрашивать;
Invoker умеет складывать команды в стопку и инициировать их выполнение по какому-то событию. 
Обратившись к Invoker можно отменить команду, пока та не выполнена.

ConcreteCommand содержит в себе запросы к Receiver, которые тот должен выполнять. 
В свою очередь Receiver содержит только набор действий (Actions), которые выполняются при обращении к ним из ConcreteCommand.

[!] В описании паттерна применяются общие понятия, такие как Класс, Объект, Абстрактный класс. Применимо к языку Go, это Пользовательский Тип, 
Значение этого Типа и Интерфейс. Также в языке Go за место общепринятого наследования используется агрегирование и встраивание.*/

// Package command - пример паттерна "Комманда".
package command

// Command обеспечивает интерфейс комманды.
type Command interface {
	Execute() string
}

// SwitchOnCommand реализует интерфейс Command.
typу SwitchOnCommand struct {
	receiver *Receiver
}

// Execute комманда.
func (s *SwitchOnCommand) Execute() string {
	return s.receiver.SwitchOn()
}

// SwitchOffCommand реализует интерфейс Command.
type SwitchOffCommand struct {
	receiver *Receiver
}

// Execute комманда.
func (s *SwitchOffCommand) Execute() string {
	return s.receiver.SwitchOff()
}

// Receiver реализация.
type Receiver struct {
}

// SwitchOn реализация.
func (r *Receiver) SwitchOn() string {
	return "Switch On"
}

// SwitchOff реализация.
func (r *Receiver) SwitchOff() string {
	return "Switch Off"
}

// Requester реализация.
type Requester struct {
	commands []Command
}

// StoreCommand добавляет комманду.
func (r *Requester) StoreCommand(command Command) {
	r.commands = append(r.commands, command)
}

// UnStoreCommand удаляет комманду.
func (r *Requester) UnStoreCommand() {
	if len(r.commands) != 0 {
		r.commands = r.commands[:len(r.commands)-1]
	}
}

// Execute выполняет все комманды.
func (r *Requester) Execute() string {
	var result string
	for _, command := range r.commands {
		result += command.Execute() + "\n"
	}
	return result
}