package common

import (
	"encoding/json"
	"fmt"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/templateMessage/request"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"my-gpt-server/conf"
	"strings"
)

func SendTemplateMessage(ctx *gin.Context, officialAccountApp officialAccount.OfficialAccount, toUser, templateId string) error {
	_, err := officialAccountApp.TemplateMessage.Send(ctx, &request.RequestTemlateMessage{
		ToUser:     toUser,
		TemplateID: templateId,
		URL:        "https://www.artisan-cloud.com/",
		Data: &power.HashMap{
			"first": &power.HashMap{
				"value": "恭喜你购买成功！",
				"color": "#173177",
			},
			"DateTime": &power.HashMap{
				"value": "2022-3-5 16:22",
				"color": "#173177",
			},
			"PayAmount": &power.HashMap{
				"value": "59.8元",
				"color": "#173177",
			},
			"Location": &power.HashMap{
				"value": "上海市长宁区",
				"color": "#173177",
			},
			"remark": &power.HashMap{
				"value": "欢迎再次购买！",
				"color": "#173177",
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateClient() (*dysmsapi20170525.Client, error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(conf.Conf.Sms.AccessKeyId),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(conf.Conf.Sms.AccessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result := &dysmsapi20170525.Client{}
	result, err := dysmsapi20170525.NewClient(config)
	return result, err
}

func SendSms(phoneNumber, templateCode, signName, templateParam string) error {
	client, err := CreateClient()
	if err != nil {
		return err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		PhoneNumbers:  tea.String(phoneNumber),
		TemplateParam: tea.String(templateParam),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, err = util.AssertAsString(error.Message)
		if err != nil {
			return err
		}
	}
	return err
}
