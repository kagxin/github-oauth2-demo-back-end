package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	// ClientID github application ClientID
	ClientID = ""
	// ClientSecret github application ClientSecret
	ClientSecret = ""
	// RedirectURI code 从定向的地址
	RedirectURI = "http://localhost:8080/oauth2"
	// GithubTokenByCode 获取 token
	GithubTokenByCode = "https://github.com/login/oauth/access_token"
	// GithubUserInfoByToken 获取用户信息
	GithubUserInfoByToken = "https://api.github.com/user"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/oauth2", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("param error"))
			return
		}
		GithubToken := getToken(code)
		userInfo := getUser(GithubToken["access_token"].(string))
		// TODO: 创建用户信息，生成用户 token，返回给前端
		c.JSON(http.StatusOK, userInfo)
	})
	router.Run("0.0.0.0:8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getToken(code string) map[string]interface{} {

	jsonStr, err := json.Marshal(map[string]string{
		"client_id":     ClientID,
		"client_secret": ClientSecret,
		"code":          code,
	})
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", GithubTokenByCode, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	rsp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}

	token := make(map[string]interface{})
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}
	return token
}

func getUser(token string) map[string]interface{} {
	client := http.Client{}
	req, err := http.NewRequest("GET", GithubUserInfoByToken, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Accept", "application/json")

	rsp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}

	userInfo := make(map[string]interface{})
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		panic(err)
	}
	return userInfo
}
