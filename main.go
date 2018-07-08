package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"log"
    "net/http"
    
	"github.com/line/line-bot-sdk-go/linebot"
	"gocv.io/x/gocv"
    "github.com/gin-gonic/gin"
)

const (
	SECRET_KEY = os.Getenv("SECRET_KEY")
	ACCEESS_TOKEN = os.Getenv("ACCEESS_TOKEN")
)


func main() {
	fmt.Println(SECRET_KEY, ACCEESS_TOKEN)

    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    router := gin.New()
    router.Use(gin.Logger())


    router.POST("/hook", func(c *gin.Context) {
        client := &http.Client{Timeout: time.Duration(15 * time.Second)}
		bot, err := linebot.New(SECRET_KEY, ACCEESS_TOKEN,, linebot.WithHTTPClient(client))
        if err != nil {
            fmt.Println(err)
            return
        }
        received, err := bot.ParseRequest(c.Request)

        for _, event := range received {
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type) {
                case *linebot.TextMessage:
                    source := event.Source
                    if source.Type == linebot.EventSourceTypeRoom {
                        if resMessage := getResMessage(message.Text); resMessage != "" {
                            postMessage := linebot.NewTextMessage(resMessage)
                            if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
                                log.Print(err)
                            }
                        }
                    }
                }
            }
        }
    })

    router.Run(":" + port)
}


func getResMessage(reqMessage string) (message string) {
    resMessages := [3]string{"わかるわかる","それで？それで？","からの〜？"}

    rand.Seed(time.Now().UnixNano())
    if rand.Intn(5) == 0 {
        if math := rand.Intn(4); math != 3 {
            message = resMessages[math];
        } else {
            message = reqMessage + "じゃねーよw"
        }
    }
    return
}