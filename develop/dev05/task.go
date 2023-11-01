/*
Утилита grep


Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).


Реализовать поддержку утилитой следующих ключей:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", напечатать номер строки
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)

// keys - структура флагов программы
type keys struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
	rExp       string
}

// readData - функция чтениия файла в срез строк
func readData(nameOfFile string) []string {
	var data []string

	file, err := os.Open(nameOfFile)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Не могу закрыть файл")
		}
	}(file)

	if err != nil {
		log.Fatal("Не могу открыть файл ", nameOfFile)
	}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		data = append(data, sc.Text())
	}

	return data
}

// showResult функция печатает результат вывода grep на экран
func showResult(res interface{}) {
	switch result := res.(type) {
	case []string:
		for _, str := range result {
			fmt.Println(str)
		}
	case []int:
		for _, num := range result {
			fmt.Println(num)
		}
	case int:
		fmt.Println(result)
	}
}

// copyToMapByInterval - функция копирует из слайса в мапу данные из заданного интервала в слайсе
func copyToMapByInterval(m map[int]string, data []string, first int, last int) {
	for k := first; k < last; k++ {
		m[k] = data[k]
	}
}

// sortMapKeys - функция возвращает слайс чисел (отсортированных ключей мапы)
func sortMapKeys(m map[int]string) []int {
	keys := make([]int, 0, 1)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// GenerateTheResult - функция, на основе промежуточной мапы и ее отсортированных ключей, формирует результат
func GenerateTheResult(m map[int]string, keys []int) []string {
	res := make([]string, 0, len(m))
	for _, key := range keys {
		res = append(res, m[key])
	}

	return res
}

// keyB для реализации ключа -B
func keyB(rExp *regexp.Regexp, f keys, data []string) map[int]string {
	mapBuf := make(map[int]string)
	for ind, str := range data {
		if rExp.MatchString(str) && !f.invert || !rExp.MatchString(str) && f.invert {
			mapBuf[ind] = str
			if ind-f.before >= 0 {
				copyToMapByInterval(mapBuf, data, ind-f.before, ind)
			} else {
				copyToMapByInterval(mapBuf, data, 0, ind)
			}
		}
	}
	return mapBuf
}

// keyA для реализации ключа -A
func keyA(rExp *regexp.Regexp, f keys, data []string) map[int]string {
	mapBuf := make(map[int]string)
	for ind, str := range data {
		if rExp.MatchString(str) && !f.invert || !rExp.MatchString(str) && f.invert {
			mapBuf[ind] = str
			if ind+f.after < len(data) {
				copyToMapByInterval(mapBuf, data, ind+1, ind+f.after+1)
			} else {
				copyToMapByInterval(mapBuf, data, ind+1, len(data))
			}
		}
	}
	return mapBuf
}

// keyC для реализации ключа -C
func keyC(rExp *regexp.Regexp, f keys, data []string) map[int]string {
	mapBuf := make(map[int]string)
	for ind, str := range data {
		if rExp.MatchString(str) && !f.invert || !rExp.MatchString(str) && f.invert {
			mapBuf[ind] = str

			if ind-f.context >= 0 {
				copyToMapByInterval(mapBuf, data, ind-f.context, ind)
			} else {
				copyToMapByInterval(mapBuf, data, 0, ind)
			}

			if ind+f.context < len(data) {
				copyToMapByInterval(mapBuf, data, ind+1, ind+f.context+1)
			} else {
				copyToMapByInterval(mapBuf, data, ind+1, len(data))
			}
		}
	}

	return mapBuf
}

// Реализует простейший grep
func grepSimple(rExp *regexp.Regexp, f keys, data []string) map[int]string {
	mapBuf := make(map[int]string)
	for ind, str := range data {
		if rExp.MatchString(str) && !f.invert || !rExp.MatchString(str) && f.invert {
			mapBuf[ind] = str
		}
	}

	return mapBuf
}

// Решает поставленную задачу (выполняетя grep с различными ключами)
func grep(data []string, f keys) (interface{}, error) {
	var (
		prefix  string
		postfix string
	)

	if f.ignoreCase {
		prefix = "(?i)"
	}

	if f.fixed {
		prefix += "^"
		postfix += "$"
	}

	rExp, err := regexp.Compile(prefix + f.rExp + postfix)
	if err != nil {
		log.Fatal("Bad regexp")
	}

	if f.count {
		cnt := 0
		if f.before > 0 {
			beforeRes := keyB(rExp, f, data)
			cnt = len(beforeRes)
		} else if f.after > 0 {
			afterRes := keyA(rExp, f, data)
			cnt = len(afterRes)
		} else if f.context > 0 {
			contextRes := keyC(rExp, f, data)
			cnt = len(contextRes)
		} else {
			simpleRes := grepSimple(rExp, f, data)
			cnt = len(simpleRes)
		}

		return cnt, nil
	} else if f.lineNum {
		var lineNumList []int

		if f.before > 0 {
			beforeRes := keyB(rExp, f, data)
			lineNumList = sortMapKeys(beforeRes)
		} else if f.after > 0 {
			afterRes := keyA(rExp, f, data)
			lineNumList = sortMapKeys(afterRes)
		} else if f.context > 0 {
			contextRes := keyC(rExp, f, data)
			lineNumList = sortMapKeys(contextRes)
		} else {
			simpleRes := grepSimple(rExp, f, data)
			lineNumList = sortMapKeys(simpleRes)
		}

		return lineNumList, nil
	} else if f.before > 0 {
		mapBuf := keyB(rExp, f, data)
		keys := sortMapKeys(mapBuf)
		res := GenerateTheResult(mapBuf, keys)

		return res, nil
	} else if f.after > 0 {
		mapBuf := keyA(rExp, f, data)
		keys := sortMapKeys(mapBuf)
		res := GenerateTheResult(mapBuf, keys)

		return res, nil
	} else if f.context > 0 {
		mapBuf := keyC(rExp, f, data)
		keys := sortMapKeys(mapBuf)
		res := GenerateTheResult(mapBuf, keys)

		return res, nil
	} else {
		mapBuf := grepSimple(rExp, f, data)
		keys := sortMapKeys(mapBuf)
		res := GenerateTheResult(mapBuf, keys)

		return res, nil
	}
}

func main() {
	flgs := keys{rExp: "\\.txt"}
	// Cчитывание файла
	data := readData("for_test.txt")
	// Выполнение операции
	res, _ := grep(data, flgs)
	// Вывод результата
	showResult(res)
}
