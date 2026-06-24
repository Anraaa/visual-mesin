-- --------------------------------------------------------
-- Host:                         10.129.78.199
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

-- Dumping structure for table datamesin.bpbl
CREATE TABLE IF NOT EXISTS `bpbl` (
  `recid` int NOT NULL AUTO_INCREMENT,
  `bc_entried` char(50) DEFAULT NULL,
  `mcncode` char(50) DEFAULT NULL,
  `probcode` char(50) DEFAULT NULL,
  `jdge` char(50) DEFAULT NULL,
  `opr` char(50) DEFAULT NULL,
  `oprname` char(50) DEFAULT NULL,
  `jdge_date` char(50) DEFAULT NULL,
  `date_shift` char(50) DEFAULT NULL,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.curtire
CREATE TABLE IF NOT EXISTS `curtire` (
  `bcd` varchar(50) DEFAULT NULL,
  `itemcode` varchar(50) DEFAULT NULL,
  `codespec` varchar(50) DEFAULT NULL,
  `mcn` varchar(50) DEFAULT NULL,
  `mcntype` varchar(50) DEFAULT NULL,
  `operator` varchar(50) DEFAULT NULL,
  `curein` varchar(50) DEFAULT NULL,
  `cureout` varchar(50) DEFAULT NULL,
  `dateshift` varchar(50) DEFAULT NULL,
  `grpshift` varchar(50) DEFAULT NULL,
  `rsn` varchar(50) DEFAULT NULL,
  `nip` varchar(50) DEFAULT NULL,
  `extemp` varchar(50) DEFAULT NULL,
  `extemp_jdg` varchar(50) DEFAULT NULL,
  `intemp` varchar(50) DEFAULT NULL,
  `intemp_jdg` varchar(50) DEFAULT NULL,
  `platen` varchar(50) DEFAULT NULL,
  `platen_jdg` varchar(50) DEFAULT NULL,
  `jacket` varchar(50) DEFAULT NULL,
  `jacket_jdg` varchar(50) DEFAULT NULL,
  `inpress_st` varchar(50) DEFAULT NULL,
  `inpress_st_jdg` varchar(50) DEFAULT NULL,
  `inpress_n2` varchar(50) DEFAULT NULL,
  `inpress_n2_jdg2` varchar(50) DEFAULT NULL,
  `curtime` varchar(50) DEFAULT NULL,
  `curtime_jdg` varchar(50) DEFAULT NULL,
  `finaljdg` varchar(50) DEFAULT NULL,
  `finaldfc` varchar(50) DEFAULT NULL,
  `desc` varchar(50) DEFAULT NULL,
  `recid` int NOT NULL AUTO_INCREMENT,
  `eventdate` datetime DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`recid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=468 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.gtentire
CREATE TABLE IF NOT EXISTS `gtentire` (
  `recid` int NOT NULL AUTO_INCREMENT,
  `bc_entried` char(50) DEFAULT NULL,
  `bld_date` char(50) DEFAULT NULL,
  `bld_mcn01` char(50) DEFAULT NULL,
  `bld_shift` char(50) DEFAULT NULL,
  `cur_in` char(50) DEFAULT NULL,
  `cur_out` char(50) DEFAULT NULL,
  `cur_shift` char(50) DEFAULT NULL,
  `cur_mcn` char(50) DEFAULT NULL,
  `cur_opr01` char(50) DEFAULT NULL,
  `jdge_date` char(50) DEFAULT NULL,
  `jdge` char(50) DEFAULT NULL,
  `probcode` char(50) DEFAULT NULL,
  `pic` char(50) DEFAULT NULL,
  `status` char(50) DEFAULT NULL,
  `whs_in` char(50) DEFAULT NULL,
  `whs_out` char(50) DEFAULT NULL,
  `serialnumb` char(50) DEFAULT NULL,
  `sarana` char(50) DEFAULT NULL,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.item_measurement
CREATE TABLE IF NOT EXISTS `item_measurement` (
  `tanggal` varchar(50) DEFAULT NULL,
  `operator` varchar(50) DEFAULT NULL,
  `grup` varchar(50) DEFAULT NULL,
  `nip` varchar(50) DEFAULT NULL,
  `shift` varchar(50) DEFAULT NULL,
  `codespec` varchar(50) DEFAULT NULL,
  `mesinl` varchar(50) DEFAULT NULL,
  `bcdl` varchar(50) DEFAULT NULL,
  `rsnl` varchar(50) DEFAULT NULL,
  `cureinl` varchar(50) DEFAULT NULL,
  `cureoutl` varchar(50) DEFAULT NULL,
  `mesinr` varchar(50) DEFAULT NULL,
  `bcdr` varchar(50) DEFAULT NULL,
  `rsnr` varchar(50) DEFAULT NULL,
  `cureinr` varchar(50) DEFAULT NULL,
  `cureoutr` varchar(50) DEFAULT NULL,
  `stepL` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `stepR` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `extempL` varchar(50) DEFAULT NULL,
  `extempR` varchar(50) DEFAULT NULL,
  `platenL` varchar(50) DEFAULT NULL,
  `platenR` varchar(50) DEFAULT NULL,
  `jacketL` varchar(50) DEFAULT NULL,
  `jacketR` varchar(50) DEFAULT NULL,
  `intempL` varchar(50) DEFAULT NULL,
  `intempR` varchar(50) DEFAULT NULL,
  `inpressN2L` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `inpressN2R` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `inpressStL` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `inpressStR` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `curtimeL` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `curtimeR` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `extemp_jdgL` varchar(50) DEFAULT NULL,
  `extemp_jdgR` varchar(50) DEFAULT NULL,
  `platen_jdgL` varchar(50) DEFAULT NULL,
  `platen_jdgR` varchar(50) DEFAULT NULL,
  `jacket_jdgL` varchar(50) DEFAULT NULL,
  `jacket_jdgR` varchar(50) DEFAULT NULL,
  `intemp_jdgL` varchar(50) DEFAULT NULL,
  `intemp_jdgR` varchar(50) DEFAULT NULL,
  `inpressN2_jdgL` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `inpressN2_jdgR` varchar(100) DEFAULT NULL,
  `inpressSt_jdgL` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `inpressSt_jdgR` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `jdg` varchar(50) DEFAULT NULL,
  `date_shift` varchar(50) DEFAULT NULL,
  `recid` int NOT NULL AUTO_INCREMENT,
  `eventdate` datetime DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.material
CREATE TABLE IF NOT EXISTS `material` (
  `recid` int NOT NULL AUTO_INCREMENT,
  `item` char(50) DEFAULT NULL,
  `bc_entried` char(50) DEFAULT NULL,
  `mcn` char(50) DEFAULT NULL,
  `opr` char(50) DEFAULT NULL,
  `txndate` char(50) DEFAULT NULL,
  `shift` char(50) DEFAULT NULL,
  `qty` char(50) DEFAULT NULL,
  `stock` char(50) DEFAULT NULL,
  `lokasi` char(50) DEFAULT NULL,
  `sarana` char(50) DEFAULT NULL,
  `jdge` char(50) DEFAULT NULL,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.monitoringtl1
CREATE TABLE IF NOT EXISTS `monitoringtl1` (
  `date_shift` varchar(50) NOT NULL,
  `ydate_shift` varchar(50) NOT NULL,
  `codesize` varchar(20) NOT NULL,
  `SpecSpeed` varchar(10) NOT NULL,
  `SpeedSP` varchar(20) NOT NULL,
  `SpeedSrv` varchar(20) NOT NULL,
  `SpeedProd` varchar(20) NOT NULL,
  `SpeedActual` varchar(20) NOT NULL,
  `ToleransiMin` varchar(20) NOT NULL,
  `ToleransiPlus` varchar(20) NOT NULL,
  `AverageSpeed` varchar(20) NOT NULL,
  `TotalDataUnder` varchar(20) NOT NULL,
  `TotalOutput` varchar(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.recorddatacyclic
CREATE TABLE IF NOT EXISTS `recorddatacyclic` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code_recipe` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `timestamp_record` datetime DEFAULT NULL,
  `comp_porkchop` varchar(50) DEFAULT NULL,
  `comp_e2middle` varchar(50) DEFAULT NULL,
  `comp_e3lower` varchar(53) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `comp_qsm60` varchar(50) DEFAULT NULL,
  `comp_qsm120` varchar(50) DEFAULT NULL,
  `preformer_name` varchar(50) DEFAULT NULL,
  `finaldie_name` varchar(50) DEFAULT NULL,
  `aspeed_e1` float DEFAULT NULL,
  `aspeed_e2` float DEFAULT NULL,
  `aspeed_e3` float DEFAULT NULL,
  `aspeed_shrink3` float DEFAULT NULL,
  `sspeed_e1` float DEFAULT NULL,
  `sspeed_e2` float DEFAULT NULL,
  `sspeed_e3` float DEFAULT NULL,
  `sspeed_shrink3` float DEFAULT NULL,
  `sberat_control` float DEFAULT NULL,
  `aberat_control` float DEFAULT NULL,
  `slebar_control` float DEFAULT NULL,
  `alebar_control` float DEFAULT NULL,
  `sberatctrl_tol--` float DEFAULT NULL,
  `sberatctrl_tol-` float DEFAULT NULL,
  `sberatctrl_tol+` float DEFAULT NULL,
  `sberatctrl_tol++` float DEFAULT NULL,
  `alebarUT` float DEFAULT NULL,
  `agapcalender_OS` float DEFAULT NULL,
  `agapcalender_DS` float DEFAULT NULL,
  `atcu_screw_pc` float DEFAULT NULL,
  `atcu_screw_e1` float DEFAULT NULL,
  `atcu_screw_e2` float DEFAULT NULL,
  `atcu_screw_e3` float DEFAULT NULL,
  `stcu_screw_pc` float DEFAULT NULL,
  `stcu_screw_e1` float DEFAULT NULL,
  `stcu_screw_e2` float DEFAULT NULL,
  `stcu_screw_e3` float DEFAULT NULL,
  `atcu_pc_barel1` float DEFAULT NULL,
  `atcu_pc_barel2` float DEFAULT NULL,
  `atcu_e1_zn1` float DEFAULT NULL,
  `atcu_e2_zn1` float DEFAULT NULL,
  `atcu_e2_zn2` float DEFAULT NULL,
  `atcu_e3_zn1` float DEFAULT NULL,
  `atcu_e3_zn2` float DEFAULT NULL,
  `stcu_pc_barel1` float DEFAULT NULL,
  `stcu_pc_barel2` float DEFAULT NULL,
  `stcu_e1_zn1` float DEFAULT NULL,
  `stcu_e2_zn1` float DEFAULT NULL,
  `stcu_e2_zn2` float DEFAULT NULL,
  `stcu_e3_zn1` float DEFAULT NULL,
  `stcu_e3_zn2` float DEFAULT NULL,
  `atcu_pc_head` float DEFAULT NULL,
  `atcu_e1_head` float DEFAULT NULL,
  `atcu_e2_head` float DEFAULT NULL,
  `atcu_e3_head` float DEFAULT NULL,
  `atemp_preformer` float DEFAULT NULL,
  `stcu_pc_head` float DEFAULT NULL,
  `stcu_e1_head` float DEFAULT NULL,
  `stcu_e2_head` float DEFAULT NULL,
  `stcu_e3_head` float DEFAULT NULL,
  `stemp_preformer` float DEFAULT NULL,
  `atemp_cal1` float DEFAULT NULL,
  `atemp_cal2` float DEFAULT NULL,
  `stemp_cal1` float DEFAULT NULL,
  `stemp_cal2` float DEFAULT NULL,
  `atcu_e120_screw` float DEFAULT NULL,
  `atcu_e120_zn1` float DEFAULT NULL,
  `atcu_e120_zn2` float DEFAULT NULL,
  `atcu_e120_head` float DEFAULT NULL,
  `stcu_e120_screw` float DEFAULT NULL,
  `stcu_e120_zn1` float DEFAULT NULL,
  `stcu_e120_zn2` float DEFAULT NULL,
  `stcu_e120_head` float DEFAULT NULL,
  `apos_d1_incline` float DEFAULT NULL,
  `apos_d2_cl1` float DEFAULT NULL,
  `apos_d3_cl2` float DEFAULT NULL,
  `apos_d4_cl3` float DEFAULT NULL,
  `apos_d5_cl4` float DEFAULT NULL,
  `apos_d6_blw` float DEFAULT NULL,
  `apos_d7_decline` float DEFAULT NULL,
  `spos_d1_incline` float DEFAULT NULL,
  `spos_d2_cl1` float DEFAULT NULL,
  `spos_d3_cl2` float DEFAULT NULL,
  `spos_d4_cl3` float DEFAULT NULL,
  `spos_d5_cl4` float DEFAULT NULL,
  `spos_d6_blw` float DEFAULT NULL,
  `spos_d7_decline` float DEFAULT NULL,
  `aspeed_d1_incline` float DEFAULT NULL,
  `aspeed_d2_cl1` float DEFAULT NULL,
  `aspeed_d3_cl2` float DEFAULT NULL,
  `aspeed_d4_cl3` float DEFAULT NULL,
  `aspeed_d5_cl4` float DEFAULT NULL,
  `aspeed_d6_blw` float DEFAULT NULL,
  `aspeed_d7_decline` float DEFAULT NULL,
  `spres_d1_incline` float DEFAULT NULL,
  `spres_d2_cl1` float DEFAULT NULL,
  `spres_d3_cl2` float DEFAULT NULL,
  `spres_d4_cl3` float DEFAULT NULL,
  `spres_d5_cl4` float DEFAULT NULL,
  `spres_d6_blw` float DEFAULT NULL,
  `spres_d7_decline` float DEFAULT NULL,
  `ampere_e1` float DEFAULT NULL,
  `ampere_e2` float DEFAULT NULL,
  `ampere_e3` float DEFAULT NULL,
  `prod_length_cooling` float DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=285947 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.recorddatapcs
CREATE TABLE IF NOT EXISTS `recorddatapcs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `recipe_code` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `waktu` datetime DEFAULT NULL,
  `sberat_finish` float DEFAULT NULL,
  `aberat_finish` float DEFAULT NULL,
  `scut_skiver` float DEFAULT NULL,
  `acut_skiver` float DEFAULT NULL,
  `slength` float DEFAULT NULL,
  `alength` float DEFAULT NULL,
  `switdhtol` float DEFAULT NULL,
  `awidthtol` float DEFAULT NULL,
  `sberat_finish--` float DEFAULT NULL,
  `sberat_finish-` float DEFAULT NULL,
  `sberat_finish+` float DEFAULT NULL,
  `sberat_finish++` float DEFAULT NULL,
  `slengthtol--` float DEFAULT NULL,
  `slengthtol-` float DEFAULT NULL,
  `slengthtol+` float DEFAULT NULL,
  `slengthtol++` float DEFAULT NULL,
  `auto_rejectStat` float DEFAULT NULL,
  `prod_OK` float DEFAULT NULL,
  `prod_NG` float DEFAULT NULL,
  `prod_Tot` float DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=457444 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rsc-pc1
CREATE TABLE IF NOT EXISTS `rsc-pc1` (
  `machine` varchar(20) NOT NULL,
  `position` varchar(40) NOT NULL,
  `code` varchar(10) NOT NULL,
  `barcode` varchar(30) NOT NULL,
  `totalseleksi` varchar(20) NOT NULL,
  `shift` varchar(10) NOT NULL,
  `starttime` varchar(40) NOT NULL,
  `stoptime` varchar(40) NOT NULL,
  KEY `key` (`machine`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rtba1
CREATE TABLE IF NOT EXISTS `rtba1` (
  `recid` int NOT NULL AUTO_INCREMENT,
  `Timestamp` varchar(50) DEFAULT NULL,
  `barcode` varchar(50) DEFAULT NULL,
  `specification` varchar(50) DEFAULT NULL,
  `pattern` varchar(100) DEFAULT NULL,
  `GT_CT` varchar(50) DEFAULT NULL,
  `preAssy_pudCT` decimal(10,1) DEFAULT NULL,
  `chafer_pudCT` decimal(10,1) DEFAULT NULL,
  `bodyplyApply_pudCT` decimal(10,1) DEFAULT NULL,
  `bodyplyStitch_pudCT` decimal(10,1) DEFAULT NULL,
  `beadTrf&BECapply_pudCT` decimal(10,1) DEFAULT NULL,
  `transfer_carcas_pudCT` decimal(10,1) DEFAULT NULL,
  `PUD_CT` decimal(10,1) DEFAULT NULL,
  `belt1_btdCT` decimal(10,1) DEFAULT NULL,
  `belt2_btdCT` decimal(10,1) DEFAULT NULL,
  `belt3_btdCT` decimal(10,1) DEFAULT NULL,
  `belt4_btdCT` decimal(10,1) DEFAULT NULL,
  `tread_btdCT` decimal(10,1) DEFAULT NULL,
  `trdstitching_btdCT` decimal(10,1) DEFAULT NULL,
  `BeltTreadtransfer_btdCT` decimal(10,1) DEFAULT NULL,
  `BTD_CT` decimal(10,1) DEFAULT NULL,
  `CarcassTRF_bddCT` decimal(10,1) DEFAULT NULL,
  `shapping_bddCT` decimal(10,1) DEFAULT NULL,
  `treadStict_bddCT` decimal(10,1) DEFAULT NULL,
  `turnUp_bddCT` decimal(10,1) DEFAULT NULL,
  `SWStict_bddCT` decimal(10,1) DEFAULT NULL,
  `gt_unload_bddCT` decimal(10,1) DEFAULT NULL,
  `BDD_CT` decimal(10,1) DEFAULT NULL,
  `treadlength` decimal(10,1) DEFAULT NULL,
  `shapping_press` decimal(10,1) DEFAULT NULL,
  `beadlock_press` decimal(10,1) DEFAULT NULL,
  `BPstc_press` decimal(10,1) DEFAULT NULL,
  `btoc` decimal(10,1) DEFAULT NULL,
  `BTB_width` decimal(10,1) DEFAULT NULL,
  `CTR_width` decimal(10,1) DEFAULT NULL,
  `pud_oc_col` decimal(10,1) DEFAULT NULL,
  `pud_oc_exp` decimal(10,1) DEFAULT NULL,
  `TreadHead_midP1` decimal(10,1) DEFAULT NULL,
  `TreadHead_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadHead_treadPos` decimal(10,1) DEFAULT NULL,
  `TreadMid_midP1` decimal(10,1) DEFAULT NULL,
  `TreadMid_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadMid_treadPos` decimal(10,1) DEFAULT NULL,
  `TreadTail_midP1` decimal(10,1) DEFAULT NULL,
  `TreadTail_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadTail_treadPos` decimal(10,1) DEFAULT NULL,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB AUTO_INCREMENT=55705 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rtbc1
CREATE TABLE IF NOT EXISTS `rtbc1` (
  `recid` int NOT NULL AUTO_INCREMENT,
  `Timestamp` varchar(50) DEFAULT NULL,
  `barcode` varchar(50) DEFAULT NULL,
  `specification` varchar(50) DEFAULT NULL,
  `pattern` varchar(100) DEFAULT NULL,
  `GT_CT` varchar(50) DEFAULT NULL,
  `preAssy_pudCT` decimal(10,1) DEFAULT NULL,
  `chafer_pudCT` decimal(10,1) DEFAULT NULL,
  `bodyplyApply_pudCT` decimal(10,1) DEFAULT NULL,
  `bodyplyStitch_pudCT` decimal(10,1) DEFAULT NULL,
  `beadTrf&BECapply_pudCT` decimal(10,1) DEFAULT NULL,
  `transfer_carcas_pudCT` decimal(10,1) DEFAULT NULL,
  `PUD_CT` decimal(10,1) DEFAULT NULL,
  `belt1_btdCT` decimal(10,1) DEFAULT NULL,
  `belt2_btdCT` decimal(10,1) DEFAULT NULL,
  `belt3_btdCT` decimal(10,1) DEFAULT NULL,
  `belt4_btdCT` decimal(10,1) DEFAULT NULL,
  `tread_btdCT` decimal(10,1) DEFAULT NULL,
  `trdstitching_btdCT` decimal(10,1) DEFAULT NULL,
  `BeltTreadtransfer_btdCT` decimal(10,1) DEFAULT NULL,
  `BTD_CT` decimal(10,1) DEFAULT NULL,
  `CarcassTRF_bddCT` decimal(10,1) DEFAULT NULL,
  `shapping_bddCT` decimal(10,1) DEFAULT NULL,
  `treadStict_bddCT` decimal(10,1) DEFAULT NULL,
  `turnUp_bddCT` decimal(10,1) DEFAULT NULL,
  `SWStict_bddCT` decimal(10,1) DEFAULT NULL,
  `gt_unload_bddCT` decimal(10,1) DEFAULT NULL,
  `BDD_CT` decimal(10,1) DEFAULT NULL,
  `treadlength` decimal(10,1) DEFAULT NULL,
  `shapping_press` decimal(10,1) DEFAULT NULL,
  `beadlock_press` decimal(10,1) DEFAULT NULL,
  `BPstc_press` decimal(10,1) DEFAULT NULL,
  `btoc` decimal(10,1) DEFAULT NULL,
  `BTB_width` decimal(10,1) DEFAULT NULL,
  `CTR_width` decimal(10,1) DEFAULT NULL,
  `pud_oc_col` decimal(10,1) DEFAULT NULL,
  `pud_oc_exp` decimal(10,1) DEFAULT NULL,
  `TreadHead_midP1` decimal(10,1) DEFAULT NULL,
  `TreadHead_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadHead_treadPos` decimal(10,1) DEFAULT NULL,
  `TreadMid_midP1` decimal(10,1) DEFAULT NULL,
  `TreadMid_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadMid_treadPos` decimal(10,1) DEFAULT NULL,
  `TreadTail_midP1` decimal(10,1) DEFAULT NULL,
  `TreadTail_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadTail_treadPos` decimal(10,1) DEFAULT NULL,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB AUTO_INCREMENT=69043 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rtbe1
CREATE TABLE IF NOT EXISTS `rtbe1` (
  `recid` int NOT NULL AUTO_INCREMENT,
  `Timestamp` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `barcode` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `specification` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `pattern` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `GT_CT` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `preAssy_pudCT` decimal(10,1) DEFAULT NULL,
  `chafer_pudCT` decimal(10,1) DEFAULT NULL,
  `bodyplyApply_pudCT` decimal(10,1) DEFAULT NULL,
  `bodyplyStitch_pudCT` decimal(10,1) DEFAULT NULL,
  `beadTrf&BECapply_pudCT` decimal(10,1) DEFAULT NULL,
  `transfer_carcas_pudCT` decimal(10,1) DEFAULT NULL,
  `PUD_CT` decimal(10,1) DEFAULT NULL,
  `belt1_btdCT` decimal(10,1) DEFAULT NULL,
  `belt2_btdCT` decimal(10,1) DEFAULT NULL,
  `belt3_btdCT` decimal(10,1) DEFAULT NULL,
  `belt4_btdCT` decimal(10,1) DEFAULT NULL,
  `tread_btdCT` decimal(10,1) DEFAULT NULL,
  `trdstitching_btdCT` decimal(10,1) DEFAULT NULL,
  `BeltTreadtransfer_btdCT` decimal(10,1) DEFAULT NULL,
  `BTD_CT` decimal(10,1) DEFAULT NULL,
  `CarcassTRF_bddCT` decimal(10,1) DEFAULT NULL,
  `shapping_bddCT` decimal(10,1) DEFAULT NULL,
  `treadStict_bddCT` decimal(10,1) DEFAULT NULL,
  `turnUp_bddCT` decimal(10,1) DEFAULT NULL,
  `SWStict_bddCT` decimal(10,1) DEFAULT NULL,
  `gt_unload_bddCT` decimal(10,1) DEFAULT NULL,
  `BDD_CT` decimal(10,1) DEFAULT NULL,
  `treadlength` decimal(10,1) DEFAULT NULL,
  `shapping_press` decimal(10,1) DEFAULT NULL,
  `beadlock_press` decimal(10,1) DEFAULT NULL,
  `BPstc_press` decimal(10,1) DEFAULT NULL,
  `btoc` decimal(10,1) DEFAULT NULL,
  `BTB_width` decimal(10,1) DEFAULT NULL,
  `CTR_width` decimal(10,1) DEFAULT NULL,
  `pud_oc_col` decimal(10,1) DEFAULT NULL,
  `pud_oc_exp` decimal(10,1) DEFAULT NULL,
  `TreadHead_midP1` decimal(10,1) DEFAULT NULL,
  `TreadHead_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadHead_treadPos` decimal(10,1) DEFAULT NULL,
  `TreadMid_midP1` decimal(10,1) DEFAULT NULL,
  `TreadMid_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadMid_treadPos` decimal(10,1) DEFAULT NULL,
  `TreadTail_midP1` decimal(10,1) DEFAULT NULL,
  `TreadTail_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadTail_treadPos` decimal(10,1) DEFAULT NULL,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB AUTO_INCREMENT=4086 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rtbe2
CREATE TABLE IF NOT EXISTS `rtbe2` (
  `recid` int NOT NULL AUTO_INCREMENT,
  `Timestamp` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `barcode` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `specification` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `pattern` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `GT_CT` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `preAssy_pudCT` decimal(10,1) DEFAULT NULL,
  `chafer_pudCT` decimal(10,1) DEFAULT NULL,
  `bodyplyApply_pudCT` decimal(10,1) DEFAULT NULL,
  `bodyplyStitch_pudCT` decimal(10,1) DEFAULT NULL,
  `beadTrf&BECapply_pudCT` decimal(10,1) DEFAULT NULL,
  `transfer_carcas_pudCT` decimal(10,1) DEFAULT NULL,
  `PUD_CT` decimal(10,1) DEFAULT NULL,
  `belt1_btdCT` decimal(10,1) DEFAULT NULL,
  `belt2_btdCT` decimal(10,1) DEFAULT NULL,
  `belt3_btdCT` decimal(10,1) DEFAULT NULL,
  `belt4_btdCT` decimal(10,1) DEFAULT NULL,
  `tread_btdCT` decimal(10,1) DEFAULT NULL,
  `trdstitching_btdCT` decimal(10,1) DEFAULT NULL,
  `BeltTreadtransfer_btdCT` decimal(10,1) DEFAULT NULL,
  `BTD_CT` decimal(10,1) DEFAULT NULL,
  `CarcassTRF_bddCT` decimal(10,1) DEFAULT NULL,
  `shapping_bddCT` decimal(10,1) DEFAULT NULL,
  `treadStict_bddCT` decimal(10,1) DEFAULT NULL,
  `turnUp_bddCT` decimal(10,1) DEFAULT NULL,
  `SWStict_bddCT` decimal(10,1) DEFAULT NULL,
  `gt_unload_bddCT` decimal(10,1) DEFAULT NULL,
  `BDD_CT` decimal(10,1) DEFAULT NULL,
  `treadlength` decimal(10,1) DEFAULT NULL,
  `shapping_press` decimal(10,1) DEFAULT NULL,
  `beadlock_press` decimal(10,1) DEFAULT NULL,
  `BPstc_press` decimal(10,1) DEFAULT NULL,
  `btoc` decimal(10,1) DEFAULT NULL,
  `BTB_width` decimal(10,1) DEFAULT NULL,
  `CTR_width` decimal(10,1) DEFAULT NULL,
  `pud_oc_col` decimal(10,1) DEFAULT NULL,
  `pud_oc_exp` decimal(10,1) DEFAULT NULL,
  `TreadHead_midP1` decimal(10,1) DEFAULT NULL,
  `TreadHead_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadHead_treadPos` decimal(10,1) DEFAULT NULL,
  `TreadMid_midP1` decimal(10,1) DEFAULT NULL,
  `TreadMid_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadMid_treadPos` decimal(10,1) DEFAULT NULL,
  `TreadTail_midP1` decimal(10,1) DEFAULT NULL,
  `TreadTail_edgeP6` decimal(10,1) DEFAULT NULL,
  `TreadTail_treadPos` decimal(10,1) DEFAULT NULL,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB AUTO_INCREMENT=2810 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rtc-tr1
CREATE TABLE IF NOT EXISTS `rtc-tr1` (
  `recid` int NOT NULL AUTO_INCREMENT,
  `Trimming_MachineNumber` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Tire_Code` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Barcode_Tire` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Start_Triming` timestamp NULL DEFAULT NULL,
  `End_Triming` timestamp NULL DEFAULT NULL,
  `Duration_TrimingProcess` time DEFAULT NULL,
  `Number_tiresTrimed` decimal(20,6) DEFAULT NULL,
  `Triming_MachineSpeed` decimal(20,6) DEFAULT NULL,
  `Pressure_Triming` decimal(20,6) DEFAULT NULL,
  `Temperature_Triming` decimal(20,6) DEFAULT NULL,
  `Triming_OperatorID` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`recid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rteex1
CREATE TABLE IF NOT EXISTS `rteex1` (
  `trxtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `date_shift` varchar(50) NOT NULL,
  `ydate_shift` varchar(50) NOT NULL,
  `machine` varchar(25) NOT NULL,
  `size` varchar(25) NOT NULL,
  `recipe` varchar(15) NOT NULL,
  `SpecWS` varchar(15) NOT NULL DEFAULT '',
  `ActWS` varchar(15) NOT NULL DEFAULT '',
  `USLWS` varchar(15) NOT NULL DEFAULT '',
  `LSLWS` varchar(15) NOT NULL DEFAULT '',
  `UCLWS` varchar(15) NOT NULL DEFAULT '',
  `LCLWS` varchar(15) NOT NULL DEFAULT '',
  `SetCutting` varchar(15) NOT NULL DEFAULT '',
  `SpecCutting` varchar(15) NOT NULL DEFAULT '',
  `LineSpeed` varchar(15) NOT NULL DEFAULT '0',
  `Deviasi` varchar(15) NOT NULL DEFAULT '0',
  `OK` varchar(15) NOT NULL DEFAULT '',
  `NG` varchar(15) NOT NULL DEFAULT '',
  `Total` varchar(15) NOT NULL DEFAULT '',
  `PersentaseOK` varchar(15) NOT NULL DEFAULT '0',
  `PersentaseNG` varchar(15) NOT NULL DEFAULT '0',
  `UCLRS` varchar(15) NOT NULL DEFAULT '',
  `LCLRS` varchar(15) NOT NULL DEFAULT '',
  `USLRS` varchar(15) NOT NULL DEFAULT '',
  `LSLRS` varchar(15) NOT NULL DEFAULT '',
  `SpecRS` varchar(15) NOT NULL DEFAULT '',
  `ActRS` varchar(15) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rteex2
CREATE TABLE IF NOT EXISTS `rteex2` (
  `mesin` varchar(20) NOT NULL,
  `trxtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `date_shift` varchar(30) NOT NULL DEFAULT '',
  `ydate_shift` varchar(30) NOT NULL DEFAULT '',
  `ng` varchar(10) NOT NULL,
  `ok` varchar(10) NOT NULL,
  `total` varchar(10) NOT NULL,
  `ngover` varchar(10) NOT NULL,
  `ngunder` varchar(10) NOT NULL,
  `lsl` varchar(10) NOT NULL,
  `spec` varchar(20) NOT NULL,
  `usl` varchar(10) NOT NULL,
  `size` varchar(20) NOT NULL,
  `actual` varchar(10) NOT NULL,
  `ngpersen` varchar(10) NOT NULL,
  `okpersen` varchar(10) NOT NULL,
  `speedline` varchar(20) NOT NULL,
  `speedbooking` varchar(20) NOT NULL,
  `recid` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`recid`),
  KEY `key1` (`trxtime`),
  KEY `key2` (`date_shift`),
  KEY `key3` (`ydate_shift`)
) ENGINE=InnoDB AUTO_INCREMENT=100108 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rteex3head
CREATE TABLE IF NOT EXISTS `rteex3head` (
  `trxtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `date_shift` varchar(25) NOT NULL,
  `ydate_shift` varchar(25) NOT NULL,
  `machine` varchar(15) NOT NULL,
  `size` varchar(15) NOT NULL,
  `ExtUpCurrent` int NOT NULL,
  `ExtMidCurrent` int NOT NULL,
  `ExtLowCurrent` int NOT NULL,
  `TempHeadExtUp` int NOT NULL,
  `TempHeadExtMid` int NOT NULL,
  `TempHeadExtLow` int NOT NULL,
  `SpeedFeedUp` int NOT NULL,
  `SpeedFeedMid` int NOT NULL,
  `SpeedFeedLow` int NOT NULL,
  `SpecRPMScrewUp` int NOT NULL,
  `SpecRPMScrewMid` int NOT NULL,
  `SpecRPMScrewLow` int NOT NULL,
  `SpeedRPMUp` int NOT NULL,
  `SpeedRPMMid` int NOT NULL,
  `SpeedRPMLow` int NOT NULL,
  `recid` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`recid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rtl-tl1
CREATE TABLE IF NOT EXISTS `rtl-tl1` (
  `tgl` varchar(50) NOT NULL,
  `jam` varchar(50) NOT NULL,
  `shift` varchar(10) NOT NULL,
  `codesize` varchar(50) NOT NULL,
  `SpecSpeed` varchar(50) NOT NULL,
  `SpeedSP` varchar(50) NOT NULL,
  `SpeedSrv` varchar(50) NOT NULL,
  `SpeedProd` varchar(50) NOT NULL,
  `SpeedA` varchar(50) NOT NULL,
  `tol--` varchar(20) NOT NULL,
  `tol++` varchar(20) NOT NULL,
  `leader` varchar(20) NOT NULL,
  `operator1` varchar(20) NOT NULL,
  `operator2` varchar(20) NOT NULL,
  `operator3` varchar(20) NOT NULL,
  `operator4` varchar(20) NOT NULL,
  `operator5` varchar(20) NOT NULL,
  `average` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.rtltl1
CREATE TABLE IF NOT EXISTS `rtltl1` (
  `date_shift` varchar(50) NOT NULL,
  `ydate_shift` varchar(50) NOT NULL,
  `txndate` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `mesin` varchar(20) NOT NULL,
  `size` varchar(25) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `TotalProd` int NOT NULL DEFAULT '0',
  `TotalOK` int NOT NULL DEFAULT '0',
  `PersenTotalOK` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `TotalReject` int NOT NULL DEFAULT '0',
  `PersenTotalReject` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `LebarOverSq` int NOT NULL DEFAULT '0',
  `PersentaseLOSQ` varchar(10) NOT NULL,
  `LebarUnSq` int NOT NULL DEFAULT '0',
  `PersentaseLUSQ` varchar(10) NOT NULL,
  `TebalOverSq` int NOT NULL DEFAULT '0',
  `PersentaseTOSQ` varchar(10) NOT NULL,
  `TebalUnderSq` int NOT NULL DEFAULT '0',
  `PersentaseTUSQ` varchar(10) NOT NULL,
  `LebarOverAp` int NOT NULL DEFAULT '0',
  `PersentaseLOAP` varchar(10) NOT NULL,
  `LebarUnderAp` int NOT NULL DEFAULT '0',
  `PersentaseLUAP` varchar(10) NOT NULL,
  `TebalOverAp` int NOT NULL DEFAULT '0',
  `PersentaseTOAP` varchar(10) NOT NULL,
  `TebalUnderAp` int NOT NULL DEFAULT '0',
  `PersentaseTUAP` varchar(10) NOT NULL,
  `Metal` int NOT NULL DEFAULT '0',
  `PersentaseMetal` varchar(10) NOT NULL,
  `LebarTotalUnder` int NOT NULL DEFAULT '0',
  `PersentaseLTU` varchar(20) NOT NULL,
  `LebarTotalOver` int NOT NULL DEFAULT '0',
  `PersentaseLTO` varchar(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

-- Dumping structure for table datamesin.trimming
CREATE TABLE IF NOT EXISTS `trimming` (
  `id` int NOT NULL AUTO_INCREMENT,
  `Machine_number` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Tirecode` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Barcode_tire` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Start_Trimming` datetime DEFAULT NULL,
  `End_Trimming` datetime DEFAULT NULL,
  `Duration_Process` time DEFAULT NULL,
  `Tires_Number` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Machine_Speed` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Pressure_Trimming` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Temperature` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `Operator_ID` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1234568 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
