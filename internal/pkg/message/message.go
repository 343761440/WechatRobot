package message

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type TextMessage struct {
	MessageCommon
	Content string
}

type ImageMessage struct {
	MessageCommon
	PicUrl  string
	MediaId string
}

type VoiceMessage struct {
	MessageCommon
	MediaId     string
	Format      string
	Recognition string
}

type VideoMessage struct {
	MessageCommon
	MediaId      string
	ThumbMediaId string
}

type ShortVideoMessage struct {
	MessageCommon
	MediaId      string
	ThumbMediaId string
}

type LocationMessage struct {
	MessageCommon
	LocationX float32
	LocationY float32
	Scale     int
	Label     string
}

type LinkMessage struct {
	MessageCommon
	Title       string
	Description string
	Url         string
}

func handleTextMessage(c *gin.Context, body []byte) {
	log.Info("handleTextMessage")
	msg := TextMessage{}
	err := xml.Unmarshal(body, &msg)
	if err != nil {
		log.Error("handleTextMessage err:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	rspmsg := TextMessage{
		Content: "你好，很高兴认识你",
		MessageCommon: MessageCommon{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      TEXT_MESSAGE,
		},
	}
	c.XML(200, rspmsg)
}

func handleImageMessage(c *gin.Context, body []byte) {
	log.Info("handleTextMessage")
	msg := ImageMessage{}
	err := xml.Unmarshal(body, &msg)
	if err != nil {
		log.Error("handleTextMessage err:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	rspmsg := TextMessage{
		Content: "图片很好看，我先收藏了",
		MessageCommon: MessageCommon{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      TEXT_MESSAGE,
		},
	}
	c.XML(200, rspmsg)
}
