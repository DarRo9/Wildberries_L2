/*Строитель (Builder)
Паттерн Builder относится к порождающим паттернам уровня объекта.

Паттерн Builder определяет процесс поэтапного построения сложного продукта. 
После того как будет построена последняя его часть, продукт можно использовать.

Возьмем одну фабрику, она производит сложный продукт, состоящий из 4 частей (крышка, бутылка, этикетка, напиток), 
которые должны быть применены в нужном порядке. Нельзя вначале взять крышку, бутылку, завинтить крышку, а потом пытаться налить туда напиток. 
Для реализации объекта, бутылки Кока-Колы, которая поставляется клиенту, нам нужен паттерн Builder.

Важно понимать, что сложный объект это не обязательно объект оперирующий несколькими другими объектами в смысле ООП. 
Например, нам нужно получить рассказ состоящий из названияа, основного текста и подписи автора. 
Наш рассказ, это сложный объект. Что бы был какой-то единый порядок составления рассказа, мы будем использовать паттерн Builder.

Требуется для реализации:

Класс Director, который будет распоряжаться строителем и отдавать ему команды в нужном порядке, а строитель будет их выполнять;
Базовый абстрактный класс Builder, который описывает интерфейс строителя, те команды, которые он обязан выполнять;
Класс ConcreteBuilder, который реализует интерфейс строителя и взаимодействует со сложным объектом;
Класс сложного объекта Product.
[!] В описании паттерна применяются общие понятия, такие как Класс, Объект, Абстрактный класс. 
Применимо к языку Go, это Пользовательский Тип, Значение этого Типа и Интерфейс. 
Также в языке Go заместо общепринятого наследования используется агрегирование и встраивание.
*/
// Package builder пример паттерна "Строитель".
package builder

// Builder обеспечивает builder интерфейс.
type Builder interface {
	MakeHeader(str string)
	MakeBody(str string)
	MakeFooter(str string)
}

// Director является менеджером
type Director struct {
	builder Builder
}

// Construct сообщает строителю, что делать и в каком порядке.
func (d *Director) Construct() {
	d.builder.MakeTitle("Title")
	d.builder.MakeMainText("MainText")
	d.builder.MakeSignature("Signature")
}

// ConcreteBuilder реализует интерфейс Builder.
type ConcreteBuilder struct {
	product *Product
}

// MakeTitle создаёт название рассказа
func (b *ConcreteBuilder) MakeTitle(str string) {
	b.product.Story += "Title of the story: " + str
}

// MakeMainText создаёт основной текст рассказа
func (b *ConcreteBuilder) MakeMainText(str string) {
	b.product.Story += "Main text: " + str
}

// MakeSignature создаёт подпись автора рассказа 
func (b *ConcreteBuilder) MakeSignature(str string) {
	b.product.Story += "Signature: " + str
}

// Реализация Product 
type Product struct {
	Story string
}

// Show возвращает результат
func (p *Product) Show() string {
	return p.Story
}