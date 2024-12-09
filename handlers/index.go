package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"local-audio-lib/constants"
	g "local-audio-lib/global"
	"local-audio-lib/types"
	"net/http"
)

func Index(c echo.Context) error {
	// 获取当前请求信息，方便拼接
	host := c.Request().Host

	// 读取缓存
	indexBytes, err := g.Rdb.Get(c.Request().Context(), constants.CacheKeyIndex).Bytes()
	if err != nil {
		g.L.Error("缓存读取失败", zap.Error(err))
		return c.String(http.StatusInternalServerError, "缓存读取失败")
	}

	// 格式化
	var privateIndex types.PrivateIndex
	if err = json.Unmarshal(indexBytes, &privateIndex); err != nil {
		g.L.Error("无法格式化索引，可能格式损坏", zap.String("index", string(indexBytes)), zap.Error(err))
		return c.String(http.StatusInternalServerError, "缓存解析失败")
	}

	// 映射为公开索引
	var publicIndex types.PublicIndex
	for fileHash, privateItem := range privateIndex {
		publicItem := types.PublicIndexItem{
			URL:    fmt.Sprintf("//%s/audio/%s%s", host, fileHash, privateItem.AudioExtension),
			Name:   privateItem.Name,
			Artist: privateItem.Artist,
			Album:  privateItem.Album,
		}

		if privateItem.HasCover {
			coverUrl := fmt.Sprintf("//%s/cover/%s", host, fileHash)
			publicItem.Cover = &coverUrl
		}

		publicIndex = append(publicIndex, publicItem)
	}

	// 以建立的公开索引响应
	return c.JSON(http.StatusOK, publicIndex)
}
