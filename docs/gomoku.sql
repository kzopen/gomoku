/*
 Navicat Premium Data Transfer

 Source Server         : 99阿里云ROOT
 Source Server Type    : MySQL
 Source Server Version : 50744
 Source Host           : 47.97.42.179:3306
 Source Schema         : gomoku

 Target Server Type    : MySQL
 Target Server Version : 50744
 File Encoding         : 65001

 Date: 27/08/2026 14:47:35
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for game_record
-- ----------------------------
DROP TABLE IF EXISTS `game_record`;
CREATE TABLE `game_record`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `room_id` bigint(20) NOT NULL COMMENT '房间号',
  `black_uid` bigint(20) NOT NULL,
  `white_uid` bigint(20) NOT NULL,
  `winner` tinyint(4) NOT NULL COMMENT '0=黑胜 1=白胜 2=和棋',
  `end_reason` tinyint(4) NOT NULL COMMENT '1=五连 2=超时 3=离开 4=和棋',
  `ai_takeover_seat` tinyint(4) NOT NULL DEFAULT -1 COMMENT '被AI接管的座位，-1=无',
  `moves` json NOT NULL COMMENT '走子序列 [{\"seat\":0,\"x\":7,\"y\":7,\"ts\":...}]',
  `started_at` datetime NOT NULL,
  `ended_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_room`(`room_id`) USING BTREE,
  INDEX `idx_black`(`black_uid`) USING BTREE,
  INDEX `idx_white`(`white_uid`) USING BTREE,
  INDEX `idx_started`(`started_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '对局记录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for player_stats
-- ----------------------------
DROP TABLE IF EXISTS `player_stats`;
CREATE TABLE `player_stats`  (
  `user_id` bigint(20) NOT NULL,
  `wins` int(11) NOT NULL DEFAULT 0,
  `losses` int(11) NOT NULL DEFAULT 0,
  `draws` int(11) NOT NULL DEFAULT 0,
  `runaways` int(11) NOT NULL DEFAULT 0 COMMENT '逃跑(掉线被AI接管)次数',
  `elo` int(11) NOT NULL DEFAULT 1200 COMMENT 'ELO 积分，初始 1200',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`) USING BTREE,
  INDEX `idx_elo`(`elo`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '玩家战绩' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of player_stats
-- ----------------------------
INSERT INTO `player_stats` VALUES (1, 0, 0, 0, 0, 1200, '2026-08-14 16:48:09', '2026-08-14 16:48:09');
INSERT INTO `player_stats` VALUES (2, 0, 0, 0, 0, 1200, '2026-08-14 16:48:09', '2026-08-14 16:48:09');
INSERT INTO `player_stats` VALUES (3, 0, 0, 0, 0, 1200, '2026-08-14 16:48:09', '2026-08-14 16:48:09');
INSERT INTO `player_stats` VALUES (4, 0, 0, 0, 0, 1200, '2026-08-15 10:26:08', '2026-08-15 10:26:08');
INSERT INTO `player_stats` VALUES (5, 0, 0, 0, 0, 1200, '2026-08-15 10:27:25', '2026-08-15 10:27:25');
INSERT INTO `player_stats` VALUES (8, 0, 0, 0, 0, 1200, '2026-08-22 22:59:18', '2026-08-22 22:59:18');
INSERT INTO `player_stats` VALUES (9, 0, 0, 0, 0, 1200, '2026-08-22 22:59:19', '2026-08-22 22:59:19');
INSERT INTO `player_stats` VALUES (10, 0, 0, 0, 0, 1200, '2026-08-22 22:59:19', '2026-08-22 22:59:19');
INSERT INTO `player_stats` VALUES (11, 0, 0, 0, 0, 1200, '2026-08-22 23:00:09', '2026-08-22 23:00:09');
INSERT INTO `player_stats` VALUES (12, 0, 0, 0, 0, 1200, '2026-08-22 23:00:11', '2026-08-22 23:00:11');
INSERT INTO `player_stats` VALUES (13, 0, 0, 0, 0, 1200, '2026-08-22 23:00:12', '2026-08-22 23:00:12');
INSERT INTO `player_stats` VALUES (14, 0, 0, 0, 0, 1200, '2026-08-22 23:03:13', '2026-08-22 23:03:13');
INSERT INTO `player_stats` VALUES (15, 0, 0, 0, 0, 1200, '2026-08-22 23:03:14', '2026-08-22 23:03:14');
INSERT INTO `player_stats` VALUES (16, 0, 0, 0, 0, 1200, '2026-08-22 23:03:26', '2026-08-22 23:03:26');
INSERT INTO `player_stats` VALUES (17, 0, 0, 0, 0, 1200, '2026-08-22 23:03:27', '2026-08-22 23:03:27');
INSERT INTO `player_stats` VALUES (18, 0, 0, 0, 0, 1200, '2026-08-22 23:03:50', '2026-08-22 23:03:50');
INSERT INTO `player_stats` VALUES (19, 0, 0, 0, 0, 1200, '2026-08-22 23:03:51', '2026-08-22 23:03:51');
INSERT INTO `player_stats` VALUES (20, 0, 0, 0, 0, 1200, '2026-08-25 15:41:20', '2026-08-25 15:41:20');
INSERT INTO `player_stats` VALUES (21, 0, 0, 0, 0, 1200, '2026-08-25 15:41:21', '2026-08-25 15:41:21');
INSERT INTO `player_stats` VALUES (22, 0, 0, 0, 0, 1200, '2026-08-25 15:43:37', '2026-08-25 15:43:37');
INSERT INTO `player_stats` VALUES (23, 0, 0, 0, 0, 1200, '2026-08-25 15:43:38', '2026-08-25 15:43:38');
INSERT INTO `player_stats` VALUES (24, 0, 0, 0, 0, 1200, '2026-08-25 15:44:54', '2026-08-25 15:44:54');
INSERT INTO `player_stats` VALUES (25, 0, 0, 0, 0, 1200, '2026-08-25 15:44:55', '2026-08-25 15:44:55');
INSERT INTO `player_stats` VALUES (26, 0, 0, 0, 0, 1200, '2026-08-25 16:21:24', '2026-08-25 16:21:24');
INSERT INTO `player_stats` VALUES (27, 0, 0, 0, 0, 1200, '2026-08-25 16:21:24', '2026-08-25 16:21:24');
INSERT INTO `player_stats` VALUES (28, 0, 0, 0, 0, 1200, '2026-08-25 16:22:00', '2026-08-25 16:22:00');
INSERT INTO `player_stats` VALUES (29, 0, 0, 0, 0, 1200, '2026-08-25 16:22:00', '2026-08-25 16:22:00');
INSERT INTO `player_stats` VALUES (30, 0, 0, 0, 0, 1200, '2026-08-25 16:22:34', '2026-08-25 16:22:34');
INSERT INTO `player_stats` VALUES (31, 0, 0, 0, 0, 1200, '2026-08-25 16:22:34', '2026-08-25 16:22:34');
INSERT INTO `player_stats` VALUES (32, 0, 0, 0, 0, 1200, '2026-08-25 16:22:49', '2026-08-25 16:22:49');
INSERT INTO `player_stats` VALUES (33, 0, 0, 0, 0, 1200, '2026-08-25 16:22:50', '2026-08-25 16:22:50');
INSERT INTO `player_stats` VALUES (34, 0, 0, 0, 0, 1200, '2026-08-25 16:40:59', '2026-08-25 16:40:59');
INSERT INTO `player_stats` VALUES (35, 0, 0, 0, 0, 1200, '2026-08-25 16:40:59', '2026-08-25 16:40:59');
INSERT INTO `player_stats` VALUES (36, 0, 0, 0, 0, 1200, '2026-08-25 16:41:16', '2026-08-25 16:41:16');
INSERT INTO `player_stats` VALUES (37, 0, 0, 0, 0, 1200, '2026-08-25 16:41:35', '2026-08-25 16:41:35');
INSERT INTO `player_stats` VALUES (38, 0, 0, 0, 0, 1200, '2026-08-25 16:41:35', '2026-08-25 16:41:35');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录名，唯一',
  `password_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'bcrypt 哈希',
  `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像 URL',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 39 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户账号' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'player1', '$2a$10$k9ugxQDekLVWJymnBEHfFenic95lgjcJhYk01Ay.TZMUfUtkd2oPO', '玩家一', '', '2026-08-14 16:48:09', '2026-08-14 16:49:02');
INSERT INTO `user` VALUES (2, 'player2', '$2a$10$k9ugxQDekLVWJymnBEHfFenic95lgjcJhYk01Ay.TZMUfUtkd2oPO', '玩家二', '', '2026-08-14 16:48:09', '2026-08-14 16:49:02');
INSERT INTO `user` VALUES (3, 'ai_test', '$2a$10$k9ugxQDekLVWJymnBEHfFenic95lgjcJhYk01Ay.TZMUfUtkd2oPO', 'AI 测试', '', '2026-08-14 16:48:09', '2026-08-14 16:49:02');
INSERT INTO `user` VALUES (4, 'tester68426', '$2a$10$J/5xVoqJKHJMyg8GErk/O.F.zfFgK.s5msNKpsCCSN3BGR1sMTaBW', 'tester68426', '', '2026-08-15 10:26:08', '2026-08-15 10:26:08');
INSERT INTO `user` VALUES (5, 'tester45071', '$2a$10$Rh5se1xiYFXHfTs1lS0RT.gU0gITRB.Nu7G5bDZIl8chj44LFngUK', 'tester45071', '', '2026-08-15 10:27:25', '2026-08-15 10:27:25');
INSERT INTO `user` VALUES (6, 'player3', '$2a$10$5XwsipnvWu0gSixxQgqXEOUO.pjMpD.ifb5AdnTvJvpd6vsS/sD02', 'player3', '', '2026-08-22 14:46:38', '2026-08-22 14:46:38');
INSERT INTO `user` VALUES (7, 'test_1787410314146', '$2a$10$R4Tjd1h5wECNoHNFRq9ziOjNUln7uRU87kq0r8LPZzI1GiYM.hBhy', 'test_1787410314146', '', '2026-08-22 22:51:53', '2026-08-22 22:51:53');
INSERT INTO `user` VALUES (8, 'match_1787410759188_a', '$2a$10$fsg3lNkOAx27Ur/bPwwvuulnaZVyRFQdU.hvdOUDzJkQm5UtJD1JG', 'match_1787410759188_a', '', '2026-08-22 22:59:18', '2026-08-22 22:59:18');
INSERT INTO `user` VALUES (9, 'match_1787410759188_b', '$2a$10$rYrduusu38bS5GJuQj.d3OtWcvY4XDWg9n2ka8ZMHDUpUh37CObBi', 'match_1787410759188_b', '', '2026-08-22 22:59:19', '2026-08-22 22:59:19');
INSERT INTO `user` VALUES (10, 'match_1787410759188_c', '$2a$10$pQJNuXAh89hdzxesvR8fXOgyUJ1J1NHEnckHNHx59csu60dcnHMNy', 'match_1787410759188_c', '', '2026-08-22 22:59:19', '2026-08-22 22:59:19');
INSERT INTO `user` VALUES (11, 'm2_1787410809793', '$2a$10$j4UcloW2..hc0BHpl6mbQ.JQb.LKeSVS0RfExlmWII9c.Dwv1ShJS', 'm2_1787410809793', '', '2026-08-22 23:00:09', '2026-08-22 23:00:09');
INSERT INTO `user` VALUES (12, 'm2a_1787410812357', '$2a$10$mrqE4.ekN8scSP9S2qlxv.NmPgApJqVI..dlD5h2DXWyIbnTVkSPO', 'm2a_1787410812357', '', '2026-08-22 23:00:11', '2026-08-22 23:00:11');
INSERT INTO `user` VALUES (13, 'm2b_1787410813194', '$2a$10$Igstt2TXGpdH0xovji42leFqORP35yOlkC5Ejdq2whx7ngFIff/te', 'm2b_1787410813194', '', '2026-08-22 23:00:12', '2026-08-22 23:00:12');
INSERT INTO `user` VALUES (14, 'g_1787410994014_a', '$2a$10$WRs8MSvk5CVZ0boKW40GBO4sAV1OqJh23RV4ODx8GV.ooghjmBe8S', 'g_1787410994014_a', '', '2026-08-22 23:03:13', '2026-08-22 23:03:13');
INSERT INTO `user` VALUES (15, 'g_1787410994014_b', '$2a$10$hwf155sYqYmQDJHiTjDkn.2NL.5XmY732I3NQijXl81Kfa4rPvTxa', 'g_1787410994014_b', '', '2026-08-22 23:03:13', '2026-08-22 23:03:13');
INSERT INTO `user` VALUES (16, 'g_1787411007241_a', '$2a$10$PGh8vwYufrWO8VEagfAE4O3M2dxnrynjnIPfifh2J9CmT/.A1bCJ6', 'g_1787411007241_a', '', '2026-08-22 23:03:26', '2026-08-22 23:03:26');
INSERT INTO `user` VALUES (17, 'g_1787411007241_b', '$2a$10$vl.Pl88OLKRu9r20B0uT/.vCCRI5GzU076UbJ8UDzJ5Iw9fNTlFH.', 'g_1787411007241_b', '', '2026-08-22 23:03:27', '2026-08-22 23:03:27');
INSERT INTO `user` VALUES (18, 'g2_1787411030914_a', '$2a$10$xOiRqkXDht8SCPp2EIWXmeuBHzAMEn/AtzqH3o14Z.d9CnL0kK.RC', 'g2_1787411030914_a', '', '2026-08-22 23:03:50', '2026-08-22 23:03:50');
INSERT INTO `user` VALUES (19, 'g2_1787411030914_b', '$2a$10$xBN0fVTNs.YDsG7w5uUb3ua4KkcAucKMGmEnGsaO/5rYwBFFe.NLq', 'g2_1787411030914_b', '', '2026-08-22 23:03:51', '2026-08-22 23:03:51');
INSERT INTO `user` VALUES (20, 'rc_1787643680896_a', '$2a$10$dF70NEotZ.NwKYK/JPk3uOs0fb3GoiShh6CphXIE7vSL3oHOGzq6i', 'rc_1787643680896_a', '', '2026-08-25 15:41:20', '2026-08-25 15:41:20');
INSERT INTO `user` VALUES (21, 'rc_1787643680896_b', '$2a$10$WgoOB7JK/uCVQu5u.nR3G.FahvSd8ZZ1QffRRIx7DpjZpfBlc0ndK', 'rc_1787643680896_b', '', '2026-08-25 15:41:21', '2026-08-25 15:41:21');
INSERT INTO `user` VALUES (22, 'rc_1787643818294_a', '$2a$10$4rfCjb8rXVF/QIMlZ9IX/e7RzN0QvrG2pT4.cdUuMbRcfTP9Jlwh2', 'rc_1787643818294_a', '', '2026-08-25 15:43:37', '2026-08-25 15:43:37');
INSERT INTO `user` VALUES (23, 'rc_1787643818294_b', '$2a$10$u5nEwcUwKfe1FOP/nnAjneDSY4PmDOyo1eMzhtqhs/Tt0bL5eyf46', 'rc_1787643818294_b', '', '2026-08-25 15:43:38', '2026-08-25 15:43:38');
INSERT INTO `user` VALUES (24, 'rc_1787643895254_a', '$2a$10$C4G8GJwaSLkfDhtB0KDu7.hE5KycYVdkqmRjPR1aLb45ZYj23fgEW', 'rc_1787643895254_a', '', '2026-08-25 15:44:54', '2026-08-25 15:44:54');
INSERT INTO `user` VALUES (25, 'rc_1787643895254_b', '$2a$10$qrNCO/AO8zPXYCLFff/4Te2GGyfA.IbhpF/TMy0jX9nZ0Z4GX5jQ2', 'rc_1787643895254_b', '', '2026-08-25 15:44:55', '2026-08-25 15:44:55');
INSERT INTO `user` VALUES (26, 'ol_1787646084784_a', '$2a$10$/30bzWDZOiklbFlNkk8SQOKpiL.qWKK8UvHu4vy4wwbsaM1BXYYY.', 'ol_1787646084784_a', '', '2026-08-25 16:21:24', '2026-08-25 16:21:24');
INSERT INTO `user` VALUES (27, 'ol_1787646084784_b', '$2a$10$ygyg1aMAvaOn0.7C/nuHUuQff2poLwbEW2fLRjJx8dRp6cvBlj35S', 'ol_1787646084784_b', '', '2026-08-25 16:21:24', '2026-08-25 16:21:24');
INSERT INTO `user` VALUES (28, 'ol_1787646120391_a', '$2a$10$jenI9GcSPq/TcXrnOJ15C.PydmqfWvexQkLVWzCVR3kEfZEb69Ocu', 'ol_1787646120391_a', '', '2026-08-25 16:21:59', '2026-08-25 16:21:59');
INSERT INTO `user` VALUES (29, 'ol_1787646120391_b', '$2a$10$2eu9h2HKOaPIug6m8VzGL.vVGDAfECSVmd3la80TRuozxxvXaS1du', 'ol_1787646120391_b', '', '2026-08-25 16:22:00', '2026-08-25 16:22:00');
INSERT INTO `user` VALUES (30, 'ol_1787646154502_a', '$2a$10$9lXs.n0VJmBASr6g5czDTuwlW92MYIBsQTAxedkuOPV93kldF.3Za', 'ol_1787646154502_a', '', '2026-08-25 16:22:34', '2026-08-25 16:22:34');
INSERT INTO `user` VALUES (31, 'ol_1787646154502_b', '$2a$10$Ez/w/ITpkrvIq/SVI62Z0u0CPjxPaenbXM8rOI9Qzmjdr/cFH4o9u', 'ol_1787646154502_b', '', '2026-08-25 16:22:34', '2026-08-25 16:22:34');
INSERT INTO `user` VALUES (32, 'ol_1787646170095_a', '$2a$10$rChxExtk.TKs1q.8qMWFB.z6TDnmISMjkdtM20tY2eQ8JNgSHlXtW', 'ol_1787646170095_a', '', '2026-08-25 16:22:49', '2026-08-25 16:22:49');
INSERT INTO `user` VALUES (33, 'ol_1787646170095_b', '$2a$10$29avmVL28pekRM.ZGciup.LtMO.NwqQk3PLFDuA1a9pR6BivCRgIy', 'ol_1787646170095_b', '', '2026-08-25 16:22:49', '2026-08-25 16:22:49');
INSERT INTO `user` VALUES (34, 'ol_1787647259467_a', '$2a$10$mXGZIXwCpAYzAfxpw4u93O8NwT.6TpchRvNh9j0o.ViJgO2IcUVhG', 'ol_1787647259467_a', '', '2026-08-25 16:40:59', '2026-08-25 16:40:59');
INSERT INTO `user` VALUES (35, 'ol_1787647259467_b', '$2a$10$Gpzb7YolxLmYyINKkynPx.R0qRapTc8I8FSFuLpbZktoIfhrPM07e', 'ol_1787647259467_b', '', '2026-08-25 16:40:59', '2026-08-25 16:40:59');
INSERT INTO `user` VALUES (36, 'zd_1787647277023', '$2a$10$qq3kM4LesflAtzqMEHo8pOdCPDvwEBAxz/4FxueGOBYG6hu/83qiK', 'zd_1787647277023', '', '2026-08-25 16:41:16', '2026-08-25 16:41:16');
INSERT INTO `user` VALUES (37, 'zf_1787647295352_a', '$2a$10$rpW5Eseo2IioNq2ApcPYtOEhG1eQ5U.Faj.5et8hf1ChWassZAFgS', 'zf_1787647295352_a', '', '2026-08-25 16:41:34', '2026-08-25 16:41:34');
INSERT INTO `user` VALUES (38, 'zf_1787647295352_b', '$2a$10$Uxj55dgfCWlak4J4lImqTeVzj7SjB.eILpPawUoVzCCVeTnlNBszu', 'zf_1787647295352_b', '', '2026-08-25 16:41:35', '2026-08-25 16:41:35');

SET FOREIGN_KEY_CHECKS = 1;
