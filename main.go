package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	SECRET_KEY   = os.Getenv("SECRET_KEY")
	ACCESS_TOKEN = os.Getenv("ACCESS_TOKEN")
	PORT         = os.Getenv("PORT")
)

func main() {
	if PORT == "" {
		PORT = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/hook", func(c *gin.Context) {
		client := &http.Client{Timeout: time.Duration(15 * time.Second)}
		bot, err := linebot.New(SECRET_KEY, ACCESS_TOKEN, linebot.WithHTTPClient(client))
		if err != nil {
			fmt.Println(err)
			return
		}
		received, err := bot.ParseRequest(c.Request)

		for _, event := range received {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if resMessage := message.Text; resMessage != "" {
						postMessage := linebot.NewTextMessage(resMessage)
						if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
							log.Print(err)
						}
					}

				case *linebot.ImageMessage:
					fmt.Println("****************************")
					fmt.Println(message.ID)
					content, err := bot.GetMessageContent(message.ID).Do()
					if err != nil {
						fmt.Println("GetMessageContent Error")
					}
					defer content.Content.Close()

					// buf := bytes.NewBuffer(content.Content)
					img, err := ioutil.ReadAll(content.Content)
					spew.Dump(img, err)

					postMessage := linebot.NewTextMessage("hoge")
					if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
						log.Print(err)
					}

				}
			}
		}
	})

	router.Run(":" + PORT)
}
