CREATE TABLE `model` (
  `id` enum('1') NOT NULL,
  `version` VARCHAR(20) NOT NULL,
  PRIMARY KEY (`id`)
);
INSERT INTO `model` (`version`) VALUES ("1.0.0");

-- Create the rest of the tables
CREATE TABLE `alert_group` (
	`id` INT NOT NULL AUTO_INCREMENT,
	`time` TIMESTAMP NOT NULL,
	`receiver` VARCHAR(100) NOT NULL,
	`status` VARCHAR(50) NOT NULL,
	`external_url` TEXT NOT NULL,
	`group_key` VARCHAR(255) NOT NULL,
	KEY `ix_time` (`time`),
        KEY `ix_status_time` (`status`, `time`),
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `group_label` (
	`id` INT NOT NULL AUTO_INCREMENT,
        `alert_group_id` INT NOT NULL,
        `group_label` VARCHAR(100) NOT NULL,
        `value` VARCHAR(1000) NOT NULL,     
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `common_label` (
	`id` INT NOT NULL AUTO_INCREMENT,
        `alert_group_id` INT NOT NULL,
        `label` VARCHAR(100) NOT NULL,
        `value` VARCHAR(1000) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `common_annotation` (
	`id` INT NOT NULL AUTO_INCREMENT,
        `alert_group_id` INT NOT NULL,
        `annotation` VARCHAR(100) NOT NULL,
        `value` VARCHAR(1000) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `alert` (
	`id` INT NOT NULL AUTO_INCREMENT,
        `alert_group_id` INT NOT NULL,
	`status` VARCHAR(50) NOT NULL,
        `starts_at` DATETIME NOT NULL,
        `ends_at` DATETIME DEFAULT NULL,
	`generator_url` TEXT NOT NULL,
	`fingerprint` VARCHAR(20) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `alert_label` (
	`id` INT NOT NULL AUTO_INCREMENT,
        `alert_id` INT NOT NULL,
        `label` VARCHAR(100) NOT NULL,
        `value` VARCHAR(1000) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `alert_annotation` (
	`id` INT NOT NULL AUTO_INCREMENT,
        `alert_id` INT NOT NULL,
    	`annotation` VARCHAR(100) NOT NULL,
    	`value` VARCHAR(1000) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
