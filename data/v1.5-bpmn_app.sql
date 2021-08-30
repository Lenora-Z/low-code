use bpmn_app;

drop table if exists `crm_client_contact`;
create table `crm_client_contact`(
    `id` int(11) unsigned not null auto_increment,
    `client_id` int(11) unsigned not null default 0 comment '客户id',
    `company_name` varchar(255)  not null default '' comment '公司名称',
    `position` varchar(255)  not null default '' comment '职务',
    `name` varchar(255) not null default '' comment '姓名',
    `mobile` varchar(20) not null default '' comment '联系电话',
    `principal` int(11) unsigned not null default 0 comment '负责人',
    `note` text comment '备注',
    `creator_id` int(11) unsigned not null default 0 comment '创建人id',
    `modifier_id` int(11) unsigned not null default 0 comment '修改人id',
    `created_at` datetime not null default current_timestamp comment '记录创建时间',
    `updated_at` datetime  comment '记录最后修改时间',
    primary key (`id`),
    key `client` (client_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='客户联系人表';

drop table if exists `crm_client`;
create table `crm_client`(
    `id` int(11) unsigned not null auto_increment,
    `channel_id` int(11) unsigned not null default 0 comment '渠道id',
    `number` int(11) unsigned not null default 0 comment '编号',
    `name` varchar(255) not null default '' comment '名称',
    `from` tinyint(2) unsigned not null default 7 comment '来源',
    `status` tinyint(2) unsigned not null default 0 comment '客户状态',
    `level` tinyint(2) unsigned not null default 0 comment '信用等级',
    `tag` varchar(255) not null default '' comment '客户标签',
    `follow_status` tinyint(2) not null default 0 comment '跟进状态',
    `principal` int(11) not null default 0 comment '负责人',
    `address` varchar(255) not null default '' comment '地址',
    `type` tinyint(2) unsigned not null default 1 comment '生态类别',
    `follower` int(11) unsigned not null  default 0 comment '跟进人',
    `creator_id` int(11) unsigned not null default 0 comment '创建人id',
    `modifier_id` int(11) unsigned not null default 0 comment '修改人id',
    `created_at` datetime not null default current_timestamp comment '记录创建时间',
    `updated_at` datetime  comment '记录最后修改时间',
    primary key (`id`),
    key `channel` (`channel_id`) using btree
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='客户表';

drop table if exists `crm_client_follow`;
create table `crm_client_follow`(
    `id` int(11) unsigned not null auto_increment,
    `client_id` int(11) unsigned not null default 0 comment '客户id',
    `contact_id` int(11) unsigned not null default 0 comment '联系人id',
    `method` tinyint(2) unsigned not null default 0 comment '跟进方式',
    `follow_at` datetime not null default current_timestamp comment '跟进时间',
    `location` varchar(255) not null default '' comment '跟进地点',
    `content` text comment '内容',
    `visitor_id` int(11) unsigned not null default 0 comment '跟进人',
    `creator_id` int(11) unsigned not null default 0 comment '创建人id',
    `modifier_id` int(11) unsigned not null default 0 comment '修改人id',
    `created_at` datetime not null default current_timestamp comment '记录创建时间',
    `updated_at` datetime  comment '记录最后修改时间',
    primary key (`id`),
    key client_contact (`client_id`,`contact_id`) using btree
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='客户跟进表';

drop table if exists `crm_channel`;
create table `crm_channel`(
    `id` int(11) unsigned not null auto_increment,
    `number` int(11) unsigned not null default 0 comment '渠道编号',
    `name` varchar(255) not null default '' comment '渠道名称',
    `principal` int(11) not null default 0 comment '负责人',
    `note` text comment '备注',
    `creator_id` int(11) unsigned not null default 0 comment '创建人id',
    `modifier_id` int(11) unsigned not null default 0 comment '修改人id',
    `created_at` datetime not null default current_timestamp comment '记录创建时间',
    `updated_at` datetime  comment '记录最后修改时间',
    primary key (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='渠道表';

drop table if exists `crm_channel_contact`;
create table `crm_channel_contact`(
    `id` int(11) unsigned not null auto_increment,
    `channel_id` int(11) unsigned not null default 0 comment '渠道id',
    `company_name` varchar(255) not null default '' comment '公司名称',
    `position` varchar(255) not null default '' comment '职务',
    `name` varchar(255) not null default '' comment '姓名',
    `mobile` varchar(20) not null default '' comment '联系电话',
    `principal` int(11) not null default 0 comment '负责人',
    `note` text comment '备注',
    `creator_id` int(11) unsigned not null default 0 comment '创建人id',
    `modifier_id` int(11) unsigned not null default 0 comment '修改人id',
    `created_at` datetime not null default current_timestamp comment '记录创建时间',
    `updated_at` datetime  comment '记录最后修改时间',
    primary key (`id`),
    key `channel` (`channel_id`) using btree
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='渠道联系人表';

drop table if exists `crm_channel_follow`;
create table `crm_channel_follow`(
    `id` int(11) unsigned not null auto_increment,
    `channel_id` int(11) unsigned not null default 0 comment '渠道id',
    `contact_id` int(11) unsigned not null default 0 comment '联系人id',
    `method` tinyint(2) unsigned not null default 0 comment '跟进方式',
    `follow_at` datetime not null default current_timestamp comment '跟进时间',
    `location` varchar(255) not null default '' comment '跟进地点',
    `content` text comment '内容',
    `visitor_id` int(11) unsigned not null default 0 comment '跟进人',
    `creator_id` int(11) unsigned not null default 0 comment '创建人id',
    `modifier_id` int(11) unsigned not null default 0 comment '修改人id',
    `created_at` datetime not null default current_timestamp comment '记录创建时间',
    `updated_at` datetime  comment '记录最后修改时间',
    primary key (`id`),
    key channel_contact (`channel_id`,`contact_id`) using btree
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='渠道跟进表';

drop table if exists `crm_project`;
create table `crm_project`(
    `id` int(11) unsigned not null auto_increment,
    `name` varchar(255) not null default '' comment '名称',
    `principal` int(11) unsigned not null default 0 comment '负责人',
    `client_id` int(11) unsigned not null default 0 comment '客户id',
    `contact_id` int(11) unsigned not null default 0 comment '联系人id',
    `amount` decimal(20,2) unsigned not null default 0 comment '合同金额',
    `start_at` datetime not null default current_timestamp comment '开始时间',
    `end_at` datetime not null default current_timestamp comment '结束时间',
    `status` tinyint(2) unsigned not null default 0 comment '完成情况',
    `note` text comment '备注',
    `creator_id` int(11) unsigned not null default 0 comment '创建人id',
    `modifier_id` int(11) unsigned not null default 0 comment '修改人id',
    `created_at` datetime not null default current_timestamp comment '记录创建时间',
    `updated_at` datetime  comment '记录最后修改时间',
    primary key (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='项目表';

drop table if exists `crm_project_schedule`;
create table `crm_project_schedule`(
    `id` int(11) unsigned not null auto_increment,
    `project_id` int(11) unsigned not null default 0 comment '项目id',
    `start_at` datetime default current_timestamp comment '开始时间',
    `note` text comment '项目记录',
    `solution` varchar(255) not null default '' comment '解决方案',
    `performer` int(11) not null default 0 comment '执行人',
    `creator_id` int(11) unsigned not null default 0 comment '创建人id',
    `modifier_id` int(11) unsigned not null default 0 comment '修改人id',
    `created_at` datetime not null default current_timestamp comment '记录创建时间',
    `updated_at` datetime  comment '记录最后修改时间',
    primary key (`id`),
    key `project` (`project_id`) using btree
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='项目进度表';

SELECT CREATE_TIME,TABLE_SCHEMA,TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES
where TABLE_SCHEMA = 'bpmn_app'
  AND TABLE_NAME like 'crm%';

SELECT TABLE_SCHEMA,TABLE_NAME,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,COLUMN_NAME,COLUMN_COMMENT FROM information_schema.COLUMNS
where TABLE_SCHEMA = 'bpmn_app'
  AND TABLE_NAME like 'crm%';

use bpmn;

insert into  `datasource_schemata` (app_id, table_schema) values (159,'bpmn_app');

INSERT INTO `datasource_table`
    (`created_at`,`app_id`,`schemata_id`,`table_schema`,`table_name`, `table_comment`)
VALUES
    ('2021-08-17 17:32:29',159,4, 'bpmn_app', 'crm_channel', '渠道表'),
    ('2021-08-17 17:32:29',159,4, 'bpmn_app', 'crm_channel_contact', '渠道联系人表'),
    ('2021-08-17 17:32:29',159,4, 'bpmn_app', 'crm_channel_follow', '渠道跟进表'),
    ('2021-08-17 17:32:29',159,4, 'bpmn_app', 'crm_client', '客户表'),
    ('2021-08-17 17:32:28',159,4, 'bpmn_app', 'crm_client_contact', '客户联系人表'),
    ('2021-08-17 17:32:29',159,4, 'bpmn_app', 'crm_client_follow', '客户跟进表'),
    ('2021-08-17 17:32:29',159,4, 'bpmn_app', 'crm_project', '项目表'),
    ('2021-08-17 17:32:29',159,4, 'bpmn_app', 'crm_project_schedule', '项目进度表');

INSERT INTO `datasource_column`
    (`app_id`,`schemata_id`,`table_schema`, `table_name`, `data_type`, `character_maximum_length`, `column_name`, `column_comment`,`show_type`,`field_type`,`is_system_field`)
VALUES
    (159,4,'bpmn_app', 'crm_channel', 'int', NULL, 'id', '','num','',1),
    (159,4,'bpmn_app', 'crm_channel', 'int', NULL, 'number', '渠道编号','num','Number',0),
    (159,4,'bpmn_app', 'crm_channel', 'varchar', 255, 'name', '渠道名称','text','Input',0),
    (159,4,'bpmn_app', 'crm_channel', 'int', NULL, 'creator_id', '创建人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_channel', 'int', NULL, 'modifier_id', '修改人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_channel', 'datetime', NULL, 'created_at', '记录创建时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_channel', 'datetime', NULL, 'updated_at', '记录最后修改时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_channel_contact', 'int', NULL, 'id', '','num','',1),
    (159,4,'bpmn_app', 'crm_channel_contact', 'int', NULL, 'channel_id', '渠道id','relation','Records',0),
    (159,4,'bpmn_app', 'crm_channel_contact', 'varchar', 255, 'company_name', '公司名称','text','Input',0),
    (159,4,'bpmn_app', 'crm_channel_contact', 'varchar', 255, 'position', '职务','text','Input',0),
    (159,4,'bpmn_app', 'crm_channel_contact', 'varchar', 255, 'name', '姓名','text','Input',0),
    (159,4,'bpmn_app', 'crm_channel_contact', 'varchar', 20, 'mobile', '联系电话','text','Phone',0),
    (159,4,'bpmn_app', 'crm_channel_contact', 'int', NULL, 'principal', '负责人','user','Member',0),
    (159,4,'bpmn_app', 'crm_channel_contact', 'text', 65535, 'note', '备注','text','Input',0),
    (159,4,'bpmn_app', 'crm_channel_contact', 'int', NULL, 'creator_id', '创建人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_channel_contact', 'int', NULL, 'modifier_id', '修改人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_channel_contact', 'datetime', NULL, 'created_at', '记录创建时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_channel_contact', 'datetime', NULL, 'updated_at', '记录最后修改时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_channel_follow', 'int', NULL, 'id', '','num','',1),
    (159,4,'bpmn_app', 'crm_channel_follow', 'int', NULL, 'channel_id', '渠道id','relation','CascadeControl',0),
    (159,4,'bpmn_app', 'crm_channel_follow', 'int', NULL, 'contact_id', '联系人id','relation','CascadeControl',0),
    (159,4,'bpmn_app', 'crm_channel_follow', 'tinyint', NULL, 'method', '跟进方式','radio','SingleChoice',0),
    (159,4,'bpmn_app', 'crm_channel_follow', 'datetime', NULL, 'follow_at', '跟进时间','date','DateTime',0),
    (159,4,'bpmn_app', 'crm_channel_follow', 'varchar', 255, 'location', '跟进地点','text','Area',0),
    (159,4,'bpmn_app', 'crm_channel_follow', 'text', 65535, 'content', '内容','text','Input',0),
    (159,4,'bpmn_app', 'crm_channel_follow', 'int', NULL, 'visitor_id', '跟进人','user','Member',0),
    (159,4,'bpmn_app', 'crm_channel_follow', 'int', NULL, 'creator_id', '创建人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_channel_follow', 'int', NULL, 'modifier_id', '修改人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_channel_follow', 'datetime', NULL, 'created_at', '记录创建时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_channel_follow', 'datetime', NULL, 'updated_at', '记录最后修改时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_client', 'int', NULL, 'id', '','num','',1),
    (159,4,'bpmn_app', 'crm_client', 'int', NULL, 'channel_id', '渠道id','relation','Records',0),
    (159,4,'bpmn_app', 'crm_client', 'int', NULL, 'number', '编号','num','Number',0),
    (159,4,'bpmn_app', 'crm_client', 'varchar', 255, 'name', '名称','text','Input',0),
    (159,4,'bpmn_app', 'crm_client', 'tinyint', NULL, 'from', '来源','radio','SingleChoice',0),
    (159,4,'bpmn_app', 'crm_client', 'tinyint', NULL, 'status', '客户状态','radio','SingleChoice',0),
    (159,4,'bpmn_app', 'crm_client', 'tinyint', NULL, 'level', '信用等级','radio','SingleChoice',0),
    (159,4,'bpmn_app', 'crm_client', 'varchar', 255, 'tag', '客户标签','text','Input',0),
    (159,4,'bpmn_app', 'crm_client', 'int', NULL, 'principal', '负责人','user','Member',0),
    (159,4,'bpmn_app', 'crm_client', 'varchar', 255, 'address', '地址','text','Area',0),
    (159,4,'bpmn_app', 'crm_client', 'tinyint', NULL, 'type', '生态类别','radio','SingleChoice',0),
    (159,4,'bpmn_app', 'crm_client', 'int', NULL, 'follower', '跟进人','user','Member',0),
    (159,4,'bpmn_app', 'crm_client', 'int', NULL, 'creator_id', '创建人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_client', 'int', NULL, 'modifier_id', '修改人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_client', 'datetime', NULL, 'created_at', '记录创建时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_client', 'datetime', NULL, 'updated_at', '记录最后修改时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_client_contact', 'int', NULL, 'id', '','num','',1),
    (159,4,'bpmn_app', 'crm_client_contact', 'int', NULL, 'client_id', '客户id','relation','Records',0),
    (159,4,'bpmn_app', 'crm_client_contact', 'varchar', 255, 'company_name', '公司名称','text','Input',0),
    (159,4,'bpmn_app', 'crm_client_contact', 'varchar', 255, 'position', '职务','text','Input',0),
    (159,4,'bpmn_app', 'crm_client_contact', 'varchar', 255, 'name', '姓名','text','Input',0),
    (159,4,'bpmn_app', 'crm_client_contact', 'varchar', 20, 'mobile', '联系电话','text','Phone',0),
    (159,4,'bpmn_app', 'crm_client_contact', 'int', NULL, 'principal', '负责人','user','Member',0),
    (159,4,'bpmn_app', 'crm_client_contact', 'text', 65535, 'note', '备注','text','Input',0),
    (159,4,'bpmn_app', 'crm_client_contact', 'int', NULL, 'creator_id', '创建人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_client_contact', 'int', NULL, 'modifier_id', '修改人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_client_contact', 'datetime', NULL, 'created_at', '记录创建时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_client_contact', 'datetime', NULL, 'updated_at', '记录最后修改时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_client_follow', 'int', NULL, 'id', '','num','',1),
    (159,4,'bpmn_app', 'crm_client_follow', 'int', NULL, 'client_id', '客户id','relation','CascadeControl',0),
    (159,4,'bpmn_app', 'crm_client_follow', 'int', NULL, 'contact_id', '联系人id','relation','CascadeControl',0),
    (159,4,'bpmn_app', 'crm_client_follow', 'tinyint', NULL, 'method', '跟进方式','radio','SingleChoice',0),
    (159,4,'bpmn_app', 'crm_client_follow', 'datetime', NULL, 'follow_at', '跟进时间','date','DateTime',0),
    (159,4,'bpmn_app', 'crm_client_follow', 'varchar', 255, 'location', '跟进地点','text','Area',0),
    (159,4,'bpmn_app', 'crm_client_follow', 'text', 65535, 'content', '内容','text','Input',0),
    (159,4,'bpmn_app', 'crm_client_follow', 'int', NULL, 'visitor_id', '跟进人','user','Member',0),
    (159,4,'bpmn_app', 'crm_client_follow', 'int', NULL, 'creator_id', '创建人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_client_follow', 'int', NULL, 'modifier_id', '修改人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_client_follow', 'datetime', NULL, 'created_at', '记录创建时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_client_follow', 'datetime', NULL, 'updated_at', '记录最后修改时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_project', 'int', NULL, 'id', '','num','',1),
    (159,4,'bpmn_app', 'crm_project', 'varchar', 255, 'name', '名称','text','Input',0),
    (159,4,'bpmn_app', 'crm_project', 'int', NULL, 'principal', '负责人','user','Member',0),
    (159,4,'bpmn_app', 'crm_project', 'int', NULL, 'client_id', '客户id','relation','CascadeControl',0),
    (159,4,'bpmn_app', 'crm_project', 'int', NULL, 'contact_id', '联系人id','relation','CascadeControl',0),
    (159,4,'bpmn_app', 'crm_project', 'decimal', NULL, 'amount', '合同金额','num','Amount',0),
    (159,4,'bpmn_app', 'crm_project', 'datetime', NULL, 'start_at', '开始时间','date','DateTime',0),
    (159,4,'bpmn_app', 'crm_project', 'datetime', NULL, 'end_at', '结束时间','date','DateTime',0),
    (159,4,'bpmn_app', 'crm_project', 'tinyint', NULL, 'status', '完成情况','radio','SingleChoice',0),
    (159,4,'bpmn_app', 'crm_project', 'text', 65535, 'note', '备注','text','Input',0),
    (159,4,'bpmn_app', 'crm_project', 'int', NULL, 'creator_id', '创建人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_project', 'int', NULL, 'modifier_id', '修改人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_project', 'datetime', NULL, 'created_at', '记录创建时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_project', 'datetime', NULL, 'updated_at', '记录最后修改时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_project_schedule', 'int', NULL, 'id', '','num','',1),
    (159,4,'bpmn_app', 'crm_project_schedule', 'int', NULL, 'project_id', '项目id','relation','Records',0),
    (159,4,'bpmn_app', 'crm_project_schedule', 'datetime', NULL, 'start_at', '开始时间','date','DateTime',0),
    (159,4,'bpmn_app', 'crm_project_schedule', 'text', 65535, 'note', '项目记录','text','Input',0),
    (159,4,'bpmn_app', 'crm_project_schedule', 'varchar', 255, 'solution', '解决方案','file','File',0),
    (159,4,'bpmn_app', 'crm_project_schedule', 'int', NULL, 'performer', '执行人','user','Number',0),
    (159,4,'bpmn_app', 'crm_project_schedule', 'int', NULL, 'creator_id', '创建人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_project_schedule', 'int', NULL, 'modifier_id', '修改人id','user','Number',1),
    (159,4,'bpmn_app', 'crm_project_schedule', 'datetime', NULL, 'created_at', '记录创建时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_project_schedule', 'datetime', NULL, 'updated_at', '记录最后修改时间','date','DateTime',1),
    (159,4,'bpmn_app', 'crm_client', 'tinyint', NULL, 'follow_status','跟进状态','radio','SingleChoice',0),
    (159,4,'bpmn_app', 'crm_channel', 'int', NULL, 'principal', '负责人','user','Member',0),
    (159,4,'bpmn_app', 'crm_channel', 'text',65535,	'note',	'备注','text','Input',0);


insert into `datasource_column_relation`
    (source_table_id, source_column_id, target_table_id, target_column_id, type)
VALUES
       (33,404,34,421,2),
       (34,421,33,404,3),
       (33,404,35,433,2),
       (35,433,33,404,3),
       (34,420,35,434,2),
       (35,434,34,420,3),
       (36,447,33,404,3),
       (33,404,36,447,2),
       (34,420,36,448,2),
       (36,448,34,420,3),
       (36,444,37,459,2),
       (37,459,36,444,3),
       (33,405,30,373,3),
       (30,373,33,405,2),
       (30,373,31,381,2),
       (31,381,30,373,3),
       (30,373,32,393,2),
       (32,393,30,373,3),
       (31,380,32,394,2),
       (32,394,31,380,3);

insert into `datasource_metadata`
    (app_id, datasource_column_id, `group`, `key`, `value`)
values
       (159,408,'crm_client_from',1,'互联网'),
       (159,408,'crm_client_from',2,'转介绍'),
       (159,408,'crm_client_from',3,'渠道'),
       (159,408,'crm_client_from',4,'独立开发'),
       (159,408,'crm_client_from',5,'展会'),
       (159,408,'crm_client_from',6,'广告'),
       (159,408,'crm_client_from',7,'员工'),
       (159,409,'crm_client_status',1,'潜在'),
       (159,409,'crm_client_status',2,'机会'),
       (159,409,'crm_client_status',3,'正式'),
       (159,409,'crm_client_status',4,'签约'),
       (159,410,'crm_client_level',1,'A'),
       (159,410,'crm_client_level',2,'B'),
       (159,410,'crm_client_level',3,'C'),
       (159,410,'crm_client_level',4,'D'),
       (159,495,'crm_client_follow_status',1,'未联系'),
       (159,495,'crm_client_follow_status',2,'初步联系'),
       (159,495,'crm_client_follow_status',3,'见面拜访'),
       (159,495,'crm_client_follow_status',4,'意向客户'),
       (159,495,'crm_client_follow_status',5,'商务洽谈'),
       (159,495,'crm_client_follow_status',6,'正式报价'),
       (159,495,'crm_client_follow_status',7,'合同签约'),
       (159,495,'crm_client_follow_status',8,'停滞客户'),
       (159,495,'crm_client_follow_status',9,'流失客户'),
       (159,495,'crm_client_follow_status',10,'无意向'),
       (159,414,'crm_client_type',1,'参股公司'),
       (159,414,'crm_client_type',2,'合作伙伴'),
       (159,414,'crm_client_type',3,'甲方客户'),
       (159,414,'crm_client_follow_method',1,'电话拜访'),
       (159,414,'crm_client_follow_method',2,'现场拜访'),
       (159,452,'crm_project_status',1,'发现机会'),
       (159,452,'crm_project_status',2,'可行性分析'),
       (159,452,'crm_project_status',3,'立项'),
       (159,452,'crm_project_status',4,'项目方案'),
       (159,452,'crm_project_status',5,'报价'),
       (159,452,'crm_project_status',6,'赢单'),
       (159,452,'crm_project_status',7,'结束'),
       (159,395,'crm_channel_follow_method',1,'电话拜访'),
       (159,395,'crm_channel_follow_method',2,'现场拜访');

