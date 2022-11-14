package parts

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

func ApplyAllModificators(words []string) []string {
	arrayOfWords := ReplaceHexAndBin(words)

	var result []string
	var resultLen, count int
	var match string

	regex_up, _ := regexp.Compile("^\\(up, \\d+\\)$")
	regex_cap, _ := regexp.Compile("^\\(cap, \\d+\\)$")
	regex_low, _ := regexp.Compile("^\\(low, \\d+\\)$")

	for i := range arrayOfWords {
		match_up := regex_up.FindString(arrayOfWords[i])
		match_cap := regex_cap.FindString(arrayOfWords[i])
		match_low := regex_low.FindString(arrayOfWords[i])

		if len(match_up)+len(match_cap)+len(match_low) != 0 {
			if len(match_cap) != 0 {
				match = match_cap
			} else if len(match_up) != 0 {
				match = match_up
			} else if len(match_low) != 0 {
				match = match_low
			}

			a := strings.Split(match[0:len(match)-1], ", ")
			num, _ := strconv.Atoi(a[1])

			for k := 0; k < num; k++ {
				if k > resultLen-1 {
					log.Fatal("Invalid input text: Number of modificator is out of range")
				}

				wordToChange := result[resultLen-1-k]

				if len(match_cap) != 0 {

					result[i-1-count-k] = strings.ToLower(wordToChange)
					result[i-1-count-k] = strings.Title(result[i-1-count-k][0:1]) + result[i-1-count-k][1:]
				} else if len(match_up) != 0 {
					result[i-1-count-k] = strings.ToUpper(wordToChange)
				} else if len(match_low) != 0 {
					result[i-1-count-k] = strings.ToLower(wordToChange)
				}
			}
			count += 1
		} else if arrayOfWords[i] == "(up)" {
			if resultLen-1 < 0 {
				count += 1
				continue
			}
			wordToChange := result[resultLen-1]
			result[i-1-count] = strings.ToUpper(wordToChange)
			count += 1
		} else if arrayOfWords[i] == "(low)" {
			if resultLen-1 < 0 {
				count += 1
				continue
			}
			wordToChange := result[resultLen-1]
			result[i-1-count] = strings.ToLower(wordToChange)
			count += 1
		} else if arrayOfWords[i] == "(cap)" {
			if resultLen-1 < 0 {
				count += 1
				continue
			}
			wordToChange := result[resultLen-1]

			result[i-1-count] = strings.ToLower(wordToChange)
			result[i-1-count] = strings.Title(result[i-1-count][0:1]) + result[i-1-count][1:]
			count += 1
		} else {
			result = append(result, arrayOfWords[i])
			resultLen += 1
		}
	}

	return result
}

func ReplaceHexAndBin(words []string) []string {
	var result []string
	var wordToChange string
	var count, base int
	change := false

	for i := range words {
		if words[i] == "(hex)" {
			change = true
			base = 16

		} else if words[i] == "(bin)" {
			change = true
			base = 2
		}
		if change {
			if i == 0 {
				log.Fatal("Invalid input: using modification as first element")
			}
			wordToChange = words[i-1]

			num, err := strconv.ParseInt(wordToChange, base, 64)
			if err != nil {
				continue
			}
			result[i-1-count] = strconv.FormatInt(int64(num), 10)
			count += 1
			change = false
		} else {
			result = append(result, words[i])
		}
	}
	return result
}
