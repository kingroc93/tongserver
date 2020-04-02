# noinspection SqlNoDataSourceInspectionForFile

-- MySQL dump 10.13  Distrib 5.7.29, for Linux (x86_64)
--
-- Host: localhost    Database: idb
-- ------------------------------------------------------
-- Server version	5.7.29-0ubuntu0.16.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `G_DATABASEURL`
--

DROP TABLE IF EXISTS `G_DATABASEURL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `G_DATABASEURL` (
  `ID` varchar(50) COLLATE utf8_bin NOT NULL,
  `DBTYPE` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `DBURL` varchar(500) COLLATE utf8_bin DEFAULT NULL,
  `USERNAME` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `PWD` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `PROJECTID` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `DBALIAS` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `G_IDS`
--

DROP TABLE IF EXISTS `G_IDS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `G_IDS` (
  `ID` varchar(50) COLLATE utf8_bin NOT NULL,
  `META` varchar(4500) COLLATE utf8_bin DEFAULT NULL,
  `PROJECTID` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `INF` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `NAME` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `DBALIAS` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `G_META`
--

DROP TABLE IF EXISTS `G_META`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `G_META` (
  `ID` varchar(45) COLLATE utf8_bin NOT NULL,
  `PROJECTID` varchar(145) COLLATE utf8_bin DEFAULT NULL,
  `NAMESPACE` varchar(145) COLLATE utf8_bin DEFAULT NULL,
  `METANAME` varchar(145) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `G_META_ITEM`
--

DROP TABLE IF EXISTS `G_META_ITEM`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `G_META_ITEM` (
  `ID` varchar(45) COLLATE utf8_bin NOT NULL,
  `META_ID` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `NAME` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `VALUE` varchar(8000) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `G_PROJECT`
--

DROP TABLE IF EXISTS `G_PROJECT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `G_PROJECT` (
  `ID` varchar(50) COLLATE utf8_bin NOT NULL,
  `PROJECTNAME` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `OWNER` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `G_SERVICE`
--

DROP TABLE IF EXISTS `G_SERVICE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `G_SERVICE` (
  `ID` varchar(50) COLLATE utf8_bin NOT NULL,
  `BODYTYPE` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `SERVICETYPE` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `NAMESPACE` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `ENABLED` int(11) DEFAULT NULL,
  `MSGLOG` int(11) DEFAULT NULL,
  `SECURITY` int(11) DEFAULT NULL,
  `META` varchar(4000) COLLATE utf8_bin DEFAULT NULL,
  `PROJECTID` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  `CONTEXT` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `G_USERPROJECT`
--

DROP TABLE IF EXISTS `G_USERPROJECT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `G_USERPROJECT` (
  `USERID` varchar(50) COLLATE utf8_bin NOT NULL,
  `PROJECTID` varchar(50) COLLATE utf8_bin NOT NULL,
  `PROJECTNAME` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`USERID`,`PROJECTID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `G_USERSERVICE`
--

DROP TABLE IF EXISTS `G_USERSERVICE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `G_USERSERVICE` (
  `ROLEID` varchar(50) CHARACTER SET utf8 NOT NULL,
  `SERVICEID` varchar(50) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`ROLEID`,`SERVICEID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `JEDA_MENU`
--

DROP TABLE IF EXISTS `JEDA_MENU`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `JEDA_MENU` (
  `MENU_ID` varchar(50) NOT NULL,
  `PARENT_MENU_ID` varchar(50) DEFAULT NULL,
  `MENU_NAME` varchar(100) DEFAULT NULL,
  `MENU_URL` varchar(500) DEFAULT NULL,
  `MENU_DESCRIPTION` varchar(1000) DEFAULT NULL,
  `MENU_IFRAME` int(11) DEFAULT NULL,
  `MENU_ICON` varchar(50) DEFAULT NULL,
  `MENU_ORDER` int(11) DEFAULT '0',
  `MENU_READ_ONLY` int(11) DEFAULT NULL,
  `MENU_OPEN_IN_HOME` int(11) DEFAULT '0',
  `MENU_VERSION` int(11) DEFAULT '0',
  `MENU_CREATOR` varchar(50) DEFAULT NULL,
  `MENU_CREATED` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `MENU_MODIFIER` varchar(50) DEFAULT NULL,
  `MENU_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`MENU_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `JEDA_ORG`
--

DROP TABLE IF EXISTS `JEDA_ORG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `JEDA_ORG` (
  `ORG_ID` varchar(50) CHARACTER SET utf8 NOT NULL,
  `ORG_NAME` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `PARENT_ORG_ID` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `ORG_DESCRIPTION` varchar(500) COLLATE utf8_bin DEFAULT NULL,
  `ORG_TEL` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `ORG_ADDRESS` varchar(500) COLLATE utf8_bin DEFAULT NULL,
  `ORG_CONTACT` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `ORG_PATH` varchar(500) COLLATE utf8_bin DEFAULT NULL,
  `ORG_LEVEL` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `ORG_ENABLED` int(11) DEFAULT '1',
  `ORG_TYPE` varchar(1) COLLATE utf8_bin DEFAULT NULL,
  `ORG_PROPERTY` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `ORG_ORDER` int(11) DEFAULT '0',
  `ORG_VERSION` int(11) DEFAULT '0',
  `ORG_CREATOR` varchar(50) COLLATE utf8_bin DEFAULT 'admin',
  `ORG_CREATED` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `ORG_MODIFIER` varchar(50) COLLATE utf8_bin DEFAULT 'admin',
  `ORG_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ORG_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `JEDA_ROLE`
--

DROP TABLE IF EXISTS `JEDA_ROLE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `JEDA_ROLE` (
  `ROLE_ID` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `ROLE_NAME` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `ROLE_DESCRIPTION` varchar(500) COLLATE utf8_bin DEFAULT NULL,
  `ROLE_TYPE` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `ROLE_ORDER` int(11) DEFAULT NULL,
  `ROLE_READ_ONLY` int(11) DEFAULT NULL,
  `ROLE_VERSION` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `JEDA_ROLE_USER`
--

DROP TABLE IF EXISTS `JEDA_ROLE_USER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `JEDA_ROLE_USER` (
  `USER_ID` varchar(50) DEFAULT NULL,
  `ROLE_ID` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `JEDA_USER`
--

DROP TABLE IF EXISTS `JEDA_USER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `JEDA_USER` (
  `USER_ID` varchar(50) NOT NULL,
  `POSITION_ID` varchar(50) DEFAULT NULL,
  `ORG_ID` varchar(50) DEFAULT NULL,
  `USER_NAME` varchar(100) DEFAULT NULL,
  `USER_PASSWORD` varchar(100) DEFAULT NULL,
  `USER_ID_NO` varchar(50) DEFAULT NULL,
  `USER_GENDER` varchar(8) DEFAULT NULL,
  `USER_EMAIL` varchar(100) DEFAULT NULL,
  `USER_BIRTHDAY` date DEFAULT NULL,
  `USER_ADDRESS` varchar(500) DEFAULT NULL,
  `USER_POST` varchar(50) DEFAULT NULL,
  `USER_TEL` varchar(50) DEFAULT NULL,
  `USER_MOBILE` varchar(50) DEFAULT NULL,
  `USER_DESCRIPTION` varchar(500) DEFAULT NULL,
  `USER_ENABLED` int(10) DEFAULT NULL,
  `USER_LOCKED` int(11) DEFAULT NULL,
  `USER_ACCOUNT_NONEXPIRED` int(11) DEFAULT NULL,
  `USER_ACCOUNT_NONLOCKED` int(11) DEFAULT NULL,
  `USER_CREDENTIALS_NONEXPIRED` int(11) DEFAULT NULL,
  `USER_ORDER` int(11) DEFAULT NULL,
  `USER_VERSION` int(11) DEFAULT NULL,
  `USER_CREATOR` varchar(50) DEFAULT NULL,
  `USER_CREATED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `USER_MODIFIER` varchar(50) DEFAULT NULL,
  `USER_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `ADDVCD` varchar(500) DEFAULT NULL,
  `VISITS` int(11) DEFAULT '0',
  `LOGIN_TIME` varchar(50) DEFAULT NULL,
  `LOGIN_NAME` varchar(60) DEFAULT NULL,
  PRIMARY KEY (`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-04-01 23:26:41
