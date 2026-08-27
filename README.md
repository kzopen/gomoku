# 五子棋（Gomoku）在线对战

基于 [Pitaya](https://github.com/topfreegames/pitaya) 游戏服务器框架实现的五子棋在线对战服务端，支持账号登录、在线匹配、实时对局、断线重连、AI 接管与对局回放。附带一个纯浏览器联调客户端（`web/`），开箱即用。

## 功能特性

- **账号体系**：登录即注册（首次登录自动建档），密码 bcrypt 哈希存储，登录态基于 token + Redis session
- **在线匹配**：Redis 队列撮合，双人自动成局，先手（黑方）随机一方
- **实时对局**：服务器权威判定——15×15 棋盘、五连胜负、制胜连线下发前端高亮
- **回合计时**：服务器权威计时，每步限时，超时自动判负；剩余时间 ≤5s 时加速推送（200ms 间隔）
- **断线重连**：`room.reconnect` 全量状态恢复（棋盘、计时、回合），对方同步收到回场通知
- **AI 接管**：玩家超时/掉线宽限期内未回来，AI 按配置等级（1 简单 / 2 普通 / 3 困难）接管落子
- **在线统计**：在线人数实时统计（Redis ZSET + 心跳续期），前端 30s 心跳刷新
- **对局入库**：终局写入 `game_record`（含完整走子序列 JSON，可用于回放）
- **战绩 / ELO**：胜负平、逃跑次数与 ELO 积分（初始 1200），登录即返回

## 技术栈

| 组件 | 说明 |
| --- | --- |
| Go 1.26.5 | 服务端语言 |
| [Pitaya v2](https://github.com/topfreegames/pitaya) | 游戏服务器框架（WebSocket 接入、组件路由、Session 管理） |
| MySQL 5.7+ | 账号、战绩、对局记录（GORM + [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)） |
| Redis | 登录态、在线状态、匹配队列、ELO 排行榜（[go-redis/v9](https://github.com/redis/go-redis)） |
| Viper | 配置加载（YAML） |
| Logrus | 结构化日志 |
| starx-wsclient.js | 浏览器端 WebSocket 客户端（`web/lib/`） |

## 目录结构

```
gomoku/
├── main.go                  # 入口：加载配置、连接 MySQL/Redis、构建 Pitaya、注册组件
├── config/
│   └── config.yaml          # 服务端配置（app / mysql / redis / game）
├── internal/
│   ├── app/
│   │   ├── auth/            # 登录/注册（auth.login）
│   │   ├── match/           # 匹配队列与撮合（match.join / match.cancel）
│   │   ├── room/            # 对局接入（room.place / room.reconnect / room.leave）
│   │   └── rank/            # 在线人数与心跳（rank.online_count / rank.ping）
│   ├── common/              # 棋盘规则（15×15、五连判定）、常量、错误码
│   ├── model/               # 配置结构、消息结构、MySQL/Redis 访问、Redis key 设计
│   ├── protocol/            # 协议定义（预留）
│   └── service/             # 房间管理器 + 对局房间（计时、AI 接管、终局入库）
├── web/                     # 浏览器联调客户端（纯 HTML/JS，无构建）
│   ├── index.html
│   └── lib/                 # starx-wsclient.js、protocol.js
└── docs/                    # 协议 proto、数据库脚本、设计规范
    ├── gomoku.proto
    ├── schema.sql           # 数据库初始化脚本（含开发种子账号）
    ├── gomoku.sql           # 数据库导出备份（含数据）
    └── *.md                 # 设计规范文档
```

## 快速开始

### 环境依赖

- Go 1.26.5+
- MySQL 5.7+（本地或远程均可）
- Redis 6+

### 1. 初始化数据库

```sql
-- 需先确保 MySQL 可用，然后执行：
mysql -uroot -p < docs/schema.sql
```

脚本会创建 `gomoku` 库及 `user` / `player_stats` / `game_record` 三张表，并插入开发种子账号（密码均为 `123456`）：

| 用户名 | 昵称 |
| --- | --- |
| player1 | 玩家一 |
| player2 | 玩家二 |
| ai_test | AI 测试 |

### 2. 修改配置

编辑 `config/config.yaml`，按实际环境调整：

```yaml
mysql:
  dsn: "root:123456@tcp(127.0.0.1:3306)/gomoku?charset=utf8mb4&parseTime=true&loc=Local"
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0
```

### 3. 启动服务端

```bash
go run main.go -c config/config.yaml
```

默认监听 `ws://127.0.0.1:3250`（WebSocket），启动成功日志：

```
starting gomoku server: type=room frontend=true port=3250 mode=standalone
```

### 4. 打开客户端联调

直接用浏览器打开 `web/index.html`，用两个窗口分别登录 `player1` / `player2`（密码 `123456`），在「对战大厅」点击匹配即可开始对局。

## 协议说明

### 路由（单机两段式）

| 路由 | 方向 | 说明 |
| --- | --- | --- |
| `auth.login` | 请求 | 登录（账号不存在则自动注册），返回 uid / token / 战绩 |
| `match.join` | 请求 | 进入匹配队列 |
| `match.cancel` | 请求 | 取消匹配 |
| `room.place` | 请求 | 落子（x, y） |
| `room.reconnect` | 请求 | 断线重连，恢复全量对局状态 |
| `room.leave` | 请求 | 主动离开（判负） |
| `rank.online_count` | 请求 | 查询在线人数 |
| `rank.ping` | 通知 | 心跳续期（前端 30s 一次） |
| `room.onMatchSuccess` | 推送 | 匹配成功，通知双方入房 |
| `room.onTurn` / `room.onTimer` | 推送 | 回合切换 / 倒计时 |
| `room.onGameOver` | 推送 | 终局（胜负、原因、走子回放、制胜连线） |
| `room.onPlayerBack` | 推送 | 对方重连回场通知 |

完整消息结构见 `internal/model/message.go` 与 `docs/gomoku.proto`。

### 终局原因

| 值 | 含义 |
| --- | --- |
| 1 | 五连 |
| 2 | 回合超时 |
| 3 | 离开（主动离开 / 掉线） |
| 4 | 和棋（预留） |

## 配置说明（config/config.yaml）

| 段 | 关键项 | 说明 |
| --- | --- | --- |
| app | `port` | WebSocket 监听端口（默认 3250） |
| app | `server_type` / `frontend` / `cluster` | 单机开发固定 `room` + `frontend=true` + `cluster=false` |
| app | `log_level` | debug / info / warn / error |
| mysql | `dsn` / `max_open` / `max_idle` | MySQL 连接串与连接池 |
| redis | `addr` / `password` / `db` | Redis 连接 |
| game | `board_size` | 棋盘大小（15） |
| game | `step_time_limit` | 每步限时（秒，默认 30） |
| game | `ai_level` | AI 难度 1 简单 / 2 普通 / 3 困难 |
| game | `forbid_rules` | 禁手规则（v1 关闭，预留） |

> 集群模式：生产环境将 `cluster` 置 `true`，按 `gate` / `room` / `rank` 分别部署，客户端使用三段式路由（如 `gate.auth.login`），详见 `docs/` 内设计规范。

## 文档

- `docs/gomoku.proto` — 协议定义
- `docs/schema.sql` — 数据库初始化脚本（含种子数据）
- `docs/gomoku.sql` — 数据库完整导出备份
- `docs/Go功能通用设计规范.md`、`docs/mysql设计规范.md` — 设计规范
- `开发文档.md` — 项目开发文档（协议、架构、演进记录）

## License

内部学习项目，暂未指定开源协议。
