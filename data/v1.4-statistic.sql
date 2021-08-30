use bpmn_app;

drop table if exists `hrm_employee_statistic`;
create table `hrm_employee_statistic`(
    `year` int(11) unsigned not null default 0 comment '',
    `month` int(11) unsigned not null default 0 comment '',
    `entry_num` int(11) unsigned not null default 0 comment '',
    `dimission_num` int(11) unsigned not null default 0 comment '',
    primary key (`year`,`month`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='离/入职人数统计表';
