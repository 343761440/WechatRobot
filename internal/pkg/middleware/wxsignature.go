package middleware

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"wxrobot/internal/app/common"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getSignParam(c *gin.Context) (signature string, timestamp string, nonce string) {
	signature = c.Query("signature")
	timestamp = c.Query("timestamp")
	nonce = c.Query("nonce")
	return
}

func SignatureMiddleware(c *gin.Context) {
	log.Debug("request:", c.Request)
	kToken := viper.GetString(common.CFG_SECRET_TOKEN)
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
		c.JSON(401, gin.H{"msg": "signature failed"})
		c.Abort()
	}
}
