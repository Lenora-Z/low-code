use `default-bpmn`;

CREATE TABLE `file` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime NOT NULL,
    `path` varchar(255) NOT NULL COMMENT '文件地址',
    `is_delete` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
    `app_id` int(11) unsigned NOT NULL COMMENT '应用id',
    `file_hash` varchar(255) NOT NULL COMMENT '文件hash值',
    `hash` varchar(50) NOT NULL COMMENT '数据hash',
    PRIMARY KEY (`id`),
    UNIQUE KEY `hash` (`hash`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

alter table `form`
    add column `type` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '表单类型 1-标准表单\n2-弹窗\n3-导航栏' after app_id,
    add column `footer` text comment '底部数据' after content;

CREATE TABLE `navigation` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL COMMENT '导航名称',
    `app_id` int(11) unsigned NOT NULL COMMENT '应用id',
    `number` varchar(100) NOT NULL COMMENT '编号',
    `content` text COMMENT '导航内容',
    `desc` varchar(255) DEFAULT NULL COMMENT '导航描述',
    `status` tinyint(1) NOT NULL COMMENT '表单状态 0-未生效 1-已生效',
    `is_online` tinyint(1) unsigned NOT NULL COMMENT '是否是上线表单 1-是',
    `created_at` datetime NOT NULL,
    `updated_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;