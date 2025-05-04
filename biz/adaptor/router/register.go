package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/xh-polaris/psych-senior/biz/adaptor/controller/chat"
	"github.com/xh-polaris/psych-senior/biz/adaptor/controller/voice"
)

func Register(r *server.Hertz) {
	root := r.Group("/", _rootMw()...)
	{
		_chat := root.Group("/chat")
		_chat.GET("/", append(_longchatMw(), chat.LongChat)...)
		_chat.GET("/history/list", chat.ListHistory)
	}
	{
		_voice := root.Group("/voice")
		_voice.GET("/asr", append(_asrMw(), voice.Asr)...)
	}
}
