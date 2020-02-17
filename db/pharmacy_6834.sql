-- --------------------------------------------------------
-- 主機:                           192.168.219.130
-- 伺服器版本:                        5.7.28-0ubuntu0.19.04.2 - (Ubuntu)
-- 伺服器作業系統:                      Linux
-- HeidiSQL 版本:                  10.3.0.5771
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- 傾印  資料表 tw_masks.pharmacy 結構
CREATE TABLE IF NOT EXISTS `pharmacy` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL DEFAULT '',
  `name` varchar(64) NOT NULL DEFAULT '',
  `tel` varchar(20) NOT NULL DEFAULT '',
  `addr` varchar(200) NOT NULL DEFAULT '',
  `latitude` decimal(14,7) NOT NULL DEFAULT '0.0000000',
  `longitude` decimal(14,7) NOT NULL DEFAULT '0.0000000',
  `adult_count` bigint(11) NOT NULL DEFAULT '0',
  `child_count` bigint(11) NOT NULL DEFAULT '0',
  `upd_time` timestamp NULL DEFAULT NULL,
  `sell_rule` varchar(200) NOT NULL DEFAULT '',
  `comment` varchar(200) NOT NULL DEFAULT '',
  `sold_out` tinyint(4) NOT NULL DEFAULT '0',
  `sold_out_date` date DEFAULT NULL,
  `disabled` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=20503 DEFAULT CHARSET=utf8mb4;

-- 取消選取資料匯出。

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
