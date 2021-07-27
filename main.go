package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"hypermark/hackerNews"
	"hypermark/utils"
)

// flags
var k string
var o bool
var s bool
var stdout bool
func init() {
	// k and s are mutually exclusive.
	flag.StringVar(&k, "k", "",
		"Save articles based on a keyword in the title.")
	flag.BoolVar(&o, "o", false,
		"Overwrite the target file instead of appending to the end.")
	flag.BoolVar(&s, "s", false, "Show all articles and exit.")
	flag.BoolVar(&stdout, "stdout", false, "Write output to stdout.")
}

func main() {
	flag.Parse()

	var outputPath *os.File
	var err error

	// outputPath is either a user-provided file or Stdout.
	outputPath, err = utils.ChooseOutputPath(flag.Args(), o, stdout)
	if err != nil {
		if err.Error() == utils.EARLY_EXIT {
			return
		}
		log.Fatal(err)
	}
	defer outputPath.Close()

	articles := hackerNews.ScrapeHN()
	if s {
		for i := 0; i < 30; i++ {
			data := articles[i].GetInfo()
			fmt.Printf("%d. %s\n%s\n%s\n\n", i+1, data[0], data[1], data[2])
		}
		return
	} else if k != "" {
		fmt.Printf("Searching for articles with '%s' in the title.\n", k)

		output := ""
		articlesFound := 0
		for i := 0; i < 30; i++ {
			if articles[i].TitleContains(k) {
				output = utils.AppendArticleTable(output, articles[i])
				articlesFound++
			}
		}
		fmt.Printf("%d articles found. Writing output to %s.\n",
			articlesFound,
			outputPath.Name())
		_, err := outputPath.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for i, article := range articles {
			fmt.Printf("%d %s\n", i+1, article.Title)
		}

		var userInput string
		fmt.Printf("\nArticles to save: (eg: 1 2 3, 1-3)\n")
		reader := bufio.NewReader(os.Stdin)
		userInput, err := reader.ReadString('\n')
		userInput = userInput[:len(userInput)-1] // remove trailing newline
		if err != nil {
			log.Fatal(err)
		}

		selections, err := utils.GetUserSelections(userInput)
		if err != nil {
			log.Fatal(err)
		}

		var output string
		for _, sel := range selections {
			output = utils.AppendArticleTable(output, articles[sel-1])
		}

		_, err = outputPath.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(
			"%d articles written to %s.\n",
			len(selections),
			outputPath.Name(),
		)
	}
}
