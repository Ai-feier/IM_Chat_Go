package service

import (
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Accept json
// @Success 200 {string} Welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	t, err := template.ParseFiles("index.html", "view/chat/head.html")
	if err != nil {
		panic(err)
	}
	t.Execute(c.Writer, "index")
}

// ToRegister
// @Tags 注册
// @Accept json
// @Success 200 {string} register
// @Router /index [get]
func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}

	// 获取 userId, token
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token

	ind.Execute(c.Writer, user)
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
