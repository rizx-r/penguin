package ws

// go get github.com/mitchellh/mapstructure
import "penguin/pkg/constants"

type (
	Msg struct {
		constants.MType `mapstructure:",mType"` // mapstructure 把 map[string]interface{} 这样的动态数据 自动解码 到 Go 的 struct 里。
		Content         string                  `mapstructure:"content"`
	}

	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:",chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		Msg                `mapstructure:"msg"`
		SendTime           int64 `mapstructure:"sendTime"`
	}

	// 消息推送
	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		constants.MType    `mapstructure:"mType"`
		Content            string `mapstructure:"content"`
	}
)
