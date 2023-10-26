/*Посетитель (Visitor)
Паттерн Visitor относится к поведенческим паттернам уровня объекта.

Паттерн Visitor позволяет обойти набор элементов (объектов) с разнородными интерфейсами, 
а также позволяет добавить новый метод в класс объекта, при этом, не изменяя сам класс этого объекта.

Требуется для реализации:

Абстрактный класс Visitor, описывающий интерфейс визитера;
Класс ConcreteVisitor, реализующий конкретного визитера. Реализует методы для обхода конкретного элемента;
Класс ObjectStructure, реализующий структуру(коллекцию), в которой хранятся элементы для обхода;
Абстрактный класс Element, реализующий интерфейс элементов структуры;
Класс ElementA, реализующий элемент структуры;
Класс ElementB, реализующий элемент структуры.
[!] В описании паттерна применяются общие понятия, такие как Класс, Объект, Абстрактный класс. 
Применимо к языку Go, это Пользовательский Тип, Значение этого Типа и Интерфейс. 
Также в языке Go за место общепринятого наследования используется агрегирование и встраивание.*/

// Package visitor - пример паттерна "Визитор".
package visitor

// Visitor обеспечивает интерфейс visitor.
type Visitor interface {
	VisitWork(p *Work) string
	VisitHome(p *Home) string
	VisitCafe(p *Cafe) string
}

// Place предоставляет интерфейс для места, которое посетитель должен посетить.
type Place interface {
	Accept(v Visitor) string
}

// People реализует интерфейс Visitor.
type People struct {
}

// VisitWork поддерживает визит в Work.
func (v *People) VisitWork(p *Work) string {
	return p.GoToWork()
}

// VisitHome реализует визит в Home.
func (v *People) VisitHome(p *Home) string {
	return p.GoToHome()
}

// VisitCafe реализует визит в Cafe.
func (v *People) VisitCafe(p *Cafe) string {
	return p.GoToCafe()
}

// City реализует коллекцию мест для визита.
type City struct {
	places []Place
}

// Add добавляет Place в коллекцию.
func (c *City) Add(p Place) {
	c.places = append(c.places, p)
}

// Accept реализует визит во все места из коллекции.
func (c *City) Accept(v Visitor) string {
	var result string
	for _, p := range c.places {
		result += p.Accept(v)
	}
	return result
}

// Work реализует интерфейс Place.
type Work struct {
}

// Accept реализация.
func (s *Work) Accept(v Visitor) string {
	return v.VisitWork(s)
}

// GoToWork реализация.
func (s *Work) GoToWork() string {
	return "Go to work"
}

// Home реализует интерфейс Place.
type Home struct {
}

// Accept реализация.
func (p *Home) Accept(v Visitor) string {
	return v.VisitHome(p)
}

// GoToHome реализация.
func (p *Home) GoToHome() string {
	return "Go to home"
}

// Cafe реализует интерфейс Place.
type Cafe struct {
}

// Accept реализация.
func (b *Cafe) Accept(v Visitor) string {
	return v.VisitCafe(b)
}

// Cafe реализация.
func (b *Cafe) GoToCafe() string {
	return "Go to cafe"
}