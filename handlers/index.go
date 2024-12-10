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
	rctx := c.Request().Context()

	// 检查公开索引是否已经缓存
	exist, err := g.Rdb.HExists(rctx, constants.CacheKeyIndexPublic, host).Result()
	if err != nil {
		// 缓存检查失败，忽略缓存继续处理
		g.L.Error("站点缓存检查失败", zap.Error(err))
	} else if exist {
		// 存在
		dataBytes, err := g.Rdb.HGet(rctx, constants.CacheKeyIndexPublic, host).Bytes()
		if err != nil {
			// 缓存读取失败，忽略缓存继续处理
			g.L.Error("站点缓存读取失败", zap.Error(err))
		} else {
			// 直接以二进制形式发送
			return c.Blob(http.StatusOK, "application/json", dataBytes)
		}
	}

	// 读取缓存
	indexBytes, err := g.Rdb.Get(rctx, constants.CacheKeyIndexPrivate).Bytes()
	if err != nil {
		g.L.Error("主索引缓存读取失败", zap.Error(err))
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
			URL:    fmt.Sprintf("//%s/audio/%s", host, fileHash),
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

	publicIndexBytes, err := json.Marshal(publicIndex)
	if err != nil {
		g.L.Error("站点缓存数据格式化失败", zap.Error(err))
		return c.String(http.StatusInternalServerError, "站点缓存数据格式化失败")
	}

	// 存入缓存
	g.Rdb.HSet(rctx, constants.CacheKeyIndexPublic, host, publicIndexBytes)

	// 直接以二进制形式发送
	return c.Blob(http.StatusOK, "application/json", publicIndexBytes)
}
