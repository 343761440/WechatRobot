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
	Content string `xml:"Content"`
	MsgId   int64  `xml:"MsgId"`
}

type ImageMessage struct {
	MessageCommon
	PicUrl  string `xml:"PicUrl"`
	MediaId string `xml:"MediaId"`
	MsgId   int64  `xml:"MsgId"`
}

type VoiceMessage struct {
	MessageCommon
	MediaId     string `xml:"MediaId"`
	Format      string `xml:"Format"`
	Recognition string `xml:"Recognition"`
	MsgId       int64  `xml:"MsgId"`
}

type VideoMessage struct {
	MessageCommon
	MediaId      string `xml:"MediaId"`
	ThumbMediaId string `xml:"ThumbMediaId"`
	MsgId        int64  `xml:"MsgId"`
}

type ShortVideoMessage struct {
	MessageCommon
	MediaId      string `xml:"MediaId"`
	ThumbMediaId string `xml:"ThumbMediaId"`
	MsgId        int64  `xml:"MsgId"`
}

type LocationMessage struct {
	MessageCommon
	LocationX float32 `xml:"Location_X"`
	LocationY float32 `xml:"Location_Y"`
	Scale     int     `xml:"Scale"`
	Label     string  `xml:"Label"`
	MsgId     int64   `xml:"MsgId"`
}

type LinkMessage struct {
	MessageCommon
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	Url         string `xml:"Url"`
	MsgId       int64  `xml:"MsgId"`
}

func NewTextMessage(content, to, from string) TextMessage {
	return TextMessage{
		Content: content,
		MessageCommon: MessageCommon{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   time.Now().Unix(),
			MsgType:      TEXT_MESSAGE,
		},
	}
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
	wxUser := GetWxUserFromCtx(c)
	nickName := ""
	if wxUser != nil {
		nickName = wxUser.NickName
	}
	rspmsg := TextMessage{
		Content: nickName + "你好，很高兴认识你",
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
