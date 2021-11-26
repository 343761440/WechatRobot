package access

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wxrobot/internal/app/common"
	"wxrobot/internal/app/utils"

	log "github.com/sirupsen/logrus"
)

type IPType int

const (
	API_IPTYPE IPType = iota
	CALLBACK_IPTYPE
)

type ipValue struct {
	urlFormat string
	redisKey  string
}

var (
	gIPurlFormat = map[IPType]*ipValue{
		API_IPTYPE: {
			urlFormat: "https://api.weixin.qq.com/cgi-bin/get_api_domain_ip?access_token=%s",
			redisKey:  "WxWhiteAPIIPList",
		},
		CALLBACK_IPTYPE: {
			urlFormat: "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=%s",
			redisKey:  "WxWhiteCBIPList",
		},
	}
)

func getIPWhiteList(ipType IPType) ([]string, error) {
	type Response struct {
		IPList []string `json:"ip_list"`
	}

	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	ipv := gIPurlFormat[ipType]
	url := fmt.Sprintf(ipv.urlFormat, token)

	rsp, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	jrsp := Response{}
	err = json.Unmarshal(rsp, &jrsp)
	if err != nil {
		return nil, err
	}

	//这里还要处理一下，有些是网段，需要转换成ip
	//偷懒一下，只处理24的网络号
	var iplist []string
	for _, ip := range jrsp.IPList {
		if strings.Contains(ip, "/") {
			srcIp := strings.Split(ip, "/")[0]
			iIp := utils.InetAtoN(srcIp)
			for i := 1; i < 255; i++ {
				rangeIp := utils.InetNtoA(iIp + int64(i))
				iplist = append(iplist, rangeIp)
			}
		} else {
			iplist = append(iplist, ip)
		}
	}

	return iplist, nil
}

func GetWhiteList(ipType IPType) ([]string, error) {
	rdb, err := common.GetRedisClient()
	if err != nil {
		return nil, err
	}

	ipv := gIPurlFormat[ipType]
	key := ipv.redisKey

	isExist, err := rdb.Exists(key).Result()
	if err != nil {
		return nil, err
	}

	if isExist {
		iplist, err := rdb.SMembers(key).Result()
		if err != nil {
			return nil, err
		}
		return iplist, nil
	}

	iplist, err := getIPWhiteList(ipType)
	if err != nil {
		return nil, err
	}

	var members []interface{}
	for _, ip := range iplist {
		members = append(members, ip)
	}
	err = rdb.SAdd(key, members...).Err()
	if err != nil {
		log.Warn("redis sadd white list failed, err=", err)
		return iplist, nil
	}
	rdb.Expire(key, time.Hour*24)
	return iplist, nil
}

func IsCallbackWhiteList(clientip string) (bool, error) {
	rdb, err := common.GetRedisClient()
	if err != nil {
		return false, err
	}

	ipv := gIPurlFormat[CALLBACK_IPTYPE]
	key := ipv.redisKey
	isExist, err := rdb.Exists(key).Result()
	if err != nil {
		return false, err
	}

	if !isExist {
		_, err := GetWhiteList(CALLBACK_IPTYPE)
		if err != nil {
			return false, nil
		}
	}

	return rdb.SIsMember(key, clientip).Result()
}
