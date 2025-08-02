package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	gofeed "github.com/mmcdole/gofeed"
	"gopkg.in/telebot.v4"
)

var (
	tokenVariableName = "TELEGRAM_BOT_TOKEN"
)

func getRssFeedsLinks() ([]string, error) {
	rssFeeds := GetEnvOrDefault("RSS_FEEDS", "")

	if len(rssFeeds) == 0 {
		return RSSFeeds, nil
	}

	var linksList []string
	err := json.Unmarshal([]byte(rssFeeds), &linksList)

	if err != nil {
		return nil, err
	}

	fmt.Printf("os.Getenv(\"RSS_FEEDS\"): %v\n", os.Getenv("RSS_FEEDS"))
	return linksList, nil
}

func parseFeed(rssLink string) (string, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(rssLink)

	if err != nil {
		log.Println(fmt.Errorf("can not request %s: %w", rssLink, err))

		return "", err
	}

	resultMessage := fmt.Sprintln(feed.Title)

	for _, item := range feed.Items {
		resultMessage += fmt.Sprintln(item.Title) + fmt.Sprintln(item.Link) + "\n"
	}

	fmt.Printf("resultMessage: %v\n", resultMessage)
	return resultMessage, nil
}

func initEnvAndCheckToken() error {
	if _, exist := os.LookupEnv(tokenVariableName); !exist {

		err := godotenv.Load(".env")
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
	logger := NewLogger()

	err := initEnvAndCheckToken()
	if err != nil {
		logger.Error("Error initialize environment variables", "error", err)
		os.Exit(1)
	}

	logger.Debug("Environment variables successfully initialized")

	pref := telebot.Settings{
		Token:  GetEnvOrDefault(tokenVariableName, ""),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	logger.Debug("Telebot preferences created")

	bot, err := telebot.NewBot(pref)

	if err != nil {
		logger.Error("error while creating bot", "error", err)
		os.Exit(1)
	}

	logger.Debug("Telebot successfully created")

	recipients, err := GetRecipients()

	if err != nil {
		logger.Error("Can not initialize messages recipients", "error", err)
		os.Exit(1)
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

	rssFeeds, err := getRssFeedsLinks()

	if err != nil {
		logger.Error("Error while parsing RSS Feeds", "error", err)
		os.Exit(1)
	}

	logger.Debug("RSS Feeds successfully parsed", "rssFeeds", rssFeeds)

	logger.Debug("Start fetching RSS Feeds")
	for _, rssFeed := range rssFeeds {
		logger.Debug("Fetching RSS Feed", "rssFeed", rssFeed)
		message, err := parseFeed(rssFeed)

		if err != nil {
			logger.Warn("Failed to fetch RSS Feed", "rssFeed", rssFeed)
		} else {
			logger.Debug("RSS Feed successfully fetched and handled", "rssFeed", rssFeed)
			logger.Debug("RSS Feed result message", "rssFeedMessage", message)

			logger.Debug("Start sending messages to recipients")
			for _, recipient := range recipients {
				bot.Send(telebot.ChatID(recipient), message)
			}
			logger.Debug("Finish sending messages to recipients")
		}

	}

	// bot.Start()
}
