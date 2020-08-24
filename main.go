package main

import (
	"encoding/json"
	"fmt"

	url "net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	GOOGLE_LOGO_URL = "https://crunchbase-production-res.cloudinary.com/image/upload/c_lpad,h_256,w_256,f_auto,q_auto:eco/v1501766036/zleq5k7rz8m8tzbajqxe.png"
	LMGTF_URL       = "http://lmgtfy.com/"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// Create link on lmgtfy with query
func getLink(query string) (string, error) {
	u, err := url.Parse(LMGTF_URL)
	if err != nil {
		return "", err
	}
	qu := u.Query()
	qu.Set("q", query)
	u.RawQuery = qu.Encode()
	return u.String(), nil
}

func main() {
	// Bot initialization
	bot, err := tb.NewBot(tb.Settings{
		Synchronous: true,
		Token:       os.Getenv("LMGTFY_ACCESS_TOKEN"),
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Panic("Exit with the error")
	}

	// Telegram commands
	HelloMsgTeml := "Hello, %s!\nThis bot is handle inline queries (example @lmgt4ybot <message>)"
	bot.Handle("/start", func(m *tb.Message) {
		log.WithFields(log.Fields{
			"user_id":  m.Sender.ID,
			"username": m.Sender.Username,
		}).Info("Received start command from the user")
		bot.Send(m.Sender, fmt.Sprintf(HelloMsgTeml, m.Sender.FirstName))
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.WithFields(log.Fields{
			"user_id":  m.Sender.ID,
			"username": m.Sender.Username,
			"message":  m.Text,
		}).Info("Received random message from the user")
		bot.Send(m.Sender, `¯\_(ツ)_/¯`)
	})

	bot.Handle(tb.OnQuery, func(q *tb.Query) {
		if len(q.Text) == 0 {
			log.WithFields(log.Fields{
				"user_id":  q.From.ID,
				"username": q.From.Username,
			}).Info("Query string is empty")
			return
		}

		link, err := getLink(q.Text)
		if err != nil {
			log.WithFields(log.Fields{
				"user_id":    q.From.ID,
				"username":   q.From.Username,
				"query_text": q.Text,
				"error":      err.Error(),
			}).Error("Can't generate a link")
		}

		err = bot.Answer(q, &tb.QueryResponse{
			Results: []tb.Result{&tb.ArticleResult{
				Title:       "Let Me Google That For You",
				Text:        link,
				Description: "Thank you for this question, click here to send a link! 😀",
				ThumbURL:    GOOGLE_LOGO_URL}},
			CacheTime: 60, // 1 min
		})

		if err != nil {
			log.WithFields(log.Fields{
				"user_id":    q.From.ID,
				"username":   q.From.Username,
				"query_text": q.Text,
				"error":      err.Error(),
			}).Error("Unable to send a message")
		}
	})

	// Setup lamda service
	log.Info("Registering gateway proxy")
	lambda.Start(func(req events.APIGatewayProxyRequest) (err error) {
		var u tb.Update
		if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
			bot.ProcessUpdate(u)
		}

		return
	})
}
