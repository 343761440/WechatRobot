package message

import (
	"encoding/xml"

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
)

type MessageCommon struct {
	XMLName      xml.Name    `xml:"xml"` // 指定最外层的标签为config
	ToUserName   string      `xml:"ToUserName"`
	FromUserName string      `xml:"FromUserName"`
	CreateTime   int64       `xml:"CreateTime"`
	MsgType      MessageType `xml:"MsgType"`
	MsgId        int64       `xml:"MsgId"`
}

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

func HandleMessage(c *gin.Context) {
	log.Info("request:", c.Request)
	// bodyReader, err := c.Request.GetBody()
	// if err != nil {
	// 	log.Info("get body failed,err=", err)
	// 	c.XML(500, "")
	// 	return
	// }

	// body, _ := ioutil.ReadAll(bodyReader)
	// log.Info("body:", body)
	msg := TextMessage{}
	c.BindXML(&msg)
	log.Info("msg:", msg)

	rspmsg := TextMessage{
		Content: "hello",
		MessageCommon: MessageCommon{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   msg.CreateTime + 1,
			MsgType:      TEXT_MESSAGE,
		},
	}
	c.XML(200, rspmsg)
}
