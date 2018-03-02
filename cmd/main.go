package cmd

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"strings"

	"github.com/spf13/cobra"
)

const DEFAULT_WORDS_FILE_NAME = "/usr/share/dict/words"
const DEFAULT_MAX_LENGTH = 64

var rootCmd = &cobra.Command{
	Use:   "pwgen",
	Short: "pwgen generates passwords from the specified (or default) file",
	Run:   run,
}
var wordsFileName string
var maxLength int
var verbose bool

func run(cmd *cobra.Command, args []string) {
	f, err := os.Open(wordsFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines, err := getLines(f)

	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	maxWords := 4
	maxTries := 20
	availableLength := maxLength - 2 // maxWords is also number of trailing dashes, 2 is the digit suffix

	var words = make([]string, 0)

	for i := 0; i < maxTries; i++ {
		index := rand.Intn(len(lines))
		word := lines[index]

		if verbose {
			fmt.Printf("availableLength: %d, candidate word: %s\n", availableLength, word)
		}

		wordLength := len(word) + 1

		if wordLength > availableLength {
			// Candidate word rejected for being too long
			continue
		}

		words = append(words, word)
		availableLength -= wordLength

		if verbose {
			fmt.Printf("Remaining availableLength: %d\n", availableLength)
		}

		if len(words) >= maxWords || availableLength <= 4 {
			break
		}
	}

	var pwd = ""

	for i, word := range words {
		if i == 0 {
			word = strings.Title(word)
		} else {
			pwd += "-"
		}

		pwd += word
	}

	var suffix = 0

	for {
		suffix = rand.Intn(100)

		// We don't want a suffix of "00", so let's skip 0 if we get it from the
		// random number generator.
		if suffix > 0 {
			break
		}
	}

	pwd = fmt.Sprintf("%s-%02d", pwd, suffix)

	fmt.Println(pwd)
}

func getLines(r io.Reader) ([]string, error) {
	var lines = make([]string, 0)
	reader := bufio.NewReader(r)

	for {
		line, err := reader.ReadString('\n')

		switch {
		case err == io.EOF:
			return lines, nil

		case err != nil:
			return nil, err
		}

		lines = append(lines, strings.TrimSpace(line))
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&wordsFileName, "file-name", "f", DEFAULT_WORDS_FILE_NAME, "File containing words from which to choose")
	rootCmd.Flags().IntVarP(&maxLength, "max-length", "m", DEFAULT_MAX_LENGTH, "Maximum length of the password")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output (useful for debugging, mostly)")
}
