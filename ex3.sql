-- --------------------------------------------------------
-- Host:                         10.132.81.54
-- Server version:               8.0.34 - MySQL Community Server - GPL
-- Server OS:                    Linux
-- HeidiSQL Version:             12.1.0.6537
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Dumping structure for table rte-ex3.alarm_history
CREATE TABLE IF NOT EXISTS `alarm_history` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `timeOn` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `timeOff` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `source` text COLLATE utf8mb4_general_ci NOT NULL,
  `message` text COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5861 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.batch_report
CREATE TABLE IF NOT EXISTS `batch_report` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `batch_id` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_name` text COLLATE utf8mb4_general_ci NOT NULL,
  `shift` text COLLATE utf8mb4_general_ci NOT NULL,
  `recipe` text COLLATE utf8mb4_general_ci NOT NULL,
  `spec_weightscale` text COLLATE utf8mb4_general_ci,
  `spec_weightscale_up` text COLLATE utf8mb4_general_ci,
  `spec_weightscale_low` text COLLATE utf8mb4_general_ci,
  `act_weightscale` text COLLATE utf8mb4_general_ci,
  `spec_lengthskiver` text COLLATE utf8mb4_general_ci,
  `act_lengthskiver` text COLLATE utf8mb4_general_ci,
  `timestamp` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=86359 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.datalog
CREATE TABLE IF NOT EXISTS `datalog` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `datetime` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `recipe1` text COLLATE utf8mb4_general_ci NOT NULL,
  `recipe2` text COLLATE utf8mb4_general_ci NOT NULL,
  `spec_speedextupper` text COLLATE utf8mb4_general_ci,
  `act_speedextupper` text COLLATE utf8mb4_general_ci,
  `act_ampextupper` text COLLATE utf8mb4_general_ci,
  `spec_speedextmiddle` text COLLATE utf8mb4_general_ci,
  `act_speedextmiddle` text COLLATE utf8mb4_general_ci,
  `act_ampextmiddle` text COLLATE utf8mb4_general_ci,
  `spec_speedextlower` text COLLATE utf8mb4_general_ci,
  `act_ampextlower` text COLLATE utf8mb4_general_ci,
  `act_speedextlower` text COLLATE utf8mb4_general_ci,
  `spec_speedline` text COLLATE utf8mb4_general_ci,
  `act_speedline` text COLLATE utf8mb4_general_ci,
  `spec_speedcalender` text COLLATE utf8mb4_general_ci,
  `act_speedcalender` text COLLATE utf8mb4_general_ci,
  `spec_runscale` text COLLATE utf8mb4_general_ci,
  `spec_runscale_up` text COLLATE utf8mb4_general_ci,
  `spec_runscale_low` text COLLATE utf8mb4_general_ci,
  `act_runscale` text COLLATE utf8mb4_general_ci,
  `act_runscale_up` text COLLATE utf8mb4_general_ci,
  `act_runscale_low` text COLLATE utf8mb4_general_ci,
  `spec_weightscale` text COLLATE utf8mb4_general_ci,
  `spec_weightscale_up` text COLLATE utf8mb4_general_ci,
  `spec_weightscale_low` text COLLATE utf8mb4_general_ci,
  `act_weightscale` text COLLATE utf8mb4_general_ci,
  `spec_tcuuuper_1` text COLLATE utf8mb4_general_ci,
  `act_tcuuuper_1` text COLLATE utf8mb4_general_ci,
  `spec_tcuuuper_2` text COLLATE utf8mb4_general_ci,
  `act_tcuuuper_2` text COLLATE utf8mb4_general_ci,
  `spec_tcuuuper_3` text COLLATE utf8mb4_general_ci,
  `act_tcuuuper_3` text COLLATE utf8mb4_general_ci,
  `spec_tcuuuper_4` text COLLATE utf8mb4_general_ci,
  `act_tcuuuper_4` text COLLATE utf8mb4_general_ci,
  `spec_tcumiddle_1` text COLLATE utf8mb4_general_ci,
  `act_tcumiddle_1` text COLLATE utf8mb4_general_ci,
  `spec_tcumiddle_2` text COLLATE utf8mb4_general_ci,
  `act_tcumiddle_2` text COLLATE utf8mb4_general_ci,
  `spec_tcumiddle_3` text COLLATE utf8mb4_general_ci,
  `act_tcumiddle_3` text COLLATE utf8mb4_general_ci,
  `spec_tcumiddle_4` text COLLATE utf8mb4_general_ci,
  `act_tcumiddle_4` text COLLATE utf8mb4_general_ci,
  `spec_tculower_1` text COLLATE utf8mb4_general_ci,
  `act_tculower_1` text COLLATE utf8mb4_general_ci,
  `spec_tculower_2` text COLLATE utf8mb4_general_ci,
  `act_tculower_2` text COLLATE utf8mb4_general_ci,
  `spec_tculower_3` text COLLATE utf8mb4_general_ci,
  `act_tculower_3` text COLLATE utf8mb4_general_ci,
  `spec_tculower_4` text COLLATE utf8mb4_general_ci,
  `act_tculower_4` text COLLATE utf8mb4_general_ci,
  `spec_tcupreformer_up` text COLLATE utf8mb4_general_ci,
  `act_tcupreformer_up` text COLLATE utf8mb4_general_ci,
  `spec_tcupreformer_down` text COLLATE utf8mb4_general_ci,
  `act_tcupreformer_down` text COLLATE utf8mb4_general_ci,
  `spec_tcucalender_up` text COLLATE utf8mb4_general_ci,
  `act_tcucalender_up` text COLLATE utf8mb4_general_ci,
  `spec_tcucalender_down` text COLLATE utf8mb4_general_ci,
  `act_tcucalender_down` text COLLATE utf8mb4_general_ci,
  `spec_tcucalenderext_1` text COLLATE utf8mb4_general_ci,
  `act_tcucalenderext_1` text COLLATE utf8mb4_general_ci,
  `spec_tcucalenderext_2` text COLLATE utf8mb4_general_ci,
  `act_tcucalenderext_2` text COLLATE utf8mb4_general_ci,
  `spec_tcucalenderext_3` text COLLATE utf8mb4_general_ci,
  `act_tcucalenderext_3` text COLLATE utf8mb4_general_ci,
  `spec_tcucalenderext_4` text COLLATE utf8mb4_general_ci,
  `act_tcucalenderext_4` text COLLATE utf8mb4_general_ci,
  `spec_gapcalender` text COLLATE utf8mb4_general_ci,
  `act_gapcalender` text COLLATE utf8mb4_general_ci,
  `spec_cuttercalender` text COLLATE utf8mb4_general_ci,
  `act_cuttercalender` text COLLATE utf8mb4_general_ci,
  `act_compoundcalender` text COLLATE utf8mb4_general_ci,
  `spec_lengthskiver` text COLLATE utf8mb4_general_ci,
  `act_lengthskiver` text COLLATE utf8mb4_general_ci,
  `spec_width` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `spec_width_up` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `spec_width_upline` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `spec_width_low` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `spec_width_lowline` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `act_width` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3918561 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.material
CREATE TABLE IF NOT EXISTS `material` (
  `id` int NOT NULL AUTO_INCREMENT,
  `compound` text COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_compound` (`compound`(255))
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.order_report
CREATE TABLE IF NOT EXISTS `order_report` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `recipe_id` int NOT NULL,
  `recipe` text COLLATE utf8mb4_general_ci NOT NULL,
  `set_qty` text COLLATE utf8mb4_general_ci NOT NULL,
  `act_qty` text COLLATE utf8mb4_general_ci,
  `shift` text COLLATE utf8mb4_general_ci NOT NULL,
  `user_name` text COLLATE utf8mb4_general_ci NOT NULL,
  `time_create` timestamp(3) NULL DEFAULT NULL,
  `time_start` timestamp(3) NULL DEFAULT NULL,
  `time_finish` timestamp(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=792 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.recipe1
CREATE TABLE IF NOT EXISTS `recipe1` (
  `id` int NOT NULL AUTO_INCREMENT,
  `size` text COLLATE utf8mb4_general_ci NOT NULL,
  `speed_extupper` text COLLATE utf8mb4_general_ci,
  `speed_extupper_low` text COLLATE utf8mb4_general_ci,
  `speed_extmiddle` text COLLATE utf8mb4_general_ci,
  `speed_extmiddle_low` text COLLATE utf8mb4_general_ci,
  `speed_extlower` text COLLATE utf8mb4_general_ci,
  `speed_extlower_low` text COLLATE utf8mb4_general_ci,
  `speed_line` text COLLATE utf8mb4_general_ci,
  `run_scale` text COLLATE utf8mb4_general_ci,
  `run_scale_up` text COLLATE utf8mb4_general_ci,
  `run_scale_low` text COLLATE utf8mb4_general_ci,
  `run_scale_act_up` text COLLATE utf8mb4_general_ci,
  `run_scale_act_low` text COLLATE utf8mb4_general_ci,
  `tcu_upper_screw` text COLLATE utf8mb4_general_ci,
  `tcu_upper_barrel1` text COLLATE utf8mb4_general_ci,
  `tcu_upper_barrel2` text COLLATE utf8mb4_general_ci,
  `tcu_upper_head` text COLLATE utf8mb4_general_ci,
  `tcu_upper_midhead` text COLLATE utf8mb4_general_ci,
  `tcu_middle_screw` text COLLATE utf8mb4_general_ci,
  `tcu_middle_barrel2` text COLLATE utf8mb4_general_ci,
  `tcu_middle_head` text COLLATE utf8mb4_general_ci,
  `tcu_lower_screw` text COLLATE utf8mb4_general_ci,
  `tcu_preformerup` text COLLATE utf8mb4_general_ci,
  `tcu_preformerdown` text COLLATE utf8mb4_general_ci,
  `tcu_upper_barrel3` text COLLATE utf8mb4_general_ci,
  `tcu_middle_barrel1` text COLLATE utf8mb4_general_ci,
  `tcu_lower_barrel1` text COLLATE utf8mb4_general_ci,
  `tcu_lower_head` text COLLATE utf8mb4_general_ci,
  `tcu_calender_extscrew` text COLLATE utf8mb4_general_ci,
  `tcu_calender_ext1` text COLLATE utf8mb4_general_ci,
  `tcu_calender_ext2` text COLLATE utf8mb4_general_ci,
  `tcu_calender_exthead` text COLLATE utf8mb4_general_ci,
  `tcu_calender_rollupper` text COLLATE utf8mb4_general_ci,
  `tcu_calender_rolllower` text COLLATE utf8mb4_general_ci,
  `speed_calender` text COLLATE utf8mb4_general_ci,
  `gap_calender` text COLLATE utf8mb4_general_ci,
  `cutter_calender` text COLLATE utf8mb4_general_ci,
  `material_extupper` text COLLATE utf8mb4_general_ci,
  `material_extmiddle` text COLLATE utf8mb4_general_ci,
  `material_extlower` text COLLATE utf8mb4_general_ci,
  `material_extcalender` text COLLATE utf8mb4_general_ci,
  `material_rubbersheet` text COLLATE utf8mb4_general_ci,
  `material_extdie` text COLLATE utf8mb4_general_ci,
  `weight_scale` text COLLATE utf8mb4_general_ci,
  `length_skiver` text COLLATE utf8mb4_general_ci,
  `weight_up` text COLLATE utf8mb4_general_ci,
  `weight_low` text COLLATE utf8mb4_general_ci,
  `width` text COLLATE utf8mb4_general_ci,
  `width_up` text COLLATE utf8mb4_general_ci,
  `width_low` text COLLATE utf8mb4_general_ci,
  `width_upline` text COLLATE utf8mb4_general_ci,
  `width_lowline` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=106 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.recipe1queue
CREATE TABLE IF NOT EXISTS `recipe1queue` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `recipe_id` int NOT NULL,
  `size` text COLLATE utf8mb4_general_ci NOT NULL,
  `qty` text COLLATE utf8mb4_general_ci NOT NULL,
  `shift` text COLLATE utf8mb4_general_ci NOT NULL,
  `user_name` text COLLATE utf8mb4_general_ci NOT NULL,
  `time_create` timestamp(3) NULL DEFAULT NULL,
  `time_start` timestamp(3) NULL DEFAULT NULL,
  `time_finish` timestamp(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_recipe2queue_recipe` (`recipe_id`) USING BTREE,
  KEY `fk_recipe1queue_recipe` (`recipe_id`) USING BTREE,
  CONSTRAINT `fk_recipe1queue_recipe1` FOREIGN KEY (`recipe_id`) REFERENCES `recipe1` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=788 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.recipe_history
CREATE TABLE IF NOT EXISTS `recipe_history` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `recipe_id` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` int NOT NULL,
  `modified_at` datetime DEFAULT NULL,
  `modified_by` int DEFAULT NULL,
  `revision` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `note` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`),
  KEY `fk_recipe_history_recipe` (`recipe_id`),
  KEY `fk_recipe_history_created_by` (`created_by`),
  KEY `fk_recipe_history_modified_by` (`modified_by`),
  CONSTRAINT `fk_recipe_history_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_recipe_history_modified_by` FOREIGN KEY (`modified_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_recipe_history_recipe` FOREIGN KEY (`recipe_id`) REFERENCES `recipe1` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.recipe_material
CREATE TABLE IF NOT EXISTS `recipe_material` (
  `id` int NOT NULL AUTO_INCREMENT,
  `recipe_id` int NOT NULL,
  `material_id` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `recipe_id` (`recipe_id`),
  KEY `material_id` (`material_id`),
  CONSTRAINT `recipe_material_ibfk_1` FOREIGN KEY (`recipe_id`) REFERENCES `recipe1` (`id`),
  CONSTRAINT `recipe_material_ibfk_2` FOREIGN KEY (`material_id`) REFERENCES `material` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

-- Dumping structure for table rte-ex3.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_name` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_password` varchar(200) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_level` enum('admin','eng','tech','prod') COLLATE utf8mb4_general_ci DEFAULT NULL,
  `timestamp` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Data exporting was unselected.

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
