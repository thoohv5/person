package middleware

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	langHeaderTag = "accept-language"
	langQueryTag  = "lang"

	userIDTag = "Userid"
)

func Params(langLabel string, userLabel string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang(langLabel, c)
		user(userLabel, c)

		c.Next()
	}
}

// lang 从请求头中获取语言信息
func lang(label string, c *gin.Context) {
	cLang := ""
	language := c.Request.Header.Get(langHeaderTag)
	headerLangList := strings.Split(language, ",")
	if len(language) != 0 && len(headerLangList) > 0 {
		cLang = headerLangList[0]
	}
	queryLang := c.Request.URL.Query().Get(langQueryTag)
	if len(queryLang) > 0 && len(cLang) == 0 {
		cLang = queryLang
	}

	if len(cLang) > 0 {
		c.Request.Header.Set(label, cLang)
	}
}

// user 从请求头中获取用户信息
func user(label string, c *gin.Context) {
	cUserID := ""

	if iUserID := c.Request.Header.Get(label); len(iUserID) > 0 {
		return
	}

	userID := c.Request.Header.Get(userIDTag)
	if len(userID) > 0 {
		cUserID = userID
	}
	// 从请求参数中获取 UserID
	queryUserId := c.Request.URL.Query().Get(userIDTag)
	if len(queryUserId) > 0 {
		cUserID = queryUserId
	}

	if len(cUserID) > 0 {
		// 这里换一下换成url编码
		decodedString, err := url.QueryUnescape(cUserID)
		if err != nil {
			return
		}
		c.Request.Header.Set(label, decodedString)
	}
}
