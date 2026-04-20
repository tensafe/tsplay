# Lesson 147 file-srv Dev vs Release

`file-srv` 至少有两种心智模型：

1. 开发态：`-serve-root .`
2. 发布态：不传 `-serve-root`，直接服务二进制内置资源

开发态更适合：

1. 边改 demo 边验证
2. 本地文件优先，内置资源兜底

发布态更适合：

1. 只发一个二进制
2. 不假设用户拿到了源码目录
