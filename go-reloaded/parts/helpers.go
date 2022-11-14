package parts

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func ArgsCheck(args []string) {
	if len(args) != 3 {
		log.Fatal("Not enough number of parameters")
	}
}

func ReadFile(filename string) string {
	file := strings.Split(filename, ".")
	if len(file) != 2 {
		log.Fatal("Could not read the file")
	}
	if file[1] != "txt" {
		log.Fatal("Could not read the file")
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Could not read the file due to this %s error \n", err)
	}
	return string(content)
}

func WriteFile(filename string, content string) {
	file := strings.Split(filename, ".")
	if len(file) != 2 {
		log.Fatal("Could not read the file")
	}
	if file[1] != "txt" {
		log.Fatal("Could not read the file")
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		log.Fatal(err2)
	}
}

/*Split words and punctuations by whitespace excpet modificators.Return string array */
func SplitByWordAndPunct(content string) []string {
	var arrayOfWords []string
	var word string
	new_word := false
	punct := ".,!?:;'"
	modificatorLength := 0
	skip := 0

	for index, symbol := range content {
		if skip != 0 {
			skip -= 1
			continue
		}
		if modificatorLength > 1 {
			new_word = true
			word += string(symbol)
			modificatorLength -= 1
			continue
		}

		if symbol == '(' {
			addWordToArray(&word, &new_word, &arrayOfWords)
			modificatorLength = GetNumberOfModificator(content[index:])
		}

		if unicode.IsSpace(symbol) {
			if new_word {
				addWordToArray(&word, &new_word, &arrayOfWords)
			}
		} else if strings.Contains(punct, string(symbol)) {
			addWordToArray(&word, &new_word, &arrayOfWords)
			if string(symbol) == "'" {
				skip = exceptionQuote(content, index, arrayOfWords)
				if skip != 0 {
					continue
				}

			}

			arrayOfWords = append(arrayOfWords, string(symbol))
		} else {
			word += string(symbol)
			new_word = true
		}

	}
	addWordToArray(&word, &new_word, &arrayOfWords)

	return arrayOfWords
}

func exceptionQuote(content string, index int, arrayOfWords []string) int {
	if len(arrayOfWords) == 0 {
		return 0
	}
	lastWord := arrayOfWords[len(arrayOfWords)-1]
	skip := 0
	if lastWord[len(lastWord)-1] == 'n' {
		if index+1 >= len(content) {
			return 0
		}
		next_symbol := content[index+1]
		skip = 1
		for unicode.IsSpace(rune(next_symbol)) {
			if index+skip >= len(content) {
				return 0
			}
			next_symbol = content[index+skip]
			skip += 1
		}
		if next_symbol == 't' {
			arrayOfWords[len(arrayOfWords)-1] += string(content[index]) + "t"
		} else {
			skip = 0
		}
		return skip
	}
	if unicode.IsLetter(rune(lastWord[len(lastWord)-1])) {
		if index+1 >= len(content) {
			return 0
		}
		next_symbol := content[index+1]
		skip = 1
		for unicode.IsSpace(rune(next_symbol)) {
			if index+skip >= len(content) {
				return 0
			}
			next_symbol = content[index+skip]
			skip += 1
		}
		if next_symbol == 's' {
			arrayOfWords[len(arrayOfWords)-1] += string(content[index]) + "s"
		} else {
			skip = 0
		}
	}

	if lastWord[len(lastWord)-1] == 'I' || lastWord[len(lastWord)-1] == 'i' {
		if index+1 >= len(content) {
			return 0
		}
		next_symbol := content[index+1]
		skip = 1
		for unicode.IsSpace(rune(next_symbol)) {
			if index+skip >= len(content) {
				return 0
			}
			next_symbol = content[index+skip]
			skip += 1
		}
		if next_symbol == 'm' {
			arrayOfWords[len(arrayOfWords)-1] += string(content[index]) + "m"
		} else {
			skip = 0
		}
		return skip
	}

	return skip
}

func addWordToArray(word *string, new_word *bool, arrayOfWords *[]string) {
	if len(*word) != 0 {

		*arrayOfWords = append(*arrayOfWords, *word)
		*word = ""
		*new_word = false
	}
}

func GetNumberOfModificator(st string) int {
	regex_up, _ := regexp.Compile("^\\(up, \\d+\\)")
	match_up := regex_up.FindStringIndex(st)

	regex_cap, _ := regexp.Compile("^\\(cap, \\d+\\)")
	match_cap := regex_cap.FindStringIndex(st)

	regex_low, _ := regexp.Compile("^\\(low, \\d+\\)")
	match_low := regex_low.FindStringIndex(st)

	if len(match_up) != 0 {
		if match_up[0] == 0 {
			return match_up[1]
		}
	}

	if len(match_cap) != 0 {
		if match_cap[0] == 0 {
			return match_cap[1]
		}
	}

	if len(match_low) != 0 {
		if match_low[0] == 0 {
			return match_low[1]
		}
	}

	return 0
}

func JoinWithPunct(arrayOfWords []string) []string {
	var array []string
	index := -1
	quoteOpen := 0
	punct := ".,!?:;'"
	skip := 0
	for i, v := range arrayOfWords {
		if skip != 0 {
			skip -= 1
			continue
		}
		if strings.Contains(punct, v) {
			if v == "'" {
				if quoteOpen == 0 {
					quoteOpen += 1
					if len(arrayOfWords) > i+1 {
						if arrayOfWords[i+1] == "'" {
							quoteOpen = 0
							array = append(array, "''")
							skip = 1
							continue
						}
					}

					if i == len(arrayOfWords)-1 {
						array = append(array, "'")
						continue
					}

				}

				if quoteOpen == 2 {

					array[index] += v
					quoteOpen = 0
				}
			} else {
				prev := 0
				for index-1 >= prev {
					if checkIfModificator(array[index-prev]) || GetNumberOfModificator(array[index-prev]) != 0 {
						prev += 1
					} else {
						break
					}
				}

				array[index-prev] += v

			}
		} else {
			if quoteOpen == 1 {
				v = "'" + v
				quoteOpen += 1
			}
			array = append(array, v)
			index += 1
		}
	}
	return array
}

func checkIfModificator(st string) bool {
	if st == "(up)" || st == "(low)" || st == "(cap)" {
		return true
	}
	return false
}

func GrammarCheck(arrayOfWords []string) {
	vowel := "aeiouh"
	for i := range arrayOfWords {
		if i == len(arrayOfWords)-1 {
			continue
		}
		if arrayOfWords[i] == "a" || arrayOfWords[i] == "A" {
			if strings.Contains(vowel, string(arrayOfWords[i+1][0])) {
				arrayOfWords[i] += "n"
			}
		} else if arrayOfWords[i] == "an" || arrayOfWords[i] == "An" {
			if !strings.Contains(vowel, string(arrayOfWords[i+1][0])) {
				arrayOfWords[i] = arrayOfWords[i][:1]
			}
		}
	}
}
