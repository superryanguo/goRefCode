
PGO 第一个版本先支持的 pprof CPU，直接读取 pprof CPU profile 文件来完成优化。

有以下两种方式：

手动指定：Go 工具链在 go build 子命令增加了-pgo=<path>参数，用于显式指定用于 PGO 构建的 profile 文件位置。
自动指定：当 Go 工具链在主模块目录下找到 default.pgo 的配置文件时，将会自动启用 PGO。
