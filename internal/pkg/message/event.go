package message

import (
	"encoding/xml"
	"net/http"
	"wxrobot/internal/pkg/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type EventType string

const (
	SUBSCIBE_EVENT    EventType = "subscribe"
	UNSUBSCRIBE_EVENT EventType = "unsubscribe"
)

type EventCommon struct {
	MessageCommon
	EventType EventType `xml:"Event"`
}

var gEventHandler = map[EventType]func(c *gin.Context, event EventCommon){
	SUBSCIBE_EVENT:    handleSubscribeEvent,
	UNSUBSCRIBE_EVENT: handleUnSubscribeEvent,
}

func handleEventMessage(c *gin.Context, body []byte) {
	log.Info("handleEventMessage")
	event := EventCommon{}
	err := xml.Unmarshal(body, &event)
	if err != nil {
		log.Error("handleEventMessage err:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	handleCb, ok := gEventHandler[event.EventType]
	if ok {
		handleCb(c, event)
	} else {
		log.Warn("Unknow Event:", event.EventType)
	}

}

func handleSubscribeEvent(c *gin.Context, event EventCommon) {
	log.Info("handleSubscribeEvent")
	log.Info("userId:", event.FromUserName)
	if err := model.CreateWxUser(event.FromUserName); err != nil {
		log.Warn("CreateUser failed, err=", err)
	}
	c.XML(200, NewTextMessage("白茶清欢无别事，我在等风也等你", event.FromUserName, event.ToUserName))
}

func handleUnSubscribeEvent(c *gin.Context, event EventCommon) {
	log.Info("handleUnSubscribeEvent")
	log.Info("userId:", event.FromUserName)
}
