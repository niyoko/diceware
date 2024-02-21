package cmd

import (
	"bytes"
	"cmp"
	"embed"
	_ "embed"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

//go:embed words/*.txt
var words embed.FS

var generateCmd = &cobra.Command{
	Use: "generate",
	Run: func(cmd *cobra.Command, args []string) {
		minEntropy := 77
		wordList := getWords()
		wordListN := len(wordList)
		perWordEntropy := math.Log2(float64(wordListN))
		wordToGenerate := int(math.Ceil(float64(minEntropy) / perWordEntropy))

		var pass strings.Builder
		rndSrc := &rndBit{}
		for i := 0; i <= wordToGenerate; i++ {
			num := rndSrc.randN(wordListN)
			if pass.Len() > 0 {
				pass.WriteRune(' ')
			}

			pass.WriteString(wordList[num])
		}

		fmt.Printf("Wordlist: %d words\n", wordListN)
		fmt.Printf("Entropy per word: %.2f bits\n", perWordEntropy)
		fmt.Printf("Total entropy: %.2f bits\n", perWordEntropy*float64(wordToGenerate))
		fmt.Println(pass.String())
	},
}

func getWords() (result []string) {
	files, _ := words.ReadDir("words")

	wordSeen := make(map[string]bool)
	for _, file := range files {
		words, _ := words.ReadFile("words/" + file.Name())
		wordsCount := getWordsCount(words)
		validWords := getValidWords(wordsCount)
		for _, word := range validWords {
			if !wordSeen[word] {
				wordSeen[word] = true
				result = append(result, word)
			}
		}
	}

	slices.SortFunc(result, func(a, b string) int { return cmp.Compare(a, b) })
	return
}

func getWordsCount(content []byte) (result map[string]uint64) {
	result = make(map[string]uint64)

	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		trimmedLine := bytes.TrimSpace(line)
		if ignoreLine(trimmedLine) {
			continue
		}

		f := bytes.Fields(trimmedLine)
		if len(f) != 2 {
			continue
		}

		word := string(f[0])
		freq, err := strconv.ParseUint(string(f[1]), 10, 64)
		if err != nil {
			continue
		}

		result[word] += freq
	}

	return
}

func ignoreLine(line []byte) bool {

	if bytes.HasPrefix(line, []byte("#")) {
		return true
	}

	runes := utf8.RuneCount(line)
	if runes == 0 || runes != len(line) {
		return true
	}

	return false
}

func getValidWords(words map[string]uint64) []string {
	result := make([]string, 0, len(words))
	freqs := make([]uint64, 0, len(words))

	for _, freq := range words {
		freqs = append(freqs, freq)
	}

	slices.SortFunc(freqs, func(a, b uint64) int { return cmp.Compare(a, b) })
	threshold := freqs[int(0.6*float64(len(freqs)))]
	for word, freq := range words {
		strLen := len(word)
		if freq >= threshold && strLen > 2 && strLen < 8 && word[0] >= 'a' && word[0] <= 'z' {
			result = append(result, word)
		}
	}

	return result
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
