-- Production Data Tables
-- Based on database-schema-migration.md section 3

-- 3.1 Building RTBA
CREATE TABLE IF NOT EXISTS `rtba1` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `Timestamp` VARCHAR(50),
    `barcode` VARCHAR(50),
    `specification` VARCHAR(50),
    `pattern` VARCHAR(100),
    `GT_CT` DECIMAL(10,1),
    `preAssy_pudCT` DECIMAL(10,1),
    `chafer_pudCT` DECIMAL(10,1),
    `bodyplyApply_pudCT` DECIMAL(10,1),
    `bodyplyStitch_pudCT` DECIMAL(10,1),
    `beadTrf_BECapply_pudCT` DECIMAL(10,1),
    `transfer_carcas_pudCT` DECIMAL(10,1),
    `PUD_CT` DECIMAL(10,1),
    `belt1_btdCT` DECIMAL(10,1),
    `belt2_btdCT` DECIMAL(10,1),
    `belt3_btdCT` DECIMAL(10,1),
    `belt4_btdCT` DECIMAL(10,1),
    `tread_btdCT` DECIMAL(10,1),
    `trdstitching_btdCT` DECIMAL(10,1),
    `BeltTreadtransfer_btdCT` DECIMAL(10,1),
    `BTD_CT` DECIMAL(10,1),
    `CarcassTRF_bddCT` DECIMAL(10,1),
    `shapping_bddCT` DECIMAL(10,1),
    `treadStict_bddCT` DECIMAL(10,1),
    `turnUp_bddCT` DECIMAL(10,1),
    `SWStict_bddCT` DECIMAL(10,1),
    `gt_unload_bddCT` DECIMAL(10,1),
    `BDD_CT` DECIMAL(10,1),
    `treadlength` DECIMAL(10,1),
    `shapping_press` DECIMAL(10,1),
    `beadlock_press` DECIMAL(10,1),
    `BPstc_press` DECIMAL(10,1),
    `btoc` DECIMAL(10,1),
    `BTB_width` DECIMAL(10,1),
    `CTR_width` DECIMAL(10,1),
    `pud_oc_col` DECIMAL(10,1),
    `pud_oc_exp` DECIMAL(10,1),
    INDEX(`barcode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `rtba2` LIKE `rtba1`;
CREATE TABLE IF NOT EXISTS `rtba3` LIKE `rtba1`;

-- 3.2 Building RTBC (same structure as RTBA)
CREATE TABLE IF NOT EXISTS `rtbc1` LIKE `rtba1`;
CREATE TABLE IF NOT EXISTS `rtbc2` LIKE `rtba1`;
CREATE TABLE IF NOT EXISTS `rtbc3` LIKE `rtba1`;
CREATE TABLE IF NOT EXISTS `rtbc4` LIKE `rtba1`;

-- 3.3 Building RTBE (same structure as RTBA)
CREATE TABLE IF NOT EXISTS `rtbe1` LIKE `rtba1`;
CREATE TABLE IF NOT EXISTS `rtbe2` LIKE `rtba1`;

-- 3.4 Extruder rteex1
CREATE TABLE IF NOT EXISTS `rteex1` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `trxtime` TIMESTAMP NULL,
    `date_shift` VARCHAR(50),
    `ydate_shift` VARCHAR(50),
    `machine` VARCHAR(25),
    `size` VARCHAR(25),
    `recipe` VARCHAR(15),
    `SpecWS` VARCHAR(15),
    `ActWS` VARCHAR(15),
    `USLWS` VARCHAR(15),
    `LSLWS` VARCHAR(15),
    `UCLWS` VARCHAR(15),
    `LCLWS` VARCHAR(15),
    `SetCutting` VARCHAR(15),
    `SpecCutting` VARCHAR(15),
    `LineSpeed` VARCHAR(15),
    `Deviasi` VARCHAR(15),
    `OK` VARCHAR(15),
    `NG` VARCHAR(15),
    `Total` VARCHAR(15),
    `PersentaseOK` VARCHAR(15),
    `PersentaseNG` VARCHAR(15),
    `UCLRS` VARCHAR(15),
    `LCLRS` VARCHAR(15),
    `USLRS` VARCHAR(15),
    `LSLRS` VARCHAR(15),
    `SpecRS` VARCHAR(15),
    `ActRS` VARCHAR(15),
    INDEX(`trxtime`),
    INDEX(`machine`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.4 Extruder rteex2
CREATE TABLE IF NOT EXISTS `rteex2` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `mesin` VARCHAR(20),
    `trxtime` TIMESTAMP NULL,
    `date_shift` VARCHAR(30),
    `ydate_shift` VARCHAR(30),
    `ok` VARCHAR(10),
    `ng` VARCHAR(10),
    `total` VARCHAR(10),
    `ngover` VARCHAR(10),
    `ngunder` VARCHAR(10),
    `okpersen` VARCHAR(10),
    `ngpersen` VARCHAR(10),
    `size` VARCHAR(20),
    `spec` VARCHAR(20),
    `lsl` VARCHAR(10),
    `usl` VARCHAR(10),
    `actual` VARCHAR(10),
    `speedline` VARCHAR(20),
    `speedbooking` VARCHAR(20)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.4 Extruder rteex3head
CREATE TABLE IF NOT EXISTS `rteex3head` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `trxtime` TIMESTAMP NULL,
    `date_shift` VARCHAR(25),
    `ydate_shift` VARCHAR(25),
    `machine` VARCHAR(15),
    `size` VARCHAR(15),
    `ExtUpCurrent` INT,
    `ExtMidCurrent` INT,
    `ExtLowCurrent` INT,
    `TempHeadExtUp` INT,
    `TempHeadExtMid` INT,
    `TempHeadExtLow` INT,
    `SpeedFeedUp` INT,
    `SpeedFeedMid` INT,
    `SpeedFeedLow` INT,
    `SpecRPMScrewUp` INT,
    `SpecRPMScrewMid` INT,
    `SpecRPMScrewLow` INT,
    `SpeedRPMUp` INT,
    `SpeedRPMMid` INT,
    `SpeedRPMLow` INT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.5 Curing curtire
CREATE TABLE IF NOT EXISTS `curtire` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `bcd` VARCHAR(50),
    `itemcode` VARCHAR(50),
    `codespec` VARCHAR(50),
    `mcn` VARCHAR(50),
    `mcntype` VARCHAR(50),
    `operator` VARCHAR(50),
    `nip` VARCHAR(50),
    `curein` VARCHAR(50),
    `cureout` VARCHAR(50),
    `dateshift` VARCHAR(50),
    `grpshift` VARCHAR(50),
    `rsn` VARCHAR(50),
    `extemp` VARCHAR(50),
    `extemp_jdg` VARCHAR(50),
    `intemp` VARCHAR(50),
    `intemp_jdg` VARCHAR(50),
    `platen` VARCHAR(50),
    `platen_jdg` VARCHAR(50),
    `jacket` VARCHAR(50),
    `jacket_jdg` VARCHAR(50),
    `inpress_st` VARCHAR(50),
    `inpress_st_jdg` VARCHAR(50),
    `inpress_n2` VARCHAR(50),
    `inpress_n2_jdg` VARCHAR(50),
    `curtime` VARCHAR(50),
    `curtime_jdg` VARCHAR(50),
    `finaljdg` VARCHAR(50),
    `finaldfc` VARCHAR(50),
    `eventdate` DATETIME,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX(`bcd`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.5 item_measurement
CREATE TABLE IF NOT EXISTS `item_measurement` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `tanggal` VARCHAR(50),
    `operator` VARCHAR(50),
    `grup` VARCHAR(50),
    `nip` VARCHAR(50),
    `shift` VARCHAR(50),
    `date_shift` VARCHAR(50),
    `codespec` VARCHAR(50),
    `mesinl` VARCHAR(50),
    `mesinr` VARCHAR(50),
    `bcdl` VARCHAR(50),
    `bcdr` VARCHAR(50),
    `rsnl` VARCHAR(50),
    `rsnr` VARCHAR(50),
    `cureinl` VARCHAR(50),
    `cureinr` VARCHAR(50),
    `cureoutl` VARCHAR(50),
    `cureoutr` VARCHAR(50),
    `step` VARCHAR(50),
    `extemp_jdgL` VARCHAR(50),
    `extemp_jdgR` VARCHAR(50),
    `platen_jdgL` VARCHAR(50),
    `platen_jdgR` VARCHAR(50),
    `jacket_jdgL` VARCHAR(50),
    `jacket_jdgR` VARCHAR(50),
    `intemp_jdgL` VARCHAR(50),
    `intemp_jdgR` VARCHAR(50),
    `inpressN2_jdgL` VARCHAR(50),
    `inpressN2_jdgR` VARCHAR(50),
    `inpressSt_jdgL` VARCHAR(50),
    `inpressSt_jdgR` VARCHAR(50),
    `jdg` VARCHAR(50),
    `eventdate` DATETIME,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.7 Trimming
CREATE TABLE IF NOT EXISTS `trimming` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `Machine_number` VARCHAR(50),
    `Tirecode` VARCHAR(50),
    `Barcode_tire` VARCHAR(50),
    `Start_Trimming` DATETIME,
    `End_Trimming` DATETIME,
    `Duration_Process` TIME,
    `Tires_Number` VARCHAR(50),
    `Machine_Speed` VARCHAR(50),
    `Pressure_Trimming` VARCHAR(50),
    `Temperature` VARCHAR(50),
    `Operator_ID` VARCHAR(50)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



-- 3.7 rtc-tr1
CREATE TABLE IF NOT EXISTS `rtc-tr1` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `Trimming_MachineNumber` VARCHAR(100),
    `Tire_Code` VARCHAR(100),
    `Barcode_Tire` VARCHAR(100),
    `Start_Triming` TIMESTAMP NULL,
    `End_Triming` TIMESTAMP NULL,
    `Duration_TrimingProcess` TIME,
    `Number_tiresTrimed` DECIMAL(20,6),
    `Triming_MachineSpeed` DECIMAL(20,6),
    `Pressure_Triming` DECIMAL(20,6),
    `Temperature_Triming` DECIMAL(20,6),
    `Triming_OperatorID` VARCHAR(50)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.8 Material
CREATE TABLE IF NOT EXISTS `material` (
    `recid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `item` VARCHAR(50),
    `bc_entried` VARCHAR(50),
    `mcn` VARCHAR(50),
    `opr` VARCHAR(50),
    `txndate` VARCHAR(50),
    `shift` VARCHAR(50),
    `qty` VARCHAR(50),
    `stock` VARCHAR(50),
    `lokasi` VARCHAR(50),
    `sarana` VARCHAR(50),
    `jdge` VARCHAR(50)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.9 monitoringtl1
CREATE TABLE IF NOT EXISTS `monitoringtl1` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `date_shift` VARCHAR(50),
    `ydate_shift` VARCHAR(50),
    `codesize` VARCHAR(20),
    `SpecSpeed` VARCHAR(20),
    `SpeedSP` VARCHAR(20),
    `SpeedSrv` VARCHAR(20),
    `SpeedProd` VARCHAR(20),
    `SpeedActual` VARCHAR(20),
    `TotalOutput` VARCHAR(20),
    `ToleransiMin` VARCHAR(20),
    `ToleransiPlus` VARCHAR(20),
    `AverageSpeed` VARCHAR(20),
    `TotalDataUnder` VARCHAR(20)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.9 rtltl1
CREATE TABLE IF NOT EXISTS `rtltl1` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `date_shift` VARCHAR(50),
    `ydate_shift` VARCHAR(50),
    `txndate` VARCHAR(100),
    `mesin` VARCHAR(20),
    `size` VARCHAR(25),
    `TotalProd` INT,
    `TotalOK` INT,
    `TotalReject` INT,
    `PersenTotalOK` VARCHAR(50),
    `PersenTotalReject` VARCHAR(50),
    `LebarOverSq` INT,
    `LebarUnSq` INT,
    `TebalOverSq` INT,
    `TebalUnderSq` INT,
    `LebarOverAp` INT,
    `LebarUnderAp` INT,
    `TebalOverAp` INT,
    `TebalUnderAp` INT,
    `Metal` INT,
    `LebarTotalUnder` INT,
    `LebarTotalOver` INT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.10 alarm_history
CREATE TABLE IF NOT EXISTS `alarm_history` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `timeOn` TIMESTAMP(3) NULL,
    `timeOff` TIMESTAMP(3) NULL,
    `source` TEXT,
    `message` TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.11 recipe1
CREATE TABLE IF NOT EXISTS `recipe1` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `size` TEXT,
    `speed_extupper` TEXT,
    `speed_extmiddle` TEXT,
    `speed_extlower` TEXT,
    `speed_line` TEXT,
    `speed_calender` TEXT,
    `run_scale` TEXT,
    `run_scale_up` TEXT,
    `run_scale_low` TEXT,
    `weight_scale` TEXT,
    `weight_up` TEXT,
    `weight_low` TEXT,
    `width` TEXT,
    `width_up` TEXT,
    `width_low` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.11 recipe1queue
CREATE TABLE IF NOT EXISTS `recipe1queue` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `order_id` VARCHAR(20),
    `recipe_id` BIGINT UNSIGNED,
    `size` TEXT,
    `qty` TEXT,
    `shift` TEXT,
    `user_name` TEXT,
    `time_create` TIMESTAMP(3) NULL,
    `time_start` TIMESTAMP(3) NULL,
    `time_finish` TIMESTAMP(3) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.11 recipe_history
CREATE TABLE IF NOT EXISTS `recipe_history` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `recipe_id` BIGINT UNSIGNED,
    `created_at` DATETIME,
    `created_by` BIGINT UNSIGNED,
    `modified_at` DATETIME NULL,
    `modified_by` BIGINT UNSIGNED NULL,
    `revision` VARCHAR(20),
    `note` TEXT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.12 order_report
CREATE TABLE IF NOT EXISTS `order_report` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `order_id` VARCHAR(20),
    `recipe_id` INT,
    `recipe` TEXT,
    `set_qty` TEXT,
    `act_qty` TEXT,
    `shift` TEXT,
    `user_name` TEXT,
    `time_create` TIMESTAMP(3) NULL,
    `time_start` TIMESTAMP(3) NULL,
    `time_finish` TIMESTAMP(3) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.12 batch_report
CREATE TABLE IF NOT EXISTS `batch_report` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `order_id` VARCHAR(20),
    `batch_id` VARCHAR(100),
    `user_name` TEXT,
    `shift` TEXT,
    `recipe` TEXT,
    `spec_weightscale` TEXT,
    `spec_weightscale_up` TEXT,
    `spec_weightscale_low` TEXT,
    `act_weightscale` TEXT,
    `spec_lengthskiver` TEXT,
    `act_lengthskiver` TEXT,
    `timestamp` TIMESTAMP(3) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.13 recorddatacyclic
CREATE TABLE IF NOT EXISTS `recorddatacyclic` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `code_recipe` VARCHAR(50),
    `timestamp_record` DATETIME,
    `comp_porkchop` VARCHAR(50),
    `comp_e2middle` VARCHAR(50),
    `comp_e3lower` VARCHAR(50),
    `comp_qsm60` VARCHAR(50),
    `comp_qsm120` VARCHAR(50),
    `preformer_name` VARCHAR(50),
    `finaldie_name` VARCHAR(50),
    `sspeed_e1` FLOAT,
    `sspeed_e2` FLOAT,
    `sspeed_e3` FLOAT,
    `sspeed_shrink3` FLOAT,
    `aspeed_e1` FLOAT,
    `aspeed_e2` FLOAT,
    `aspeed_e3` FLOAT,
    `aspeed_shrink3` FLOAT,
    `sberat_control` FLOAT,
    `aberat_control` FLOAT,
    `slebar_control` FLOAT,
    `alebar_control` FLOAT,
    `ampere` FLOAT,
    `prod_length_cooling` FLOAT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.13 recorddatapcs
CREATE TABLE IF NOT EXISTS `recorddatapcs` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `recipe_code` VARCHAR(50),
    `waktu` DATETIME,
    `sberat_finish` FLOAT,
    `aberat_finish` FLOAT,
    `scut_skiver` FLOAT,
    `acut_skiver` FLOAT,
    `slength` FLOAT,
    `alength` FLOAT,
    `switdhtol` FLOAT,
    `awidthtol` FLOAT,
    `auto_rejectStat` FLOAT,
    `prod_OK` FLOAT,
    `prod_NG` FLOAT,
    `prod_Tot` FLOAT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.14 datalog
CREATE TABLE IF NOT EXISTS `datalog` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `datetime` TIMESTAMP(3) NULL,
    `recipe1` TEXT,
    `recipe2` TEXT,
    `spec_speedextupper` TEXT,
    `spec_speedextmiddle` TEXT,
    `spec_speedextlower` TEXT,
    `act_speedextupper` TEXT,
    `act_speedextmiddle` TEXT,
    `act_speedextlower` TEXT,
    `act_ampextupper` TEXT,
    `act_ampextmiddle` TEXT,
    `act_ampextlower` TEXT,
    `spec_speedline` TEXT,
    `spec_speedcalender` TEXT,
    `act_speedline` TEXT,
    `act_speedcalender` TEXT,
    `spec_runscale` TEXT,
    `spec_runscale_up` TEXT,
    `spec_runscale_low` TEXT,
    `act_runscale` TEXT,
    `act_runscale_up` TEXT,
    `act_runscale_low` TEXT,
    `spec_weightscale` TEXT,
    `spec_weightscale_up` TEXT,
    `spec_weightscale_low` TEXT,
    `act_weightscale` TEXT,
    `spec_gapcalender` TEXT,
    `act_gapcalender` TEXT,
    `spec_cuttercalender` TEXT,
    `act_cuttercalender` TEXT,
    `act_compoundcalender` TEXT,
    `spec_lengthskiver` TEXT,
    `act_lengthskiver` TEXT,
    `spec_width` TEXT,
    `spec_width_up` TEXT,
    `spec_width_upline` TEXT,
    `spec_width_low` TEXT,
    `spec_width_lowline` TEXT,
    `act_width` TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3.15 rsc_pc1
CREATE TABLE IF NOT EXISTS `rsc_pc1` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `machine` VARCHAR(20),
    `position` VARCHAR(40),
    `code` VARCHAR(10),
    `barcode` VARCHAR(30),
    `totalseleksi` VARCHAR(20),
    `shift` VARCHAR(10),
    `starttime` VARCHAR(40),
    `stoptime` VARCHAR(40)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Undocumented Curing tables (rtci1, rtctr1) - basic structure
CREATE TABLE IF NOT EXISTS `rtci1` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `recid` BIGINT UNSIGNED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `rtctr1` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `recid` BIGINT UNSIGNED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add missing resource group items (curtire, trimming, rtc-tr1)
INSERT IGNORE INTO resource_group_items (group_id, resource_name, label, sort_order) VALUES
(2, 'curtire', 'curtire', 3),
(6, 'trimming', 'trimming', 2),
(6, 'rtc-tr1', 'rtc-tr1', 3);
