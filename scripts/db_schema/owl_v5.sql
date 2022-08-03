-- MySQL dump 10.13  Distrib 5.7.27, for Linux (x86_64)
--
-- Host: localhost    Database: owl_v5
-- ------------------------------------------------------
-- Server version	5.7.27-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `action`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action`
(
    `id`               int(10) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_id`      bigint(20) unsigned NOT NULL,
    `type`             int(1) unsigned NOT NULL,
    `kind`             int(1) NOT NULL,
    `alarm_subject`    varchar(255) COLLATE utf8_unicode_ci  NOT NULL,
    `alarm_template`   varchar(1024) COLLATE utf8_unicode_ci NOT NULL,
    `restore_subject`  varchar(255) COLLATE utf8_unicode_ci  NOT NULL DEFAULT '',
    `restore_template` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `script_id`        int(10) unsigned NOT NULL,
    `begin_time`       time                                  NOT NULL DEFAULT '00:00:00',
    `end_time`         time                                  NOT NULL DEFAULT '23:59:59',
    `time_period`      int(4) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY                `fk_action_strategy_id` (`strategy_id`),
    KEY                `script_id` (`script_id`),
    CONSTRAINT `action_ibfk_1` FOREIGN KEY (`script_id`) REFERENCES `scripts` (`id`),
    CONSTRAINT `fk_action_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `action_bak`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action_bak`
(
    `id`               int(10) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_id`      bigint(20) unsigned NOT NULL,
    `type`             int(1) unsigned NOT NULL,
    `kind`             int(1) NOT NULL,
    `alarm_subject`    varchar(255) COLLATE utf8_unicode_ci  NOT NULL,
    `alarm_template`   varchar(1024) COLLATE utf8_unicode_ci NOT NULL,
    `restore_subject`  varchar(255) COLLATE utf8_unicode_ci  NOT NULL DEFAULT '',
    `restore_template` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `script_id`        int(10) unsigned NOT NULL,
    `begin_time`       time                                  NOT NULL DEFAULT '00:00:00',
    `end_time`         time                                  NOT NULL DEFAULT '23:59:59',
    `time_period`      int(4) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY                `fk_action_strategy_id` (`strategy_id`),
    KEY                `script_id` (`script_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `action_result`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action_result`
(
    `strategy_event_id` bigint(20) unsigned NOT NULL,
    `count`             int(4) NOT NULL DEFAULT '1',
    `action_id`         bigint(20) unsigned NOT NULL,
    `action_type`       tinyint(1) unsigned NOT NULL,
    `action_kind`       tinyint(1) NOT NULL,
    `script_id`         int(1) unsigned NOT NULL,
    `user_id`           int(10) unsigned NOT NULL,
    `username`          char(255) COLLATE utf8_unicode_ci    NOT NULL,
    `phone`             varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `mail`              varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `weixin`            varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `subject`           varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `content`           text COLLATE utf8_unicode_ci         NOT NULL,
    `success`           tinyint(1) unsigned NOT NULL,
    `response`          text COLLATE utf8_unicode_ci         NOT NULL,
    KEY                 `fk_action_result_strategy_event_id` (`strategy_event_id`),
    KEY                 `index_strategy_event_id_count` (`strategy_event_id`,`count`),
    CONSTRAINT `fk_action_result_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `action_user_group`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action_user_group`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `action_id`     int(10) unsigned NOT NULL,
    `user_group_id` int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_action_id_user_group_id` (`action_id`,`user_group_id`),
    KEY             `fk_action_user_group_user_group_id` (`user_group_id`),
    CONSTRAINT `fk_action_user_group_action_id` FOREIGN KEY (`action_id`) REFERENCES `action` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_action_user_group_user_group_id` FOREIGN KEY (`user_group_id`) REFERENCES `user_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chart`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chart`
(
    `id`        int(10) unsigned NOT NULL AUTO_INCREMENT,
    `title`     varchar(255) NOT NULL,
    `creator`   varchar(255) NOT NULL DEFAULT '',
    `span`      int(255) NOT NULL DEFAULT '12',
    `height`    int(255) NOT NULL DEFAULT '100',
    `create_at` datetime     NOT NULL,
    `type`      varchar(255) NOT NULL,
    `panel_id`  int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    KEY         `chart_panel_id_fk` (`panel_id`),
    CONSTRAINT `chart_panel_id_fk` FOREIGN KEY (`panel_id`) REFERENCES `panel` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chart_element`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chart_element`
(
    `metric`   varchar(1024) NOT NULL,
    `tags`     varchar(2048) NOT NULL,
    `chart_id` int(10) unsigned NOT NULL,
    KEY        `chart_element_chart_id_fk` (`chart_id`),
    CONSTRAINT `chart_element_chart_id_fk` FOREIGN KEY (`chart_id`) REFERENCES `chart` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `host`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host`
(
    `id`            char(32) COLLATE utf8_unicode_ci     NOT NULL,
    `name`          varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `ip`            varchar(16) COLLATE utf8_unicode_ci  NOT NULL DEFAULT '',
    `hostname`      varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `agent_version` varchar(10) COLLATE utf8_unicode_ci  NOT NULL DEFAULT 'Unknown',
    `status`        enum('0','1','2','3') COLLATE utf8_unicode_ci NOT NULL DEFAULT '3',
    `create_at`     timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at`     timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `mute_time`     timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `uptime`        float                                NOT NULL DEFAULT '0',
    `idle_pct`      float                                NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_host_id` (`id`),
    UNIQUE KEY `idx_hostname_ip` (`hostname`,`ip`,`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `host_group`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host_group`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(255) COLLATE utf8_unicode_ci  NOT NULL,
    `description` text COLLATE utf8_unicode_ci          NOT NULL,
    `product_id`  int(10) unsigned NOT NULL,
    `creator`     varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `create_at`   datetime                              NOT NULL,
    `update_at`   datetime                              NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name_product_id` (`name`,`product_id`) USING BTREE,
    KEY           `group_product_id_fk` (`product_id`),
    CONSTRAINT `group_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `host_group_host`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host_group_host`
(
    `id`            int(10) unsigned NOT NULL AUTO_INCREMENT,
    `host_id`       char(32) COLLATE utf8_unicode_ci NOT NULL,
    `host_group_id` int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_host_id_group_id` (`host_id`,`host_group_id`),
    KEY             `fk_host_group_group_id` (`host_group_id`),
    CONSTRAINT `fk_host_group_group_id` FOREIGN KEY (`host_group_id`) REFERENCES `host_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_host_group_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `host_group_plugin`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host_group_plugin`
(
    `id`        int(10) unsigned NOT NULL AUTO_INCREMENT,
    `group_id`  int(10) unsigned NOT NULL,
    `plugin_id` int(10) unsigned NOT NULL,
    `args`      varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `interval`  int(1) NOT NULL DEFAULT '60',
    `timeout`   int(1) NOT NULL DEFAULT '10',
    `comment`   varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    KEY         `group_id` (`group_id`),
    KEY         `plugin_id` (`plugin_id`),
    CONSTRAINT `fk_group_id` FOREIGN KEY (`group_id`) REFERENCES `host_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_plugin_id` FOREIGN KEY (`plugin_id`) REFERENCES `plugin` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `host_plugin`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host_plugin`
(
    `id`        int(10) unsigned NOT NULL AUTO_INCREMENT,
    `host_id`   char(32) COLLATE utf8_unicode_ci      NOT NULL,
    `plugin_id` int(10) unsigned NOT NULL,
    `args`      varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `interval`  int(1) NOT NULL DEFAULT '60',
    `timeout`   int(1) NOT NULL DEFAULT '10',
    `comment`   varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    KEY         `fk_host_plugin_plugin_id` (`plugin_id`),
    KEY         `fk_host_plugin_host_id` (`host_id`),
    CONSTRAINT `fk_host_plugin_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_host_plugin_plugin_id` FOREIGN KEY (`plugin_id`) REFERENCES `plugin` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `metric`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `metric`
(
    `id`        bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `host_id`   char(32) COLLATE utf8_unicode_ci     NOT NULL,
    `metric`    varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'cpu.idle/cpu=all',
    `tags`      varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `dt`        enum('GAUGE','DERIVE','COUNTER') COLLATE utf8_unicode_ci NOT NULL,
    `cycle`     int(10) NOT NULL,
    `create_at` datetime                             NOT NULL,
    `update_at` datetime                             NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`metric`,`host_id`,`tags`) USING BTREE,
    KEY         `fk_metric_host` (`host_id`),
    KEY         `idx_metric` (`metric`),
    CONSTRAINT `fk_metric_host` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `operations`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `operations`
(
    `ip`       char(15) COLLATE utf8_unicode_ci  NOT NULL,
    `operator` char(255) COLLATE utf8_unicode_ci NOT NULL,
    `method`   char(10) COLLATE utf8_unicode_ci  NOT NULL,
    `api`      text COLLATE utf8_unicode_ci      NOT NULL,
    `body`     text COLLATE utf8_unicode_ci      NOT NULL,
    `result`   tinyint(1) unsigned NOT NULL,
    `time`     timestamp                         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY        `operation_time` (`time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `panel`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `panel`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `product_id` int(10) unsigned NOT NULL,
    `name`       varchar(255) NOT NULL,
    `creator`    varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY          `fk_panel_product_id` (`product_id`),
    CONSTRAINT `fk_panel_product_id` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `plugin`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `plugin`
(
    `id`        int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(255) COLLATE utf8_unicode_ci  NOT NULL,
    `path`      varchar(255) COLLATE utf8_unicode_ci  NOT NULL,
    `args`      varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `interval`  int(1) NOT NULL DEFAULT '60',
    `timeout`   int(1) NOT NULL DEFAULT '10',
    `checksum`  varchar(255) COLLATE utf8_unicode_ci  NOT NULL,
    `update_at` datetime                              NOT NULL,
    `create_at` datetime                              NOT NULL,
    `creator`   varchar(255) COLLATE utf8_unicode_ci  NOT NULL DEFAULT '',
    `comment`   text COLLATE utf8_unicode_ci          NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `product`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(255)  NOT NULL,
    `description` varchar(1024) NOT NULL DEFAULT '',
    `creator`     varchar(1024) NOT NULL DEFAULT '',
    `create_at`   datetime      NOT NULL,
    `is_delete`   tinyint(1) unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `product_name_uindex` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `product_host`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_host`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `product_id` int(10) unsigned NOT NULL,
    `host_id`    char(32) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `product_host_product_id_host_id_unique` (`product_id`,`host_id`),
    KEY          `product_host_host_id_fk` (`host_id`),
    KEY          `host_id_fk` (`product_id`),
    CONSTRAINT `product_host_host_id_fk` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `product_host_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `product_user`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_user`
(
    `product_id` int(10) unsigned NOT NULL,
    `user_id`    int(10) unsigned NOT NULL,
    UNIQUE KEY `product_user_product_id_user_id_unique` (`product_id`,`user_id`),
    KEY          `product_user_product_id_fk` (`product_id`),
    KEY          `product_user_user_id_fk` (`user_id`),
    CONSTRAINT `product_user_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `product_user_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `scripts`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scripts`
(
    `id`        int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(255) NOT NULL,
    `file_path` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `product_id`  int(10) unsigned NOT NULL,
    `name`        varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `priority`    int(1) unsigned NOT NULL,
    `alarm_count` int(4) unsigned NOT NULL DEFAULT '0',
    `cycle`       int(10) unsigned NOT NULL DEFAULT '5',
    `expression`  varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `description` varchar(255) COLLATE utf8_unicode_ci          DEFAULT '',
    `enable`      tinyint(1) unsigned NOT NULL DEFAULT '0',
    `user_id`     int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `product_id` (`product_id`,`name`) USING BTREE,
    KEY           `fk_strategy_product_id` (`product_id`),
    CONSTRAINT `fk_strategy_product_id` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy_bak`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy_bak`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `product_id`  int(10) unsigned NOT NULL,
    `name`        varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `priority`    int(1) unsigned NOT NULL,
    `alarm_count` int(4) unsigned NOT NULL DEFAULT '0',
    `cycle`       int(10) unsigned NOT NULL DEFAULT '5',
    `expression`  varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `description` varchar(255) COLLATE utf8_unicode_ci          DEFAULT '',
    `enable`      tinyint(1) unsigned NOT NULL DEFAULT '0',
    `user_id`     int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `product_id` (`product_id`,`name`) USING BTREE,
    KEY           `fk_strategy_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy_event`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy_event`
(
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `product_id`     int(10) unsigned NOT NULL,
    `strategy_id`    bigint(20) unsigned NOT NULL,
    `strategy_name`  varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `priority`       int(1) unsigned NOT NULL,
    `cycle`          int(10) unsigned NOT NULL,
    `alarm_count`    int(4) unsigned NOT NULL,
    `expression`     varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `create_time`    timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `aware_end_time` timestamp                            NOT NULL DEFAULT '1979-12-31 16:00:00',
    `count`          int(4) unsigned NOT NULL DEFAULT '1',
    `status`         int(1) unsigned NOT NULL DEFAULT '1',
    `host_id`        char(32) COLLATE utf8_unicode_ci     NOT NULL,
    `host_name`      varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `ip`             varchar(16) COLLATE utf8_unicode_ci  NOT NULL,
    PRIMARY KEY (`id`),
    KEY              `idx_host_id_strategy_id_status` (`host_id`,`strategy_id`,`status`),
    KEY              `fk_strategy_event_product_id` (`product_id`),
    CONSTRAINT `fk_strategy_event_product_id` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy_event_failed`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy_event_failed`
(
    `strategy_id` bigint(20) unsigned NOT NULL,
    `host_id`     varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `status`      int(11) NOT NULL,
    `message`     varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `create_time` timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `idx_strategy_id_host_id` (`strategy_id`,`host_id`) USING BTREE,
    KEY           `strategy_id` (`strategy_id`),
    KEY           `stragety_event_failed_host_id_fk` (`host_id`),
    CONSTRAINT `stragety_event_failed_host_id_fk` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `strategy_event_failed_strategy_id_fk` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy_event_process`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy_event_process`
(
    `strategy_event_id` bigint(20) unsigned NOT NULL,
    `process_status`    int(255) DEFAULT NULL,
    `process_user`      varchar(255) DEFAULT NULL,
    `process_comments`  varchar(255) DEFAULT NULL,
    `process_time`      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    KEY                 `strategy_event_id` (`strategy_event_id`),
    CONSTRAINT `fk_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy_event_record`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy_event_record`
(
    `strategy_event_id` bigint(20) unsigned NOT NULL,
    `count`             int(4) unsigned NOT NULL DEFAULT '1',
    `strategy_id`       bigint(20) unsigned NOT NULL,
    `strategy_name`     varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `priority`          int(1) unsigned NOT NULL,
    `cycle`             int(10) unsigned NOT NULL,
    `alarm_count`       int(4) unsigned NOT NULL,
    `expression`        varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `create_time`       timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`       timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `aware_end_time`    timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `status`            int(1) unsigned NOT NULL DEFAULT '1',
    `host_id`           char(32) COLLATE utf8_unicode_ci     NOT NULL,
    `host_name`         varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `ip`                varchar(16) COLLATE utf8_unicode_ci  NOT NULL,
    KEY                 `idx_strategy_event_record_strategy_event_id` (`strategy_event_id`) USING BTREE,
    KEY                 `strategy_event_record_host_id_fk` (`host_id`),
    CONSTRAINT `strategy_event_record_host_id_fk` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy_group`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy_group`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_id` bigint(20) unsigned NOT NULL,
    `group_id`    int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_strategy_id_group_id` (`strategy_id`,`group_id`),
    KEY           `fk_strategy_group_group_id` (`group_id`),
    CONSTRAINT `fk_strategy_group_group_id` FOREIGN KEY (`group_id`) REFERENCES `host_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_strategy_group_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy_host_exclude`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy_host_exclude`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_id` bigint(20) unsigned NOT NULL,
    `host_id`     char(32) COLLATE utf8_unicode_ci NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_strategy_id_host_id` (`strategy_id`,`host_id`),
    KEY           `idx_strategy_host_host_id` (`host_id`),
    CONSTRAINT `fk_strategy_host_host_id` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_strategy_host_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy_template`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy_template`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(255)  NOT NULL,
    `alarm_count` int(4) unsigned NOT NULL,
    `cycle`       int(10) unsigned NOT NULL,
    `expression`  varchar(255)  NOT NULL,
    `description` varchar(1024) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trigger`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trigger`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_id` bigint(20) unsigned NOT NULL,
    `metric`      varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `tags`        varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `number`      int(10) unsigned NOT NULL DEFAULT '0',
    `index`       char(10) COLLATE utf8_unicode_ci     NOT NULL,
    `method`      varchar(10) COLLATE utf8_unicode_ci  NOT NULL,
    `symbol`      char(5) COLLATE utf8_unicode_ci      NOT NULL,
    `threshold`   double(255, 2
) NOT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_strategy_id_index_index` (`strategy_id`,`index`),
  CONSTRAINT `fk_strategy_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `strategy` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trigger_bak`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trigger_bak`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_id` bigint(20) unsigned NOT NULL,
    `metric`      varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `tags`        varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `number`      int(10) unsigned NOT NULL DEFAULT '0',
    `index`       char(10) COLLATE utf8_unicode_ci     NOT NULL,
    `method`      varchar(10) COLLATE utf8_unicode_ci  NOT NULL,
    `symbol`      char(5) COLLATE utf8_unicode_ci      NOT NULL,
    `threshold`   double(255, 2
) NOT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_strategy_id_index_index` (`strategy_id`,`index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trigger_event`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trigger_event`
(
    `strategy_event_id` bigint(20) unsigned NOT NULL,
    `index`             char(10) COLLATE utf8_unicode_ci NOT NULL,
    `metric`            varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `tags`              varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `number`            int(10) unsigned NOT NULL DEFAULT '0',
    `aggregate_tags`    varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `current_threshold` double(255, 2
) NOT NULL,
  `method` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `symbol` char(5) COLLATE utf8_unicode_ci NOT NULL,
  `threshold` double(255,2) NOT NULL,
  `triggered` tinyint(1) NOT NULL DEFAULT '0',
  KEY `fk_trigger_event_strategy_event_id` (`strategy_event_id`),
  CONSTRAINT `fk_trigger_event_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trigger_event_record`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trigger_event_record`
(
    `strategy_event_id` bigint(20) unsigned NOT NULL,
    `count`             int(4) unsigned NOT NULL DEFAULT '1',
    `index`             char(10) COLLATE utf8_unicode_ci NOT NULL,
    `metric`            varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `tags`              varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `number`            int(10) unsigned NOT NULL DEFAULT '0',
    `aggregate_tags`    varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `current_threshold` double(255, 2
) NOT NULL,
  `method` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `symbol` char(5) COLLATE utf8_unicode_ci NOT NULL,
  `threshold` double(255,2) NOT NULL,
  `triggered` tinyint(1) NOT NULL DEFAULT '0',
  KEY `strategy_event_id` (`strategy_event_id`),
  CONSTRAINT `fk_trigger_event_record_strategy_event_id` FOREIGN KEY (`strategy_event_id`) REFERENCES `strategy_event_record` (`strategy_event_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trigger_template`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trigger_template`
(
    `id`                   int(10) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_template_id` int(10) unsigned NOT NULL,
    `metric`               varchar(255) NOT NULL,
    `tags`                 varchar(255) NOT NULL,
    `number`               int(10) NOT NULL,
    `index`                char(10)     NOT NULL,
    `method`               varchar(10)  NOT NULL,
    `symbol`               char(5)      NOT NULL,
    `threshold`            double(255, 2
) NOT NULL,
  `description` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `strategy_template_id` (`strategy_template_id`),
  CONSTRAINT `trigger_template_ibfk_1` FOREIGN KEY (`strategy_template_id`) REFERENCES `strategy_template` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user`
(
    `id`           int(10) unsigned NOT NULL AUTO_INCREMENT,
    `username`     char(255) COLLATE utf8_unicode_ci    NOT NULL,
    `display_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `password`     char(255) COLLATE utf8_unicode_ci    NOT NULL DEFAULT '',
    `role`         tinyint(1) NOT NULL DEFAULT '0',
    `phone`        varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `mail`         varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `wechat`       varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `type`         varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `status`       tinyint(1) NOT NULL DEFAULT '0',
    `create_at`    timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at`    timestamp                            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_group`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_group`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(255) COLLATE utf8_unicode_ci  NOT NULL,
    `product_id`  int(10) unsigned NOT NULL,
    `description` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `creator`     varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`,`product_id`) USING BTREE,
    KEY           `user_group_product_id_fk` (`product_id`),
    CONSTRAINT `user_group_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_group_user`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_group_user`
(
    `id`            int(10) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`       int(10) unsigned NOT NULL,
    `user_group_id` int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id_group_id` (`user_id`,`user_group_id`),
    KEY             `fk_user_usergroup_user_group_id` (`user_group_id`),
    CONSTRAINT `fk_user_user_group_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_user_usergroup_user_group_id` FOREIGN KEY (`user_group_id`) REFERENCES `user_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-07-04 11:18:45
