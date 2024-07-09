package service

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GetOneAnswer(c *gin.Context, msg string) {
	url := "https://api.openai-hk.com/v1/chat/completions"
	apiKey := conf.Conf.Gpt.Key

	payload := map[string]interface{}{
		"max_tokens":       1200,
		"model":            "gpt-3.5-turbo",
		"temperature":      0.8,
		"top_p":            1,
		"presence_penalty": 1,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible.",
			},
			{
				"role":    "user",
				"content": msg,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		common.ResError(c, err.Error())
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}
	defer resp.Body.Close()

	// 处理响应
	// 请根据实际需求解析和处理响应数据
	respBody, err := io.ReadAll(resp.Body)

	reg := regexp.MustCompile(`( )+|(\n)+`)
	after := reg.ReplaceAllString(string(respBody), "$1$2")
	var resData model.GptResponse
	_ = json.Unmarshal([]byte(strings.ReplaceAll(after, "\n", "")), &resData)
	_, err = conf.Mysql.Insert(model.GptQuestion{
		Question:  msg,
		Answer:    resData.Choices[0].Message.Content,
		CreatedAt: int(time.Now().Unix()),
	})
	common.ResOk(c, "ok", resData)
}

func QuestionList(c *gin.Context, search string, page, pageSize int) {
	var questions []*model.GptQuestion
	sess := conf.Mysql.NewSession()
	if search != "" {
		sess.Where("question LIKE ?", "%"+search+"%")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("id DESC").FindAndCount(&questions)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}
	for _, question := range questions {
		question.CreatedTime = time.Unix(int64(question.CreatedAt), 0).Format("2006-01-02 15:04:05")
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: questions})
}

func QuestionDetail(c *gin.Context, id int) {
	var question model.GptQuestion
	_, err := conf.Mysql.Where("id = ?", id).Get(&question)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}
	question.CreatedTime = time.Unix(int64(question.CreatedAt), 0).Format("2006-01-02 15:04:05")
	common.ResOk(c, "ok", question)
}
