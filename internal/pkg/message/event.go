package message

import (
	"encoding/xml"
	"net/http"
	"wxrobot/internal/pkg/model"

	"wxrobot/internal/pkg/log"

	"github.com/gin-gonic/gin"
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

	userName := ""
	wxUser, _ := model.GetWxUser(event.FromUserName)
	if wxUser != nil {
		userName = wxUser.NickName
	} else {
		user := model.WxUser{
			UserId: event.FromUserName,
		}
		nextUser := getNextUser()
		if nextUser != nil {
			log.Info("Get Next User Success:", nextUser)
			user.NickName = nextUser.NickName
			user.UserType = nextUser.UserType
			user.Birthday = nextUser.Birthday
		}

		if err := model.CreateWxUser(&user); err != nil {
			log.ErrorWithRecord("CreateUser failed, err=", err)
		}
	}

	content := "白茶清欢无别事，我在等风也等你"
	if len(userName) > 0 {
		content += "\n"
		content += "终于等到你了 " + userName + "\n"
		content += "回复m即可查看主菜单"
	}
	c.XML(200, NewTextMessage(content, c))
}

func handleUnSubscribeEvent(c *gin.Context, event EventCommon) {
	log.Info("handleUnSubscribeEvent")
	log.Info("userId:", event.FromUserName)
}
