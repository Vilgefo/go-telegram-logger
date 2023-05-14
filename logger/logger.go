package logger

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(bot *tgbotapi.BotAPI, logChannel int64) *zap.Logger {
	l, err := zap.NewProduction()
	cfg := zap.NewProductionConfig()
	if err != nil {
		panic("couln't initialize logger")
	}
	logger := zap.New(&LoggerCore{
		Core:       l.Core(),
		bot:        bot,
		logChannel: logChannel,
		Encoder:    zapcore.NewConsoleEncoder(cfg.EncoderConfig),
	})
	return logger
}

type LoggerCore struct {
	Encoder zapcore.Encoder
	zapcore.Core
	bot        *tgbotapi.BotAPI
	logChannel int64
}

func (c *LoggerCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c *LoggerCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	//refactor and remove double log encode
	buf, err := c.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		panic(err)
	}
	text := "<pre>" + entry.Level.String() + ": " + entry.Message + "\n\n"

	text += "encoded message: \n" + string(buf.Bytes()) + " </pre>"
	msg := tgbotapi.NewMessage(c.logChannel, text)
	msg.ParseMode = tgbotapi.ModeHTML
	c.bot.Send(msg)

	return c.Core.Write(entry, fields)
}
