package common

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/templateMessage/request"
	"github.com/gin-gonic/gin"
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
