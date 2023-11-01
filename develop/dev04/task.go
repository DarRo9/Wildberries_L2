/*
Написать функцию поиска всех множеств анаграмм по словарю. 


Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.


Требования:
Входные данные для функции: ссылка на массив, каждый элемент которого - слово на русском языке в кодировке utf8
Выходные данные: ссылка на мапу множеств анаграмм
Ключ - первое встретившееся в словаре слово из множества. Значение - ссылка на массив, каждый элемент которого, 
слово из множества.
Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру. 
В результате каждое слово должно встречаться только один раз.
*/

package main

import (
	"fmt"
	"sort"
	"strings"
)

// anagramSearch производит поиск анаграмм
func anagramSearch(arr []string) map[string][]string {
	inf := make(map[string][]string)
	keys := make(map[string]string)

	for _, s1 := range arr {
		s1 = strings.ToLower(s1)
		s2 := []rune(s1)

		sort.Slice(s2, func(i, j int) bool {
			return s2[i] < s2[j]
		})

		if _, ok := keys[string(s2)]; !ok {
			inf[s1] = append(inf[s1], s1)
			keys[string(s2)] = s1
			continue
		}

		key := keys[string(s2)]
		inf[key] = append(inf[key], s1)
	}

	for i := range inf {
		if len(inf[i]) < 2 {
			delete(inf, i)
		} else {
			sort.Strings(inf[i])
		}
	}

	return inf
}

func main() {
	arr := []string{"пятак", "пятак", "пятка", "тяпка", "листок", "слиток", "столик"}

	data := anagramSearch(arr)

	for _, i := range data {
		fmt.Println(i)
	}
}
