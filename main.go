package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var adjectiveMap = make(map[string]int)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	// app := echo.New()
	// app.Logger.Fatal(
	// 	app.Start(os.Getenv("PORT")),
	// )

	//var adjectiveMap map[string]int
	//adjectiveMap = make(map[string]int)

	loadData()
	scrape()
}

func loadData() {
	//open local text file
	file, err := os.Open("wordmap.txt")
	if err != nil {
		fmt.Println("Error opening wordmap.txt", err)
		return
	}
	defer file.Close()

	//scan file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//skip blank lines
		if line == "" {
			continue
		}
		//split the line 'word:score'
		split := strings.Split(line, ":")
		//trim white space
		trimmed := strings.TrimSpace(split[1])
		//convert score to integer
		num, err := strconv.Atoi(trimmed)
		if err != nil {
			fmt.Println(err)
			return
		}
		//map each word:score
		word := split[0]
		score := num
		adjectiveMap[word] = score
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

type Article struct {
	Words        int
	wordsChecked int
	Paragraphs   int
	Score        int
	URL          string
}

func scrape() {
	wordCount := 0
	paragraphCount := 0
	wordsChecked := 0
	score := 0

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})
	//CNN: ".vossi-paragraph"
	//FOX: ".article-body p"
	c.OnHTML(".paragraph", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		split := strings.Split(e.Text, " ")
		//we need _ here because golang is expecting us to care about the index, we dont here so we use _
		for _, currentWord := range split {
			value, exists := adjectiveMap[currentWord]
			if exists {
				wordsChecked++
				score += value
			}
		}
		wordCount++
		paragraphCount++
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	//CNN
	c.Visit("https://www.cnn.com/2024/11/14/politics/elon-musk-doge-trump/index.html")

	//FOX
	//c.Visit("https://www.foxnews.com/us/special-education-teacher-resigns-apologizes-viral-video-threatening-trump-voters-sparks-backlash")

	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	e.Request.Visit(e.Attr("href"))
	// })

	fmt.Println("Total words: ", wordCount)
	fmt.Println("Total paragraphs: ", paragraphCount)
	fmt.Println("Words measured: ", wordsChecked)
	fmt.Println("Score: ", score)
}
