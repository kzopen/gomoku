# web/lib — 官方 starx JS 客户端（原样引入 + 一处必要补丁）

> 本目录两个文件来自 **starx 官方 JavaScript 客户端**（starx 为 Pitaya 的前身，二者网络协议兼容），
> 从上游仓库复制。除下述**一处必要补丁**外未做其他修改；升级时重新应用该补丁。

| 文件 | 作用 |
| --- | --- |
| `protocol.js` | pomelo/starx 二进制协议编解码（Package 分帧 + Message 消息头），暴露全局 `window.Protocol` |
| `starx-wsclient.js` | 连接/握手/心跳/重连封装，暴露全局 `window.starx`，提供 `init` / `request` / `notify` / `on` / `disconnect` API |

- 上游仓库：<https://github.com/nano-ecosystem/starx-client-websocket>
- 许可：MIT License（Copyright (c) 2012-2013 NetEase, Inc. and other contributors，见文件头部注释）
- 参考实现：[pomelo-jsclient-websocket](https://github.com/pomelonode/pomelo-jsclient-websocket)
- 注意：`topfreegames/pitaya-client-websocket` 链接已失效（404），官方维护的 JS 客户端即上述仓库。

## ⚠️ 必要补丁（protocol.js）

**问题**：Pitaya 服务端会主动下发心跳探测包（`type=3`、`bodyLen=0`）。上游 `Package.decode`
对 `bodyLen==0` 的包会执行 `copyArray(null, ...)` 直接抛错，导致客户端收到第一个心跳即崩溃
（浏览器/Node 均复现）。

**补丁**（`Package.decode` 内）：`if (body) { copyArray(body, 0, bytes, offset, length); }`
——空 body 包不再复制，`body` 保持 `null`（心跳处理逻辑不读 body，无影响）。

## 序列化说明

该客户端**默认使用 JSON 序列化**（`JSON.stringify` + UTF-8），开箱即用；如需 protobuf，
可在 `starx.init({ encode, decode })` 中注入基于 `protobufjs` 的编解码函数（见《开发文档》§7.5）。
服务端对应切换 Pitaya 的 `pitaya.serializertype` 配置（`1`=JSON 默认 / `2`=protobuf），详见《开发文档》§7.5。
