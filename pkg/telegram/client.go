package telegram

import (
	"time"

	config "github.com/jaredallard/mm/pkg/config"
	log "github.com/jaredallard/mm/pkg/logger"
	tb "gopkg.in/tucnak/telebot.v2"
)

var bot *tb.Bot

// Channel shim
type Channel struct {
	ID string
}

// Recipient returns the channel id provideed
func (c *Channel) Recipient() string {
	return c.ID
}

// Setup the telgram bot.
func Setup(cfg *config.ConfigurationFile) {
	var err error
	bot, err = tb.NewBot(tb.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal("Failed to start Telegram client.", err.Error())
	}

	log.Info("Telegram client created")
}

// SendToChannel a message
func SendToChannel(channelID string, message string) error {
	c := Channel{
		ID: channelID,
	}

	_, err := bot.Send(&c, message, &tb.ReplyMarkup{})
	if err != nil {
		return err
	}

	return nil
}

// Start the polling agent.
func Start() {
	bot.Start()
}
