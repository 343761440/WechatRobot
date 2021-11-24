package wxmanager

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	kToken  = "serendipitymeeting6883"
	kAESkey = "WEhj6yTiFyuMshRfoZJdNbEoqNi5TyJvfswvgFnxVTW"
)

func getSignParam(c *gin.Context) (signature string, timestamp string, nonce string) {
	signature = c.Query("signature")
	timestamp = c.Query("timestamp")
	nonce = c.Query("nonce")
	return
}

func SignatureMiddleware(c *gin.Context) {
	log.Info("request:", c.Request)

	sign, timestamp, nonce := getSignParam(c)
	params := []string{}
	params = append(params, kToken, timestamp, nonce)
	sort.Strings(params)
	waitEncryStr := ""
	for _, p := range params {
		waitEncryStr += p
	}
	h := sha1.New()
	h.Write([]byte(waitEncryStr))
	sha1String := hex.EncodeToString(h.Sum([]byte{}))

	if sha1String == sign {
		c.Next()
	} else {
		log.Error("signature failed")
		c.JSON(401, gin.H{"msg": "signature failed"})
		c.Abort()
	}
}

func Signature(c *gin.Context) {
	echostr := c.Query("echostr")
	c.Writer.Write([]byte(echostr))
}

func RecvMessage(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ok"})
}
