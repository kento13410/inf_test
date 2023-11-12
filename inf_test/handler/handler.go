package handler

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"inf_test/db"
	"inf_test/model"
	"net/http"
	"strings"
)

func Signup(c *gin.Context) {
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account creation failed",
			"cause":   "required user_id and password",
		})
		return
	}

	if !isValidLength(user) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account creation failed",
			"cause":   "**User ID:** 6-20 characters\n**Password:** 8-20 characters",
		})
		return
	}

	if !isAlphaNumeric(user.UserID) || !isAlphaNumeric(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account creation failed",
			"cause":   "required user_id and password alpha numeric",
		})
		return
	}

	if err := db.AddUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account creation failed",
			"cause":   "already same user_id is used",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account successfully created",
		"user": gin.H{
			"user_id":  user.UserID,
			"password": user.Password,
		},
	})
}

func isValidLength(user model.User) bool {
	lenID, lenPass := len(user.UserID), len(user.Password)
	return 6 <= lenID && lenID <= 20 && 8 <= lenPass && lenPass <= 20
}

func isAlphaNumeric(str string) bool {
	for _, char := range str {
		if (char < '0' || char > '9') && (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return false
		}
	}
	return true
}

func GetUser(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	userID, password, err := DecodeAuthorizationHeader(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication Failed"})
		return
	}
	user, err := db.SelectUser(userID, password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No User found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User details by user_id",
		"user":    user,
	})
}

func DecodeAuthorizationHeader(header string) (username, password string, err error) {
	// Authorization ヘッダーから "Basic" を除去する
	header = strings.TrimPrefix(header, "Basic ")

	// Base64 でデコードする
	decoded, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return "", "", err
	}

	// ユーザー名とパスワードを取得する
	username, password = strings.Split(string(decoded), ":")[0], strings.Split(string(decoded), ":")[1]

	return username, password, nil
}

func UpdateUser(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	userID, password, err := DecodeAuthorizationHeader(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication Failed"})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User updation failed",
			"cause":   "required nickname or comment",
		})
		return
	}

	if user.UserID != "" || user.Password != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Can not change user_id and password"})
		return
	}

	user.UserID, user.Password = userID, password

	if err := db.UpdateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No User found"})
		return
	}

	user.UserID, user.Password = "", ""
	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully updated",
		"recipe":  user,
	})
}

func DeleteUser(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	userID, password, err := DecodeAuthorizationHeader(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication Failed"})
		return
	}

	if err := db.DeleteUser(userID, password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deletion Failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account and user successfully removed"})
}
