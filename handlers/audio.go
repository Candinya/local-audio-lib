package handlers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"local-audio-lib/constants"
	g "local-audio-lib/global"
	"net/http"
)

func Audio(c echo.Context) error {
	id := c.Param("id")

	audioPath, err := g.Rdb.HGet(c.Request().Context(), constants.CacheKeyAudioFile, id).Result()
	if err != nil {
		g.L.Error("音频路径读取失败", zap.String("id", id), zap.Error(err))
		return c.String(http.StatusInternalServerError, "缓存读取失败")
	}

	// 设置超大缓存（一周）
	c.Response().Header().Add("Cache-Control", "max-age=604800")

	// 返回文件
	return c.File(audioPath)
}
