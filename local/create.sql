-- 2022-07-01 09:13:51



CREATE TABLE IF NOT EXISTS `d_cfg` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `value` varchar(10000) DEFAULT '',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `d_history` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0',
  `project_id` int(11) NOT NULL DEFAULT '0',
  `node_id` int(11) NOT NULL DEFAULT '0',
  `node` text,
  `version` varchar(100) NOT NULL DEFAULT '',
  `desc` varchar(1000) NOT NULL DEFAULT '',
  `project` mediumtext,
  `status` varchar(100) DEFAULT '',
  `log` longtext,
  `pid` int(11) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  `ci` varchar(1000) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_project_id` (`project_id`)
) ENGINE=InnoDB AUTO_INCREMENT=162 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `d_node` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '',
  `desc` varchar(100) DEFAULT '',
  `ip` varchar(100) NOT NULL DEFAULT '',
  `ssh_port` varchar(10) DEFAULT '22',
  `ssh_key` text,
  `ssh_username` varchar(100) DEFAULT '',
  `ssh_password` varchar(100) DEFAULT '',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='宿主机表';

CREATE TABLE IF NOT EXISTS `d_project` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `desc` varchar(100) DEFAULT '',
  `uid` int(11) NOT NULL DEFAULT '0',
  `ding` varchar(200) DEFAULT '',
  `docker` mediumtext,
  `native` mediumtext,
  `cronjob` varchar(100) DEFAULT '',
  `mode` varchar(10) DEFAULT 'native',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='虚拟项目表';

CREATE TABLE IF NOT EXISTS `d_project_cronjob` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(11) NOT NULL DEFAULT '0',
  `uid` int(11) NOT NULL DEFAULT '0',
  `ci` varchar(1000) DEFAULT '',
  `node_id` int(11) DEFAULT '0',
  `node` text,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `d_tag` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `desc` varchar(500) NOT NULL DEFAULT '',
  `uid` int(11) NOT NULL DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `d_tag_rel` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `tid` int(11) NOT NULL DEFAULT '0',
  `pid` int(11) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `d_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL DEFAULT '',
  `password` varchar(100) NOT NULL DEFAULT '',
  `real_name` varchar(100) NOT NULL DEFAULT '',
  `avatar_url` varchar(100) DEFAULT '',
  `email` varchar(100) DEFAULT '',
  `adm` tinyint(1) DEFAULT '0',
  `token` varchar(100) DEFAULT '',
  `token_created_at` datetime DEFAULT NULL,
  `token_expired_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';