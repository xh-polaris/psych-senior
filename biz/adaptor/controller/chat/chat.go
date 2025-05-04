package chat

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xh-polaris/gopkg/util/log"
	"github.com/xh-polaris/psych-senior/biz/adaptor"
	"github.com/xh-polaris/psych-senior/biz/application/service"
)

// LongChat 开启一轮长对话
// @router /chat/ [GET]
func LongChat(ctx context.Context, c *app.RequestContext) {
	// 尝试升级协议, 并处理
	err := adaptor.UpgradeWs(ctx, c, service.ChatHandler)
	if err != nil {
		log.Error(err.Error())
	}
}
