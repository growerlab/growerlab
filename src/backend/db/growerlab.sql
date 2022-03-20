-- MySQL dump 10.13  Distrib 8.0.26, for Linux (aarch64)
--
-- Host: localhost    Database: growerlab
-- ------------------------------------------------------
-- Server version	8.0.26-0ubuntu0.21.04.3

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `activate_code`
--

DROP TABLE IF EXISTS `activate_code`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `activate_code` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `code` varchar(16) NOT NULL DEFAULT '',
  `created_at` bigint NOT NULL,
  `used_at` bigint NOT NULL,
  `expired_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户激活码';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `namespace`
--

DROP TABLE IF EXISTS `namespace`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `namespace` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '路径',
  `owner_id` int NOT NULL COMMENT '命名空间所有者（用户）',
  `type` tinyint NOT NULL COMMENT '1用户 2组织',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_path` (`path`),
  KEY `unq_owner` (`owner_id`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='命名空间';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `permission`
--

DROP TABLE IF EXISTS `permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `permission` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `namespace_id` int NOT NULL COMMENT '权限所属命名空间（不依赖组织或个人，而是依赖命名空间）',
  `context_type` int NOT NULL COMMENT '上下文类型',
  `context_param_1` int NOT NULL COMMENT '上下文参数1',
  `context_param_2` int NOT NULL COMMENT '上下文参数2',
  `user_domain_type` int NOT NULL COMMENT '用户域类型',
  `user_domain_param` int NOT NULL COMMENT '用户域参数',
  `created_at` int NOT NULL COMMENT '创建时间',
  `deleted_at` int DEFAULT NULL COMMENT '删除时间',
  `code` int NOT NULL COMMENT '权限代号',
  PRIMARY KEY (`id`),
  KEY `idx_ctx` (`code`,`context_type`,`context_param_1`,`context_param_2`),
  KEY `idx_namespace` (`namespace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `repository`
--

DROP TABLE IF EXISTS `repository`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `repository` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '仓库uuid（fork仓库相同）',
  `fork` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否fork项目',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '仓库路径',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '仓库名',
  `namespace_id` int NOT NULL COMMENT '仓库所属命名空间id',
  `owner_id` int NOT NULL COMMENT '仓库创建者,fork后不变',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '仓库描述',
  `created_at` bigint NOT NULL,
  `server_id` int NOT NULL COMMENT '所属服务器id',
  `server_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '在服务器中的物理路径',
  `public` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否公开',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_namespace_path` (`namespace_id`,`path`),
  KEY `idx_server` (`server_id`),
  KEY `idx_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='仓库表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `server`
--

DROP TABLE IF EXISTS `server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `server` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `summary` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '说明备注',
  `host` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '服务器host',
  `port` int NOT NULL DEFAULT '9000',
  `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '服务器状态 0关闭 1正常 2暂停',
  `created_at` bigint NOT NULL,
  `deleted_at` bigint DEFAULT NULL COMMENT '是否被删除',
  PRIMARY KEY (`id`),
  KEY `idx_host` (`host`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `session`
--

DROP TABLE IF EXISTS `session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `session` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `owner_id` int NOT NULL,
  `token` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `created_at` bigint NOT NULL,
  `expired_at` bigint NOT NULL,
  `client_ip` varchar(46) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户当前登录的ip',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_owner` (`owner_id`,`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户邮箱',
  `encrypted_password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户密码',
  `username` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '唯一性用户名（将用在url中）',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户昵称',
  `public_email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '公开的邮箱地址',
  `last_login_ip` varchar(46) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '最后的登录ip（兼容ipv6长度）',
  `last_login_at` int DEFAULT NULL,
  `created_at` bigint NOT NULL,
  `deleted_at` int DEFAULT NULL,
  `verified_at` int DEFAULT NULL,
  `register_ip` varchar(46) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '注册ip',
  `is_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否管理员',
  `namespace_id` int NOT NULL COMMENT '用户的用户域id',
  PRIMARY KEY (`id`),
  KEY `unq_email` (`email`),
  KEY `unq_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-11-30 20:53:49
