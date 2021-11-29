package message

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type MessageType string

const (
	TEXT_MESSAGE        MessageType = "text"
	IMAGE_MESSAGE       MessageType = "image"
	VOICE_MESSAGE       MessageType = "voice"
	VIDEO_MESSAGE       MessageType = "video"
	SHORT_VIDEO_MESSAGE MessageType = "shortvideo"
	LOCATION_MESSAGE    MessageType = "location"
	LINK_MESSAGE        MessageType = "link"
	EVENT_MESSAGE       MessageType = "event"
)

type MessageCommon struct {
	XMLName      xml.Name    `xml:"xml"` // 指定最外层的标签为config
	ToUserName   string      `xml:"ToUserName"`
	FromUserName string      `xml:"FromUserName"`
	CreateTime   int64       `xml:"CreateTime"`
	MsgType      MessageType `xml:"MsgType"`
	MsgId        int64       `xml:"MsgId"`
}

var gMsgHandler = map[MessageType]func(c *gin.Context, body []byte){
	TEXT_MESSAGE:  handleTextMessage,
	IMAGE_MESSAGE: handleImageMessage,
	EVENT_MESSAGE: handleEventMessage,
}

func proxyMessage(c *gin.Context, body []byte, mtype MessageType) {
	handleCb := gMsgHandler[mtype]
	handleCb(c, body)
}

func HandleMessage(c *gin.Context) {
	log.Debug("request:", c.Request)

	body, _ := ioutil.ReadAll(c.Request.Body)

	msg := MessageCommon{}
	err := xml.Unmarshal(body, &msg)
	if err != nil {
		log.Error("HandleMessage err:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	proxyMessage(c, body, msg.MsgType)
}
