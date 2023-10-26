/*Фасад (Facade)
Паттерн Facade относится к структурным паттернам уровня объекта.

Паттерн Facade предоставляет высокоуровневый унифицированный интерфейс в виде набора имен методов к набору взаимосвязанных классов или объектов некоторой подсистемы, что облегчает ее использование.

Разбиение сложной системы на подсистемы позволяет упростить процесс разработки, а также помогает максимально снизить зависимости одной подсистемы от другой. Однако использовать такие подсистемы становиться довольно сложно. Один из способов решения этой проблемы является паттерн Facade. Наша задача, сделать простой, единый интерфейс, через который можно было бы взаимодействовать с подсистемами.

В качестве примера можно привести интерфейс автомобиля. Современные автомобили имеют унифицированный интерфейс для водителя, под которым скрывается сложная подсистема. Благодаря применению навороченной электроники, делающей большую часть работы за водителя, тот может с лёгкостью управлять автомобилем, не задумываясь, как там все работает.

Требуется для реализации:

Класс Facade предоставляющий унифицированный доступ для классов подсистемы;
Класс подсистемы SubSystemA;
Класс подсистемы SubSystemB;
Класс подсистемы SubSystemC.
Заметьте, что фасад не является единственной точкой доступа к подсистеме, он не ограничивает возможности, которые могут понадобиться "продвинутым" пользователям, желающим работать с подсистемой напрямую.

[!] В описании паттерна применяются общие понятия, такие как Класс, Объект, Абстрактный класс. Применимо к языку Go, это Пользовательский Тип, Значение этого Типа и Интерфейс. Также в языке Go за место общепринятого наследования используется агрегирование и встраивание.*/

// Package facade - пример паттерна "Фасад".
package facade

import (
	"strings"
)

// NewDeveloper создаёт Developer.
func NewDeveloper() *Developer {
	return &Developer{
		codesyntax: &CodeSyntax{},
		algorithms:  &Algorithms{},
		projects: &Projects{},
	}
}

// Developer реализует Developer и фасад.
type Developer struct {
	codesyntax *CodeSyntax
	algorithms  *Algorithms
	projects *Projects
}

// Todo возвращает то, что  Developer должен делать.
func (d *Developer) Todo() string {
	result := []string{
		d.codesyntax.Know(),
		d.algorithms.Understand(),
		d.projects.Experience(),
	}
	return strings.Join(result, "\n")
}

// CodeSyntax реализует подсистему "CodeSyntax"
type CodeSyntax struct {
}

// Know реализация.
func (c *CodeSyntax) Know() string {
	return "Know the code syntax"
}

// Algorithms реализует подсистему "Algorithms"
type Algorithms struct {
}

// Understand реализация.
func (a *Algorithms) Understand() string {
	return "Understand algorithms"
}

// Projects реализует подсистему "Projects"
type Projects struct {
}

// Experience реализация.
func (p *Projects) Experience() string {
	return "Experience with projects"
}