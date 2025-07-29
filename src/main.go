package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	gofeed "github.com/mmcdole/gofeed"
	"gopkg.in/telebot.v4"
)

type RSSFeed struct {
	name, link string
}

var (
	RSSFeeds          = []RSSFeed{{link: "https://cprss.s3.amazonaws.com/golangweekly.com.xml", name: "Golang weekly"}, {link: "https://cprss.s3.amazonaws.com/nodeweekly.com.xml", name: "NodeJS weekly"}, {link: "https://cprss.s3.amazonaws.com/javascriptweekly.com.xml", name: "Javascript weekly"}, {link: "https://feedpress.me/cssweekly", name: "CSS weekly"}, {link: "https://cprss.s3.amazonaws.com/frontendfoc.us.xml", name: "Frontend focus"}, {link: "https://us5.campaign-archive.com/feed?u=ea228d7061e8bbfa8639666ad&id=104d6bcc2d", name: "Web tools"}, {link: "https://www.sitepoint.com/sitepoint.rss", name: "Sitepoint"}}
	tokenVariableName = "TELEGRAM_BOT_TOKEN"
)

func parseFeed(rssFeed RSSFeed) (string, error) {
	type MessageItem struct {
		title string
		link  string
	}

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(rssFeed.link)

	if err != nil {
		log.Println(fmt.Errorf("can not request %s: %w", rssFeed, err))

		return "", err
	}

	// items := make([]MessageItem, 0, feed.Len())
	resultMessage := fmt.Sprintln(feed.Title)

	for _, item := range feed.Items {
		resultMessage += fmt.Sprintln(item.Title) + fmt.Sprintln(item.Link) + "\n"

		// messageItem := MessageItem{title: fmt.Sprintln(item.Title), link: fmt.Sprintln(item.Link)}

		// items = append(items, messageItem)
	}

	fmt.Printf("resultMessage: %v\n", resultMessage)
	return resultMessage, nil
}

func initEnvAndCheckToken() error {
	if _, exist := os.LookupEnv(tokenVariableName); !exist {
		entries, err := os.ReadDir(".")
		if err != nil {
			return fmt.Errorf("error while reading directory: %w", err)
		}

		for _, entry := range entries {
			fmt.Println(entry.Name())
		}

		err = godotenv.Load("../.env")
		if err != nil {
			return fmt.Errorf(".env file load error: %w", err)
		}

		_, exist = os.LookupEnv(tokenVariableName)
		if !exist {
			return fmt.Errorf("token variable %s not found in .env file", tokenVariableName)
		}
	}

	return nil
}

func main() {
	err := initEnvAndCheckToken()
	if err != nil {
		log.Panic(err)
	}

	pref := telebot.Settings{
		Token:  os.Getenv(tokenVariableName),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)

	if err != nil {
		log.Panic(fmt.Errorf("error while creating bot: %w", err))
	}

	// bot.Handle("/start", func(c telebot.Context) error {
	// 	fp := gofeed.NewParser()

	// 	feed, err := fp.ParseURL("https://cprss.s3.amazonaws.com/javascriptweekly.com.xml")

	// 	if err != nil {
	// 		log.Println(fmt.Errorf("can not request javascriptweekly.com: %w", err))
	// 	}

	// 	item := feed.Items[0]

	// 	description, err := htmltomarkdown.ConvertString(item.Description)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Printf("c.Chat().ID: %v\n", c.Chat().ID)

	// 	return c.Send(description[0:4096])

	// 	return c.Send("bot started")
	// })

	for _, rssFeed := range RSSFeeds {
		message, err := parseFeed(rssFeed)

		if err != nil {
			log.Println(err)
		} else {
			bot.Send(telebot.ChatID(306650763), message)
		}

	}

	// bot.Start()

}
