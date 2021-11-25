package menu

import (
	"fmt"
	"wxrobot/internal/pkg/access"

	log "github.com/sirupsen/logrus"
)

const (
	kCreateMenuUrlFormat = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"
)

type MenuButton struct {
}

func CreateMenu() error {

	token, err := access.GetAccessToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(kCreateMenuUrlFormat, token)
	log.Info(url)
	return nil
}
