package message

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"wxrobot/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

/*
 业务逻辑控制主要在这个文件
 1、管理员命名新用户
	1.1 管理员获取所有userId
	1.2 管理员设置人员类型
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
	}
	gImportantHandler = map[string]KeyWordInfo{
		"1": {"真心话环节", qaHandler},
		"2": {"大冒险环节", todoListHandler},
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

func qaHandler(c *gin.Context, args ...string) {
	cmd := args[0]
	if cmd == "1" {
		qlist, err := model.ListQuestions()
		if err != nil {
			c.XML(http.StatusOK, NewTextMessage("抱歉，目前出了点状况，请联系管理员", c))
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
			c.XML(http.StatusOK, NewTextMessage("抱歉，目前出了点状况，请联系管理员", c))
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
					c.XML(http.StatusOK, NewTextMessage("我暂时出了点问题，请联系一下管理员", c))
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
		for cmd, hinfo := range gImportantHandler {
			content += "\t" + cmd + "." + hinfo.description + "\n"
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
