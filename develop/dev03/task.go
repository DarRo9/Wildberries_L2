/*
Утилита cut

Реализовать утилиту аналог консольной команды cut (man cut). Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
*/

/*
Утилита sort
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры): на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительно

Реализовать поддержку утилитой следующих ключей:

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учетом суффиксов
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type flags struct {
	filename     string
	sortColumn   int
	sortByNum    bool
	reversedSort bool
	uniqueValues bool
}

// reverse обращает срез в обратном порядке
func reverse(data []string) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// sliceOfKeys возвращает срез ключей отображения
func sliceOfKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

/* сaseOrder исправляет порядок слов в отсортированном лексикографически срезе так, чтобы
слова с одинковой первой буквой, но с разным регистром находились согласно логике выполнение sort в linux */
func сaseOrder(data []string) {
	for i := 0; i < len(data)-1; i++ {
		str1 := []rune(data[i])
		str2 := []rune(data[i+1])

		if unicode.ToLower(str1[0]) == unicode.ToLower(str2[0]) &&
			unicode.IsUpper(str1[0]) &&
			unicode.IsLower(str2[0]) {
			buf := data[i]
			data[i] = data[i+1]
			data[i+1] = buf
		}
	}
}

// SortString осуществляет быструю сортировку строк
func SortString(data []string, start, end int, byNum bool) {
	if start < end {
		base := data[start]

		left := start
		right := end

		for left < right {

			if !byNum {
				for left < right && strings.ToLower(data[right]) >= strings.ToLower(base) {
					right--
				}
			} else {
				r, err := strconv.Atoi(data[right])
				b, err := strconv.Atoi(base)
				for left < right && r >= b {
					if err != nil {
						log.Fatal("Not number:", data[right])
					}

					right--

					r, err = strconv.Atoi(data[right])
					b, err = strconv.Atoi(base)
				}
			}

			if left < right {
				data[left] = data[right]
				left++
			}

			if !byNum {
				for left < right && strings.ToLower(data[left]) <= strings.ToLower(base) {
					left++
				}
			} else {
				l, err := strconv.Atoi(data[left])
				b, err := strconv.Atoi(base)
				for left < right && l <= b {
					if err != nil {
						log.Fatal("Not number:", data[left])
					}
					left++
					l, err = strconv.Atoi(data[left])
					b, err = strconv.Atoi(base)
				}
			}

			if left < right {
				data[right] = data[left]
				right--
			}
		}

		data[left] = base

		SortString(data, start, left-1, byNum)
		SortString(data, left+1, end, byNum)
	}
}

// uniqueS оставляет только уникальные строки в срезе
func uniqueS(data []string) []string {
	res := make([]string, 0, len(data))
	m := make(map[string]bool)
	for _, str := range data {
		if _, ok := m[str]; !ok {
			m[str] = true
			res = append(res, str)
		}
	}
	return res
}

// sortByColumn сортирует срез строк по колонкам в строке (колонки по умолчанию делятся пробелом)
func sortByColumn(data []string, fgs flags) {
	srcMap := make(map[string]string)
	for _, str := range data {
		columns := strings.Split(str, " ")
		if len(columns) > fgs.sortColumn {
			srcMap[columns[fgs.sortColumn]] = str
		} else {
			srcMap[columns[0]] = str
		}
	}
	keysToBeSorted := sliceOfKeys(srcMap)
	SortString(keysToBeSorted, 0, len(keysToBeSorted)-1, fgs.sortByNum)
	for ind, key := range keysToBeSorted {
		data[ind] = srcMap[key]
	}
}

// sort сортирует срез в соответствии с флагами
func sort(data []string, fgs flags) []string {
	res := make([]string, len(data))
	copy(res, data)

	if fgs.sortColumn >= 0 {
		sortByColumn(res, fgs)
		сaseOrder(res)
	} else {
		SortString(res, 0, len(res)-1, false)
		сaseOrder(res)
	}

	if fgs.uniqueValues {
		SortString(res, 0, len(res)-1, false)
		сaseOrder(res)
		res = uniqueS(res)
	}

	if fgs.reversedSort {
		reverse(res)
	}

	return res
}

// readF чтения файла в срез строк
func readF(filename string) []string {
	var rows []string

	file, err := os.Open(filename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Ошибка закрытия файла")
		}
	}(file)

	if err != nil {
		log.Fatal("Ошибка открытия файла", filename)
	}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		rows = append(rows, sc.Text())
	}

	return rows
}

// writeToF записи среза строк в файл
func writeToF(filename string, data []string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Ошибка создания файла")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal("Ошибка закрытия файла")
		}
	}(f)

	for i := 0; i < len(data); i++ {
		_, err := fmt.Fprintln(f, data[i])
		if err != nil {
			log.Fatal("Ошибка записи в файл")
		}
	}
}

func main() {
	// Установка флагов
	flags := flags{}
	// Cчитывание файла
	text := readF("for_sorting.txt")
	// Сортировка среза
	sortedText := sort(text, flags)
	// Запись отсортированного среза в файл
	writeToF("for_sorting.txt", sortedText)
	fmt.Println(sortedText)
}
