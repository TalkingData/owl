/*
 Navicat Premium Data Transfer

 Source Server         : 172.28.5.2
 Source Server Type    : MySQL
 Source Server Version : 50713
 Source Host           : 172.28.5.2:3333
 Source Schema         : owl

 Target Server Type    : MySQL
 Target Server Version : 50713
 File Encoding         : 65001

 Date: 23/03/2018 18:46:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for action
-- ----------------------------
DROP TABLE IF EXISTS `action`;
CREATE TABLE `action` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `type` int(1) unsigned NOT NULL,
  `kind` int(1) NOT NULL,
  `alarm_subject` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `alarm_template` varchar(1024) COLLATE utf8_unicode_ci NOT NULL,
  `restore_subject` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `restore_template` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `script_id` int(10) unsigned NOT NULL,
  `begin_time` time NOT NULL DEFAULT '00:00:00',
  `end_time` time NOT NULL DEFAULT '23:59:59',
  `time_period` int(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_action_strategy_id` (`strategy_id`),
  KEY `script_id` (`script_id`),
  CONSTRAINT `action_ibfk_1` FOREIGN KEY (`script_id`) REFERENCES `scripts` (`id`),
  CONSTRAINT `fk_action_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for action_result
-- ----------------------------
DROP TABLE IF EXISTS `action_result`;
CREATE TABLE `action_result` (
  `strategy_event_id` bigint(20) unsigned NOT NULL,
  `count` int(4) NOT NULL DEFAULT '1',
  `action_id` bigint(20) unsigned NOT NULL,
  `action_type` tinyint(1) unsigned NOT NULL,
  `action_kind` tinyint(1) NOT NULL,
  `script_id` int(1) unsigned NOT NULL,
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
  KEY `index_strategy_event_id_count` (`strategy_event_id`,`count`),
  CONSTRAINT `fk_action_result_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for action_user_group
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for chart
-- ----------------------------
DROP TABLE IF EXISTS `chart`;
CREATE TABLE `chart` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `creator` varchar(255) NOT NULL DEFAULT '',
  `span` int(255) NOT NULL DEFAULT '12',
  `height` int(255) NOT NULL DEFAULT '100',
  `create_at` datetime NOT NULL,
  `type` varchar(255) NOT NULL,
  `panel_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `chart_panel_id_fk` (`panel_id`),
  CONSTRAINT `chart_panel_id_fk` FOREIGN KEY (`panel_id`) REFERENCES `panel` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for chart_element
-- ----------------------------
DROP TABLE IF EXISTS `chart_element`;
CREATE TABLE `chart_element` (
  `metric` varchar(1024) NOT NULL,
  `tags` varchar(2048) NOT NULL,
  `chart_id` int(10) unsigned NOT NULL,
  KEY `chart_element_chart_id_fk` (`chart_id`),
  CONSTRAINT `chart_element_chart_id_fk` FOREIGN KEY (`chart_id`) REFERENCES `chart` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for host
-- ----------------------------
DROP TABLE IF EXISTS `host`;
CREATE TABLE `host` (
  `id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `ip` varchar(16) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `hostname` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `agent_version` varchar(10) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'Unknown',
  `status` enum('0','1','2','3') COLLATE utf8_unicode_ci NOT NULL DEFAULT '3',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mute_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `uptime` float NOT NULL DEFAULT '0',
  `idle_pct` float NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_host_id` (`id`),
  UNIQUE KEY `idx_hostname_ip` (`hostname`,`ip`,`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for host_group
-- ----------------------------
DROP TABLE IF EXISTS `host_group`;
CREATE TABLE `host_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `description` text COLLATE utf8_unicode_ci NOT NULL,
  `product_id` int(10) unsigned NOT NULL,
  `creator` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_at` datetime NOT NULL,
  `update_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name_product_id` (`name`,`product_id`) USING BTREE,
  KEY `group_product_id_fk` (`product_id`),
  CONSTRAINT `group_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for host_group_host
-- ----------------------------
DROP TABLE IF EXISTS `host_group_host`;
CREATE TABLE `host_group_host` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `host_group_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_host_id_group_id` (`host_id`,`host_group_id`),
  KEY `fk_host_group_group_id` (`host_group_id`),
  CONSTRAINT `fk_host_group_group_id` FOREIGN KEY (`host_group_id`) REFERENCES `host_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_host_group_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for host_group_plugin
-- ----------------------------
DROP TABLE IF EXISTS `host_group_plugin`;
CREATE TABLE `host_group_plugin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_id` int(10) unsigned NOT NULL,
  `plugin_id` int(10) unsigned NOT NULL,
  `args` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `interval` int(1) NOT NULL DEFAULT '60',
  `timeout` int(1) NOT NULL DEFAULT '10',
  PRIMARY KEY (`id`),
  KEY `group_id` (`group_id`),
  KEY `plugin_id` (`plugin_id`),
  CONSTRAINT `fk_group_id` FOREIGN KEY (`group_id`) REFERENCES `host_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_plugin_id` FOREIGN KEY (`plugin_id`) REFERENCES `plugin` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for host_plugin
-- ----------------------------
DROP TABLE IF EXISTS `host_plugin`;
CREATE TABLE `host_plugin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `plugin_id` int(10) unsigned NOT NULL,
  `args` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `interval` int(1) NOT NULL DEFAULT '60',
  `timeout` int(1) NOT NULL DEFAULT '10',
  PRIMARY KEY (`id`),
  KEY `fk_host_plugin_plugin_id` (`plugin_id`),
  KEY `fk_host_plugin_host_id` (`host_id`),
  CONSTRAINT `fk_host_plugin_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_host_plugin_plugin_id` FOREIGN KEY (`plugin_id`) REFERENCES `plugin` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for metric
-- ----------------------------
DROP TABLE IF EXISTS `metric`;
CREATE TABLE `metric` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `metric` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'cpu.idle/cpu=all',
  `tags` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `dt` enum('GAUGE','DERIVE','COUNTER') COLLATE utf8_unicode_ci NOT NULL,
  `cycle` int(10) NOT NULL,
  `create_at` datetime NOT NULL,
  `update_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`metric`,`host_id`,`tags`) USING BTREE,
  KEY `fk_metric_host` (`host_id`),
  CONSTRAINT `fk_metric_host` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for operations
-- ----------------------------
DROP TABLE IF EXISTS `operations`;
CREATE TABLE `operations` (
  `ip` char(15) COLLATE utf8_unicode_ci NOT NULL,
  `operator` char(255) COLLATE utf8_unicode_ci NOT NULL,
  `method` char(10) COLLATE utf8_unicode_ci NOT NULL,
  `api` text COLLATE utf8_unicode_ci NOT NULL,
  `body` text COLLATE utf8_unicode_ci NOT NULL,
  `result` tinyint(1) unsigned NOT NULL,
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY `operation_time` (`time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for panel
-- ----------------------------
DROP TABLE IF EXISTS `panel`;
CREATE TABLE `panel` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(10) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `creator` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_panel_product_id` (`product_id`),
  CONSTRAINT `fk_panel_product_id` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for plugin
-- ----------------------------
DROP TABLE IF EXISTS `plugin`;
CREATE TABLE `plugin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `path` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `args` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `interval` int(1) NOT NULL DEFAULT '60',
  `timeout` int(1) NOT NULL DEFAULT '10',
  `checksum` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `update_at` datetime NOT NULL,
  `create_at` datetime NOT NULL,
  `creator` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for product
-- ----------------------------
DROP TABLE IF EXISTS `product`;
CREATE TABLE `product` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(1024) NOT NULL DEFAULT '',
  `creator` varchar(1024) NOT NULL DEFAULT '',
  `create_at` datetime NOT NULL,
  `is_delete` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `product_name_uindex` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for product_host
-- ----------------------------
DROP TABLE IF EXISTS `product_host`;
CREATE TABLE `product_host` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(10) unsigned NOT NULL,
  `host_id` char(32) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `product_host_product_id_host_id_unique` (`product_id`,`host_id`),
  KEY `product_host_host_id_fk` (`host_id`),
  KEY `host_id_fk` (`product_id`),
  CONSTRAINT `product_host_host_id_fk` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `product_host_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for product_user
-- ----------------------------
DROP TABLE IF EXISTS `product_user`;
CREATE TABLE `product_user` (
  `product_id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  UNIQUE KEY `product_user_product_id_user_id_unique` (`product_id`,`user_id`),
  KEY `product_user_product_id_fk` (`product_id`),
  KEY `product_user_user_id_fk` (`user_id`),
  CONSTRAINT `product_user_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`),
  CONSTRAINT `product_user_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for scripts
-- ----------------------------
DROP TABLE IF EXISTS `scripts`;
CREATE TABLE `scripts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `file_path` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for strategy
-- ----------------------------
DROP TABLE IF EXISTS `strategy`;
CREATE TABLE `strategy` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(10) unsigned NOT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `priority` int(1) unsigned NOT NULL,
  `alarm_count` int(4) unsigned NOT NULL DEFAULT '0',
  `cycle` int(10) unsigned NOT NULL DEFAULT '5',
  `expression` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT '',
  `enable` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `user_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `product_id` (`product_id`,`name`) USING BTREE,
  KEY `fk_strategy_product_id` (`product_id`),
  CONSTRAINT `fk_strategy_product_id` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for strategy_event
-- ----------------------------
DROP TABLE IF EXISTS `strategy_event`;
CREATE TABLE `strategy_event` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(10) unsigned NOT NULL,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `strategy_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `priority` int(1) unsigned NOT NULL,
  `cycle` int(10) unsigned NOT NULL,
  `alarm_count` int(4) unsigned NOT NULL,
  `expression` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `aware_end_time` timestamp NOT NULL DEFAULT '1980-01-01 00:00:00',
  `count` int(4) unsigned NOT NULL DEFAULT '1',
  `status` int(1) unsigned NOT NULL DEFAULT '1',
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `host_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `ip` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_host_id_strategy_id_status` (`host_id`,`strategy_id`,`status`),
  KEY `fk_strategy_event_product_id` (`product_id`),
  CONSTRAINT `fk_strategy_event_product_id` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for strategy_event_failed
-- ----------------------------
DROP TABLE IF EXISTS `strategy_event_failed`;
CREATE TABLE `strategy_event_failed` (
  `strategy_id` bigint(20) unsigned NOT NULL,
  `host_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `status` int(11) NOT NULL,
  `message` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `idx_strategy_id_host_id` (`strategy_id`,`host_id`) USING BTREE,
  KEY `strategy_id` (`strategy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for strategy_event_process
-- ----------------------------
DROP TABLE IF EXISTS `strategy_event_process`;
CREATE TABLE `strategy_event_process` (
  `strategy_event_id` bigint(20) unsigned NOT NULL,
  `process_status` int(255) DEFAULT NULL,
  `process_user` varchar(255) DEFAULT NULL,
  `process_comments` varchar(255) DEFAULT NULL,
  `process_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  KEY `strategy_event_id` (`strategy_event_id`),
  CONSTRAINT `fk_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for strategy_event_record
-- ----------------------------
DROP TABLE IF EXISTS `strategy_event_record`;
CREATE TABLE `strategy_event_record` (
  `strategy_event_id` bigint(20) unsigned NOT NULL,
  `count` int(4) unsigned NOT NULL DEFAULT '1',
  `strategy_id` bigint(20) unsigned NOT NULL,
  `strategy_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `priority` int(1) unsigned NOT NULL,
  `cycle` int(10) unsigned NOT NULL,
  `alarm_count` int(4) unsigned NOT NULL,
  `expression` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `aware_end_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` int(1) unsigned NOT NULL DEFAULT '1',
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  `host_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `ip` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  KEY `idx_strategy_event_record_strategy_event_id` (`strategy_event_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for strategy_group
-- ----------------------------
DROP TABLE IF EXISTS `strategy_group`;
CREATE TABLE `strategy_group` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `group_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_strategy_id_group_id` (`strategy_id`,`group_id`),
  KEY `fk_strategy_group_group_id` (`group_id`),
  CONSTRAINT `fk_strategy_group_group_id` FOREIGN KEY (`group_id`) REFERENCES `host_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_strategy_group_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for strategy_host_exclude
-- ----------------------------
DROP TABLE IF EXISTS `strategy_host_exclude`;
CREATE TABLE `strategy_host_exclude` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `host_id` char(32) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_strategy_id_host_id` (`strategy_id`,`host_id`),
  KEY `idx_strategy_host_host_id` (`host_id`),
  CONSTRAINT `fk_strategy_host_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_strategy_host_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for strategy_template
-- ----------------------------
DROP TABLE IF EXISTS `strategy_template`;
CREATE TABLE `strategy_template` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `alarm_count` int(4) unsigned NOT NULL,
  `cycle` int(10) unsigned NOT NULL,
  `expression` varchar(255) NOT NULL,
  `description` varchar(1024) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for trigger
-- ----------------------------
DROP TABLE IF EXISTS `trigger`;
CREATE TABLE `trigger` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_id` bigint(20) unsigned NOT NULL,
  `metric` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `tags` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `number` int(10) unsigned NOT NULL DEFAULT '0',
  `index` char(10) COLLATE utf8_unicode_ci NOT NULL,
  `method` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `symbol` char(5) COLLATE utf8_unicode_ci NOT NULL,
  `threshold` double(255,2) NOT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_strategy_id_index_index` (`strategy_id`,`index`),
  CONSTRAINT `fk_strategy_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for trigger_event
-- ----------------------------
DROP TABLE IF EXISTS `trigger_event`;
CREATE TABLE `trigger_event` (
  `strategy_event_id` bigint(20) unsigned NOT NULL,
  `index` char(10) COLLATE utf8_unicode_ci NOT NULL,
  `metric` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
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
-- Table structure for trigger_event_record
-- ----------------------------
DROP TABLE IF EXISTS `trigger_event_record`;
CREATE TABLE `trigger_event_record` (
  `strategy_event_id` bigint(20) unsigned NOT NULL,
  `count` int(4) unsigned NOT NULL DEFAULT '1',
  `index` char(10) COLLATE utf8_unicode_ci NOT NULL,
  `metric` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `tags` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `number` int(10) unsigned NOT NULL DEFAULT '0',
  `aggregate_tags` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `current_threshold` double(255,2) NOT NULL,
  `method` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `symbol` char(5) COLLATE utf8_unicode_ci NOT NULL,
  `threshold` double(255,2) NOT NULL,
  `triggered` tinyint(1) NOT NULL DEFAULT '0',
  KEY `strategy_event_id` (`strategy_event_id`),
  CONSTRAINT `fk_trigger_event_record_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event_record` (`strategy_event_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for trigger_template
-- ----------------------------
DROP TABLE IF EXISTS `trigger_template`;
CREATE TABLE `trigger_template` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `strategy_template_id` int(10) unsigned NOT NULL,
  `metric` varchar(255) NOT NULL,
  `tags` varchar(255) NOT NULL,
  `number` int(10) NOT NULL,
  `index` char(10) NOT NULL,
  `method` varchar(10) NOT NULL,
  `symbol` char(5) NOT NULL,
  `threshold` double(255,2) NOT NULL,
  `description` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `strategy_template_id` (`strategy_template_id`),
  CONSTRAINT `trigger_template_ibfk_1` FOREIGN KEY (`strategy_template_id`) REFERENCES `strategy_template` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` char(255) COLLATE utf8_unicode_ci NOT NULL,
  `display_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `password` char(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `role` tinyint(1) NOT NULL DEFAULT '0',
  `phone` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `mail` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `wechat` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for user_group
-- ----------------------------
DROP TABLE IF EXISTS `user_group`;
CREATE TABLE `user_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `product_id` int(10) unsigned NOT NULL,
  `description` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `creator` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`,`product_id`) USING BTREE,
  KEY `user_group_product_id_fk` (`product_id`),
  CONSTRAINT `user_group_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for user_group_user
-- ----------------------------
DROP TABLE IF EXISTS `user_group_user`;
CREATE TABLE `user_group_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `user_group_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id_group_id` (`user_id`,`user_group_id`),
  KEY `fk_user_usergroup_user_group_id` (`user_group_id`),
  CONSTRAINT `fk_user_user_group_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_user_usergroup_user_group_id` FOREIGN KEY (`user_group_id`) REFERENCES `user_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

ALTER TABLE `strategy_event_failed`
ADD CONSTRAINT `stragety_event_failed_host_id_fk` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT `strategy_event_failed_strategy_id_fk` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE  `strategy_event_record`
ADD CONSTRAINT `strategy_event_record_host_id_fk` FOREIGN KEY (`host_id`) REFERENCES  `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- 增加用户类型
ALTER TABLE `owl`.`user`
ADD COLUMN `type` varchar(255) NOT NULL DEFAULT '' AFTER `wechat`;

-- 插件增加备注
ALTER TABLE `host_group_plugin` ADD COLUMN `comment` varchar(1024) NOT NULL DEFAULT '' AFTER `timeout`;
ALTER TABLE `host_plugin` ADD COLUMN `comment` varchar(1024) NOT NULL DEFAULT '' AFTER `timeout`;
ALTER TABLE `plugin` ADD COLUMN `comment` text NOT NULL AFTER `creator`;
SET FOREIGN_KEY_CHECKS = 1;
