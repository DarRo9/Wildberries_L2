/*Фабричный метод (FactoryMethod)
Паттерн Factory Method относится к порождающим паттернам уровня класса и сфокусирован только на отношениях между классами.

Паттерн Factory Method полезен, когда система должна оставаться легко расширяемой путем добавления объектов новых типов. 
Этот паттерн является основой для всех порождающих паттернов и может легко трансформироваться под нужды системы. 
По этому, если перед разработчиком стоят не четкие требования для продукта или не ясен способ организации взаимодействия между продуктами, 
то для начала можно воспользоваться паттерном Factory Method, пока полностью не сформируются все требования.

Паттерн Factory Method применяется для создания объектов с определенным интерфейсом, реализации которого предоставляются потомками. 
Другими словами, есть базовый абстрактный класс фабрики, который говорит, что каждая его наследующая фабрика должна реализовать 
такой-то метод для создания своих продуктов.

Реализация фабричного метода может быть разной, в большинстве случаем это зависит от языка реализации. 
Это может быть полиморфизм или параметризированный метод.

Пример: К нам приходят файлы трех расширений .txt, .png, .doc. 
В зависимости от расширения файла мы должны сохранять его в одном из каталогов /file/txt/, /file/png/ и /file/doc/. 
Значит, у нас будет файловая фабрика с параметризированным фабричным методом, принимающим путь к файлу, который нам нужно сохранить в одном из каталогов. 
Этот фабричный метод возвращает нам объект, используя который мы можем манипулировать с нашим файлом (сохранить, посмотреть тип и каталог для сохранения).
Заметьте, мы никак не указываем какой экземпляр объекта-продукта нам нужно получить, 
это делает фабричный метод путем определения расширения файла и на его основе выбора подходящего класса продукта. 
Тем самым, если наша система будет расширяться и доступных расширений файлов станет, например 25, 
то нам всего лишь нужно будет изменить фабричный метод и реализовать классы продуктов.

Требуется для реализации:

Базовый абстрактный класс Creator, описывающий интерфейс, который должна реализовать конкретная фабрика для производства продуктов. 
Этот базовый класс описывает фабричный метод.
Базовый класс Product, описывающий интерфейс продукта, который возвращает фабрика. 
Все продукты возвращаемые фабрикой должны придерживаться единого интерфейса.
Класс конкретной фабрики по производству продуктов ConcreteCreator. Этот класс должен реализовать фабричный метод;
Класс реального продукта ConcreteProductFirst;
Класс реального продукта ConcreteProductSecond;
Класс реального продукта ConcreteProductThird.
Factory Method отличается от Abstract Factory, тем, что Abstract Factory производит семейство объектов, 
эти объекты разные, обладают разными интерфейсами, 
но взаимодействуют между собой. В то время как Factory Method производит продукты придерживающиеся одного интерфейса и эти продукты не связаны между собой,
 не вступают во взаимодействие.

[!] В описании паттерна применяются общие понятия, такие как Класс, Объект, Абстрактный класс. Применимо к языку Go, это Пользовательский Тип, 
Значение этого Типа и Интерфейс. Также в языке Go за место общепринятого наследования используется агрегирование и встраивание.
*/
// Пакет factory_method является примером шаблона factory Method.
package factory_method

import (
	"log"
)

// action помогает клиентам узнать доступные действия.
type action string

const (
	First action = "First"
	Second action = "Second"
	Third action = "Third"
)

// Creator предоставляет заводской интерфейс.
type Creator interface {
	CreateProduct(action action) Product // Factory Method
}

/* Product поддерживает product интерфейс.
All products returned by factory must provide a single interface.*/
type Product interface {
	Use() string // Каждый продукт должен иметь возможность для использования
}

// ConcreteCreator реализует Creator interface.
type ConcreteCreator struct{}

// NewCreator это конструктор для ConcreteCreator.
func NewCreator() Creator {
	return &ConcreteCreator{}
}

// CreateProduct - это Factory Method.
func (p *ConcreteCreator) CreateProduct(action action) Product {
	var product Product

	switch action {
	case First:
		product = &ConcreteProductFirst{string(action)}
	case Second:
		product = &ConcreteProductSecond{string(action)}
	case C:
		product = &ConcreteProductThird{string(action)}
	default:
		log.Fatalln("Unknown Action")
	}

	return product
}

// ConcreteProductFirst реализует product "First".
type ConcreteProductFirst struct {
	action string
}

// Use возвращает знвчение поля структуры product action.
func (p *ConcreteProductFirst) Use() string {
	return p.action
}

// ConcreteProductSecond реализует product "Second".
type ConcreteProductSecond struct {
	action string
}

// Use возвращает знвчение поля структуры product action.
func (p *ConcreteProductSecond) Use() string {
	return p.action
}

// ConcreteProductThird реализует product "Third".
type ConcreteProductThird struct {
	action string
}

// Use возвращает знвчение поля структуры product action.
func (p *ConcreteProductThird) Use() string {
	return p.action
}