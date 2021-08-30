use `bpmn`;
alter table form
    add column `page_type` tinyint(1) not null default 1 comment '页面类型 1-表单 2-展示页' after type,
    add column `from` int(11) unsigned not null default 0 comment '原表' after app_id;

alter table `user`
    modify column `true_name` varchar(255) not null default '' comment '真实姓名';

alter table `flow`
    add column assignee varchar(255) not null default '' comment '处理用户' after service_group;

create table `field_multi_form`(
    `form_id` int(11) unsigned not null default 0 comment '父表单id',
    `field_id` int(11) unsigned not null default 0 comment '控件id',
    `child_id` int(11) unsigned not null default 0 comment '子表单id',
    `method` tinyint(1) unsigned not null default 1 comment '表单引入方式 1-关联 2-原表复制',
    `mode` tinyint(1) unsigned not null default 1 comment '呈现样式 1-单条记录 2-多条记录 3-表格',

    primary key (form_id,child_id),
    key `field` (field_id) using btree
)engine=InnoDB default CHARSET=utf8 comment '多表单控件属性表';

create table `flow_assignee`(
    `user_id` int(11) unsigned not null default 0 comment '用户id',
    `activity` varchar(255) not null default '' comment '任务id',
    `flow_id` int(11) unsigned not null default 0 comment '流程id'
)engine=InnoDB default CHARSET=utf8 comment '用户流程关联表';

create table `field_linkage`(
    `form_id` int(11) unsigned not null default 0 comment '表单id',
    `field_id` int(11) unsigned not null default 0 comment '级联控件id',
    `first_id` int(11) unsigned not null default 0 comment '一级表单id',
    `key` varchar(50) not null default '' comment '一级控件key',
    `content` text comment '级联关系,格式:表单id#控件key,以&连接',
    key (form_id) using btree
)engine=InnoDB default CHARSET=utf8 comment '级联控件属性表';

alter table `flow`
    add column notifier varchar(255) not null default '' comment '抄送用户列表' after assignee;

create table `flow_notifier`(
    `user_id` int(11) unsigned not null default 0 comment '用户id',
    `activity` varchar(255) not null default '' comment '任务id',
    `flow_id` int(11) unsigned not null default 0 comment '流程id'
)engine=InnoDB default CHARSET=utf8 comment '抄送用户与流程关联表';

update application set `status` = 2;

update version set `status` = 2;

update navigation set `status` = 0,`is_online` = 0;

update form set `status` = 0,`is_online` = 0,`page_type` = 0;

update flow set `status` = 0,`is_online` = 0;
