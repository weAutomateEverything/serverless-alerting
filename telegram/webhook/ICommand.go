package main

import (
	"gopkg.in/telegram-bot-api.v4"
)

type Command interface {
	CommandIdentifier() string
	CommandDescription() string
	RestrictToAuthorised() bool
	Show(chat uint32) bool
	Execute(update tgbotapi.Update)
}
