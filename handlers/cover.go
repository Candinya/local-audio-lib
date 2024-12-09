package handlers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"local-audio-lib/constants"
	g "local-audio-lib/global"
	"net/http"
)

func Cover(c echo.Context) error {
	id := c.Param("id")

	imageBinary, err := g.Rdb.HGet(c.Request().Context(), constants.CacheKeyCoverContent, id).Bytes()
	if err != nil {
		g.L.Error("封面内容读取失败", zap.String("id", id), zap.Error(err))
		return c.String(http.StatusInternalServerError, "缓存读取失败")
	}

	mimeType, err := g.Rdb.HGet(c.Request().Context(), constants.CacheKeyCoverMimeType, id).Result()
	if err != nil {
		g.L.Error("封面格式读取失败", zap.String("id", id), zap.Error(err))

		// 读取失败就运行时解析
		mimeType = http.DetectContentType(imageBinary)
	}

	// 返回二进制信息
	return c.Blob(http.StatusOK, mimeType, imageBinary)
}
