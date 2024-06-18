package worker

import (
	"telegram-bot/pkg/logger"
	"telegram-bot/pkg/smap"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type TelBot struct {
	l     logger.Logger
	bot   *tgbotapi.BotAPI
	users *smap.STMap
}

func NewTelBot(l logger.Logger, users *smap.STMap, bot *tgbotapi.BotAPI) *TelBot {
	return &TelBot{
		l:     l,
		bot:   bot,
		users: users,
	}
}

func (b *TelBot) InitTelBotConsumer() {
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upChan, err := b.bot.GetUpdatesChan(ucfg)
	if err != nil {
		b.l.WithFields(
			b.l.TraceWrap(err),
		).Panic(err)
	}
	var text string
	for update := range upChan {
		ID := update.Message.Chat.ID

		if ok := b.users.Get(ID); ok {
			b.l.Infof("User already exist ID: %d", ID)
			text = "Я тебя уже знаю"
		} else {
			b.users.Set(ID)
			b.l.Infof("Add new user ID: %d", ID)
			text = "Будем знакомы"
		}

		msg := tgbotapi.NewMessage(ID, text)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.l.WithFields(
				b.l.TraceWrap(err),
			).Error(err)
		}
	}
}

func (b *TelBot) InitTelBotProducer(ch <-chan string) {
	for msg := range ch {
		for _, chat := range b.users.GetAllKeys() {
			ms := tgbotapi.NewMessage(chat, msg)
			_, err := b.bot.Send(ms)
			if err != nil {
				b.l.WithFields(
					b.l.TraceWrap(err),
				).Error(err)
			}
			b.l.Infof("Success sended! msg: %s ID: %d", msg, chat)
		}
	}
}
