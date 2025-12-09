package types

// PrivateIndexItem : 私有索引，用于系统内部存储
type PrivateIndexItem struct {
	// 元信息
	Name   string  `json:"name"`
	Artist *string `json:"artist,omitempty"`
	Album  *string `json:"album,omitempty"`

	// 附加信息
	HasCover        bool   `json:"has_cover"`
	CoverThemeColor string `json:"cover_theme,omitempty"` // 封面的主题颜色
	//HasLyrics bool `json:"has_lyrics"`
}

type PrivateIndex = map[string]PrivateIndexItem

// PublicIndexItem : 公开索引，用于暴露给客户端使用
type PublicIndexItem struct {
	URL string `json:"url"` // 音频文件， /audio/id.ext

	Name   string  `json:"name"`             // 音频名
	Artist *string `json:"artist,omitempty"` // 艺术家
	Album  *string `json:"album,omitempty"`  // 专辑名

	Cover *string `json:"cover,omitempty"` // 封面 /cover/id.ext
	Theme *string `json:"theme,omitempty"` // 提取封面主色作为主题颜色
	//Lyrics *string `json:"lrc,omitempty"`
}

type PublicIndex = []PublicIndexItem
