-- MySQL dump 10.13  Distrib 8.0.36, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: asap
-- ------------------------------------------------------
-- Server version	8.0.36

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `company`
--

DROP TABLE IF EXISTS `company`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `company` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `owner` int NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`),
  KEY `owner` (`owner`),
  CONSTRAINT `company_ibfk_1` FOREIGN KEY (`owner`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `company`
--

LOCK TABLES `company` WRITE;
/*!40000 ALTER TABLE `company` DISABLE KEYS */;
INSERT INTO `company` VALUES (1,'1',1,'2024-11-27 16:04:54'),(3,'2',1,'2024-11-27 16:06:30'),(7,'3',1,'2024-11-27 16:09:41'),(8,'6',1,'2024-11-27 16:13:27'),(15,'4',1,'2024-11-27 16:56:26'),(18,'43',3,'2024-11-28 10:09:02'),(19,'56',3,'2024-11-28 10:21:33'),(20,'123',1,'2024-11-28 12:35:46'),(23,'1234',1,'2024-12-03 09:37:28'),(29,'12345',1,'2024-12-04 13:57:16'),(30,'123456',1,'2024-12-04 13:57:20'),(31,'1234567',1,'2024-12-04 13:57:23'),(33,'–π—Ü—É–∫',1,'2024-12-04 13:57:28'),(34,'—Ñ—ã–≤–∞',1,'2024-12-04 13:57:30'),(38,'q',1,'2024-12-04 15:12:33');
/*!40000 ALTER TABLE `company` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invites`
--

DROP TABLE IF EXISTS `invites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `invites` (
  `userID` int NOT NULL,
  `companyID` int NOT NULL,
  `AdditionalInfo` varchar(255) DEFAULT NULL,
  KEY `userID` (`userID`),
  KEY `companyID` (`companyID`),
  CONSTRAINT `invites_ibfk_1` FOREIGN KEY (`userID`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `invites_ibfk_2` FOREIGN KEY (`companyID`) REFERENCES `company` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invites`
--

LOCK TABLES `invites` WRITE;
/*!40000 ALTER TABLE `invites` DISABLE KEYS */;
/*!40000 ALTER TABLE `invites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `projects`
--

DROP TABLE IF EXISTS `projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `projects` (
  `name` varchar(255) NOT NULL,
  `leaderID` int NOT NULL,
  `companyID` int NOT NULL,
  `created` datetime NOT NULL,
  `status` tinyint(1) NOT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  KEY `leaderID` (`leaderID`),
  KEY `companyID` (`companyID`),
  CONSTRAINT `projects_ibfk_1` FOREIGN KEY (`leaderID`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `projects_ibfk_2` FOREIGN KEY (`companyID`) REFERENCES `company` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projects`
--

LOCK TABLES `projects` WRITE;
/*!40000 ALTER TABLE `projects` DISABLE KEYS */;
INSERT INTO `projects` VALUES ('123',1,20,'2024-11-28 13:05:26',0,1),('1',1,15,'2024-11-28 13:50:42',0,2),('2',1,15,'2024-11-28 13:50:44',0,3),('3',1,15,'2024-11-28 13:50:45',0,4),('4',1,7,'2024-11-28 14:29:11',0,5),('1',3,19,'2024-11-29 10:03:22',0,6),('123',1,23,'2024-12-03 09:37:45',0,7),('1',1,7,'2024-12-11 19:11:07',0,9),('6',1,7,'2024-12-11 19:12:02',0,10),('123',4,7,'2024-12-11 21:17:57',1,11),('1234',4,7,'2024-12-11 21:18:04',1,12),('5',1,15,'2024-12-12 16:05:14',0,13),('4',1,20,'2024-12-12 16:06:27',0,14);
/*!40000 ALTER TABLE `projects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sessions`
--

DROP TABLE IF EXISTS `sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sessions` (
  `token` char(43) NOT NULL,
  `data` blob NOT NULL,
  `expiry` timestamp(6) NOT NULL,
  PRIMARY KEY (`token`),
  KEY `sessions_expiry_idx` (`expiry`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sessions`
--

LOCK TABLES `sessions` WRITE;
/*!40000 ALTER TABLE `sessions` DISABLE KEYS */;
INSERT INTO `sessions` VALUES ('0R6M2ECVLZK0u2txx7BL9PNBBEBVhbiGyKNBHhAAeko',_binary '%ˇÄ\0DeadlineˇÇ\0ValuesˇÑ\0\0\0ˇÅTimeˇÇ\0\0\0\'ˇÉmap[string]interface {}ˇÑ\0\0\02ˇÄ\0\0\0\ﬁ\Ìbn\Ê,ˇˇauthenticatedUserIDint\0\0','2024-12-12 20:00:45.024045'),('IqOfZNA0TDx1NqxKcXO9COnK0LTPoca7vyJ4CVUDGsE',_binary '%ˇÄ\0DeadlineˇÇ\0ValuesˇÑ\0\0\0ˇÅTimeˇÇ\0\0\0\'ˇÉmap[string]interface {}ˇÑ\0\0\02ˇÄ\0\0\0\ﬁ\Ì\€\ÌGaˇˇauthenticatedUserIDint\0\0','2024-12-13 04:40:29.424108');
/*!40000 ALTER TABLE `sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tasks`
--

DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tasks` (
  `name` varchar(255) NOT NULL,
  `category` enum('Urgent','High Priority','Low Priority') NOT NULL,
  `created` datetime NOT NULL,
  `expired` datetime NOT NULL,
  `isDone` tinyint(1) NOT NULL,
  `whocomplete` varchar(255) DEFAULT NULL,
  `projID` int NOT NULL,
  `companyID` int NOT NULL,
  KEY `projID` (`projID`),
  KEY `companyID` (`companyID`),
  CONSTRAINT `tasks_ibfk_1` FOREIGN KEY (`projID`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
  CONSTRAINT `tasks_ibfk_2` FOREIGN KEY (`companyID`) REFERENCES `company` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tasks`
--

LOCK TABLES `tasks` WRITE;
/*!40000 ALTER TABLE `tasks` DISABLE KEYS */;
INSERT INTO `tasks` VALUES ('2','High Priority','2024-12-03 10:10:00','2024-12-05 13:07:00',1,'3',7,23),('3','Urgent','2024-12-03 10:39:53','2024-12-10 16:30:00',0,NULL,7,23),('6','Low Priority','2024-12-03 10:43:27','2024-12-03 15:45:00',0,NULL,7,23),('23','Low Priority','2024-12-03 10:49:53','2024-12-04 15:52:00',0,NULL,7,23),('2','Urgent','2024-12-03 12:19:14','2024-12-03 17:19:00',1,'2',5,7),('1','Urgent','2024-12-04 13:38:45','2024-12-04 19:40:00',0,NULL,3,15),('2','Urgent','2024-12-04 13:44:54','2024-12-04 19:46:00',1,'2',2,15),('123','High Priority','2024-12-12 13:36:31','2024-12-12 16:36:00',1,'2',5,7),('567','Low Priority','2024-12-12 13:36:45','2024-12-12 16:36:00',1,'2',5,7),('5','High Priority','2024-12-12 13:39:20','2024-11-29 18:42:00',1,'2',5,7),('67','High Priority','2024-12-12 13:39:36','2024-12-07 19:42:00',0,NULL,5,7),('435','High Priority','2024-12-12 13:40:15','2024-12-12 19:43:00',0,NULL,5,7),('45','High Priority','2024-12-12 13:41:19','2024-11-30 19:44:00',0,NULL,5,7),('123','High Priority','2024-12-12 19:03:37','2024-11-30 00:05:00',1,'qwerty123',1,20);
/*!40000 ALTER TABLE `tasks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `hashed_password` char(60) NOT NULL,
  `created` datetime NOT NULL,
  `role` enum('admin','user') DEFAULT 'user',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2','q@q.q','$2a$12$oA9kLI9qWY7Hic6QxSPS0ukz3yz6QmKZ7F/NP3.z3rixrbql5tTni','2024-11-27 15:38:55','admin'),(2,'3','2@1.1','$2a$12$0mlo/dVVzxtuErtXcWd.pe5tBe1NGP73G56iIh/kGIbqx4jSlG6Cm','2024-11-27 15:42:35','user'),(3,'qwerty','bob@example.com','$2a$12$d49QV6uO0Qd.um0qWkrhhezB8o2tjY1FdJE3p/B9dxynHIcDwef/a','2024-11-28 10:08:47','user'),(4,'123','123@123.s','$2a$12$4Z/CzDmBxINoWKF.aUxx1OaYOOeUYkgDJ3zwjrENljHNbzmTBiMTu','2024-11-28 12:28:53','user'),(5,'1234','11@1.1','$2a$12$Rd4SPz6MF16YU5drGU6wgeDmGHanMkiz.thZ1tedIkso5odh8MCNa','2024-11-28 12:30:40','user');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `usersincompanies`
--

DROP TABLE IF EXISTS `usersincompanies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `usersincompanies` (
  `user` int NOT NULL,
  `company` int NOT NULL,
  `role` varchar(255) NOT NULL,
  KEY `user` (`user`),
  KEY `company` (`company`),
  CONSTRAINT `usersincompanies_ibfk_1` FOREIGN KEY (`user`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `usersincompanies_ibfk_2` FOREIGN KEY (`company`) REFERENCES `company` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usersincompanies`
--

LOCK TABLES `usersincompanies` WRITE;
/*!40000 ALTER TABLE `usersincompanies` DISABLE KEYS */;
INSERT INTO `usersincompanies` VALUES (1,7,'1'),(3,18,'owner'),(3,19,'owner'),(1,20,'owner'),(1,23,'owner'),(1,18,'worker'),(1,19,'worker'),(1,29,'owner'),(1,30,'owner'),(1,31,'owner'),(1,33,'owner'),(1,34,'owner'),(1,38,'owner');
/*!40000 ALTER TABLE `usersincompanies` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-12-13  1:11:58
