CREATE DATABASE IF NOT EXISTS `baetyl_cloud` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `baetyl_cloud`;

CREATE TABLE IF NOT EXISTS `baetyl_application_history` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'app名称',
  `version` varchar(36) NOT NULL DEFAULT '' COMMENT 'app版本',
  `is_deleted` smallint(6) NOT NULL DEFAULT '0' COMMENT '删除标记1:删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `content` mediumtext COMMENT 'app详情',
  PRIMARY KEY (`id`),
  KEY `idx_app_history` (`namespace`,`name`,`version`),
  KEY `idx_app_date` (`namespace`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='application历史信息表';


CREATE TABLE IF NOT EXISTS `baetyl_batch` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '批号',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `description` varchar(1024) NOT NULL DEFAULT '' COMMENT '型号描述信息',
  `quota_num` int(11) NOT NULL DEFAULT '200' COMMENT '数量',
  `enable_whitelist` int(11) NOT NULL DEFAULT '1' COMMENT '是否启用白名单',
  `security_type` varchar(32) NOT NULL DEFAULT 'Token' COMMENT '安全等级 None/Token/Cert/Dongle',
  `security_key` varchar(64) NOT NULL DEFAULT '' COMMENT 'null/token/cert_id/dongle_id',
  `callback_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'callback name',
  `labels` varchar(2048) NOT NULL DEFAULT '{}' COMMENT '标签，json格式字符串,会设置到激活的node上',
  `fingerprint` varchar(1024) NOT NULL DEFAULT '{}' COMMENT '设备指纹信息，json格式字符串，包含类型等数据',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`namespace`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='批号管理';


CREATE TABLE IF NOT EXISTS `baetyl_batch_record` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '注册记录uuid',
  `batch_name` varchar(128) NOT NULL DEFAULT '' COMMENT '批号',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `fingerprint_value` varchar(255) NOT NULL DEFAULT '' COMMENT '注册序列号,多种组合采用逗号隔开',
  `active` int(1) NOT NULL DEFAULT '0' COMMENT '是否激活',
  `node_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'node资源名称',
  `active_ip` varchar(64) NOT NULL DEFAULT '0.0.0.0' COMMENT '激活ip',
  `active_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '激活时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_fingerprint` (`namespace`,`batch_name`,`fingerprint_value`),
  UNIQUE KEY `unique_name` (`namespace`,`batch_name`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='注册记录管理';

CREATE TABLE IF NOT EXISTS `baetyl_callback` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'callback uuid',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `method` varchar(36) NOT NULL DEFAULT 'GET' COMMENT 'Get/Post/Put/Delete',
  `params` varchar(2048) NOT NULL DEFAULT '{}' COMMENT 'query params',
  `body` varchar(2048) NOT NULL DEFAULT '{}' COMMENT 'body',
  `header` varchar(1024) NOT NULL DEFAULT '{}' COMMENT 'header',
  `url` varchar(1024) NOT NULL DEFAULT '' COMMENT 'url',
  `description` varchar(1024) NOT NULL DEFAULT '' COMMENT '描述信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_ns_name` (`namespace`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='回调';

CREATE TABLE IF NOT EXISTS `baetyl_index_application_config` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `application` varchar(128) NOT NULL DEFAULT '' COMMENT 'app名称',
  `config` varchar(128) NOT NULL DEFAULT '' COMMENT 'config名称',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_application` (`namespace`,`application`),
  KEY `idx_config` (`namespace`,`config`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='应用与配置索引表';

CREATE TABLE IF NOT EXISTS `baetyl_index_application_node` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `application` varchar(128) NOT NULL DEFAULT '' COMMENT 'app名称',
  `node` varchar(128) NOT NULL DEFAULT '' COMMENT 'node名称',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_application` (`namespace`,`application`),
  KEY `idx_node` (`namespace`,`node`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='应用与节点索引表';

CREATE TABLE IF NOT EXISTS `baetyl_index_application_secret` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `application` varchar(128) NOT NULL DEFAULT '' COMMENT 'app名称',
  `secret` varchar(128) NOT NULL DEFAULT '' COMMENT 'secret名称',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_application` (`namespace`,`application`),
  KEY `idx_secret` (`namespace`,`secret`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='应用与secret索引表';

CREATE TABLE IF NOT EXISTS `baetyl_system_config` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `type` varchar(128) NOT NULL DEFAULT '' COMMENT '配置类别',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '配置的键',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  `value` text NOT NULL COMMENT '配置的值',
  PRIMARY KEY (`id`),
  KEY `idx_type_key` (`type`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='系统配置表';

CREATE TABLE IF NOT EXISTS `baetyl_node_shadow` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'node影子名称',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `report` text COMMENT '上报内容',
  `desire` text COMMENT '期望内容',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`namespace`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='节点影子';

CREATE TABLE IF NOT EXISTS `baetyl_certificate` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `cert_id` varchar(128) NOT NULL DEFAULT '' COMMENT '证书id',
  `parent_id` varchar(64) NOT NULL DEFAULT '' COMMENT '上级证书id',
  `type` varchar(64) NOT NULL DEFAULT '' COMMENT '证书类型（根 CA、二级 CA、节点客户端证书、模块客户端证书、模块服务端证书）',
  `common_name` varchar(128) NOT NULL DEFAULT '' COMMENT '常用名',
  `description` varchar(256) NOT NULL DEFAULT '' COMMENT '描述信息',
  `not_before` datetime NOT NULL DEFAULT '2017-01-01 00:00:00' COMMENT '证书生效时间',
  `not_after` datetime NOT NULL DEFAULT '2017-01-01 00:00:00' COMMENT '记录失效时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  `csr` text COMMENT 'csr请求生成证书的信息',
  `content` text COMMENT '证书内容',
  `private_key` text COMMENT '根证书private_key信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_cert_id` (`cert_id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='证书表';

CREATE TABLE IF NOT EXISTS `baetyl_property` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `key` varchar(128) NOT NULL DEFAULT '',
  `value` text NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP ,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_type_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='System configuration property table';

COMMIT;