package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Token      string
	RSSFeeds   []string
	Recipients []int
	LogLevel   string
	RateLimit  time.Duration
	MaxRetries int
	MaxItems   int
	MaxWorkers int
}

// DefaultConfig returns default configuration values
func DefaultConfig() *Config {
	return &Config{
		LogLevel:   "INFO",
		RateLimit:  time.Second,
		MaxRetries: 3,
		MaxItems:   5,
		MaxWorkers: 3,
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	// Load token
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}
	config.Token = token

	// Load RSS feeds
	rssFeeds, err := GetRssFeedsLinks()
	if err != nil {
		return nil, fmt.Errorf("failed to get RSS feeds: %w", err)
	}
	config.RSSFeeds = rssFeeds

	// Load recipients
	recipients, err := GetRecipients()
	if err != nil {
		return nil, fmt.Errorf("failed to get recipients: %w", err)
	}
	config.Recipients = recipients

	// Load log level
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	// Load rate limit
	if rateLimitStr := os.Getenv("RATE_LIMIT_MS"); rateLimitStr != "" {
		if rateLimit, err := time.ParseDuration(rateLimitStr + "ms"); err == nil {
			config.RateLimit = rateLimit
		}
	}

	// Load max retries
	if maxRetriesStr := os.Getenv("MAX_RETRIES"); maxRetriesStr != "" {
		if maxRetries, err := parseInt(maxRetriesStr); err == nil && maxRetries > 0 {
			config.MaxRetries = maxRetries
		}
	}

	// Load max items per feed
	if maxItemsStr := os.Getenv("MAX_ITEMS_PER_FEED"); maxItemsStr != "" {
		if maxItems, err := parseInt(maxItemsStr); err == nil && maxItems > 0 {
			config.MaxItems = maxItems
		}
	}

	// Load max workers
	if maxWorkersStr := os.Getenv("MAX_WORKERS"); maxWorkersStr != "" {
		if maxWorkers, err := parseInt(maxWorkersStr); err == nil && maxWorkers > 0 {
			config.MaxWorkers = maxWorkers
		}
	}

	return config, nil
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	if len(c.RSSFeeds) == 0 {
		return fmt.Errorf("at least one RSS feed is required")
	}
	if len(c.Recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	if c.MaxRetries <= 0 {
		return fmt.Errorf("max retries must be positive")
	}
	if c.MaxItems <= 0 {
		return fmt.Errorf("max items must be positive")
	}
	if c.MaxWorkers <= 0 {
		return fmt.Errorf("max workers must be positive")
	}
	return nil
}

// String returns a string representation of the configuration (without sensitive data)
func (c *Config) String() string {
	return fmt.Sprintf("Config{RSSFeeds: %d, Recipients: %d, LogLevel: %s, RateLimit: %v, MaxRetries: %d, MaxItems: %d, MaxWorkers: %d}",
		len(c.RSSFeeds), len(c.Recipients), c.LogLevel, c.RateLimit, c.MaxRetries, c.MaxItems, c.MaxWorkers)
}
