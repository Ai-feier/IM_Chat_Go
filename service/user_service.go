package service

import (
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// CreateUser
// @Summary 创建用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [POST]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致",
		})
	}
	user.PassWord = password
	models.CreateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "用户添加成功",
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param name query string false "用户名"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [DELETE]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)

	models.CreateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "用户添加成功",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param name formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [PUT]
func UpdateUser(c *gin.Context) {
	user := &models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")

	models.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"message": "用户添加成功",
	})
}
