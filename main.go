package main

import (
	"time"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
	"fmt"
	url "net/url"
	"os"
)

func main() {
	bot, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("LMGTFY_ACCESS_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/start", func(m *tb.Message) {
		bot.Send(m.Sender, fmt.Sprintf("Hello, %s!\nThis bot handle inline queries from your messages in chat (example @lmgt4ybot <message>)", m.Sender.FirstName))
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		bot.Send(m.Sender, `¯\_(ツ)_/¯`)
	})

	bot.Handle(tb.OnQuery, func(q *tb.Query) {
		if len(q.Text) == 0 {
			return
		}

		u, err := url.Parse("http://lmgtfy.com/")
		if err != nil{
			log.Fatal(err)
		}
		qu := u.Query()
		qu.Set("q", q.Text)
		u.RawQuery = qu.Encode()

		err = bot.Answer(q, &tb.QueryResponse{
			Results:   []tb.Result {&tb.ArticleResult{
				Title:"Let Me Google That For You",
				Text: u.String(),
				Description:"Thank you for this question, click here to send link! =)",
				ThumbURL:"https://crunchbase-production-res.cloudinary.com/image/upload/c_lpad,h_256,w_256,f_auto,q_auto:eco/v1501766036/zleq5k7rz8m8tzbajqxe.png"}},
			CacheTime: 60, // a minute
		})

		if err != nil {
			log.Fatal(err)
		}
	})

	bot.Start()
}
