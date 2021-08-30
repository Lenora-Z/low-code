use `default-bpmn`;
drop table if exists `router`;
create table `router`(
    `id` int(11) unsigned not null auto_increment,
    `created_at` datetime NOT NULL,
    `updated_at` datetime DEFAULT NULL,
    `nav_id` int(11) unsigned not null default 0 comment '路由id',
    `title` varchar(100) not null default '' comment '路由名称',
    `parent_id` int(11) unsigned not null default 0 comment '父级路由',
    `key` varchar(20) not null default 0 comment 'key值',
    `form_id` int(11) unsigned not null default 0 comment '对应的表单id',
    `icon` varchar(255) DEFAULT NULL COMMENT '路由icon',
    `status` tinyint(1) unsigned not null default 0 comment '生效状态 0-未生效 1-已生效',

    primary key (`id`),
    key `nav` (`nav_id`) using btree ,
    key `route_key` (`key`) using btree
)engine=InnoDB default charset =utf8 comment '路由表';


update field set type = 'Input' where type = 'Text';

alter table `user`
    add column `pwd_status` tinyint(2) not null default 0 comment '1已重置' after group_id;