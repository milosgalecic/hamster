-- MySQL dump 10.13  Distrib 8.0.39, for Linux (x86_64)
--
-- Host: localhost    Database: hamster
-- ------------------------------------------------------
-- Server version	8.0.39-0ubuntu0.22.04.1

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
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `license_number` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone_number` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
  `average_consumption` float DEFAULT NULL,
  `km_traveled` float NOT NULL DEFAULT '0',
  `active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `drivers`
--

LOCK TABLES `drivers` WRITE;
/*!40000 ALTER TABLE `drivers` DISABLE KEYS */;
INSERT INTO `drivers` VALUES (1,'John Doe','XYZ123456','555-1234','2024-09-30 13:31:20','active',12.5,20000.8,1),(2,'John Doe','DL123456','555-1234','2024-10-01 14:35:01','active',15.2,10000,1),(3,'Jane Smith','DL654321','555-5678','2024-10-01 14:35:01','active',16.5,8500,1),(4,'Mike Johnson','DL987654','555-9876','2024-10-01 14:35:01','on leave',14.8,12000,1),(5,'Alice Williams','DL456789','555-4567','2024-10-01 14:35:01','active',17.3,7500,1),(6,'Bob Brown','DL111222','555-1122','2024-10-01 14:35:01','nije retartded',18,13000,0),(7,'Milos Galecic','333-555','0655368030','2024-10-07 09:28:42','DR',69,69,0),(8,'Natasa','333-555','1232343425','2024-10-07 10:45:48','active',1,111,0),(9,'Natasa Sojic','333-555','1232343425','2024-10-07 10:46:10','active',1,111,0),(10,'Fedor Galecic','1111','1234235','2024-10-07 10:47:30','pending',1234,1234,0);
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
  `description` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'Unknown',
  `driver_id` int DEFAULT NULL,
  `truck_id` int NOT NULL,
  `scheduled_date` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `distance` float NOT NULL DEFAULT '0',
  `package_size` float NOT NULL DEFAULT '0',
  `scheduled_arrival_time` timestamp NULL DEFAULT NULL,
  `client_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'Klijent nije naveden',
  `start_location` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'Nepoznato',
  `destination_location` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'Nepoznato',
  `package_weight` float DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_driver` (`driver_id`),
  KEY `fk_truck` (`truck_id`),
  CONSTRAINT `fk_driver` FOREIGN KEY (`driver_id`) REFERENCES `drivers` (`id`),
  CONSTRAINT `fk_truck` FOREIGN KEY (`truck_id`) REFERENCES `trucks` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jobs`
--

LOCK TABLES `jobs` WRITE;
/*!40000 ALTER TABLE `jobs` DISABLE KEYS */;
INSERT INTO `jobs` VALUES (37,'Delivery of electronics',1,1,'2024-10-08 20:56:00','2024-10-07 13:32:33','active',1001.5,16.2,'2024-10-15 21:56:00','Client A','Belgrade','Sombor',50.5),(38,'Furniture transport',2,2,'2024-10-24 18:49:00','2024-10-07 13:32:33','pending',2,30,'2024-10-16 20:47:00','Client B','Nis','Sremska Mitrovica',75),(39,'Clothing shipment',3,3,'2024-10-07 13:32:33','2024-10-07 13:32:33','pending',150,12.5,'2024-10-08 13:32:33','Client C','Kragujevac','Subotica',40),(40,'Transport of machinery',4,1,'2024-10-07 13:32:33','2024-10-07 13:32:33','pending',500,100,'2024-10-12 13:32:33','Client D','Zrenjanin','Kraljevo',120),(41,'Food delivery',5,4,'2024-10-07 13:32:33','2024-10-07 13:32:33','pending',100,10,'2024-10-08 13:32:33','Client E','Sombor','Belgrade',20),(42,'Transport of construction materials',6,5,'2024-10-07 13:32:33','2024-10-07 13:32:33','pending',400,50,'2024-10-14 13:32:33','Client F','Cacak','Nis',90),(43,'Medical supplies shipment',7,6,'2024-10-07 13:32:33','2024-10-07 13:32:33','pending',120,25,'2024-10-09 13:32:33','Client G','Pancevo','Vrsac',30),(44,'Transport of books',1,7,'2024-10-07 13:32:33','2024-10-07 13:32:33','pending',250,20,'2024-10-11 13:32:33','Client H','Zajecar','Sremska Mitrovica',45),(45,'Delivery of home appliances',2,1,'2024-10-07 13:32:33','2024-10-07 13:32:33','pending',180,15.5,'2024-10-09 13:32:33','Client I','Leskovac','Novi Pazar',60),(46,'Fragile goods transport',3,2,'2024-10-07 13:32:33','2024-10-07 13:32:33','pending',320,40,'2024-10-10 13:32:33','Client J','Pristina','Belgrade',85);
/*!40000 ALTER TABLE `jobs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `trucks`
--

DROP TABLE IF EXISTS `trucks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `trucks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `model` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `license_plate` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'available',
  `km_traveled` float NOT NULL DEFAULT '0',
  `average_consumption` float DEFAULT NULL,
  `active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `trucks`
--

LOCK TABLES `trucks` WRITE;
/*!40000 ALTER TABLE `trucks` DISABLE KEYS */;
INSERT INTO `trucks` VALUES (1,'Ford F-150','Loshmi','2024-09-30 13:29:53','available',15000.5,15.8,0),(2,'Ford F-150','ABC123','2024-10-01 14:34:16','available',10000,15.5,1),(3,'Chevy Silverado','XYZ789','2024-10-01 14:34:16','available',8000,16.2,1),(4,'Ram 1500','DEF456','2024-10-01 14:34:16','in repair',12000,14.8,1),(5,'Toyota Tacoma','GHI012','2024-10-01 14:34:16','available',5000,17.3,1),(6,'Nissan Frontier','JKL345','2024-10-01 14:34:16','in use',3000,18,1),(7,'Masina','555-333','2024-10-07 11:05:29','free',100000,30,0),(8,'Milos','solim','2024-10-08 12:38:47','active',0,0,0),(9,'Dorfe','Registracija','2024-10-08 13:10:20','pending',999,123,0);
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

-- Dump completed on 2024-10-10  9:21:04
