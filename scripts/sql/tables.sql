CREATE DATABASE IF NOT EXISTS `baetyl_cloud` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `baetyl_cloud`;
CREATE TABLE IF NOT EXISTS `baetyl_batch` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '批号',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `description` varchar(1024) NOT NULL DEFAULT '' COMMENT '型号描述信息',
  `quota_num` int(11) NOT NULL DEFAULT '200' COMMENT '数量',
  `enable_whitelist` int(11) NOT NULL DEFAULT '1' COMMENT '是否启用白名单',
  `cluster` int(11) NOT NULL DEFAULT '0' COMMENT '是否支持集群部署',
  `security_type` varchar(32) NOT NULL DEFAULT 'Token' COMMENT '安全等级 None/Token/Cert/Dongle',
  `security_key` varchar(64) NOT NULL DEFAULT '' COMMENT 'null/token/cert_id/dongle_id',
  `callback_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'callback name',
  `labels` varchar(2048) NOT NULL DEFAULT '{}' COMMENT '标签，json格式字符串,会设置到激活的node上',
  `accelerator` varchar(32) NOT NULL DEFAULT '' COMMENT 'AI加速器',
  `sys_apps` varchar(1024) NOT NULL DEFAULT '' COMMENT '可选官方模块',
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

CREATE TABLE IF NOT EXISTS `baetyl_node_shadow` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'node影子名称',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `report` text COMMENT '上报内容',
  `desire` text COMMENT '期望内容',
  `report_meta` text COMMENT '上报内容元数据',
  `desire_meta` text COMMENT '期望内容元数据',
  `desire_version` varchar(36) NOT NULL DEFAULT '' COMMENT 'desire版本号，用于CAS',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`namespace`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='节点影子';

CREATE TABLE IF NOT EXISTS `baetyl_certificate` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `cert_id` varchar(128) NOT NULL DEFAULT '' COMMENT '证书id',
  `parent_id` varchar(64) NOT NULL DEFAULT '' COMMENT '上级证书id',
  `type` varchar(64) NOT NULL DEFAULT '' COMMENT '证书类型（根 CA、二级 CA、节点客户端证书、模块客户端证书、模块服务端证书）',
  `common_name` varchar(128) NOT NULL DEFAULT '' COMMENT '常用名',
  `description` varchar(1024) NOT NULL DEFAULT '' COMMENT '描述信息',
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
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'name',
  `value` text NOT NULL COMMENT 'value',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='System configuration property table';

CREATE TABLE IF NOT EXISTS `baetyl_module` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '模块名称',
  `version` varchar(36) NOT NULL DEFAULT '' COMMENT '版本',
  `image` varchar(1024) NOT NULL DEFAULT '' COMMENT '镜像',
  `programs` varchar(2048) NOT NULL DEFAULT '' COMMENT '进程模式程序包',
  `type` varchar(36) NOT NULL DEFAULT '0' COMMENT '应用类型，user:用户类型模块，runtime_user: 函数运行时类型的用户模块，system:系统类型模块，opt_system: 可选系统类型模块',
  `flag` int(10) NOT NULL DEFAULT '0' COMMENT '应用标识',
  `is_latest` int(1) NOT NULL DEFAULT '0' COMMENT '是否是最新版本',
  `description` varchar(1024) NOT NULL DEFAULT '' COMMENT '描述',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`,`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='module table';

CREATE TABLE IF NOT EXISTS `baetyl_cron_app` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID,主键',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '应用名称',
  `namespace` varchar(64) NOT NULL DEFAULT '' COMMENT '命名空间',
  `selector` varchar(2048) NOT NULL DEFAULT '' COMMENT '边缘节点选择器',
  `cron_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'cron time',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`namespace`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='cron app table';
COMMIT;
