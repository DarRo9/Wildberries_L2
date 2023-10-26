/*Стратегия (Strategy)
Паттерн Strategy относится к поведенческим паттернам уровня объекта.

Паттерн Strategy определяет набор алгоритмов схожих по роду деятельности, инкапсулирует их в отдельный класс и делает их подменяемыми. 
Паттерн Strategy позволяет подменять алгоритмы без участия клиентов, которые используют эти алгоритмы.

Требуется для реализации:

Класс Context, представляющий собой контекст выполнения той или иной стратегии;
Абстрактный класс Strategy, определяющий интерфейс различных стратегий;
Класс ConcreteStrategyA, реализует одну из стратегий представляющую собой алгоритмы, направленные на достижение определенной цели;
Класс ConcreteStrategyB, реализует одно из стратегий представляющую собой алгоритмы, направленные на достижение определенной цели.
[!] В описании паттерна применяются общие понятия, такие как Класс, Объект, Абстрактный класс. Применимо к языку Go, это Пользовательский Тип, 
Значение этого Типа и Интерфейс. Также в языке Go за место общепринятого наследования используется агрегирование и встраивание.*/

// Package strategy - пример паттерна "Стратегия".
package strategy

// GreetingStrategy - интерфейс для приветствий.
type GreetingStrategy interface {
	Greeting()
}

// OfficialGreeting реализует оффициальное приветствие.
type OfficialGreeting struct {
}

// Greeting() приветствует оффициально.
func (o *OfficialGreeting) Greeting() string {
	return "Здравствуйте!"
}

// FriendlyGreeting реализует дружеское приветствие.
type FriendlyGreeting struct {
}

// Greeting() приветствует дружески.
func (f *FriendlyGreeting) Greeting() string {
	return "Привет!"
}

// Context обеспечивает context для вежливого приветствия.
type Context struct {
	strategy GreetingStrategy
}

// Algorithm замещает стратегии.
func (c *Context) Algorithm(a GreetingStrategy) {
	c.strategy = a
}

// Greeting приветствует в зависимости от выбранной стратегии.
func (c *Context) Greeting() string {
	return c.strategy.Greeting()
}