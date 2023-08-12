package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [GET]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 8)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

// FindByNameAndPwd
// @Summary 通过用户名和密码查询
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findByNameAndPwd [get]
func FindByNameAndPwd(c *gin.Context) {
	name := c.Request.FormValue("name")
	data := models.FindByName(name)
	if data.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}

	newPwd := c.Request.FormValue("password")
	flag := utils.ValidPassword(newPwd, data.PassWord, data.Salt)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "登录成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 创建用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{
		LoginTime:     time.Now(),
		HeartbeatTime: time.Now(),
		LoginOutTime:  time.Now(),
	}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")
	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "两次密码不一致！",
			"data":    user,
		})
		return
	}

	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "用户名或密码不能为空！",
			"data":    user,
		})
		return
	}

	// 查询用户名是否存在
	data := models.FindByName(user.Name)
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "用户名已注册！",
			"data":    user,
		})
		return
	}

	// 校验手机号与邮箱
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"messsage": "手机号或邮箱不匹配",
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致",
		})
		return
	}
	// 加随机数校验
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	user.PassWord = utils.MakePassword(password, salt)
	models.CreateUser(&user)
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "新增用户成功！",
		"data":    user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := &models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "删除用户成功！",
		"data":    user,
	})

}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param email formData string false "邮箱"
// @param phone formData string false "手机号"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	user.Phone = c.PostForm("phone")
	user.PassWord = c.PostForm("password")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "修改参数不符合要求！",
			"data":    user,
		})
		return
	}
	models.UpdateUser(&user)
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "修改用户成功！",
		"data":    user,
	})
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	upgrade, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "发送错误",
		})
	}
	defer func() {
		err2 := upgrade.Close()
		if err2 != nil {
			fmt.Println(err2)
			return
		}
	}()
	MsgHandle(c, upgrade)
}

func MsgHandle(c *gin.Context, upgrade *websocket.Conn) {
	// 死循环用于服务器能够持续监听，发送消息
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
		}
		t := time.Now().Format("2003-03-19 07:03:19")
		s := fmt.Sprintf("%s: %s", t, msg)
		err = upgrade.WriteMessage(1, []byte(s))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
