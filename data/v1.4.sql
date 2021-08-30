alter table `file`
    add column `name` varchar(255) not null default '' comment '文件名称' after path;

alter table `router`
    add column `app_id` int(11) unsigned not null default 0 comment 'app id' after `nav_id`,
    add column `order_num` int(5) unsigned not null default 0 comment '排序号',
    add column `action` tinyint(3) unsigned not null default 1 comment '触发事件',
    add column `action_content` text comment '事件内容',
    modify column `key` varchar(50) NOT NULL DEFAULT '' COMMENT 'key值' ;

alter table `form`
    add column `is_delete` tinyint(1) unsigned not null default 0 comment '删除状态 0-未删除 1-已删除' after is_online,
    add column `datasource_table_id` int(11) unsigned not null default 0 comment '数据表id' after app_id;

alter table `field`
    add column `datasource_column_id` int(11) unsigned not null default 0 comment '数据表字段id' after form_id;

alter table `field_linkage`
    add column `id` int(11) unsigned auto_increment not null auto_increment first,
    drop column `first_id`,
    drop column `key`,
    add primary key (`id`);

drop table if exists `field_table`;
CREATE TABLE `field_table` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `form_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表单id',
    `field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '控件id',
    `datasource_table_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '数据表id',
    `is_export` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '导出配置 0: 关闭导出 1:允许导出',
    `is_filter` tinyint(1) unsigned not null default 0 comment '数据过滤配置  0:关闭过滤  1:允许过滤',
    PRIMARY KEY (`id`),
    KEY `form_field_idx` (`form_id`,`field_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='列表控件表';

drop table if exists `field_button`;
CREATE TABLE `field_button` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `form_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表单id',
    `field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '控件id',
    `flow_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '流程id',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '按钮名称',
    `event` tinyint(3) NOT NULL DEFAULT '1' COMMENT '触发事件类型',
    PRIMARY KEY (`id`),
    KEY `form_field_idx` (`form_id`,`field_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='按钮控件表';

drop table if exists `field_records`;
CREATE TABLE `field_records` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `form_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表单id',
    `field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '控件id',
    `datasource_column_relation_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表字段关联关系id',
    `datasource_column_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '关联字段id列表,逗号分分割',
    `mode` tinyint(3) NOT NULL DEFAULT '1' COMMENT '呈现方式 1-卡片 2-回填',
    `count_type` tinyint(3) NOT NULL DEFAULT '1' COMMENT '关联记录数量 1-单条 2-多条',
    `detail_status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '详情是否展示 0-关闭 1-开启',
    `columns` varchar(255) NOT NULL DEFAULT '' COMMENT '显示字段,关联字段id列表,逗号分分割',
    PRIMARY KEY (`id`),
    KEY `form_field_idx` (`form_id`,`field_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='关联控件表';

drop table if exists `field_table_button`;
CREATE TABLE `field_table_button` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `form_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表单id',
    `field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '控件id',
    `flow_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '流程id',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '按钮名称',
    `event` tinyint(3) NOT NULL DEFAULT '1' COMMENT '触发事件类型',
    PRIMARY KEY (`id`),
    KEY `form_field_idx` (`form_id`,`field_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='列表控件显示字段配置表';

drop table if exists `field_table_column`;
CREATE TABLE `field_table_column` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `form_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表单id',
    `field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '控件id',
    `datasource_column_relation_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表字段关联表id,0表示本数据表字段',
    `datasource_column_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表字段id',
    `show_name` varchar(255) NOT NULL DEFAULT '' COMMENT '表头显示名称',
    `is_condition` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否作为筛选条件 0-关闭 1-开启',
    PRIMARY KEY (`id`),
    KEY `form_field_idx` (`form_id`,`field_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='列表控件显示字段配置表';


drop table if exists `field_table_filter`;
CREATE TABLE `field_table_filter` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `form_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表单id',
    `field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '控件id',
    `datasource_column_relation_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表字段关联表id,0表示本数据表字段',
    `datasource_column_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表字段id',
    `field_type` varchar(20) NOT NULL COMMENT '控件类型',
    `field_type_condition` tinyint(3) NOT NULL DEFAULT '1' COMMENT '控件筛选选项条件',
    `field_type_condition_value` varchar(255) NOT NULL DEFAULT '' COMMENT '控件筛选选项值',
    PRIMARY KEY (`id`),
    KEY `form_field_idx` (`form_id`,`field_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='列表控件数据过滤配置表';

alter table `service`
    add column `type` tinyint(3) unsigned not null default 1 comment '服务类型 1-上链服务 2-发送邮件 3-代码包';

alter table `param`
    add column `check` tinyint(1) unsigned not null default 1 comment '是否必填',
    add column `mode` tinyint(1) unsigned not null default 1 comment '入/出参 1-入 2-出',
    add column `type` tinyint(3) unsigned not null default 1 comment '参数类型 1-string 2-uint 3-int 4-float 5-double 6-list',
    add column `default_value` varchar(255) default '' comment '默认值',
    add column `fixed_value` varchar(255) default '' comment '固定值',
    add column `is_visual` tinyint(1) default 1 comment '配置端是否可见 0-不可见 1-可见';

alter table `flow`
    add column `desc` text comment '表单描述' after `name`,
    add column `user_id` int(11) unsigned not null default 0 comment '用户id' after `desc`,
    add column `json` text comment '工作流json' after `xml`,
    add column `is_delete` tinyint(1) unsigned not null default 0 comment '删除状态' after `is_online` ;

drop table if exists `flow_activity`;
create table `flow_activity`(
    `id` int(11) unsigned not null auto_increment,
    `flow_id` int(11) unsigned not null default 0 comment '流程id',
    `service_id` int(11) unsigned not null default 0 comment '服务id',
    `node_id` bigint(20) unsigned not null default 0 comment '前端生成的节点id',
    `name` varchar(255) not null default '' comment '节点名称',
    `desc` varchar(255) not null default '' comment '节点描述',
    `type` tinyint(3) unsigned not null default 1 comment '节点类型 1-开始节点-数据表字段时间类型;2-开始节点-固定时间类型;3-数据操作类型4-外部服务类型',
    primary key (`id`)
) engine = InnoDB default charset=utf8 comment='节点表';

drop table  if exists `flow_activity_start_table`;
create table `flow_activity_start_table`(
    `id` int(11) unsigned not null auto_increment,
    `flow_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '流程id',
    `flow_activity_id` int(11) unsigned not null default 0 comment '节点id',
    `datasource_table_id` int(11) unsigned not null default 0 comment '表id',
    `datasource_column_id` int(11) unsigned not null default 0 comment '表字段id',
    `trigger_interval` int(5) not null default 0 comment '触发间隔 0：当前时间 -n:提前n小时 n:延后n小时',
    primary key (`id`)
)engine = InnoDB default charset=utf8 comment='开始节点数据表字段时间类型';

drop table  if exists `flow_activity_start_time`;
CREATE TABLE `flow_activity_start_time` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `flow_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '流程id',
    `flow_activity_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '节点id',
    `trigger_time` datetime DEFAULT NULL COMMENT '触发时间',
    `trigger_interval` int(5) NOT NULL DEFAULT '0' COMMENT '触发间隔 0：当前时间 -n:提前n小时 n:延后n小时',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='开始节点固定时间类型';

drop table  if exists `flow_activity_service`;
CREATE TABLE `flow_activity_service` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `app_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '应用id',
    `flow_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '流程id',
    `flow_activity_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '当前流程节点id',
    `param_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '参数id',
    `input_node_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '入参-前端生成的节点id',
    `input_field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '入参来源 控件id',
    `input_param_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '入参来源 服务出参id',
    `input_lowcode_user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '入参来源 低代码用户id',
    `input_text` text COMMENT '入参来源 输入文本',
    `input_field_table_column_id` int(11) unsigned not null default 0 comment '入参来源 列表控件显示字段表id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='外部服务类型';

drop table  if exists `flow_activity_data`;
CREATE TABLE `flow_activity_data` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `app_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '应用id',
    `flow_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '流程id',
    `flow_activity_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '当前流程节点id',
    `datasource_table_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '数据表id',
    `datasource_column_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '数据表字段id',
    `type` tinyint(3) NOT NULL DEFAULT '1' COMMENT '参数字段类型 1:操作字段 2:条件字段',
    `expression` tinyint(3) NOT NULL DEFAULT '1' COMMENT '操作符 1:是/等于 2:不是/不等于',
    `op` tinyint(3) NOT NULL DEFAULT '1' COMMENT '关联符 0：无 1:且 2:或',
    `op_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '关联本表主键id',
    `group_id` int(11) unsigned not null default 0 comment '组编号',
    `input_node_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '入参-前端生成的节点id',
    `input_field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '入参来源 控件id',
    `input_param_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '入参来源 服务出参id',
    `input_text` varchar(255) NOT NULL DEFAULT '' COMMENT '入参来源 输入文本',
    `input_field_table_column_id` int(11) unsigned not null default 0 comment '入参来源 列表控件显示字段表id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='数据操作类型';

drop table  if exists `flow_activity_gateway`;
CREATE TABLE `flow_activity_gateway` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `flow_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '流程id',
    `flow_activity_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '当前流程节点id',
    `node_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '前端生成的节点id',
    `left_input_field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '左值入参来源 控件id',
    `left_input_param_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '左值入参来源 服务出参id',
    `left_input_field_table_column_id` int(11) unsigned not null default 0 comment '左值入参来源 列表控件显示字段表id',
    `expression` tinyint(3) NOT NULL DEFAULT '1' COMMENT '操作符 1:是/等于 2:不是/不等于',
    `op` tinyint(3) NOT NULL DEFAULT '1' COMMENT '关联符 0：无 1:且 2:或',
    `op_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '关联本表主键id 0：无 ',
    `group_id` int(11) unsigned not null default 0 comment '组编号',
    `right_input_field_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '右值入参来源 控件id',
    `right_input_param_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '右值入参来源 服务出参id',
    `right_input_text` varchar(255)  NOT NULL DEFAULT '' COMMENT '右值入参来源 手动输入文本',
    `right_input_field_table_column_id` int(11) unsigned not null default 0 comment '右值入参来源 列表控件显示字段表id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网关分支类型';

drop table  if exists `datasource_table`;
CREATE TABLE `datasource_table` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `app_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '应用id',
    `schemata_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '数据源id',
    `table_schema` varchar(64) NOT NULL DEFAULT '' COMMENT '数据库名称',
    `table_name` varchar(64) NOT NULL DEFAULT '' COMMENT '数据表名称',
    `table_comment` varchar(2048) NOT NULL DEFAULT '' COMMENT '数据库表备注',
    PRIMARY KEY (`id`),
    KEY `app_id_idx` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据表';


drop table  if exists `datasource_schemata`;
CREATE TABLE `datasource_schemata` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `app_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '应用id',
    `table_schema` varchar(64) NOT NULL DEFAULT '' COMMENT '数据库名称',
    PRIMARY KEY (`id`),
    KEY `app_id_idx` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据源';

drop table  if exists `datasource_metadata`;
CREATE TABLE `datasource_metadata` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `app_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '应用id',
    `datasource_column_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '表字段id',
    `group` varchar(64) NOT NULL DEFAULT '' COMMENT '分组组名',
    `key` varchar(64) NOT NULL DEFAULT '' COMMENT 'key键',
    `value` varchar(64) NOT NULL DEFAULT '' COMMENT 'value值',
    PRIMARY KEY (`id`),
    KEY `app_id_idx` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据字典';


drop table  if exists `datasource_column_relation`;
CREATE TABLE `datasource_column_relation` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `source_table_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '源数据表id',
    `source_column_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '源表字段id',
    `target_table_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '目标数据表id',
    `target_column_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '目标表字段id',
    `type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '关联类型.1:1对1，2:1对多,3:多对1,4:多对多',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='表字段关联表';

drop table  if exists `datasource_column`;
CREATE TABLE `datasource_column` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `app_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '应用id',
    `schemata_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '数据源id',
    `table_schema` varchar(64) NOT NULL DEFAULT '' COMMENT '数据库名称',
    `table_name` varchar(64) NOT NULL DEFAULT '' COMMENT '数据表名称',
    `data_type` varchar(64) NOT NULL DEFAULT '' COMMENT '字段数据类型',
    `character_maximum_length` bigint(21) unsigned DEFAULT NULL COMMENT '字段长度',
    `column_name` varchar(64) NOT NULL DEFAULT '' COMMENT '字段名称',
    `column_comment` varchar(1024) NOT NULL DEFAULT '' COMMENT '字段备注',
    `show_type` varchar(20) NOT NULL COMMENT '字段显示类型(业务类型) text:文本，num: 数字 date: 日期 radio:单选 checkbox:多选 user:成员 org:组织 relation: 关联 file: 文件',
    `field_type` varchar(20) NOT NULL COMMENT '控件类型',
    `is_system_field` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为系统字段0 不是系统字段 1是系统字段',
    PRIMARY KEY (`id`),
    KEY `app_id_idx` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='表字段';

truncate table `field`;
truncate table `field_button`;
truncate table `field_linkage`;
truncate table `field_multi_form`;
truncate table `field_records`;
truncate table `field_table`;
truncate table `field_table_button`;
truncate table `field_table_column`;
truncate table `field_table_filter`;
truncate table `file`;
truncate table `flow`;
truncate table `flow_assignee`;
truncate table `flow_mapping`;
truncate table `flow_notifier`;
truncate table `form`;
truncate table `navigation`;
truncate table `org_role`;
truncate table `organization`;
truncate table `param`;
truncate table `param_rely`;
truncate table `permission`;
truncate table `role`;
truncate table `role_permission`;
truncate table `router`;
truncate table `service`;
truncate table `user`;
truncate table `version`;

update application set `status` = 0 where `status` = 1;

insert into `service`
    (name, type,created_at)
    value
    ('上链服务',1,CURRENT_TIME),
    ('发送邮件',2,CURRENT_TIME),
    ('生成签名',3,CURRENT_TIME),
    ('合同归档',3,current_time);


INSERT INTO `param`
    ( `name`, `service_id`, `created_at`,  `check`, `mode`, `type`, `default_value`, `fixed_value`, `is_visual`)
VALUES
    ( 'value', 1, '2021-03-26 21:35:09', 1, 1, 6, '', '', 1),
    ( 'result', 1, '2021-08-02 16:55:00',  1, 2, 1, '', '', 1),
    ( 'sender', 2, '2021-08-03 17:10:54',  1, 1, 1, '', '', 1),
    ( 'receivers', 2, '2021-08-03 17:10:54',  1, 1, 6, '', '', 1),
    ( 'subject', 2, '2021-08-03 17:10:54',  1, 1, 1, '', '', 1),
    ( 'content', 2, '2021-08-03 17:10:54',  1, 1, 1, '', '', 1),
    ( 'access_key', 3, '2021-08-03 17:10:54',  1, 1, 1, '', '', 1),
    ( 'method', 3, '2021-08-03 17:10:54',  1, 1, 1, '', '', 1),
    ( 'path', 3, '2021-08-03 17:10:54',  1, 1, 1, '', '', 1),
    ( 'secret_key', 3, '2021-08-03 17:10:54',  1, 1, 1, '', '', 1),
    ( 'date', 3, '2021-08-03 17:10:54',  1, 2, 1, '', '', 1),
    ( 'rand', 3, '2021-08-03 17:10:54',  1, 2, 2, '0', '', 1),
    ( 'signature', 3, '2021-08-03 17:10:54',  1, 2, 1, '', '', 1),
    ( 'id', 4, '2021-08-16 17:27:14',  1, 1, 2, '0', '', 1);