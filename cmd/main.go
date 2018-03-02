package cmd

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const FILE_NAME = "/usr/share/dict/words"

var rootCmd = &cobra.Command{
	Use:   "pwgen",
	Short: "pwgen generates passwords from " + FILE_NAME,
	Run:   ff,
}

func ff(cmd *cobra.Command, args []string) {
	f, err := os.Open(FILE_NAME)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines, err := getLines(f)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(lines))
	fmt.Printf("Count: %d\n", len(lines))
	fmt.Printf("Line at index %d: %s\n", index, lines[index])
}

func getLines(r io.Reader) ([]string, error) {
	var a = make([]string, 0)
	reader := bufio.NewReader(r)

	for {
		line, err := reader.ReadString('\n')

		switch {
		case err == io.EOF:
			return a, nil

		case err != nil:
			return nil, err
		}

		a = append(a, line)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//rootCmd.AddCommand(rootCmd)
}
