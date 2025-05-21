package model

type CommonWarningStruct struct {
	MsgType string             `json:"msg_type"`
	Text    TextWarningContent `json:"text"`
	//MarkDown string `json:"mark_down"`
}

type TextWarningContent struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}
