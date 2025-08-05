package main

import (
	"time"

	tb "gopkg.in/telebot.v4"
)

type Telebot struct {
	bot *tb.Bot
}

func NewTelebot(token string) (*Telebot, error) {
	pref := tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tb.NewBot(pref)

	if err != nil {
		return nil, err
	}

	_telebot := Telebot{bot}

	return &_telebot, nil
}

func (telebot *Telebot) Send(id int, message string) {
	telebot.bot.Send(tb.ChatID(id), message)
}
