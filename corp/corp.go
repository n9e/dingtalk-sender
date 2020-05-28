package corp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/toolkits/pkg/logger"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const dingTimeOut = time.Second * 1

// Err
type Err struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// Client
type Client struct {
	Mobiles []string
	token   string
	openUrl string
	IsAtAll bool
}

// Result 发送消息返回结果
type Result struct {
	Err
}

// New
func New(token string, mobiles []string, isAtAll bool) *Client {
	c := new(Client)
	c.openUrl = "https://oapi.dingtalk.com/robot/send?access_token="
	c.token = token
	c.Mobiles = mobiles
	c.IsAtAll = isAtAll
	return c
}

func (c Client) GetToken() string {
	return c.token
}

// Send 发送信息
func (c *Client) Send(token, msg string) error {

	postData := c.generateData(msg)
	if c.GetToken() != "" {
		// 配置了token 说明采用配置文件的token
		token = c.GetToken()
	}
	url := c.openUrl + token

	resultByte, err := jsonPost(url, postData)
	if err != nil {
		return fmt.Errorf("invoke send api fail: %v", err)
	}

	result := Result{}
	err = json.Unmarshal(resultByte, &result)
	if err != nil {
		return fmt.Errorf("parse send api response fail: %v", err)
	}

	if result.ErrCode != 0 || result.ErrMsg != "ok" {
		err = fmt.Errorf("invoke send api return ErrCode = %d, ErrMsg = %s ", result.ErrCode, result.ErrMsg)
	}

	return err
}

func jsonPost(url string, data interface{}) ([]byte, error) {
	jsonBody, err := encodeJSON(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		logger.Info("ding talk new post request err =>", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("ding talk post request err =>", err)
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func encodeJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) generateData(msg string) interface{} {
	postData := make(map[string]interface{})
	postData["msgtype"] = "text"
	sendContext := make(map[string]interface{})
	sendContext["content"] = msg
	postData["text"] = sendContext

	at := make(map[string]interface{})
	if !c.IsAtAll && len(c.Mobiles) > 0 && c.token != "" {
		at["atMobiles"] = c.Mobiles // 根据手机号@指定人
	} else {
		c.IsAtAll = true
	}

	at["isAtAll"] = c.IsAtAll // s是否@所有人
	postData["at"] = at

	return postData
}

func getClient() *http.Client {
	// 通过http.Client 中的 DialContext 可以设置连接超时和数据接受超时 （也可以使用Dial, 不推荐）
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.DialTimeout(network, addr, dingTimeOut) // 设置建立链接超时
				if err != nil {
					return nil, err
				}
				_ = conn.SetDeadline(time.Now().Add(dingTimeOut)) // 设置接受数据超时时间
				return conn, nil
			},
			ResponseHeaderTimeout: dingTimeOut, // 设置服务器响应超时时间
		},
	}
}
