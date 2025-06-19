## 如何使用

配置文件模板为 config.yml.example ，请复制为 config.yml 以方便程序使用。您可以根据您的实际部署方案进行对应的调整。

默认只加载 .mp3 文件，如果要加载所有扩展名的文件，可以这样设置 `library.extensions` 项：

```yml
library:
  extensions:
    - "*"
```

以下说明针对使用 Docker Compose 部署，参考 docker-compose.yml 文件：

1. 将音频文件放置在 data 目录的 audio 目录下，可以使用子目录（例如将 mp3 文件放在 mp3 子目录下）
2. 启动程序，等待程序生成音频封面的缓存（默认配置下会生成在 data 目录的 cover 目录中
3. 可以通过预设的端口访问，例如 `127.0.0.1:1323` ，您可以自行设置反向代理

如果遇到问题，请将 `system.debug` 设置为 `true` 以收集较为详细的日志。
