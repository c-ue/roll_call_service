package pageHandle

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"html/template"
	"roll_call_service/server/config"
	"roll_call_service/server/logger"
	"runtime"
	"strconv"
)

func Error(ctx *fasthttp.RequestCtx, serverConf config.Config) {
	var ConnID = strconv.FormatUint(ctx.ConnID(), 10)
	var log *zap.Logger = logger.Console()

	// -------------------------------------------------------
	// 处理 HTTP 响应数据
	// -------------------------------------------------------
	// HTTP header 构造
	ctx.Response.Header.SetStatusCode(200)
	ctx.Response.Header.SetConnectionClose() // 关闭本次连接, 这就是短连接 HTTP
	ctx.Response.Header.SetBytesKV([]byte("Content-Type"), []byte("text/html; charset=utf8"))
	ctx.Response.Header.SetBytesKV([]byte("TransactionID"), []byte(ConnID))
	// HTTP payload 设置
	// 这里 HTTP payload 是 []byte
	//ctx.Response.SetBody(payload.Bytes())

	// -------------------------------------------------------
	// 处理逻辑开始
	// -------------------------------------------------------
	templateFileName := "template/Error.tmpl"
	t := template.Must(template.ParseFiles(templateFileName))
	if err := t.Execute(ctx, string(ConnID)); err != nil {
		_, file, _, _ := runtime.Caller(1)
		log.Debug("---------------- Template Produce Error [" + file + ";" + err.Error() + "]-------------")
		return
	}

}
