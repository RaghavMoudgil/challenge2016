Basic Instrcution 
to run the code create 3 tables 
----------------------------------------------- Create Table Queries------------------------------------
CREATE TABLE `cities` (
  `City Code` varchar(50) NOT NULL,
  `Province Code` varchar(10) NOT NULL,
  `Country Code` varchar(10) NOT NULL,
  `City Name` varchar(100) NOT NULL,
  `Province name` varchar(100) NOT NULL,
  `Country Name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


CREATE TABLE `distributors` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `parent_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `parent_id` (`parent_id`),
  CONSTRAINT `distributors_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `distributors` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


CREATE TABLE `distributor_permissions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `distributor_id` int NOT NULL,
  `permission_type` enum('INCLUDE','EXCLUDE') NOT NULL,
  `region_code` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `distributor_id` (`distributor_id`),
  CONSTRAINT `distributor_permissions_ibfk_1` FOREIGN KEY (`distributor_id`) REFERENCES `distributors` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

----------------------------------------------- END Create Table Queries------------------------------------
----------------------------------------------Use this query to insert DUMMY DATA-----------------------------
insert into `distributors` (`id`, `name`, `parent_id`) values('4','DISTRIBUTOR1',NULL);
insert into `distributors` (`id`, `name`, `parent_id`) values('5','DISTRIBUTOR2','4');
insert into `distributors` (`id`, `name`, `parent_id`) values('6','DISTRIBUTOR3','5');



insert into `distributor_permissions` (`id`, `distributor_id`, `permission_type`, `region_code`) values('7','4','INCLUDE','CHICAGO-ILLINOIS-UNITEDSTATES');
insert into `distributor_permissions` (`id`, `distributor_id`, `permission_type`, `region_code`) values('9','4','INCLUDE','CHENNAI-TAMILNADU-INDIA');
insert into `distributor_permissions` (`id`, `distributor_id`, `permission_type`, `region_code`) values('10','4','INCLUDE','BANGALORE-KARNATAKA-INDIA');
insert into `distributor_permissions` (`id`, `distributor_id`, `permission_type`, `region_code`) values('11','5','INCLUDE','BANGALORE-KARNATAKA-INDIA');
insert into `distributor_permissions` (`id`, `distributor_id`, `permission_type`, `region_code`) values('12','5','INCLUDE','MUMBAI-MAHARASHTRA-INDIA');
insert into `distributor_permissions` (`id`, `distributor_id`, `permission_type`, `region_code`) values('13','6','INCLUDE','HUBLI-KARNATAKA-INDIA');




----------------------------------------------END---------------------------------------------------




Project Overview

This project uses MySQL to manage city data, distributors, and their permissions.

Database Tables
cities
Stores information about cities including:

City Code (primary key)

Province Code and Name

Country Code and Name

distributors
Stores distributor details and their hierarchy:

Each distributor has a unique ID and name.

A distributor can have a parent distributor (parent_id) forming a hierarchy.

If a parent distributor is deleted, the child’s parent is set to NULL.

distributor_permissions
Stores permissions for distributors by region:

Links to distributors via distributor_id.

Permissions are of type INCLUDE or EXCLUDE.

region_code defines which region the permission applies to.( it is the combination of CITY-PROVINCE-COUNTRY)

When a distributor is deleted, all related permissions are removed automatically.

Relationships and Flow
Distributors are organized in a parent-child hierarchy.

Permissions define where distributors can operate based on included or excluded regions.

City data provides the regions (provinces, countries, and cities) related to these permissions.