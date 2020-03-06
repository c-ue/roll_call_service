package main

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/url"
	"roll_call_service/server/logger"
	"strconv"
)

func debugRequestHandle(ctx *fasthttp.RequestCtx) {
	var ConnID = strconv.FormatUint(ctx.ConnID(), 10)
	var log *zap.Logger = logger.Console()
	var uri = ctx.Path()
	{
		// 取出 URI
		log.Debug("---------------- HTTP URI -------------")
		log.Debug(" HTTP 请求 URL 原始数据 > ", zap.String("request", ctx.String()))
		log.Debug("在 URI 中的路徑 > " + string(uri))
	}

	// 取出 web client 请求的 URL/URI 中的参数部分
	{
		log.Debug("---------------- HTTP URI 参数 -------------")
		var uri = ctx.URI().QueryString()
		log.Debug("在 URI 中的原始数据 > " + string(uri))
		log.Debug("---------------- HTTP URI 每一个键值对 -------------")
		ctx.URI().QueryArgs().VisitAll(func(key, value []byte) {
			log.Debug(ConnID, zap.String("key", url.QueryEscape(string(key))), zap.String("value", url.QueryEscape(string(value))))
		})
	}
	// -------------------------------------------------------
	// 注意对比一下, 下面的代码段, 与 web client  中几乎一样
	// -------------------------------------------------------
	{
		// 取出 web client 请求中的 HTTP header
		{
			log.Debug("---------------- HTTP header 每一个键值对-------------")
			ctx.Request.Header.VisitAll(func(key, value []byte) {
				log.Debug(ConnID, zap.String("key", url.QueryEscape(string(key))), zap.String("value", url.QueryEscape(string(value))))
			})

		}
		// 取出 web client 请求中的 HTTP payload
		{
			log.Debug("---------------- HTTP payload -------------")
			log.Debug(ConnID, zap.String("http payload", url.QueryEscape(string(ctx.Request.Body()))))
		}
	}
}
