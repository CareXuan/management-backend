package controller

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/gin-gonic/gin"
	"io"
	"management-backend/common"
	"management-backend/conf"
	"net/http"
	"sort"
	"strings"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.ResError(c, "文件上传失败")
		return
	}
	fileMaxSize := 4 << 20 //4M
	if int(file.Size) > fileMaxSize {
		common.ResError(c, "文件不允许大小于32KB")
		return
	}
	reader, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	b := make([]byte, 512)
	reader.Read(b)
	contentType := http.DetectContentType(b)
	if contentType != "image/jpeg" && contentType != "image/png" {
		common.ResError(c, "只能上传jpeg/png")
		return
	}
	fileName := "YQ" + common.GetOneNewCard(18) + ".png"
	dst := "./upload/" + fileName
	c.SaveUploadedFile(file, dst)
	common.ResOk(c, "ok", conf.Conf.Upload.Url+"/"+fileName)
}

func WechatCheck(c *gin.Context) {
	fmt.Println(conf.Conf.Wechat.Warning.TestWarning)
	err := common.SendTemplateMessage(c, *conf.WechatApp, "oWUHY7STySg4-oE-ZufbgwmwyBeY", conf.Conf.Wechat.Warning.TestWarning, &power.HashMap{
		"thing1": &power.HashMap{
			"value": "测试机",
		},
		"character_string3": &power.HashMap{
			"value": "220V",
		},
	})
	if err != nil {
		fmt.Println(err)
		common.ResError(c, "err")
		return
	}
	common.ResOk(c, "ok", nil)
}

type WechatMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"` // 文本消息的内容
	Event        string   `xml:"Event"`   // 事件类型
}

func WechatHandler(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	// 1. 验证签名
	if !validateSignature(signature, timestamp, nonce) {
		c.String(http.StatusForbidden, "signature invalid")
		return
	}

	// 2. 如果是 GET 请求，说明是首次验证
	if c.Request.Method == http.MethodGet {
		c.String(http.StatusOK, echostr)
		return
	}

	// 3. 如果是 POST 请求，说明是消息推送
	var msg WechatMessage
	if err := c.ShouldBindXML(&msg); err != nil {
		c.String(http.StatusBadRequest, "invalid xml")
		return
	}

	// 4. 根据消息类型或事件类型进行处理
	if msg.MsgType == "event" && msg.Event == "subscribe" {
		// 用户关注事件
		reply := generateTextReply(msg, "欢迎关注我们公众号！")
		fmt.Println(reply)
		c.XML(http.StatusOK, reply)
		return
	}

	if msg.MsgType == "text" {
		// 用户发送文本消息
		reply := generateTextReply(msg, "你发送的是："+msg.Content)
		c.XML(http.StatusOK, reply)
		return
	}

	// 其他情况默认回复
	reply := generateTextReply(msg, "暂不支持该消息类型")
	c.XML(http.StatusOK, reply)
}

// 验证微信签名
func validateSignature(signature, timestamp, nonce string) bool {
	strs := []string{conf.Conf.Wechat.Token, timestamp, nonce}
	sort.Strings(strs)
	str := strings.Join(strs, "")
	h := sha1.New()
	io.WriteString(h, str)
	hashcode := fmt.Sprintf("%x", h.Sum(nil))
	return hashcode == signature
}

// 构造文本回复消息
func generateTextReply(msg WechatMessage, replyText string) WechatMessage {
	return WechatMessage{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   msg.CreateTime,
		MsgType:      "text",
		Content:      replyText,
	}
}
