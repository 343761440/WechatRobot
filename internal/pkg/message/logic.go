package message

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"wxrobot/internal/app/common"
	"wxrobot/internal/app/utils"
	"wxrobot/internal/pkg/log"
	"wxrobot/internal/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*
 业务逻辑控制主要在这个文件
 1、小林同学命名新用户
	1.1 小林同学获取所有userId
	1.2 小林同学设置人员类型
	1.3
 2、重要人的QA
 3、和重要人的待完成事项
 4、普通用户的默认返回
 5、重要人员不能识别命令的返回(随机土味之类的，或者再找找有没有小冰的SDK)
 6、订阅其他用户的留言:
	6.1 记录留言
	6.2 查看留言
*/

type KeyWordHandle func(c *gin.Context, args ...string)

type NextUser struct {
	NextName string //下一个注册NickName
	UserType int    // important / ""
	sync.Mutex
}

type KeyWordInfo struct {
	description string
	handler     KeyWordHandle
}

var (
	gRootHandlers = map[string]KeyWordHandle{
		"new": newUserHandler,
		"err": lastErrorHandler,
	}
	gImportantHandler = map[string]KeyWordInfo{
		"1": {"真心话环节", qaHandler},
		"2": {"要做的xx件事", todoListHandler},
		"3": {"实时天气", weatherHandler},
		"4": {"电影推荐", movieRecoHandler},
		"5": {"最近在听", musicRecoHandler},
		"6": {"今日份冷笑话", coldjokeHandler},
	}
	gNextUser = NextUser{}
	gUserType = map[string]int{
		"important": model.USER_IMPORTANT,
		"normal":    model.USER_NORMAL,
		"friend":    model.USER_FRIEND,
	}
)

func (nu *NextUser) clear() {
	nu.NextName = ""
	nu.UserType = model.USER_NORMAL
}

func matchUserType(userType string) (int, bool) {
	tp, ok := gUserType[userType]
	return tp, ok
}

func newUserHandler(c *gin.Context, args ...string) {
	gNextUser.Lock()
	defer gNextUser.Unlock()

	if len(args) != 2 {
		c.XML(http.StatusOK, NewTextMessage("格式输入有误，正确格式：new username userType", c))
		return
	}
	utype, ok := matchUserType(args[1])
	if !ok {
		c.XML(http.StatusOK, NewTextMessage("无法识别的userType", c))
		return
	}

	userName := args[0]
	gNextUser.NextName = userName
	gNextUser.UserType = utype
	c.XML(http.StatusOK, NewTextMessage(args[1]+"用户设置成功", c))
}

func lastErrorHandler(c *gin.Context, args ...string) {
	lasterr := log.GetLastError()
	if len(lasterr) == 0 {
		lasterr = "NoError"
	}
	c.XML(http.StatusOK, NewTextMessage(lasterr, c))
}

func getNextUser() *model.WxUser {
	gNextUser.Lock()
	defer gNextUser.Unlock()

	if len(gNextUser.NextName) > 0 {
		res := &model.WxUser{
			NickName: gNextUser.NextName,
			UserType: gNextUser.UserType,
		}
		gNextUser.clear()
		return res
	} else {
		return nil
	}
}

//TODO:每天只能看5个
func qaHandler(c *gin.Context, args ...string) {
	cmd := args[0]
	if cmd == "1" {
		qlist, err := model.ListQuestions()
		if err != nil {
			log.ErrorWithRecord("ListQuestions failed, err=", err)
			c.XML(http.StatusOK, NewTextMessage("抱歉，目前出了点状况，请联系小林同学", c))
			return
		}
		if len(qlist) == 0 {
			c.XML(http.StatusOK, NewTextMessage("当前没有可以查看的Q&A了~", c))
			return
		}
		content := "当前有以下Q&A可供查看，回复问题前的数字即可查看详细内容(例如1001\n"
		for _, q := range qlist {
			content += q.QuestionID + "." + q.Question + "\n"
		}
		content += "Tips：所有问题的回复只能查看一次，查看后就无法再次查看了哦"
		c.XML(http.StatusOK, NewTextMessage(content, c))
		return
	} else {
		q, err := model.GetQuestion(cmd)
		if err != nil {
			c.XML(http.StatusOK, NewTextMessage("未找到当前数字对应的Question", c))
			return
		}
		content := "Q" + q.QuestionID + ".:" + q.Question + "\n"
		content += q.Answer
		c.XML(http.StatusOK, NewTextMessage(content, c))
	}
}

func todoListHandler(c *gin.Context, args ...string) {
	cmd := args[0]
	if cmd == "2" {
		todolist, err := model.ListTodoItems(model.TODO_ALL)
		if err != nil {
			log.ErrorWithRecord("ListTodoItems failed, err=", err)
			c.XML(http.StatusOK, NewTextMessage("抱歉，目前出了点状况，请联系小林同学", c))
			return
		}
		if len(todolist) == 0 {
			c.XML(http.StatusOK, NewTextMessage("当前还没有待做的事项清单哦~", c))
			return
		}
		nofinishCount := 0
		finishCount := 0
		content := "当前有以下待做清单：\n"
		for _, todo := range todolist {
			if todo.FinishState > 0 {
				finishCount++
			} else {
				nofinishCount++
			}
			num := fmt.Sprint(todo.ID)
			finish := fmt.Sprintf("(%d/1)", todo.FinishState)
			content += "\t" + num + "." + todo.ItemInfo + " " + finish + "\n"
		}
		content += "\n- - - - - - - - - - - - - - - - - - - - \n"
		content += fmt.Sprintf("目前已经完成了%d件待做事项，还有%d件待做事项等待完成哦，一起加油吧~\n", finishCount, nofinishCount)
		content += "\nTips：回复21+待做事项可以补充清单\n"
		content += "(例如：21 一起看日出"
		c.XML(http.StatusOK, NewTextMessage(content, c))
		return
	} else {
		substrs := strings.Split(args[0], " ")
		subcmd := substrs[0]
		if subcmd == "21" {
			if len(substrs) >= 2 {
				info := ""
				for i := 1; i < len(substrs); i++ {
					info += substrs[i] + " "
				}
				if err := model.CreateTodoItems(model.TodoItem{ItemInfo: info}); err != nil {
					log.ErrorWithRecord("CreateTodoItems failed, err=", err)
					c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下小林同学", c))
				} else {
					c.XML(http.StatusOK, NewTextMessage("增加到待做清单成功！", c))
				}
				return
			} else {
				c.XML(http.StatusOK, NewTextMessage("正确格式：21 待做事项 (21与待做事项间的空格不要漏哦)", c))
				return
			}
		}
	}
	c.XML(http.StatusOK, NewTextMessage("o(╥﹏╥)o我当前还无法消化这个信息", c))
}

func weatherHandler(c *gin.Context, args ...string) {
	urlformat := "https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%d&extensions=base"

	type Lives struct {
		Province      string `json:"province"`
		City          string `json:"city"`
		Adcode        string `json:"adcode"`
		Weather       string `json:"weather"`
		Temperature   string `json:"temperature"`
		Winddirection string `json:"winddirection"`
		Windpower     string `json:"windpower"`
		Humidity      string `json:"humidity"`
		Reporttime    string `json:"reporttime"`
	}

	type Result struct {
		Status   string  `json:"status"`
		Count    string  `json:"count"`
		Info     string  `json:"info"`
		InfoCode string  `json:"infocode"`
		Lives    []Lives `json:"lives"`
	}

	// 默认滨江区
	// 杭州市		330100
	// 杭州市市辖区	330101
	// 上城区		330102
	// 下城区		330103
	// 江干区		330104
	// 拱墅区		330105
	// 西湖区		330106
	// 滨江区		330108
	// 萧山区		330109
	// 余杭区		330110
	url := fmt.Sprintf(urlformat, "d1c4ce5567b24fee573915a2d3d8110e", 330108)
	resp, err := utils.HttpGet(url)
	if err != nil {
		log.ErrorWithRecord("weather HttpGet failed, err=", err)
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下小林同学~", c))
		return
	}

	var result Result
	err = json.Unmarshal(resp, &result)
	if err != nil {
		log.ErrorWithRecord("weather json unmarshal failed, err=", err, " res:", string(resp))
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下小林同学~", c))
		return
	}

	logrus.Info("result:", result)

	if result.InfoCode != "10000" {
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下小林同学~", c))
		return
	}

	content := "当前地区：" + result.Lives[0].City + "\n"
	content += "天气：" + result.Lives[0].Weather + "\n"
	content += "实时气温：" + result.Lives[0].Temperature + "℃\n"
	content += "空气湿度" + result.Lives[0].Humidity + "\n"
	content += "风向描述：" + result.Lives[0].Winddirection + "\n"
	content += "风力级别：" + result.Lives[0].Windpower + "级\n"
	content += "\n- - - - - - - - - - - - - - - - - - - - \n"
	content += "发布时间：" + result.Lives[0].Reporttime + "\n"

	//得加一些文案，表示提醒，例如升温，降温等
	c.XML(http.StatusOK, NewTextMessage(content, c))
}

func getGlobalMovieIncrbyNum() (int64, error) {
	rdb, err := common.GetRedisClient()
	if err != nil {
		return 0, err
	}

	key := "movieGloabNum"
	num, err := rdb.Incr(key).Result()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func movieRecoHandler(c *gin.Context, args ...string) {
	type Data struct {
		ID          string `json:"id"`
		Poster      string `json:"poster"`
		Name        string `json:"name"`
		Genre       string `json:"genre"`
		Description string `json:"description"`
		Language    string `json:"language"`
		Country     string `json:"country"`
		Lang        string `json:"lang"`
		Movie       string `json:"movie"`
	}

	type Object struct {
		Datas      []Data `json:"data"`
		OriginName string `json:"originalName"`
		DoubanRate string `json:"doubanRating"`
		Alias      string `json:"alias"`
		Year       string `json:"year"`
	}

	urlformat := "https://api.wmdb.tv/api/v1/top?type=Douban&skip=%d&limit=1&lang=Cn"
	num, err := getGlobalMovieIncrbyNum()
	if err != nil {
		log.ErrorWithRecord("getGlobalMovieIncrbyNum failed,err=", err)
		num = int64(rand.Intn(1000) + 1)
	}

	url := fmt.Sprintf(urlformat, num)
	resp, err := utils.HttpGet(url)
	if err != nil {
		log.ErrorWithRecord("HttpGet movie failed, err=", err, " res:", string(resp))
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下小林同学~", c))
		return
	}

	var objects []Object
	err = json.Unmarshal(resp, &objects)
	if err != nil {
		log.ErrorWithRecord("Unmarshal movie failed, err=", err, " res:", string(resp))
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下小林同学~", c))
		return
	}
	movie := objects[0]
	content := "电影名：《" + movie.Datas[0].Name + "》\n"
	content += "类型：" + movie.Datas[0].Genre + "\n"
	content += "语言：" + movie.Datas[0].Language + "\n"
	content += "国家：" + movie.Datas[0].Country + "\n"
	content += "海报：" + movie.Datas[0].Poster + "\n"
	content += "豆瓣评分：" + movie.DoubanRate + "\n"
	content += "影片概要：" + movie.Datas[0].Description + "\n"
	c.XML(http.StatusOK, NewTextMessage(content, c))
}

func musicRecoHandler(c *gin.Context, args ...string) {

}

func coldjokeHandler(c *gin.Context, args ...string) {

}

func rootCmdAnalyze(c *gin.Context, content string) {
	strs := strings.Split(content, " ")
	if len(strs) > 0 {
		handler, ok := gRootHandlers[strs[0]]
		if ok {
			handler(c, strs[1:]...)
			return
		}
	}
	cmds := ""
	for cmd := range gRootHandlers {
		cmds += cmd + "\n"
	}
	c.XML(http.StatusOK, NewTextMessage("无法识别的指令,当前支持的指令:"+cmds, c))
}

func importantCmdAnalyze(c *gin.Context, content string) {
	userName := ""
	user := GetWxUserFromCtx(c)
	if user != nil {
		userName = user.NickName
	}
	handler, ok := gImportantHandler[content[0:1]]
	if !ok {
		content := "Hi~ " + userName + " 回复前面的数字即可进入下列选项哦\n"
		var keys []string
		for key := range gImportantHandler {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			content += "\t" + key + "." + gImportantHandler[key].description + "\n"
		}
		c.XML(http.StatusOK, NewTextMessage(content, c))
		return
	}
	handler.handler(c, content)
}

func friendCmdAnalyze(c *gin.Context, content string) {
	c.XML(http.StatusOK, NewTextMessage("谢谢你的留言。", c))
}

func defaultResponse(c *gin.Context) {
	c.XML(http.StatusOK, NewTextMessage("本订阅号正在维护中...", c))
}

func HandleLogicMessage(c *gin.Context, content string) {
	user := GetWxUserFromCtx(c)
	switch user.UserType {
	case model.USER_ADMIN:
		rootCmdAnalyze(c, content)
	case model.USER_FRIEND:
		friendCmdAnalyze(c, content)
	case model.USER_IMPORTANT:
		importantCmdAnalyze(c, content)
	case model.USER_NORMAL:
		defaultResponse(c)
	default:
		defaultResponse(c)
	}
}
