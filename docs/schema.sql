-- ============================================================
-- 五子棋数据库初始化脚本
-- 目标库: gomoku (远程 MySQL 47.97.42.179:3306)
-- 执行: mysql -h47.97.42.179 -uroot -p < schema.sql
-- ============================================================

CREATE DATABASE IF NOT EXISTS `gomoku`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `gomoku`;

-- ------------------------------------------------------------
-- 账号 / 角色基础数据
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT,
  `username`      VARCHAR(64)  NOT NULL COMMENT '登录名，唯一',
  `password_hash` VARCHAR(128) NOT NULL COMMENT 'bcrypt 哈希',
  `nickname`      VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar`        VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像 URL',
  `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE = InnoDB COMMENT ='用户账号';

-- ------------------------------------------------------------
-- 战绩与积分
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `player_stats` (
  `user_id`    BIGINT      NOT NULL,
  `wins`       INT         NOT NULL DEFAULT 0,
  `losses`     INT         NOT NULL DEFAULT 0,
  `draws`      INT         NOT NULL DEFAULT 0,
  `runaways`   INT         NOT NULL DEFAULT 0 COMMENT '逃跑(掉线被AI接管)次数',
  `elo`        INT         NOT NULL DEFAULT 1200 COMMENT 'ELO 积分，初始 1200',
  `created_at` DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  KEY `idx_elo` (`elo`)
) ENGINE = InnoDB COMMENT ='玩家战绩';

-- ------------------------------------------------------------
-- 对局记录（含回放）
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_record` (
  `id`                BIGINT   NOT NULL AUTO_INCREMENT,
  `room_id`           BIGINT   NOT NULL COMMENT '房间号',
  `black_uid`         BIGINT   NOT NULL,
  `white_uid`         BIGINT   NOT NULL,
  `winner`            TINYINT  NOT NULL COMMENT '0=黑胜 1=白胜 2=和棋',
  `end_reason`        TINYINT  NOT NULL COMMENT '1=五连 2=超时 3=离开 4=和棋',
  `ai_takeover_seat`  TINYINT  NOT NULL DEFAULT -1 COMMENT '被AI接管的座位，-1=无',
  `moves`             JSON     NOT NULL COMMENT '走子序列 [{"seat":0,"x":7,"y":7,"ts":...}]',
  `started_at`        DATETIME NOT NULL,
  `ended_at`          DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_room` (`room_id`),
  KEY `idx_black` (`black_uid`),
  KEY `idx_white` (`white_uid`),
  KEY `idx_started` (`started_at`)
) ENGINE = InnoDB COMMENT ='对局记录';

-- ------------------------------------------------------------
-- 种子数据（开发用，便于登录调试）
-- 密码均为 123456 (bcrypt 预生成)
-- ------------------------------------------------------------
INSERT IGNORE INTO `user` (`id`, `username`, `password_hash`, `nickname`, `avatar`) VALUES
  (1, 'player1', '$2a$10$k9ugxQDekLVWJymnBEHfFenic95lgjcJhYk01Ay.TZMUfUtkd2oPO', '玩家一', ''),
  (2, 'player2', '$2a$10$k9ugxQDekLVWJymnBEHfFenic95lgjcJhYk01Ay.TZMUfUtkd2oPO', '玩家二', ''),
  (3, 'ai_test', '$2a$10$k9ugxQDekLVWJymnBEHfFenic95lgjcJhYk01Ay.TZMUfUtkd2oPO', 'AI 测试', '');

INSERT IGNORE INTO `player_stats` (`user_id`, `wins`, `losses`, `draws`, `runaways`, `elo`) VALUES
  (1, 0, 0, 0, 0, 1200),
  (2, 0, 0, 0, 0, 1200),
  (3, 0, 0, 0, 0, 1200);
