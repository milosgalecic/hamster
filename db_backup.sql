-- MySQL dump 10.13  Distrib 8.0.39, for Linux (x86_64)
--
-- Host: localhost    Database: hamster
-- ------------------------------------------------------
-- Server version	8.0.39-0ubuntu0.24.04.2

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
-- Table structure for table `drivers`
--

DROP TABLE IF EXISTS `drivers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `drivers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `license_number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone_number` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
  `average_consumption` float DEFAULT NULL,
  `km_traveled` float NOT NULL DEFAULT '0',
  `active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `drivers`
--

LOCK TABLES `drivers` WRITE;
/*!40000 ALTER TABLE `drivers` DISABLE KEYS */;
INSERT INTO `drivers` VALUES (12,'John Doe','A1234567','123-456-7890','2024-11-13 17:27:22','аvailable',8.5,12000,0),(13,'Jane Smith','B2345678','234-567-8901','2024-11-13 17:27:22','inactive',7.2,15000,0),(14,'Mark Taylor','C3456789','345-678-9012','2024-11-13 17:27:22','active',6.8,8000,1),(15,'Lucy Adams','D4567890','456-789-0123','2024-11-13 17:27:22','active',7.5,6000,1),(16,'Tom Reed','E5678901','567-890-1234','2024-11-13 17:27:22','inactive',8,11000,0),(17,'Anna Green','F6789012','678-901-2345','2024-11-13 17:27:22','active',7,20000,1),(18,'Gary White','G7890123','789-012-3456','2024-11-13 17:27:22','active',6.5,5000,1),(19,'Lily Brown','H8901234','890-123-4567','2024-11-13 17:27:22','inactive',7.8,9500,0),(20,'Charlie Black','I9012345','901-234-5678','2024-11-13 17:27:22','active',6.9,13000,1),(21,'Sarah Blue','J0123456','012-345-6789','2024-11-13 17:27:22','inactive',7.1,12000,0);
/*!40000 ALTER TABLE `drivers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jobs`
--

DROP TABLE IF EXISTS `jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jobs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `description` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'Unknown',
  `driver_id` int DEFAULT NULL,
  `truck_id` int NOT NULL,
  `scheduled_date` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` enum('active','pending','completed','issue','canceled') COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  `distance` float NOT NULL DEFAULT '0',
  `package_size` float NOT NULL DEFAULT '0',
  `scheduled_arrival_time` timestamp NULL DEFAULT NULL,
  `client_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'Klijent nije naveden',
  `start_location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'Nepoznato',
  `destination_location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'Nepoznato',
  `package_weight` float DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_driver` (`driver_id`),
  KEY `fk_truck` (`truck_id`),
  CONSTRAINT `fk_driver` FOREIGN KEY (`driver_id`) REFERENCES `drivers` (`id`),
  CONSTRAINT `fk_truck` FOREIGN KEY (`truck_id`) REFERENCES `trucks` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=80 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jobs`
--

LOCK TABLES `jobs` WRITE;
/*!40000 ALTER TABLE `jobs` DISABLE KEYS */;
INSERT INTO `jobs` VALUES (70,'Delivery of electronics',12,10,'2024-11-15 07:00:00','2024-11-13 17:34:07','pending',200.5,1.5,'2024-11-15 11:00:00','ElectroStore','City A','City B',500),(71,'Furniture transport',13,11,'2024-11-16 08:00:00','2024-11-13 17:34:07','pending',150,2.5,'2024-11-16 12:00:00','HomeDecor','Warehouse A','Showroom B',800),(72,'Grocery delivery',14,12,'2024-11-17 05:30:00','2024-11-13 17:34:07','active',120,0.8,'2024-11-17 09:00:00','FreshMarket','Distribution Center','Retail Store',300),(73,'Machinery parts transport',15,13,'2024-11-18 13:00:00','2024-11-13 17:34:07','completed',250,3.2,'2024-11-18 17:30:00','HeavyWorks','City C','City D',1200),(74,'Medical supplies',16,14,'2024-11-19 10:00:00','2024-11-13 17:34:07','issue',100,0.5,'2024-11-19 13:30:00','MediCare','Pharmacy A','Clinic B',200),(75,'Automobile parts',17,15,'2024-11-20 12:00:00','2024-11-13 17:34:07','active',300,2,'2024-11-20 18:00:00','AutoWorld','Plant X','Service Center',750),(76,'Perishable goods',18,16,'2024-11-21 04:00:00','2024-11-13 17:34:07','pending',180,1,'2024-11-21 08:30:00','FarmFresh','Cold Storage','Supermarket',350),(77,'Chemical supplies',19,17,'2024-11-22 06:30:00','2024-11-13 17:34:07','pending',400,3.5,'2024-11-22 13:00:00','ChemSafe','Factory A','Lab B',1000),(78,'Building materials',20,18,'2024-11-23 09:00:00','2024-11-13 17:34:07','completed',220,5,'2024-11-23 14:30:00','BuildCo','Site 1','Construction Zone',1600),(79,'Textile delivery',21,19,'2024-11-24 07:45:00','2024-11-13 17:34:07','active',90,1.2,'2024-11-24 10:45:00','TextilesHub','Manufacturing Unit','Store',400);
/*!40000 ALTER TABLE `jobs` ENABLE KEYS */;
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
/*!40000 ALTER TABLE `sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `trucks`
--

DROP TABLE IF EXISTS `trucks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `trucks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `model` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `license_plate` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'available',
  `km_traveled` float NOT NULL DEFAULT '0',
  `average_consumption` float NOT NULL DEFAULT '0',
  `active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `trucks`
--

LOCK TABLES `trucks` WRITE;
/*!40000 ALTER TABLE `trucks` DISABLE KEYS */;
INSERT INTO `trucks` VALUES (10,'Volvo FH16','ABC-123','2024-11-13 17:27:29','аvailable',50000,25.5,0),(11,'Scania R500','DEF-456','2024-11-13 17:27:29','maintenance',60000,27,1),(12,'Mercedes-Benz Actros','GHI-789','2024-11-13 17:27:29','available',55000,26.3,1),(13,'MAN TGX','JKL-012','2024-11-13 17:27:29','in use',70000,28,1),(14,'Iveco Stralis','MNO-345','2024-11-13 17:27:29','available',45000,24.8,1),(15,'DAF XF','PQR-678','2024-11-13 17:27:29','in use',48000,26.5,1),(16,'Renault T High','STU-901','2024-11-13 17:27:29','available',62000,27.2,1),(17,'Volvo FMX','VWX-234','2024-11-13 17:27:29','maintenance',51000,25,1),(18,'Scania G450','YZA-567','2024-11-13 17:27:29','available',54000,24.6,1),(19,'Mercedes-Benz Arocs','BCD-890','2024-11-13 17:27:29','in use',67000,28.3,1);
/*!40000 ALTER TABLE `trucks` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-11-13 18:55:34
