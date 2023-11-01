/*
Задача на распаковку
Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

Дополнительно 
Реализовать поддержку escape-последовательностей. 
Например:
qwe\4\5 => qwe45 (*)
qwe\45 => qwe44444 (*)
qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка, функция должна возвращать ошибку. Написать unit-тесты.
*/
package main

import (
	"fmt"
	"strconv"
	"strings"
	"errors"
)

func is_digit(r rune) bool {
	if _, err := strconv.Atoi(string(r)); err == nil {
		return true
	} else {
		return false
	}
}

// is_alpha проверяет является ли символ буквой
func is_alpha(r rune) bool {
	if !is_digit(r) && string(r) != `\` {
		return true
	}
	return false
}

// add записывает руну r в builder
func add(r rune, b *strings.Builder, cnt int) {
	for i := 0; i < cnt; i++ {
		b.WriteString(string(r))
	}
}

func string_unpacking(str string) (string, error) {
	arr := []rune(str)

	builder := strings.Builder{}

	ind := 0
	for ind < len(arr) {
		r := arr[ind]
		if is_alpha(r) {
			builder.WriteString(string(r))
		} else if is_digit(r) {
			cnt := 0

			j := ind
			prevInd := ind
			for j < len(arr) && is_digit(arr[j]) {
				buf, _ := strconv.Atoi(string(arr[j]))
				cnt = cnt*10 + buf
				j++
				ind++
			}
			if prevInd > 0 {
				add(arr[prevInd-1], &builder, cnt-1)
			} else {
				return "", errors.New("invalid string")
			}
			continue
		} else if string(r) == `\` {
			if ind < len(arr)-1 {
				builder.WriteString(string(arr[ind+1]))
				ind++
			} else {
				return "", errors.New("invalid string")
			}
		}
		ind++
	}

	return builder.String(), nil
}

func main() {
	example1, err1 := string_unpacking(`a\`)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println(example1)
	}
	
	example2, err2 := string_unpacking(`a4bc2d5e`)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println(example2)
	}

	example3, err3 := string_unpacking(`abcd`)
	if err3 != nil {
		fmt.Println(err3)
	} else {
		fmt.Println(example3)
	}

	example4, err4 := string_unpacking(`a11b`)
	if err4 != nil {
		fmt.Println(err4)
	} else {
		fmt.Println(example4)
	}

	example5, err5 := string_unpacking(`a12`)
	if err5 != nil {
		fmt.Println(err5)
	} else {
		fmt.Println(example5)
	}

	example6, err6 := string_unpacking(`45`)
	if err6 != nil {
		fmt.Println(err6)
	} else {
		fmt.Println(example6)
	}

	example7, err7 := string_unpacking(``)
	if err7 != nil {
		fmt.Println(err7)
	} else {
		fmt.Println(example7)
	}

	example8, err8 := string_unpacking(`qwe\4\5`)
	if err8 != nil {
		fmt.Println(err8)
	} else {
		fmt.Println(example8)
	}

	example9, err9 := string_unpacking(`qwe\45`)
	if err9 != nil {
		fmt.Println(err9)
	} else {
		fmt.Println(example9)
	}

	example10, err10 := string_unpacking(`qwe\\5`)
	if err10 != nil {
		fmt.Println(err10)
	} else {
		fmt.Println(example10)
	}

	example11, err11 := string_unpacking(`qwe\\`)
	if err11 != nil {
		fmt.Println(err11)
	} else {
		fmt.Println(example11)
	}
}
