package main

import (
	"fmt"
	"go-reloaded/parts"
	"os"
	"strings"
)

func main() {
	args := os.Args

	parts.ArgsCheck(args)

	readFile := args[1]
	writeFile := args[2]

	fileContent := parts.ReadFile(readFile)

	arrayOfWordsAndPunct := parts.SplitByWordAndPunct(fileContent)

	arrayOfWords := parts.JoinWithPunct(arrayOfWordsAndPunct)
	fmt.Println(arrayOfWords)

	arrayOfWordsWithMod := parts.ApplyAllModificators(arrayOfWords)

	parts.GrammarCheck(arrayOfWordsWithMod)

	joinedString := strings.Join(arrayOfWordsWithMod, " ")

	parts.WriteFile(writeFile, joinedString)
}

// harold don ' t wilson (cap, 2) : ' I ' m a optimist ,but a optimist who carries a raincoat....  ' a.
