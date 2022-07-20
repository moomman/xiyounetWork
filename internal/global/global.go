package global

import (
	"ttms/internal/model/config"
	pay "ttms/internal/pkg/alipay"
	"ttms/internal/pkg/app"
	"ttms/internal/pkg/goroutine/work"
	"ttms/internal/pkg/logger"
	"ttms/internal/pkg/mangerFunc"
	"ttms/internal/pkg/snowflake"
	"ttms/internal/pkg/token"
)

var (
	Logger       *logger.Log           // 日志
	Settings     config.All            // 全局配置
	Maker        token.Maker           // 操作token
	Snowflake    *snowflake.Snowflake  // 生成ID
	RootDir      string                // 项目跟路径
	Page         *app.Page             // 分页
	Worker       *work.Worker          // 工作池
	MangerFunc   mangerFunc.MangerFunc // manger函数
	AliPayClient *pay.Client           // 支付宝
)
