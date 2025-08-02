package main

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	RSSFeeds = []string{"https://cprss.s3.amazonaws.com/golangweekly.com.xml", "https://cprss.s3.amazonaws.com/nodeweekly.com.xml", "https://cprss.s3.amazonaws.com/javascriptweekly.com.xml", "https://feedpress.me/cssweekly", "https://cprss.s3.amazonaws.com/frontendfoc.us.xml", "https://us5.campaign-archive.com/feed?u=ea228d7061e8bbfa8639666ad&id=104d6bcc2d", "https://www.sitepoint.com/sitepoint.rss"}
)

// getEnvOrDefault returns the value of the environment variable named by the key.
// If the environment variable is not set, it returns the default value.
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetRssFeedsLinks() ([]string, error) {
	rssFeeds := GetEnvOrDefault("RSS_FEEDS", "")

	if len(rssFeeds) == 0 {
		return RSSFeeds, nil
	}

	var linksList []string
	err := json.Unmarshal([]byte(rssFeeds), &linksList)

	if err != nil {
		return nil, err
	}

	return linksList, nil
}

func GetRecipients() ([]int32, error) {
	recipientsIdsJson := GetEnvOrDefault("RECIPIENTS", "")

	if len(recipientsIdsJson) == 0 {
		return nil, errors.New("Can not find RECIPIENTS env")
	}

	var recipients []int32
	err := json.Unmarshal([]byte(recipientsIdsJson), &recipients)

	if err != nil {
		return nil, err
	}

	return recipients, nil
}
