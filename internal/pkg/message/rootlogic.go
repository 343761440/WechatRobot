package message

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"wxrobot/internal/pkg/log"
	"wxrobot/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

var (
	gRootHandlers = map[string]KeyWordHandle{
		"new":     newUserHandler,
		"err":     lastErrorHandler,
		"lsuser":  listUsersHandler,
		"todo":    todoController,
		"sub":     subscribeController,
		"remeber": rememberController,
		"joke":    addJokeController,
	}
	gImportantHandlers = map[string]KeyWordHandle{
		"jc": jcController,
	}
	gNextUser = NextUser{}
	gUserType = map[string]int{
		"Important": model.USER_IMPORTANT,
		"Normal":    model.USER_NORMAL,
		"Friend":    model.USER_FRIEND,
	}
)

func getNextUser() *model.WxUser {
	gNextUser.Lock()
	defer gNextUser.Unlock()

	if len(gNextUser.User.NickName) > 0 {
		res := &gNextUser.User
		gNextUser.clear()
		return res
	} else {
		return nil
	}
}

func (nu *NextUser) clear() {
	nu.User = model.WxUser{}
}

func matchUserType(userType string) (int, bool) {
	tp, ok := gUserType[userType]
	return tp, ok
}

func newUserHandler(c *gin.Context, args ...string) {
	gNextUser.Lock()
	defer gNextUser.Unlock()

	if len(args) < 2 {
		c.XML(http.StatusOK, NewTextMessage("格式输入有误，正确格式：new username userType birthday(可选)", c))
		return
	}
	utype, ok := matchUserType(args[1])
	if !ok {
		c.XML(http.StatusOK, NewTextMessage("无法识别的userType:"+args[1], c))
		return
	}

	userName := args[0]
	gNextUser.User.NickName = userName
	gNextUser.User.UserType = utype
	if len(args) >= 3 {
		gNextUser.User.Birthday = args[2]
	}
	c.XML(http.StatusOK, NewTextMessage(args[1]+"用户设置成功", c))
}

func lastErrorHandler(c *gin.Context, args ...string) {
	lasterr := log.GetLastError()
	if len(lasterr) == 0 {
		c.XML(http.StatusOK, NewTextMessage("NoError Happend", c))
	}
	content := "Error List:"
	for _, e := range lasterr {
		content += e + "\n"
	}
	c.XML(http.StatusOK, NewTextMessage(content, c))
}

func listUsersHandler(c *gin.Context, args ...string) {
	users, err := model.ListWxUsers()
	if err != nil {
		c.XML(http.StatusOK, NewTextMessage(err.Error(), c))
		return
	}

	DescribeUser := func(wu model.WxUser) string {
		content := "ID：" + fmt.Sprint(wu.ID) + "\n"
		content += "NickName：" + wu.NickName + "\n"
		content += "Birthday：" + wu.Birthday + "\n"
		content += "UserType：" + wu.GetUserType() + "\n"
		return content
	}

	res := ""
	for _, user := range users {
		res += DescribeUser(user)
		res += "\n- - - - - - - - - - - - - - - - - - - - \n"
	}
	c.XML(http.StatusOK, NewTextMessage(res, c))
}

func todoController(c *gin.Context, args ...string) {
	if len(args) == 0 {
		c.XML(http.StatusOK, NewTextMessage("todo needs ls/add/fin/del/mod", c))
		return
	}

	if args[0] == "ls" {
		items, err := model.ListTodoItems(model.TODO_ALL)
		if err != nil {
			c.XML(http.StatusOK, NewTextMessage("ListTodoItems failer,err="+err.Error(), c))
			return
		}
		content := describeTodos(items)
		c.XML(http.StatusOK, NewTextMessage(content, c))
		return
	} else if args[0] == "fin" {
		if len(args) < 2 {
			c.XML(http.StatusOK, NewTextMessage("todo fin+id", c))
			return
		}
		id, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			c.XML(http.StatusOK, NewTextMessage("todo parse arg2 failed,err="+err.Error(), c))
			return
		}
		err = model.UpdateTodoFinish(int(id), 1)
		if err != nil {
			c.XML(http.StatusOK, NewTextMessage("todo update finish state failed,err="+err.Error(), c))
			return
		}
		c.XML(http.StatusOK, NewTextMessage("todo finish success", c))
	} else if args[0] == "add" {
		if len(args) >= 2 {
			info := ""
			for i := 1; i < len(args); i++ {
				info += args[i] + " "
			}
			if err := model.CreateTodoItems(model.TodoItem{ItemInfo: info}); err != nil {
				c.XML(http.StatusOK, NewTextMessage("add todo failed,err="+err.Error(), c))
			} else {
				c.XML(http.StatusOK, NewTextMessage("add todo success!", c))
			}
			return
		} else {
			c.XML(http.StatusOK, NewTextMessage("正确格式：todo add xxxxx", c))
			return
		}
	}
}

//订阅某人的最近输入
//目前只针对Important
func subscribeController(c *gin.Context, args ...string) {
	if len(args) == 0 {
		c.XML(http.StatusOK, NewTextMessage("sub needs name, like:sub 张三", c))
		return
	}
}

//查看记仇小本本
func rememberController(c *gin.Context, args ...string) {
	if len(args) == 0 {
		c.XML(http.StatusOK, NewTextMessage("remeber ls", c))
		return
	}
}

//查看love joke
func addJokeController(c *gin.Context, args ...string) {
	if len(args) == 0 {
		c.XML(http.StatusOK, NewTextMessage("joke ls/add", c))
		return
	}
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
