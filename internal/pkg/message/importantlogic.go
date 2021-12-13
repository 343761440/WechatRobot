package message

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
	"wxrobot/internal/app/common"
	"wxrobot/internal/app/utils"
	"wxrobot/internal/pkg/log"
	"wxrobot/internal/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	imporQACmd     = "0"
	imporTodoCmd   = "1"
	imporWeathCmd  = "2"
	imporMovieCmd  = "3"
	imporRandomCmd = "4"
	imporRmbCmd    = "5"
	imporBookCmd   = "6"
)

var (
	gImportantDescribes = map[string]KeyWordInfo{
		//imporQACmd:     {"真心话环节", qaHandler},
		imporTodoCmd:   {"要做的xx件事", todoListHandler},
		imporWeathCmd:  {"实时天气", weatherHandler},
		imporMovieCmd:  {"电影推荐", movieRecoHandler},
		imporRandomCmd: {"随机匣子", coldjokeHandler},
		imporRmbCmd:    {"记仇小本本", jcDescriber},
		//imporBookCmd: {"书单推荐", bookDescriber},
	}
)

/* ------------------------ 以下为Important用户的handler ----------------------------*/

//TODO:每天只能看5个
func qaHandler(c *gin.Context, args ...string) {
	cmd := args[0]
	if cmd == imporQACmd {
		qlist, err := model.ListQuestions()
		if err != nil {
			log.ErrorWithRecord("ListQuestions failed, err=", err)
			c.XML(http.StatusOK, NewTextMessage("抱歉，目前出了点状况，请联系林同学", c))
			return
		}
		if len(qlist) == 0 {
			content := "你已经看完了所有的Q&A啦~\n"
			c.XML(http.StatusOK, NewTextMessage(content, c))
			return
		}
		content := "当前有以下Q&A可供查看，回复数字即可查看详细内容(例如：1001)\n"
		for _, q := range qlist {
			content += q.QuestionID + "." + q.Question + "\n"
		}
		content += "\n- - - - - - - - - - - - - - - - - - - - \n"
		content += "Tips：所有问题的回复只能查看一次，查看后就无法再次查看了哦。"
		c.XML(http.StatusOK, NewTextMessage(content, c))
		return
	} else {
		q, err := model.GetQuestion(cmd)
		if err != nil {
			c.XML(http.StatusOK, NewTextMessage("未找到当前数字对应的Question", c))
			return
		}
		content := "Q" + q.QuestionID + ":" + q.Question + "\n"
		content += "\n- - - - - - - - - - - - - - - - - - - - \n"
		content += q.Answer
		c.XML(http.StatusOK, NewTextMessage(content, c))
		model.UpdateQuestionReadness(cmd)
	}
}

func describeTodo(item model.TodoItem) string {
	num := fmt.Sprint(item.ID)
	finish := fmt.Sprintf("(%d/1)", item.FinishState)
	content := "\t" + num + "." + item.ItemInfo + " " + finish + "\n"
	return content
}

func describeTodos(items []model.TodoItem) string {
	sort.Slice(items, func(i, j int) bool {
		if items[i].FinishState > items[j].FinishState {
			return true
		} else if items[i].FinishState < items[j].FinishState {
			return false
		} else {
			return items[i].ID < items[j].ID
		}
	})
	nofinishCount := 0
	finishCount := 0
	seperator := false
	content := "当前有以下待做清单：\n"
	for _, todo := range items {
		if todo.FinishState > 0 {
			finishCount++
			seperator = false
		} else {
			nofinishCount++
			if !seperator {
				content += "\n- - - - - - - - - - - - - - - - - - - - \n"
			}
			seperator = true
		}
		content += describeTodo(todo)
	}
	content += "\n- - - - - - - - - - - - - - - - - - - - \n"
	content += fmt.Sprintf("目前已经完成了%d件待做事项，还有%d件待做事项等待完成哦，一起加油吧~\n", finishCount, nofinishCount)
	return content
}

func todoListHandler(c *gin.Context, args ...string) {
	cmd := args[0]
	if cmd == imporTodoCmd {
		todolist, err := model.ListTodoItems(model.TODO_ALL)
		if err != nil {
			log.ErrorWithRecord("ListTodoItems failed, err=", err)
			c.XML(http.StatusOK, NewTextMessage("抱歉，目前出了点状况，请联系林同学", c))
			return
		}
		if len(todolist) == 0 {
			c.XML(http.StatusOK, NewTextMessage("当前还没有待做的事项清单哦~", c))
			return
		}
		content := describeTodos(todolist)
		content += "\nTips：回复" + imporTodoCmd + "1" + "+待做事项可以补充清单\n"
		content += "(例如：" + imporTodoCmd + "1" + " 一起看日出"
		c.XML(http.StatusOK, NewTextMessage(content, c))
		return
	} else {
		substrs := strings.Split(args[0], " ")
		subcmd := substrs[0]
		if subcmd == imporTodoCmd+"1" {

			if len(substrs) >= 2 {
				info := ""
				for i := 1; i < len(substrs); i++ {
					info += substrs[i] + " "
				}
				if err := model.CreateTodoItems(model.TodoItem{ItemInfo: info}); err != nil {
					log.ErrorWithRecord("CreateTodoItems failed, err=", err)
					c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学", c))
				} else {
					c.XML(http.StatusOK, NewTextMessage("增加到待做清单成功！", c))
				}
				return
			} else {
				log.ErrorWithRecord("add todo failed by wrong format, msg=", cmd)
				c.XML(http.StatusOK, NewTextMessage("正确格式："+subcmd+" 待做事项 ("+subcmd+"与待做事项间的空格不要漏哦)", c))
				return
			}
		}
	}
	log.ErrorWithRecord("todoListHandler unknow message, msg=", cmd)
	c.XML(http.StatusOK, NewTextMessage("o(╥﹏╥)o我当前还无法消化这个信息", c))
}

func weatherHandler(c *gin.Context, args ...string) {
	urlformat := "https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%d"

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

	type Cast struct {
		Date         string `json:"date"`
		Week         string `json:"week"`
		DayWeather   string `json:"dayweather"`
		NightWeather string `json:"nightweather"`
		DayTemp      string `json:"daytemp"`
		NightTemp    string `json:"nighttemp"`
	}

	type Forecast struct {
		City       string `json:"city"`
		Reporttime string `json:"reporttime"`
		Casts      []Cast `json:"casts"`
	}

	type Result struct {
		Status   string    `json:"status"`
		Count    string    `json:"count"`
		Info     string    `json:"info"`
		InfoCode string    `json:"infocode"`
		Lives    []Lives   `json:"lives"`
		Forecast *Forecast `json:"forecast"`
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
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}

	var result Result
	err = json.Unmarshal(resp, &result)
	if err != nil {
		log.ErrorWithRecord("weather json unmarshal failed, err=", err, " res:", string(resp))
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}

	logrus.Info("result:", result)

	if result.InfoCode != "10000" {
		log.ErrorWithRecord("code not 10000, resp:", result)
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}

	if len(result.Lives) == 0 {
		log.ErrorWithRecord("lives is 0, resp:", result)
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}

	content := "当前地区：" + result.Lives[0].City + "\n"
	content += "天气：" + result.Lives[0].Weather + "\n"
	content += "实时气温：" + result.Lives[0].Temperature + "℃\n"
	content += "空气湿度" + result.Lives[0].Humidity + "\n"
	content += "风向描述：" + result.Lives[0].Winddirection + "\n"
	content += "风力级别：" + result.Lives[0].Windpower + "级\n"
	content += "发布时间：" + result.Lives[0].Reporttime + "\n"
	content += "\n- - - - - - - - - - - - - - - - - - - - \n"
	if result.Forecast != nil {
		if len(result.Forecast.Casts) > 0 {
			cast := result.Forecast.Casts[0]
			content += "白天温度：" + cast.DayTemp + "℃\n"
			content += "晚上温度：" + cast.NightTemp + "℃\n"
			content += "白天天气：" + cast.DayWeather + "\n"
			content += "晚上天气：" + cast.NightWeather + "\n"
			content += "\n- - - - - - - - - - - - - - - - - - - - \n"
		}
	}

	if result.Lives[0].Temperature < "10" {
		content += "当前温度只有" + result.Lives[0].Temperature + "℃" + "，出门一定一定要做好保暖工作哦"
	} else if result.Lives[0].Temperature >= "10" && result.Lives[0].Temperature < "20" {
		content += "当前温度为" + result.Lives[0].Temperature + "℃" + "，要好好穿衣服，不要感冒了"
	} else {
		content += "当前温度为" + result.Lives[0].Temperature + "℃" + "，稍有回暖，但不能松懈"
	}

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
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}

	var objects []Object
	err = json.Unmarshal(resp, &objects)
	if err != nil {
		log.ErrorWithRecord("Unmarshal movie failed, err=", err, " res:", string(resp))
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}
	if len(objects) == 0 {
		log.ErrorWithRecord("get movie objects failed, len is zero")
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}
	if len(objects[0].Datas) == 0 {
		log.ErrorWithRecord("get movie Datas failed, len is zero")
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}
	movie := objects[0]
	content := "电影名：《" + movie.Datas[0].Name + "》\n"
	content += "类型：" + movie.Datas[0].Genre + "\n"
	content += "语言：" + movie.Datas[0].Language + "\n"
	content += "国家：" + movie.Datas[0].Country + "\n"
	content += "发行时间：" + movie.Year + "\n"
	content += "海报：" + movie.Datas[0].Poster + "\n"
	content += "豆瓣评分：" + movie.DoubanRate + "\n"
	content += "影片概要：" + movie.Datas[0].Description + "\n"
	c.XML(http.StatusOK, NewTextMessage(content, c))
}

func coldjokeHandler(c *gin.Context, args ...string) {
	lovejoke, err := model.GetALoveJoke()
	if err != nil {
		log.ErrorWithRecord("get coldjokeHandler failed, err=", err)
		c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下林同学~", c))
		return
	}
	if len(lovejoke) == 0 {
		c.XML(http.StatusOK, NewTextMessage("请你看星星：我好想你", c))
		return
	}
	c.XML(http.StatusOK, NewTextMessage(lovejoke, c))
}

//记仇小本本
func jcDescriber(c *gin.Context, args ...string) {
	name := GetUserNameFromCtx(c)
	content := name + "打开了记仇小本本......\n"
	content += "是谁！惹" + name + "不高兴了！！！\n"
	content += "记仇小本本操作方法:"
	content += "\n- - - - - - - - - - - - - - - - - - - - \n"
	content += "写小本本：jc add xxxxx，例如：jc add 2021.12.13,小林审美不行\n"
	content += "\n- - - - - - - - - - - - - - - - - - - - \n"
	content += "显示小本本：jc ls\n"
	content += "\n- - - - - - - - - - - - - - - - - - - - \n"
	content += "修改小本本：jc mod 事件id 事件内容,例如：jc mod 3 2021.12.13,小林审美是真的不行\n"
	c.XML(http.StatusOK, NewTextMessage(content, c))
}

func bookDescriber(c *gin.Context, args ...string) {

}

func jcController(c *gin.Context, args ...string) {

}

//书单推荐
func bookController(c *gin.Context, args ...string) {

}

func importantCmdAnalyze(c *gin.Context, content string) {
	userName := ""
	isBirth := false
	user := GetWxUserFromCtx(c)
	if user != nil {
		userName = user.NickName
		isBirth = user.IsBirthday()
	}

	//有没有命中关键词命令
	strs := strings.Split(content, " ")
	if len(strs) > 0 {
		handler, ok := gImportantHandlers[strs[0]]
		if ok {
			handler(c, strs[1:]...)
			return
		}
	}

	//没命中则进入主菜单
	handler, ok := gImportantDescribes[content[0:1]]
	if !ok {
		content := "Hi~ " + userName + "\n"
		if isBirth {
			log.Info("time:", time.Now().Format("2006-01-02"))
			content += "今天是" + time.Now().Format("2006-01-02") + " ,是你的破蛋日，祝你生日快乐~\n"
		}
		content += "回复前面的数字即可进入下列选项，可重复回复哦\n"
		var keys []string
		for key := range gImportantDescribes {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			content += "\t" + key + "." + gImportantDescribes[key].description + "\n"
		}
		c.XML(http.StatusOK, NewTextMessage(content, c))
		return
	}
	handler.handler(c, content)
}
