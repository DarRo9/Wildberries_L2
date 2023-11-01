/*
Утилита cut

Реализовать утилиту аналог консольной команды cut (man cut). Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var fields string
var delimiter string
var separated bool

// регулярные выражения для processingFields
var re0 = regexp.MustCompile(`^\d+$`)
var re1 = regexp.MustCompile(`^\d+\-\d+$`)
var re2 = regexp.MustCompile(`^\-\d+$`)
var re3 = regexp.MustCompile(`^\d+\-$`)

var errNumbering = errors.New("поля нумеруются с 1")
var errDiapason = errors.New("неверный уменьшающийся диапазон")

// cutPrint печатает поля из диапазонов fs строки s.
func cutPrint(str string, fs [][]int, dr rune) {

	// поиск подстрок, разделенных символом dr
	var sliceOfStrings []string
	runes := []rune(str)
	st := 0
	congruence := false
	for i, r := range runes {
		if r == dr {
			congruence = true
			if i == st {
				sliceOfStrings = append(sliceOfStrings, "")
			} else {
				sliceOfStrings = append(sliceOfStrings, string(runes[st:i]))
			}
			st = i + 1
		} else if i == len(runes)-1 {
			sliceOfStrings = append(sliceOfStrings, string(runes[st:]))
		}
	}

	// строка не содержит разделителя
	if !congruence {
		if separated {
			return
		}
		fmt.Println(str)
		return
	}

	// если строка оканчивается символом-разделителем
	if lr, _ := utf8.DecodeLastRuneInString(str); lr == dr {
		sliceOfStrings = append(sliceOfStrings, "")
	}

	// выбор подстрок (полей), входящих в диапазоны fs
	var outSlice []string
	for i, f := range sliceOfStrings {
		if indexInSection(i+1, fs) {
			outSlice = append(outSlice, f)
		}
	}
	fmt.Println(strings.Join(outSlice, string(dr)))
}

// indexInSection возвращает true, если i принадлежит одному из отрезков в section
func indexInSection(i int, section [][]int) bool {
	res := section[0][0] <= i && i <= section[0][1]
	for j := 1; j < len(section); j++ {
		res = res || (section[j][0] <= i && i <= section[j][1])
	}
	return res
}

/* processingFields парсит строку, переданную через флаг -f
и возвращает слайс интервалов, отображаемых полей */ 
func processingFields(str string) (fs [][]int, err error) {
	strings1 := strings.FieldsFunc(str, func(r rune) bool { return r == ',' })

	// если не указано ни одного поля: ","
	if len(strings1) == 0 {
		return nil, errNumbering
	}

	// указано одно поле (диапазон) перед или после запятой: ",2", "2-4,"
	contaisComma, err := regexp.MatchString(",", str)
	if err != nil {
		return nil, err
	}
	if len(strings1) == 1 && contaisComma {
		return nil, errNumbering
	}
	for _, rString := range strings1 {

		// если подстрока соответствует регулярному выражению, диапазон заносится в результат
		switch {
		case re0.MatchString(rString):
			a, _ := strconv.Atoi(rString)
			if a < 1 {
				return nil, errNumbering
			}
			fs = append(fs, []int{a, a})
		case re1.MatchString(rString):
			ds := strings.FieldsFunc(rString, func(r rune) bool { return r == '-' })
			a, _ := strconv.Atoi(ds[0])
			if a < 1 {
				return nil, errNumbering
			}
			b, _ := strconv.Atoi(ds[1])
			if a > b {
				return nil, errDiapason
			}
			fs = append(fs, []int{a, b})
		case re2.MatchString(rString):
			b, _ := strconv.Atoi(rString[len("-"):])
			if b < 1 {
				return nil, errDiapason
			}
			fs = append(fs, []int{1, b})
		case re3.MatchString(rString):
			a, _ := strconv.Atoi(rString[:len(rString)-len("-")])
			if a < 1 {
				return nil, errNumbering
			}
			fs = append(fs, []int{a, math.MaxInt})
		default:
			return nil, errors.New("недопустимое значение поля")
		}
	}
	return fs, nil
}

func main() {
	if fields == "" {
		fmt.Fprintln(os.Stderr, "поля не указаны")
		return
	}

	fs, err := processingFields(fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Проверка, что в -d передан один символ-руна
	if dn := utf8.RuneCountInString(delimiter); dn != 1 {
		fmt.Fprintln(os.Stderr, "разделитель должен быть одним символом")
		return
	}
	dr, _ := utf8.DecodeRuneInString(delimiter)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		cutPrint(sc.Text(), fs, dr)
	}
}
