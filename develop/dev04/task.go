package main

import (
	"slices"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func Anagrams(words []string) map[string][]string {
	words = removeDuplicates(words)

	temp := make(map[string][]string)

	for _, word := range words {
		key := getKey(word)
		temp[key] = append(temp[key], word)
	}

	result := make(map[string][]string)
	for _, anagrams := range temp {
		if len(anagrams) < 2 {
			continue
		}
		result[anagrams[0]] = anagrams
		slices.Sort(anagrams)
	}
	return result

}

func removeDuplicates(words []string) []string {
	set := make(map[string]struct{}, len(words))
	result := make([]string, 0, len(words))
	for _, word := range words {
		word = strings.ToLower(word)
		if _, ok := set[word]; ok || len(word) == 0 {
			continue
		}
		result = append(result, word)
		set[word] = struct{}{}
	}
	return result
}

func getKey(word string) string {
	r := []rune(word)
	slices.Sort(r)
	return string(r)
}
