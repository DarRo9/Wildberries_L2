Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

package main
 
import (
    "fmt"
    "os"
)
 
func Foo() error {
    var err *os.PathError = nil
    return err
}
 
func main() {
    err := Foo()
    fmt.Println(err)
    fmt.Println(err == nil)
}

Ответ:
<nil>
false

Интерфейс равен nil, так как один из полей интерфейса, которым являются ошибки в Go, равен nil, а именно - (interface table). 
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
Вся информация по методам, типам и прочему скрывается в tab - интерфейсной таблицей (interface table).
Но, после присвоения этого интерфейса переменной err, поле data перестаёт быть nil и приобретает указатель на конкретный тип *os.PathError. Поэтому при сравнении с nil результат false. 
При этом из-за пустого значения interface table при выводе значения переменной выводится nil.
