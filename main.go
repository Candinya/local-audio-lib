package main

import (
	"go.uber.org/zap"
	g "local-audio-lib/global"
	"local-audio-lib/inits"
	"log"
)

func main() {
	// 读取配置
	cfg, err := inits.Config()
	if err != nil {
		log.Fatalf("配置文件读取失败: %v", err)
	}

	// 初始化 zap
	if g.L, err = inits.Logger(cfg.System.Debug); err != nil {
		log.Fatalf("zap 初始化失败: %v", err)
	}

	g.L.Debug("zap 初始化成功，切换为主日志系统")

	// 连接 redis
	if g.Rdb, err = inits.Redis(cfg.System.RedisConn); err != nil {
		g.L.Fatal("Redis 初始化失败", zap.Error(err))
	}

	// 创建资源库
	if err = inits.Library(cfg.Library); err != nil {
		g.L.Fatal("资源库初始化失败", zap.Error(err))
	}

	// 启动服务器
	if err = inits.WebServer(cfg.System.Listen); err != nil {
		g.L.Fatal("服务启动失败", zap.Error(err))
	}
}
