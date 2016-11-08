/*
 Source Server Type    : MySQL
 Source Server Version : 50713
 Source Database       : owl

 Target Server Type    : MySQL
 Target Server Version : 50713
 File Encoding         : utf-8

 Date: 11/08/2016 10:46:50 AM
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `action`
-- ----------------------------
DROP TABLE IF EXISTS `action`;
CREATE TABLE `action` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `type` tinyint(1) unsigned NOT NULL,
  `file_path` text COLLATE utf8_unicode_ci NOT NULL,
  `alarm_subject` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `restore_subject` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `alarm_template` text COLLATE utf8_unicode_ci NOT NULL,
  `restore_template` text COLLATE utf8_unicode_ci NOT NULL,
  `time_out` int(4) unsigned NOT NULL DEFAULT '30',
  `send_type` int(1) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_action_strategy_id` (`strategy_id`),
  CONSTRAINT `fk_action_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=452 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `action_result`
-- ----------------------------
DROP TABLE IF EXISTS `action_result`;
CREATE TABLE `action_result` (
  `strategy_event_id` bigint(20) unsigned NOT NULL,
  `action_id` bigint(20) unsigned NOT NULL,
  `action_type` tinyint(1) unsigned NOT NULL,
  `action_send_type` int(1) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `username` char(255) COLLATE utf8_unicode_ci NOT NULL,
  `phone` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `mail` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `weixin` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `subject` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `content` text COLLATE utf8_unicode_ci NOT NULL,
  `success` tinyint(1) unsigned NOT NULL,
  `response` text COLLATE utf8_unicode_ci NOT NULL,
  KEY `fk_action_result_strategy_event_id` (`strategy_event_id`),
  CONSTRAINT `fk_action_result_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `action_user`
-- ----------------------------
DROP TABLE IF EXISTS `action_user`;
CREATE TABLE `action_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `action_id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_action_id_user_id` (`action_id`,`user_id`),
  KEY `fk_action_user_user_id` (`user_id`),
  CONSTRAINT `fk_action_user_action_id` FOREIGN KEY (`action_id`) REFERENCES `action` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_action_user_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `action_user_group`
-- ----------------------------
DROP TABLE IF EXISTS `action_user_group`;
CREATE TABLE `action_user_group` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `action_id` int(10) unsigned NOT NULL,
  `user_group_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_action_id_user_group_id` (`action_id`,`user_group_id`),
  KEY `fk_action_user_group_user_group_id` (`user_group_id`),
  CONSTRAINT `fk_action_user_group_action_id` FOREIGN KEY (`action_id`) REFERENCES `action` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_action_user_group_user_group_id` FOREIGN KEY (`user_group_id`) REFERENCES `user_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=636 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `chart`
-- ----------------------------
DROP TABLE IF EXISTS `chart`;
CREATE TABLE `chart` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `size` int(1) unsigned NOT NULL,
  `refer_count` int(1) unsigned zerofill NOT NULL,
  `thumbnail` longtext COLLATE utf8_unicode_ci NOT NULL,
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `fk_chart_user` (`user_id`),
  CONSTRAINT `fk_chart_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=98 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `chart_element`
-- ----------------------------
DROP TABLE IF EXISTS `chart_element`;
CREATE TABLE `chart_element` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `chart_id` int(10) unsigned DEFAULT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `metric` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `tags` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_chart_element_chart_id` (`chart_id`),
  CONSTRAINT `fk_chart_element_chart_id` FOREIGN KEY (`chart_id`) REFERENCES `chart` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=140 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `group`
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `host`
-- ----------------------------
DROP TABLE IF EXISTS `host`;
CREATE TABLE `host` (
  `id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `ip` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  `sn` varchar(128) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `hostname` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `agent_version` varchar(10) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'Unknown',
  `status` enum('0','1','2','3') COLLATE utf8_unicode_ci NOT NULL DEFAULT '3',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_hostname_ip` (`hostname`,`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `host_group`
-- ----------------------------
DROP TABLE IF EXISTS `host_group`;
CREATE TABLE `host_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `group_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_host_id_group_id` (`host_id`,`group_id`),
  KEY `fk_host_group_group_id` (`group_id`),
  CONSTRAINT `fk_host_group_group_id` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_host_group_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=460 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `host_plugin`
-- ----------------------------
DROP TABLE IF EXISTS `host_plugin`;
CREATE TABLE `host_plugin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `plugin_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_host_id_plugin_id` (`host_id`,`plugin_id`),
  KEY `fk_host_plugin_plugin_id` (`plugin_id`),
  CONSTRAINT `fk_host_plugin_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_host_plugin_plugin_id` FOREIGN KEY (`plugin_id`) REFERENCES `plugin` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=217 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `metric`
-- ----------------------------
DROP TABLE IF EXISTS `metric`;
CREATE TABLE `metric` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'cpu.idle/cpu=all',
  `dt` enum('GAUGE','DRIVER','COUNTER') COLLATE utf8_unicode_ci NOT NULL,
  `cycle` int(10) NOT NULL,
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`,`host_id`),
  KEY `fk_metric_host` (`host_id`),
  CONSTRAINT `fk_metric_host` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=19813 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `metric_index`
-- ----------------------------
DROP TABLE IF EXISTS `metric_index`;
CREATE TABLE `metric_index` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `metric` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_metric` (`metric`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=465 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `operations`
-- ----------------------------
DROP TABLE IF EXISTS `operations`;
CREATE TABLE `operations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `operation_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ip` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  `operator` char(255) COLLATE utf8_unicode_ci NOT NULL,
  `content` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `operation_result` tinyint(1) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `operation_time` (`operation_time`)
) ENGINE=InnoDB AUTO_INCREMENT=222220 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `panel`
-- ----------------------------
DROP TABLE IF EXISTS `panel`;
CREATE TABLE `panel` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `favor` int(1) NOT NULL DEFAULT '0',
  `thumbnail` longtext COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `fk_panel_user` (`user_id`),
  CONSTRAINT `fk_panel_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=67 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `panel_chart`
-- ----------------------------
DROP TABLE IF EXISTS `panel_chart`;
CREATE TABLE `panel_chart` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `chart_id` int(10) unsigned NOT NULL,
  `panel_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_chart_id_panel_id` (`chart_id`,`panel_id`),
  KEY `fk_panel_chart_panel_id` (`panel_id`),
  CONSTRAINT `fk_panel_chart_chart_id` FOREIGN KEY (`chart_id`) REFERENCES `chart` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_panel_chart_panel_id` FOREIGN KEY (`panel_id`) REFERENCES `panel` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=60 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `plugin`
-- ----------------------------
DROP TABLE IF EXISTS `plugin`;
CREATE TABLE `plugin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `args` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `interval` int(1) NOT NULL DEFAULT '30',
  `timeout` int(1) NOT NULL DEFAULT '10',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `strategy`
-- ----------------------------
DROP TABLE IF EXISTS `strategy`;
CREATE TABLE `strategy` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `priority` int(1) unsigned NOT NULL,
  `type` int(1) unsigned NOT NULL,
  `pid` int(10) unsigned NOT NULL,
  `alarm_count` int(4) unsigned NOT NULL DEFAULT '0',
  `cycle` int(10) unsigned NOT NULL DEFAULT '5',
  `expression` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `group_id` int(10) unsigned NOT NULL DEFAULT '0',
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT '',
  `enable` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=193 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `strategy_event`
-- ----------------------------
DROP TABLE IF EXISTS `strategy_event`;
CREATE TABLE `strategy_event` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `strategy_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `strategy_type` int(1) unsigned NOT NULL,
  `priority` int(1) unsigned NOT NULL,
  `cycle` int(10) unsigned NOT NULL,
  `alarm_count` int(4) unsigned NOT NULL,
  `expression` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `count` int(4) unsigned NOT NULL DEFAULT '1',
  `status` int(1) unsigned NOT NULL DEFAULT '1',
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `host_cname` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `host_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `ip` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  `sn` varchar(128) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `process_user` char(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `process_comments` text COLLATE utf8_unicode_ci,
  `process_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5513 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `strategy_group`
-- ----------------------------
DROP TABLE IF EXISTS `strategy_group`;
CREATE TABLE `strategy_group` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `group_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_strategy_id_group_id` (`strategy_id`,`group_id`),
  KEY `fk_strategy_group_group_id` (`group_id`),
  CONSTRAINT `fk_strategy_group_group_id` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_strategy_group_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=298 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `strategy_host`
-- ----------------------------
DROP TABLE IF EXISTS `strategy_host`;
CREATE TABLE `strategy_host` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_strategy_id_host_id` (`strategy_id`,`host_id`),
  KEY `fk_strategy_host_host_id` (`host_id`),
  CONSTRAINT `fk_strategy_host_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_strategy_host_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=440 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `tag_index`
-- ----------------------------
DROP TABLE IF EXISTS `tag_index`;
CREATE TABLE `tag_index` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `tagk` varchar(255) NOT NULL,
  `tagv` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tagk_tagv` (`tagk`,`tagv`)
) ENGINE=InnoDB AUTO_INCREMENT=1119 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `trigger`
-- ----------------------------
DROP TABLE IF EXISTS `trigger`;
CREATE TABLE `trigger` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `metric` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `tags` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `number` int(10) unsigned NOT NULL DEFAULT '0',
  `index` char(10) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `method` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `symbol` char(5) COLLATE utf8_unicode_ci NOT NULL,
  `threshold` double(255,2) NOT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_strategy_id_index_index` (`strategy_id`,`index`),
  CONSTRAINT `fk_strategy_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=259 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `trigger_event`
-- ----------------------------
DROP TABLE IF EXISTS `trigger_event`;
CREATE TABLE `trigger_event` (
  `strategy_event_id` bigint(20) unsigned NOT NULL,
  `index` char(10) COLLATE utf8_unicode_ci NOT NULL,
  `metric` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `tags` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `number` int(10) unsigned NOT NULL DEFAULT '0',
  `aggregate_tags` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `current_threshold` double(255,2) NOT NULL,
  `method` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `symbol` char(5) COLLATE utf8_unicode_ci NOT NULL,
  `threshold` double(255,2) NOT NULL,
  `triggered` tinyint(1) NOT NULL DEFAULT '0',
  KEY `fk_trigger_event_strategy_event_id` (`strategy_event_id`),
  CONSTRAINT `fk_trigger_event_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `user`
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` char(255) COLLATE utf8_unicode_ci NOT NULL,
  `password` char(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `role` tinyint(1) NOT NULL DEFAULT '2',
  `phone` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `mail` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `weixin` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=190 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `user_group`
-- ----------------------------
DROP TABLE IF EXISTS `user_group`;
CREATE TABLE `user_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `user_user_group`
-- ----------------------------
DROP TABLE IF EXISTS `user_user_group`;
CREATE TABLE `user_user_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `user_group_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id_group_id` (`user_id`,`user_group_id`),
  KEY `fk_user_usergroup_user_group_id` (`user_group_id`),
  CONSTRAINT `fk_user_user_group_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_user_usergroup_user_group_id` FOREIGN KEY (`user_group_id`) REFERENCES `user_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=273 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
