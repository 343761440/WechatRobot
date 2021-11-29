package message

import "github.com/gin-gonic/gin"

type EventType string

const (
	SUBSCIBE_EVENT    EventType = "subscribe"
	UNSUBSCRIBE_EVENT EventType = "unsubscribe"
)

func handleEventMessage(c *gin.Context, body []byte) {

}
