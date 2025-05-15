package common

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/templateMessage/request"
	"github.com/gin-gonic/gin"
)

func SendTemplateMessage(ctx *gin.Context, officialAccountApp officialAccount.OfficialAccount, toUser, templateId string, data *power.HashMap) error {
	_, err := officialAccountApp.TemplateMessage.Send(ctx, &request.RequestTemlateMessage{
		ToUser:     toUser,
		TemplateID: templateId,
		Data:       data,
	})
	if err != nil {
		return err
	}
	return nil
}
