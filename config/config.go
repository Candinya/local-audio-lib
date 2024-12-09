package config

type Config struct {
	System  SystemConfig  `yaml:"system"`
	Library LibraryConfig `yaml:"library"`
}

type SystemConfig struct {
	Debug     bool   `yaml:"debug"`      // 调试模式
	Listen    string `yaml:"listen"`     // web server 监听的端口
	RedisConn string `yaml:"redis_conn"` // redis 连接字符串
}

type LibraryConfig struct {
	Path struct {
		Audio string `yaml:"audio"`
		Cover string `yaml:"cover"`
	} `yaml:"path"` // 内容库的路径
	Extensions []string `yaml:"extensions"` // 需要索引的后缀名
}
