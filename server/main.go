package main

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"roll_call_service/server/config"
	"roll_call_service/server/logger"
	"roll_call_service/server/pageHandle"
)

func main() {
	var log *zap.Logger = logger.Console()
	var conf config.Config = config.ReadConfig("config/server_config.toml")
	var server_address = conf.SERVER.IP + ":" + conf.SERVER.PORT

	// -------------------------------------------------------
	//  fasthttp 的 handler 处理函数
	// -------------------------------------------------------
	var requestHandler = func(ctx *fasthttp.RequestCtx) {
		if conf.SERVER.DEBUG {
			debugRequestHandle(ctx)
		}

		// 取出 web 访问的 URL/URI
		//var ConnID = strconv.FormatUint(ctx.ConnID(),10)
		var uriPath = string(ctx.Path())
		switch {
		// 如果访问的 URI 路由是 /uri 开头 , 则进行下面这个响应
		case uriPath == "/index.html":
			{
				pageHandle.Index(ctx, conf, log)
				return
			}

			// 访问路踊不是 /uri 的其他响应
		default:
			{
				pageHandle.Error(ctx, conf, log)
				return
			}
		}

		return

	}

	// -------------------------------------------------------
	// 创建 fasthttp 服务器
	// -------------------------------------------------------
	// Create custom server.
	s := &fasthttp.Server{
		Handler: requestHandler,     // 注意这里
		Name:    "Roll Call Server", // 服务器名称
	}
	// -------------------------------------------------------
	// 运行服务端程序
	// -------------------------------------------------------
	log.Debug("------------------ Roll Call Server 服务器尝试启动------ ")

	if err := s.ListenAndServe(server_address); err != nil {
		log.Fatal("ERROR in ListenAndServe", zap.Error(err))
	}
}
