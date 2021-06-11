package main

import (
	"log"
	"net/http"
	"promnotify/prometheusalert"
	"promnotify/webhook"
	yamlconfig "promnotify/yamlconfig"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	yamlconfig.CreateYamlFactory()
	//监听文件变化载配置文件 非文件链接
	yamlconfig.ConfigYml.ConfigFileChangeListen()
	// yamlconfig.ConfigYml.yamlConfig.viper.WatchConfig()

	// 监听告警消息
	go webhook.WebhookWorker()

	router := gin.Default()
	router.POST("/notify", notify)
	router.POST("/reload", reload)
	router.Run(":11122")
}

func notify(ctx *gin.Context) {
	// 获取消息类型
	notifyType := ctx.DefaultQuery("type", "")
	notifyAisle := ctx.DefaultQuery("aisle", "")
	notifyGroup := ctx.DefaultQuery("group", "")
	apiKey := ctx.DefaultQuery("api_key", "")
	if apiKey != yamlconfig.ConfigYml.GetString("api_key") {
		ctx.JSON(http.StatusForbidden, gin.H{
			"err": -1,
		})
		return
	}

	var p prometheusalert.Prometheus
	if err := ctx.ShouldBindJSON(&p); err != nil {
		log.Panicln(err)
	}
	webhook.Hierarchy(&p, notifyType, notifyAisle, notifyGroup)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

func reload(ctx *gin.Context) {
	apiKey := ctx.DefaultQuery("api_key", "")
	if apiKey == yamlconfig.ConfigYml.GetString("api_key") || ctx.ClientIP() == "127.0.0.1" {
		yamlconfig.ConfigYml.ReloadCache()
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}
	ctx.JSON(http.StatusForbidden, gin.H{
		"err": -1,
	})
}
