package inits

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"local-audio-lib/config"
	"local-audio-lib/constants"
	g "local-audio-lib/global"
	"local-audio-lib/types"
	"local-audio-lib/utils"
	"os"
	"path"
	"path/filepath"

	"github.com/bogem/id3v2"
	"go.uber.org/zap"
)

// processFile : 对每一个文件计算摘要得到 ID ，根据 ID 检查旧索引，有则直接迁移入新索引，无则创建新的数据
func processFile(fileName string, fileExt string, coverLibPath string, oldIndex *types.PrivateIndex, newIndex *types.PrivateIndex) error {
	// 打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("无法打开文件 %s: %v", fileName, err)
	}
	defer f.Close()

	// 计算摘要
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("无法计算摘要 %s: %v", fileName, err)
	}
	fileHash := hex.EncodeToString(h.Sum(nil))

	// 设置文件路径
	g.Rdb.HSet(context.Background(), constants.CacheKeyAudioFile, fileHash, fileName)

	// 检查匹配
	if val, exist := (*oldIndex)[fileHash]; exist {
		// 继承
		(*newIndex)[fileHash] = val
		delete(*oldIndex, fileHash) // 删除以确保在之后的清理步骤中不会被影响

		// 对这个文件处理完了，可以返回
		return nil
	} // else 不存在，需要新增

	// 归位读取指针，用于读取指定长度
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("指针归位失败 %s: %v", fileName, err)
	}

	// 读取 id3 标签
	tag, err := id3v2.ParseReader(f, id3v2.Options{
		Parse: true,
	})
	if err != nil {
		return fmt.Errorf("文件 id3 标签读取失败 %s: %v", fileName, err)
	}

	// 读取封面
	hasCover := false
	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	for _, f := range pictures {
		pic, ok := f.(id3v2.PictureFrame)
		if ok {
			// 找到图片了，记录图片
			hasCover = true

			coverPath := path.Join(coverLibPath, fileHash)
			err = os.WriteFile(coverPath, pic.Picture, 0644)
			if err != nil {
				return fmt.Errorf("封面保存失败 %s: %v", fileName, err)
			}

			break // 找到了，就不用再处理其他文件了
		}
	}

	// 创建索引
	indexItem := types.PrivateIndexItem{
		HasCover: hasCover,
	}
	if tagTitle := tag.Title(); tagTitle != "" {
		indexItem.Name = tagTitle
	} else {
		// 使用文件名作为 fallback，避免为空
		baseFileName := filepath.Base(fileName)
		indexItem.Name = baseFileName[:len(baseFileName)-len(fileExt)]
	}
	if tagArtist := tag.Artist(); tagArtist != "" {
		indexItem.Artist = &tagArtist
	}
	if tagAlbum := tag.Album(); tagAlbum != "" {
		indexItem.Album = &tagAlbum
	}

	(*newIndex)[fileHash] = indexItem

	return nil
}

// Library : 扫描路径下的所有音频文件，生成索引
func Library(cfg config.LibraryConfig) error {
	// 检查封面目录
	if _, err := os.Stat(cfg.Path.Cover); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(cfg.Path.Cover, os.ModePerm)
			if err != nil {
				return fmt.Errorf("封面路径不存在且无法被创建: %v", err)
			}
		} else {
			return fmt.Errorf("无法确认封面路径是否存在: %v", err)
		}
	} // else 已经存在

	// 检查是否有旧的索引
	var oldIndex types.PrivateIndex

	if exist, err := g.Rdb.Exists(context.Background(), constants.CacheKeyIndexPrivate).Result(); err != nil {
		return fmt.Errorf("无法检测是否存在旧索引: %v", err)
	} else if exist > 0 {
		// 存在旧索引，尝试读取
		if oldIndexBytes, err := g.Rdb.Get(context.Background(), constants.CacheKeyIndexPrivate).Bytes(); err != nil {
			return fmt.Errorf("无法读取旧索引: %v", err)
		} else if err = json.Unmarshal(oldIndexBytes, &oldIndex); err != nil {
			return fmt.Errorf("无法格式化旧索引，可能格式损坏: %v", err)
		}
	} else {
		// 不存在旧索引，初始化为空
		oldIndex = make(types.PrivateIndex)
	}

	// 准备新的索引
	newIndex := make(types.PrivateIndex)

	// 读取文件列表
	allFiles, err := utils.ListFileRecursive(cfg.Path.Audio)
	if err != nil {
		return fmt.Errorf("无法列出目录: %v", err)
	}

	// 检查是否为处理所有文件
	processAnyExtension := utils.ArrayHas(cfg.Extensions, "*")

	// 处理文件
	for _, fileName := range allFiles {
		fileExt := filepath.Ext(fileName)
		if processAnyExtension || utils.ArrayHas(cfg.Extensions, fileExt) {
			// 处理文件
			if err = processFile(fileName, fileExt, cfg.Path.Cover, &oldIndex, &newIndex); err != nil {
				// return fmt.Errorf("无法处理文件: %v", err)
				g.L.Error("无法处理文件", zap.String("fileName", fileName), zap.Error(err))
				// 忽略这个文件
			}
		} // else 忽略
	}

	// 存储新索引
	newIndexBytes, err := json.Marshal(newIndex)
	if err != nil {
		return fmt.Errorf("新索引格式化失败: %v", err)
	}

	g.Rdb.Set(context.Background(), constants.CacheKeyIndexPrivate, newIndexBytes, 0)

	// 删除旧索引中剩余的记录
	for oldFileHash, oldItem := range oldIndex {
		g.Rdb.HDel(context.Background(), constants.CacheKeyAudioFile, oldFileHash)
		if oldItem.HasCover {
			_ = os.Remove(path.Join(cfg.Path.Cover, oldFileHash)) // 忽略错误
		}
	}

	// 清空旧的站点索引
	g.Rdb.Del(context.Background(), constants.CacheKeyIndexPublic)

	// 处理完成，返回
	return nil
}
