create schema `bpmn_app` collate utf8_general_ci;

use `bpmn_app`;

create table bpmn_app.hrm_applicant
(
    id int(11) unsigned auto_increment comment '应聘者id'
        primary key,
    name char(50) default '' not null comment '姓名',
    gender tinyint(1) unsigned default 0 not null comment '性别(1:男;2:女)',
    birthday datetime null comment '出生年月',
    education tinyint(1) unsigned default 0 not null comment '学历(1:博士;2:硕士;3:本科;4:大专;5:大专以下)',
    college char(30) default '' not null comment '毕业学校',
    major char(30) default '' not null comment '专业',
    email char(50) default '' not null comment '邮箱',
    recruit_type tinyint(1) unsigned default 0 not null comment '求职类型(1:实习;2:兼职;3:正式)',
    recruit_job char(30) default '' not null comment '求职岗位',
    resume_path varchar(255) default '' not null comment '简历附件地址',
    status tinyint(1) unsigned default 1 not null comment '招聘状态(1:未面试;2:面试中;3:不通过;4:待发offer;5:已发offer;6:待入职;7:拒绝offer;8:未入职;9:已入职)',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment '应聘者表' charset=utf8mb4;

create table bpmn_app.hrm_department
(
    id int(11) unsigned auto_increment comment '部门id'
        primary key,
    name char(20) default '' not null comment '部门名称',
    establish_time datetime null comment '成立时间',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment '部门表' charset=utf8mb4;

create table bpmn_app.hrm_employee
(
    id int(11) unsigned auto_increment comment '员工id'
        primary key,
    applicant_id char(100) default '0' not null comment '应聘者id',
    job_number char(10) default '' not null comment '工号',
    name char(50) default '' not null comment '姓名',
    gender tinyint(1) unsigned default 0 not null comment '性别(1:男;2:女)',
    department_id char(100) default '' not null comment '部门id',
    job char(30) default '' not null comment '岗位',
    grade char(20) default '' not null comment '职级',
    entry_time datetime null comment '入职时间',
    probation char(30) default '' not null comment '试用期',
    phone_number char(15) default '' not null comment '联系电话',
    identity_number char(20) default '' not null comment '身份证号',
    political_status char(20) default '' not null comment '政治面貌',
    become_time datetime null comment '转正时间',
    is_married tinyint(1) unsigned default 0 not null comment '婚姻状况(1:未婚;2:已婚)',
    education tinyint(1) unsigned default 0 not null comment '学历(1:博士;2:硕士;3:本科;4:大专;5:大专以下)',
    college char(30) default '' not null comment '毕业学校',
    major char(30) default '' not null comment '专业',
    hukou_type tinyint(1) unsigned default 4 not null comment '户口性质(1:农业户口;2:城镇户口;3:集体户口;4:其他)',
    insurance_base int(11) unsigned default 0 not null comment '社保基数',
    fund_base int(11) unsigned default 0 not null comment '公积金基数',
    identity_address varchar(255) default '' not null comment '身份证地址',
    identity_expired_time datetime null comment '身份证有效期',
    address varchar(255) default '' not null comment '现地址',
    email char(50) default '' not null comment '邮箱',
    bank_name char(30) default '' not null comment '工资卡开户行',
    bank_account char(30) default '' not null comment '银行账户',
    certificate_name varchar(255) default '' not null comment '各类资质证书名称(多个用‘,’分隔)',
    certificate_path varchar(255) default '' not null comment '各类资质证书地址(多个用‘,’分隔)',
    status tinyint(1) unsigned default 3 not null comment '入职后状态(1:实习;2:兼职;3:正式;4:试用期;5:待离职;6:已离职)',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment '员工表' charset=utf8mb4;

create table bpmn_app.hrm_employee_statistic
(
    year int(11) unsigned default 0 not null,
    month int(11) unsigned default 0 not null,
    entry_num int(11) unsigned default 0 not null,
    dimission_num int(11) unsigned default 0 not null,
    primary key (year, month)
)
    comment '离/入职人数统计表' charset=utf8;

create table bpmn_app.hrm_interview
(
    id int(11) unsigned auto_increment comment '面试记录id'
        primary key,
    applicant_id char(100) default '' not null comment '应聘者id',
    interview_time datetime null comment '面试时间',
    interviewer_id char(100) default '' not null comment '面试官id',
    evaluation varchar(255) default '' not null comment '面试评价',
    is_next tinyint(1) unsigned default 1 not null comment '是否有下一轮面试(0:否;1:是)',
    status tinyint(1) unsigned default 1 not null comment '面试结果(1:待填写;2:不通过;3:通过)',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment '面试记录表' charset=utf8mb4;

create table bpmn_app.hrm_labor_contract
(
    id int(11) unsigned auto_increment comment '劳动合同记录id'
        primary key,
    employee_id char(100) default '' not null comment '员工id',
    code char(20) default '' not null comment '合同编号',
    sum int(11) unsigned default 2 not null comment '存档份数，默认为2',
    sign_time datetime null comment '签订时间',
    expired_time datetime null comment '到期时间',
    receive_time datetime null comment '领取时间',
    is_received tinyint(1) unsigned default 0 not null comment '是否领取(0:待领取;1:已领取)',
    status tinyint(1) unsigned default 1 not null comment '合同状态(1:正常;2:将到期;3:已过期;4:中止)',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment '劳动合同记录表' charset=utf8mb4;

create table bpmn_app.hrm_offer
(
    id int(11) unsigned auto_increment comment 'offer记录id'
        primary key,
    applicant_id char(100) default '' not null comment '应聘者id',
    department_id char(100) default '' not null comment '拟录用部门id',
    job char(30) default '' not null comment '拟录用岗位',
    manager_id char(100) default '' not null comment '汇报对象id',
    work_place char(100) default '' not null comment '工作地点',
    before_salary int default 0 not null comment '试用期薪资总包',
    before_basic_salary int default 0 not null comment '试用期基本工资',
    before_merits_salary int default 0 not null comment '试用期绩效工资',
    probation char(30) default '' not null comment '试用期',
    after_salary int default 0 not null comment '转正后薪资总包',
    after_basic_salary int default 0 not null comment '转正后基本工资',
    after_merits_salary int default 0 not null comment '转正后绩效工资',
    retain_time datetime null comment '保留日期',
    board_time datetime null comment '到岗日期',
    send_time datetime null comment '发送时间',
    sender_id char(100) default '' not null comment '发送人id',
    status tinyint(1) unsigned default 1 not null comment 'offer状态(1:待发放;2:待回复;3:已接受;4:已拒绝)',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment 'offer管理表' charset=utf8mb4;

create table bpmn_app.hrm_resignation
(
    id int(11) unsigned auto_increment comment '离职记录id'
        primary key,
    employee_id char(100) default '' not null comment '员工id',
    apply_time datetime null comment '申请时间',
    resign_time datetime null comment '离职时间',
    reason varchar(255) default '' not null comment '离职原因',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment '离职记录表' charset=utf8mb4;

create table bpmn_app.hrm_talent_pool
(
    id int(11) unsigned auto_increment comment '人才id'
        primary key,
    applicant_id char(100) default '' not null comment '应聘者id',
    level char(20) default '' not null comment '人才等级(一般人才、重要人才)',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment '人才库表' charset=utf8mb4;

create table bpmn_app.hrm_train
(
    id int(11) unsigned auto_increment comment '培训记录id'
        primary key,
    name char(50) default '' not null comment '培训名称',
    code char(20) default '' not null comment '培训编号',
    start_time datetime null comment '培训开始时间',
    end_time datetime null comment '培训结束时间',
    dept_org_id char(100) default '' not null comment '组织培训的部门id',
    dept_join_id char(100) default '' not null comment '参与培训的部门id',
    budget int(11) unsigned default 0 not null comment '培训预算(元)',
    goal varchar(255) default '' not null comment '培训宗旨',
    material_path varchar(255) default '' not null comment '培训资料附件地址',
    creator_id int(11) unsigned not null comment '记录创建人id',
    comment varchar(255) default '' not null comment '备注',
    created_at datetime default CURRENT_TIMESTAMP not null comment '数据创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '数据更新时间'
)
    comment '培训记录表' charset=utf8mb4;

insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (1,'董事长（总裁now()办公室',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (2,'投融资部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (3,'综合资源部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (4,'财务部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (5,'人力资源部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (6,'法务部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (7,'品牌宣传部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (8,'平台运营中心',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (9,'渠道部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (10,'市场部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (11,'解决方案中心',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (12,'研究院',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (13,'数字政法事业部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (14,'研发一部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (15,'研发二部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (16,'研发三部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (17,'研发五部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (18,'数据与信息维护中心',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (19,'数字农业事业部',now(),19);
insert into bpmn_app.hrm_department (id,name,establish_time,creator_id) values (20,'国际与创新事业部',now(),19);

insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0001','高航','董事长',1,'18606525678',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0006','俞学劢','CEO',1,'13575481734',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0026','陆静怡','董秘',2,'15801722182',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0093','倪帆环','董事长助理',1,'13777470377',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0455','董倩','董事长助理',2,'13588186969',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0463','孔清扬','董事长助理（实习）',2,'13718572250',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0239','唐斌','总裁',1,'18606998338',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0309','蒋俊青','副总裁',1,'13857171613',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0329','帅磊','总监助理',2,'18258122415',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (1,'0252','濮笑威','资源管理经理',2,'18758231921',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (2,'0389','蒋利峰','投融资VP',1,'13811500036',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0246','张云霆','综合资源总经理',1,'13868008119',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0333','俞斐','综合资源副总经理',2,'13386521312',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0255','肖靖','品牌运营经理',1,'18458132839',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0003','张苏丽','行政总监',2,'13588230107',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0109','洪静芬','行政助理',2,'18767136016',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0471','陈若钰','前台',2,'19818532302',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0501','黄路妹','人事行政专员',2,'15068864912',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0514','章冰心','前台',2,'18297693526',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0119','李方华','保洁',2,'13515710615',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (3,'0009','周陈','司机',1,'15824452380',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (4,'0005','汪令','财务总监',2,'15757131672',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (4,'0027','高银霖','会计',2,'15005810906',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (4,'0153','侯雁翎','会计',2,'18758322584',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (4,'0151','郑慧琴','会计助理',2,'18623077366',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (4,'0335','周怡铃','出纳',2,'15857166514',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (4,'0370','韩天','财务经理',1,'15858178915',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (5,'0376','陈怡坚','人事总监',1,'18857876339',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (5,'0033','金晶','人事经理',2,'13757185932',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (5,'0373','李旭玲','招聘专员',2,'15700151759',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (5,'0388','柳治昊','人事专员',1,'13221802287',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (5,'0438','周一颖','人事经理',2,'13588359983',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (5,'0466','周锦','HRBP',2,'18867111060',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (6,'0314','高秀春','法务',2,'13777382226',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (6,'0475','戴淑娟','法务助理',2,'18895319339',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (7,'0117','侯鲁阳','市场运营',1,'18037966957',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (7,'0273','赵莉','社群运营',2,'15757124820',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (7,'0380','袁雯芬','文案策划',2,'13732256675',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (7,'0473','王涛','平面设计师',1,'17682302562',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (7,'0506','李雅诗','内容运营',2,'18958007052',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (8,'0010','俞佳楠','负责人',2,'18268196637',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (8,'0097','张志波','市场运营',1,'18072700096',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (8,'0420','孙百灵','新媒体运营',2,'17762098066',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (8,'0450','单寒雪','UI设计师',2,'18201188727',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (9,'0022','冯晔','副总裁',1,'13601026747',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (9,'0243','张策','解决方案专家',1,'18500093895',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (9,'0440','程昊','渠道经理',1,'13675796731',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (9,'0428','周洋','渠道经理',1,'13634101220',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (9,'0476','安彬','总裁助理',1,'18167106981',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (9,'0339','冠景曦','商务助理',1,'13500737606',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (9,'0488','王居钦','实习生',1,'18258288411',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (10,'0011','沈海华','商务总监',2,'13777818200',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (10,'0204','杨慧梁','商务经理',1,'13858175992',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (10,'0247','余俊杰','商务经理',1,'18072857826',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (10,'0467','吴雷军','商务经理',1,'15158069360',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (10,'0491','苏立滔','商务经理',1,'13575775722',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (10,'0025','孙明明（鉴信）','商务经理',1,'15868862009',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (11,'0300','杨柳军','解决方案专家',1,'13656651841',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (11,'0299','宋慧然','售前工程师',1,'18989893675',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (11,'0170','金晖','售前工程师',1,'18667181992',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (11,'0441','吴伟钟','售前工程师',1,'18969142822',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (12,'0015','张金琳','科学家',2,'13735531523',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (12,'0016','季姝','研究员',2,'18266967500',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (12,'0350','岳高','研究员',1,'18740593317',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (12,'0439','李峰','知识产权专员',1,'13645710378',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (12,'0485','王继莲','研究院助理（实习）',2,'18872943530',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0014','陈豪鸣','副总裁',1,'18758303304',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0394','吴元子','VP（兼职）',1,'18999906133',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0395','吴雨薇','商务经理（兼职）',2,'13262971917',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0512','高畅','商务经理（兼职）',2,'18916178670',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0056','马露鹏','产品经理',1,'15990055807',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0072','胡兰心','产品经理',2,'18883992815',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0321','陈婉菲','产品助理',2,'15927651749',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0186','骆嘉隽','产品经理',1,'15957126884',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0477','顾志鹏','产品经理',1,'15088706079',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0497','陈如意','产品助理实习',2,'17865572628',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0500','陈路杰','产品助理实习',1,'19858186613',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0381','裴宏达','高级运营专家',1,'18072804897',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0127','楼莹','运营主管',2,'13758204655',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0216','毋玉针','运营助理',2,'17788569475',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0268','彭珍','运营支持',2,'15906120409',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0343','张真真','运营支持',2,'17805808852',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0365','韩媛','运营支持',2,'18758143367',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0414','喻龙','运营支持',1,'17737664620',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0415','钮航','运营支持',2,'17764572085',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0421','高音璐','运营支持',2,'13372527170',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0470','李辉','运营支持',1,'17606507054',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0495','王文杰','运营支持',1,'13073694279',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0241','肖军','直销主管',1,'15088620457',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0169','朱飞祥','销售经理',1,'15557167320',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0221','张知火','销售经理',1,'15979904706',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0276','吴秋','销售经理',1,'15088604015',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0326','卓峰','销售经理',1,'13958114421',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0410','陈志伟','商务经理',1,'15988467406',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0177','郝玉超','政法商务',1,'18610938500',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0442','郑建成','商务经理',1,'13732256361',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (13,'0294','李怡','商务经理',2,'18601349468',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0031','张文勇','负责人',1,'18167127094',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0081','杨大敏','区块链研发工程师',1,'13918450498',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0259','邱璐','前端工程师',2,'18779105742',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0260','张学政','Golong开发',1,'13253538613',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0213','褚德权','前端工程师',1,'18815287483',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0219','曹可磊','后端开发',1,'18868803292',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0184','张翔','php开发',2,'13588852454',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0083','钟小飞','测试工程师',2,'18858502363',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0073','李文韬','UI设计师',1,'15906644473',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0318','方菲','UI设计师',2,'18119607520',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0330','屠翔','数据产品经理',1,'18958482327',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0349','孙帅帅','前端工程师',1,'15055938613',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0352','沈惠良','前端工程师',1,'18758020025',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0353','史祥宇','node.js开发',1,'17826833657',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0363','董源灏','产品助理',1,'13905712235',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0347','王玉霞','大数据开发工程师',2,'19957010696',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0379','朱李姣','测试工程师',2,'13971240376',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0416','严清鑫','测试工程师',2,'19947884532',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0423','刘思琦','产品经理',2,'19560196277',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0426','张亚朋','测试工程师',1,'18939461688',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0462','冯杰','产品经理',1,'13957141918',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0472','孙奕铖','后端开发',1,'19855036641',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0478','刘浩','后端开发',1,'17612770767',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0482','王小凤','前端开发',2,'18226618639',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0486','潘高峰','java开发',1,'13989472823',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0490','王芑','后端开发',1,'18995942890',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0493','程伟森','后端开发',1,'15330238219',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0494','张梓铭','测试开发',1,'18857115981',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0507','童家美','测试开发',2,'15258834522',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0516','王志彬','java开发',1,'15057156571',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0266','苏波','解决方案专家',1,'13611307375',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0148','谢振元','后端开发',1,'13381029469',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0227','薄尊旭','后端开发',1,'18211037059',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (14,'0238','王守伟','区块链开发',1,'13521059527',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0055','何锐','技术总监',1,'13810178103',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0108','孙亮','php开发工程师',1,'15810302468',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0124','李雪毓','前端开发工程师',2,'13141171998',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0233','张天乐','后端PHP开发',1,'13717857221',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0456','汤继康','后端PHP开发',1,'15611752059',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0253','王仁杰','测试工程师',2,'13691332631',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0237','马学超','UI设计师',1,'15869179380',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0214','李现平','后端工程师',1,'15669983656',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0188','李昂','后端工程师',1,'18506819175',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0279','曹慧慧','产品经理',2,'15715727764',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0336','王杰','前端工程师',1,'17367073121',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0358','陈厚斌','项目经理',1,'15336540995',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0372','李磊','运维工程师',1,'13858054558',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0391','邓珂','测试工程师',1,'18237700291',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0403','徐超超','测试工程师',1,'18658176223',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0404','王世豪','JAVA开发工程师',1,'17639845335',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0407','唐宗胜','JAVA开发工程师',1,'18715116356',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0408','周馨远','产品经理',2,'15291830587',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0409','左晨宇','测试工程师',1,'17355280829',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0413','胡文林','前端工程师',1,'15515305037',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0424','杜梵','前端工程师',1,'15068057219',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0444','朱秋月','UI设计师',2,'13588176553',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0460','何浩','大数据工程师',1,'18658163713',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0468','孙革兵','架构师',1,'13320038505',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0499','王峰','后端开发',1,'17764544711',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0502','杨文寿','前端开发' ,1,'15372000059',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (15,'0508','白森','Java开发',1,'13123915439',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0320','王广磊','技术主管',1,'18057182101',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0035','屈朋辉','区块链研发',1,'18106500602',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0360','唐小立','后端工程师',1,'13588801445',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0367','黄忠成','JAVA开发工程师',1,'18658109550',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0332','陈威镇','运维工程师',1,'18956851661',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0342','王坤杰','运维工程师',1,'17798689096',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0023','刘尧尧','区块链研发',1,'13145205611',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0012','吴燕','前端工程师',2,'18368729972',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0013','田柳','UI设计师',2,'18868815710',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0059','孙宽慰','爬虫工程师',1,'15957156851',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0061','韩换飞','爬虫工程师',1,'17521001703',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0224','李宗辉','前端工程师',1,'15136065755',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0265','禹金鹏','python工程师',1,'15088735720',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0280','殷彦香','测试工程师',2,'15858261021',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0297','徐文斌','前端工程师',1,'13666604580',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0175','谢磊','php开发转java',1,'18368033334',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0341','王彦彬','安卓工程师',1,'13606524556',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0348','方明旺','ios开发工程师',1,'18668157573',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0251','罗巧巧','前端实习',2,'15659622723',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0374','刘超','架构师',1,'15958166963',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0393','王倩','前端工程师',2,'18969147338',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0401','仰红霞','测试工程师',2,'13516722602',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0412','翟君慧','项目经理',2,'15618969149',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0417','李涛涛','项目经理',1,'19905814214',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0418','陈伟根','java工程师',1,'18668133519',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0419','刘洋洋','php开发工程师',1,'15936061620',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0429','陈昱灯','java工程师',1,'13758587061',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0430','刘锦祥','前端工程师',1,'15957195869',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0446','潘婷','UI设计师',2,'15088995205',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0448','郭艳','测试工程师',2,'18600464071',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0449','陈燕青','测试工程师',2,'17605880738',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0451','傅豪朋','python工程师',1,'15137728278',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0459','李科','java工程师',1,'17630605618',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0479','梁奇峰','java工程师',1,'13588795841',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0483','刘加圣','大数据实施工程师',1,'17605810844',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0489','杨统','前端实习',1,'13989539748',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0510','张铁宝','java工程师',1,'15500180069',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0515','何林','python工程师',1,'18680325804',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0503','任奕柯','后端工程师',2,'13651629961',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0504','宋博','前端工程师',1 ,'19939209995',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (16,'0517','许长海','测试工程师',1,'15618385363',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0432','赵冲','负责人',1,'15157010666',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0433','何明飞','php开发工程师',1,'18367803500',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0434','李晓伟','php开发工程师',1,'15958042925',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0435','郭鑫杰','java工程师',1,'17682308036',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0436','高波','前端、产品、项目',1,'18367826026',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0437','丁天立','产品经理、项目、商务',1,'13958178831',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0457','郭晶璐','ui设计师',2,'15035731735',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0458','马世根','ui设计师',1,'13003682811',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0464','许佳辰','java工程师',1,'17600792030',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0474','徐鹏元','Java工程师',1,'17760715873',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0480','柴春梅','前端工程师',2,'18868191902',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0481','熊小亮','前端工程师',1,'15279216013',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0484','刘明','前端工程师',1,'15958043011',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0498','刘西菩','后端开发',2,'13925210349',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0505','刘培武','测试工程师',1,'17603841625',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0509','焦涵','产品经理',1,'18741270380',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0511','王磊','Java工程师',1,'13263112656',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (17,'0513','周波','测试工程师',1,'18629436016',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (18,'0077','张凯','运维工程师',1,'18297959138',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (18,'0288','李鹏程','运维工程师',1,'18701320197',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (18,'0315','孙海洋','项目管理',1,'17767130930',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (18,'0445','赵杏嶂','运维工程师',1,'18238674490',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (19,'0378','崔伟','首席科学家、VP',1,'13601300572',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (19,'0245','齐林杰','产品经理',1,'13279206067',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (19,'0218','周佳丽','副总裁助理',2,'15258290380',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (19,'0461','叶前程','商务经理',1,'18800303278',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (19,'0496','潘钟瑞','助理实习生',1,'18873086687',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (19,'0196','曾凡','商务经理',1,'15011223779',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (20,'0209','马超','国际业务助理',1,'13109570606',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (20,'0469','沈也','产品经理',2,'18268022277',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (20,'0492','胡周平','产品经理',1,'18268070310',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (20,'0047','钟建','运维主管（兴义矿场）',1,'18072836356',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (20,'0051','何顺岗','运维专员（后旗）',1,'13967021335',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (20,'0066','罗勇','机房电工（兴义）',1,'18788760226',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (20,'0080','刘仕兵','运维专员（兴义）',1,'18785984307',null,null,null,19);
insert into `bpmn_app`.hrm_employee (department_id,job_number,name,job,gender,phone_number,entry_time,become_time,identity_expired_time,creator_id) values (20,'','张静宇','运维助理',2,'15024963914',null,null,null,19);

INSERT INTO `bpmn`.datasource_schemata (id, app_id, table_schema) VALUES (1, 80, 'bpmn_app');

INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (19, 80, 1, 'bpmn_app', 'hrm_applicant', '应聘者表');
INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (20, 80, 1, 'bpmn_app', 'hrm_department', '部门表');
INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (21, 80, 1, 'bpmn_app', 'hrm_employee', '员工表');
INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (22, 80, 1, 'bpmn_app', 'hrm_interview', '面试记录表');
INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (23, 80, 1, 'bpmn_app', 'hrm_labor_contract', '劳动合同记录表');
INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (24, 80, 1, 'bpmn_app', 'hrm_offer', 'offer管理表');
INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (25, 80, 1, 'bpmn_app', 'hrm_resignation', '离职记录表');
INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (26, 80, 1, 'bpmn_app', 'hrm_talent_pool', '人才库表');
INSERT INTO `bpmn`.datasource_table (id, app_id, schemata_id, table_schema, table_name, table_comment) VALUES (27, 80, 1, 'bpmn_app', 'hrm_train', '培训记录表');

INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (15, 80, 125, 'hrm_applicant_gender', '1', '男');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (16, 80, 127, 'hrm_applicant_education', '1', '博士');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (17, 80, 130, 'hrm_applicant_recruit_type', '1', '实习');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (18, 80, 133, 'hrm_applicant_status', '1', '未面试');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (19, 80, 149, 'hrm_employee_gender', '1', '男');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (20, 80, 158, 'hrm_employee_is_married', '1', '未婚');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (21, 80, 159, 'hrm_employee_education', '1', '博士');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (22, 80, 162, 'hrm_employee_hukou_type', '1', '农业户口');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (23, 80, 173, 'hrm_employee_status', '1', '实习');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (24, 80, 183, 'hrm_interview_is_next', '0', '否');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (25, 80, 184, 'hrm_interview_status', '1', '待填写');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (26, 80, 196, 'hrm_interview_is_received', '0', '待领取');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (27, 80, 197, 'hrm_labor_contract_status', '1', '正常');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (28, 80, 210, 'hrm_offer_status', '1', '待发放');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (68, 80, 125, 'hrm_applicant_gender', '2', '女');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (69, 80, 127, 'hrm_applicant_education', '2', '硕士');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (70, 80, 127, 'hrm_applicant_education', '3', '本科');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (71, 80, 127, 'hrm_applicant_education', '4', '大专');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (72, 80, 127, 'hrm_applicant_education', '5', '大专以下');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (73, 80, 130, 'hrm_applicant_recruit_type', '2', '兼职');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (74, 80, 130, 'hrm_applicant_recruit_type', '3', '正式');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (75, 80, 133, 'hrm_applicant_status', '2', '面试中');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (76, 80, 133, 'hrm_applicant_status', '3', '不通过');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (77, 80, 133, 'hrm_applicant_status', '4', '待发offer');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (78, 80, 133, 'hrm_applicant_status', '5', '已发offer');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (79, 80, 133, 'hrm_applicant_status', '6', '待入职');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (80, 80, 133, 'hrm_applicant_status', '7', '拒绝offer');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (81, 80, 133, 'hrm_applicant_status', '8', '未入职');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (82, 80, 133, 'hrm_applicant_status', '9', '已入职');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (83, 80, 149, 'hrm_employee_gender', '2', '女');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (84, 80, 158, 'hrm_employee_is_married', '2', '已婚');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (85, 80, 159, 'hrm_employee_education', '2', '硕士');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (86, 80, 159, 'hrm_employee_education', '3', '本科');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (87, 80, 159, 'hrm_employee_education', '4', '大专');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (88, 80, 159, 'hrm_employee_education', '5', '大专以下');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (89, 80, 162, 'hrm_employee_hukou_type', '2', '城镇户口');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (90, 80, 162, 'hrm_employee_hukou_type', '3', '集体户口');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (91, 80, 162, 'hrm_employee_hukou_type', '4', '其他');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (92, 80, 173, 'hrm_employee_status', '2', '兼职');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (93, 80, 173, 'hrm_employee_status', '3', '正式');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (94, 80, 173, 'hrm_employee_status', '4', '试用期');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (95, 80, 173, 'hrm_employee_status', '5', '待离职');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (96, 80, 173, 'hrm_employee_status', '6', '已离职');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (97, 80, 183, 'hrm_interview_is_next', '1', '是');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (98, 80, 184, 'hrm_interview_status', '2', '不通过');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (99, 80, 184, 'hrm_interview_status', '3', '通过');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (100, 80, 196, 'hrm_interview_is_received', '1', '已领取');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (101, 80, 197, 'hrm_labor_contract_status', '2', '将到期');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (102, 80, 197, 'hrm_labor_contract_status', '3', '已过期');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (103, 80, 197, 'hrm_labor_contract_status', '4', '中止');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (104, 80, 210, 'hrm_offer_status', '2', '待回复');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (105, 80, 210, 'hrm_offer_status', '3', '已接受');
INSERT INTO `bpmn`.datasource_metadata (id, app_id, datasource_column_id, `group`, `key`, value) VALUES (106, 80, 210, 'hrm_offer_status', '4', '已拒绝');

INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (123, 80, 1, 'bpmn_app', 'hrm_applicant', 'int', null, 'id', '应聘者id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (124, 80, 1, 'bpmn_app', 'hrm_applicant', 'char', 50, 'name', '姓名', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (125, 80, 1, 'bpmn_app', 'hrm_applicant', 'tinyint', null, 'gender', '性别', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (126, 80, 1, 'bpmn_app', 'hrm_applicant', 'datetime', null, 'birthday', '出生年月', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (127, 80, 1, 'bpmn_app', 'hrm_applicant', 'tinyint', null, 'education', '学历', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (128, 80, 1, 'bpmn_app', 'hrm_applicant', 'char', 30, 'college', '毕业学校', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (129, 80, 1, 'bpmn_app', 'hrm_applicant', 'char', 30, 'major', '专业', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (130, 80, 1, 'bpmn_app', 'hrm_applicant', 'tinyint', null, 'recruit_type', '求职类型', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (131, 80, 1, 'bpmn_app', 'hrm_applicant', 'char', 30, 'recruit_job', '求职岗位', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (132, 80, 1, 'bpmn_app', 'hrm_applicant', 'varchar', 255, 'resume_path', '简历附件地址', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (133, 80, 1, 'bpmn_app', 'hrm_applicant', 'tinyint', null, 'status', '招聘状态', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (134, 80, 1, 'bpmn_app', 'hrm_applicant', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (135, 80, 1, 'bpmn_app', 'hrm_applicant', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (136, 80, 1, 'bpmn_app', 'hrm_applicant', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (137, 80, 1, 'bpmn_app', 'hrm_applicant', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (369, 80, 1, 'bpmn_app', 'hrm_applicant', 'char', 50, 'email', '邮箱', 'text', 'Mail', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (138, 80, 1, 'bpmn_app', 'hrm_department', 'int', null, 'id', '部门id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (139, 80, 1, 'bpmn_app', 'hrm_department', 'char', 20, 'name', '部门名称', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (140, 80, 1, 'bpmn_app', 'hrm_department', 'datetime', null, 'establish_time', '成立时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (141, 80, 1, 'bpmn_app', 'hrm_department', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (142, 80, 1, 'bpmn_app', 'hrm_department', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (143, 80, 1, 'bpmn_app', 'hrm_department', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (144, 80, 1, 'bpmn_app', 'hrm_department', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (145, 80, 1, 'bpmn_app', 'hrm_employee', 'int', null, 'id', '员工id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (146, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 100, 'applicant_id', '应聘者id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (147, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 10, 'job_number', '工号', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (148, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 50, 'name', '姓名', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (149, 80, 1, 'bpmn_app', 'hrm_employee', 'tinyint', null, 'gender', '性别', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (150, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 100, 'department_id', '部门id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (151, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 30, 'job', '岗位', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (152, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 20, 'grade', '职级', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (153, 80, 1, 'bpmn_app', 'hrm_employee', 'datetime', null, 'entry_time', '入职时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (154, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 30, 'probation', '试用期', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (155, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 15, 'phone_number', '联系电话', 'text', 'Phone', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (156, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 20, 'identity_number', '身份证号', 'text', 'Certificates', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (157, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 20, 'political_status', '政治面貌', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (158, 80, 1, 'bpmn_app', 'hrm_employee', 'tinyint', null, 'is_married', '婚姻状况', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (159, 80, 1, 'bpmn_app', 'hrm_employee', 'tinyint', null, 'education', '学历', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (160, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 30, 'college', '毕业学校', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (161, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 30, 'major', '专业', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (162, 80, 1, 'bpmn_app', 'hrm_employee', 'tinyint', null, 'hukou_type', '户口性质', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (163, 80, 1, 'bpmn_app', 'hrm_employee', 'int', null, 'insurance_base', '社保基数', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (164, 80, 1, 'bpmn_app', 'hrm_employee', 'int', null, 'fund_base', '公积金基数', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (165, 80, 1, 'bpmn_app', 'hrm_employee', 'varchar', 255, 'identity_address', '身份证地址', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (166, 80, 1, 'bpmn_app', 'hrm_employee', 'datetime', null, 'identity_expired_time', '身份证有效期', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (167, 80, 1, 'bpmn_app', 'hrm_employee', 'varchar', 255, 'address', '现地址', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (168, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 50, 'email', '邮箱', 'text', 'Mail', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (169, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 30, 'bank_name', '工资卡开户行', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (170, 80, 1, 'bpmn_app', 'hrm_employee', 'char', 30, 'bank_account', '银行账户', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (171, 80, 1, 'bpmn_app', 'hrm_employee', 'varchar', 255, 'certificate_name', '各类资质证书名称', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (172, 80, 1, 'bpmn_app', 'hrm_employee', 'varchar', 255, 'certificate_path', '各类资质证书地址', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (173, 80, 1, 'bpmn_app', 'hrm_employee', 'tinyint', null, 'status', '入职后状态', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (174, 80, 1, 'bpmn_app', 'hrm_employee', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (175, 80, 1, 'bpmn_app', 'hrm_employee', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (176, 80, 1, 'bpmn_app', 'hrm_employee', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (177, 80, 1, 'bpmn_app', 'hrm_employee', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (372, 80, 1, 'bpmn_app', 'hrm_employee', 'datetime', null, 'become_time', '转正时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (178, 80, 1, 'bpmn_app', 'hrm_interview', 'int', null, 'id', '面试记录id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (179, 80, 1, 'bpmn_app', 'hrm_interview', 'char', 100, 'applicant_id', '应聘者id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (180, 80, 1, 'bpmn_app', 'hrm_interview', 'datetime', null, 'interview_time', '面试时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (181, 80, 1, 'bpmn_app', 'hrm_interview', 'char', 100, 'interviewer_id', '面试官id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (182, 80, 1, 'bpmn_app', 'hrm_interview', 'varchar', 255, 'evaluation', '面试评价', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (183, 80, 1, 'bpmn_app', 'hrm_interview', 'tinyint', null, 'is_next', '是否有下一轮面试', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (184, 80, 1, 'bpmn_app', 'hrm_interview', 'tinyint', null, 'status', '面试结果', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (185, 80, 1, 'bpmn_app', 'hrm_interview', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (186, 80, 1, 'bpmn_app', 'hrm_interview', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (187, 80, 1, 'bpmn_app', 'hrm_interview', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (188, 80, 1, 'bpmn_app', 'hrm_interview', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (189, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'int', null, 'id', '劳动合同记录id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (190, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'char', 100, 'employee_id', '员工id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (191, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'char', 20, 'code', '合同编号', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (192, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'int', null, 'sum', '存档份数', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (193, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'datetime', null, 'sign_time', '签订时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (194, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'datetime', null, 'expired_time', '到期时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (195, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'datetime', null, 'receive_time', '领取时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (196, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'tinyint', null, 'is_received', '是否领取', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (197, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'tinyint', null, 'status', '合同状态', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (198, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (199, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (200, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (201, 80, 1, 'bpmn_app', 'hrm_labor_contract', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (202, 80, 1, 'bpmn_app', 'hrm_offer', 'int', null, 'id', 'offer记录id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (203, 80, 1, 'bpmn_app', 'hrm_offer', 'char', 100, 'applicant_id', '应聘者id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (204, 80, 1, 'bpmn_app', 'hrm_offer', 'char', 100, 'department_id', '拟录用部门id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (205, 80, 1, 'bpmn_app', 'hrm_offer', 'char', 30, 'job', '拟录用岗位', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (206, 80, 1, 'bpmn_app', 'hrm_offer', 'char', 100, 'manager_id', '汇报对象id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (207, 80, 1, 'bpmn_app', 'hrm_offer', 'char', 100, 'work_place', '工作地点', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (208, 80, 1, 'bpmn_app', 'hrm_offer', 'datetime', null, 'send_time', '发送时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (209, 80, 1, 'bpmn_app', 'hrm_offer', 'char', 100, 'sender_id', '发送人id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (210, 80, 1, 'bpmn_app', 'hrm_offer', 'tinyint', null, 'status', 'offer状态', 'radio', 'SingleChoice', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (211, 80, 1, 'bpmn_app', 'hrm_offer', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (212, 80, 1, 'bpmn_app', 'hrm_offer', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (213, 80, 1, 'bpmn_app', 'hrm_offer', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (214, 80, 1, 'bpmn_app', 'hrm_offer', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (468, 80, 1, 'bpmn_app', 'hrm_offer', 'int', null, 'before_salary', '试用期薪资总包', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (469, 80, 1, 'bpmn_app', 'hrm_offer', 'int', null, 'before_basic_salary', '试用期基本工资', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (470, 80, 1, 'bpmn_app', 'hrm_offer', 'int', null, 'before_merits_salary', '试用期绩效工资', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (471, 80, 1, 'bpmn_app', 'hrm_offer', 'char', 30, 'probation', '试用期', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (472, 80, 1, 'bpmn_app', 'hrm_offer', 'int', null, 'after_salary', '转正后薪资总包', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (473, 80, 1, 'bpmn_app', 'hrm_offer', 'int', null, 'after_basic_salary', '转正后基本工资', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (474, 80, 1, 'bpmn_app', 'hrm_offer', 'int', null, 'after_merits_salary', '转正后绩效工资', 'num', 'Number', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (475, 80, 1, 'bpmn_app', 'hrm_offer', 'datetime', null, 'retain_time', '保留日期', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (476, 80, 1, 'bpmn_app', 'hrm_offer', 'datetime', null, 'board_time', '到岗日期', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (215, 80, 1, 'bpmn_app', 'hrm_resignation', 'int', null, 'id', '离职记录id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (216, 80, 1, 'bpmn_app', 'hrm_resignation', 'char', 100, 'employee_id', '员工id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (217, 80, 1, 'bpmn_app', 'hrm_resignation', 'datetime', null, 'apply_time', '申请时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (218, 80, 1, 'bpmn_app', 'hrm_resignation', 'datetime', null, 'resign_time', '离职时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (219, 80, 1, 'bpmn_app', 'hrm_resignation', 'varchar', 255, 'reason', '离职原因', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (220, 80, 1, 'bpmn_app', 'hrm_resignation', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (221, 80, 1, 'bpmn_app', 'hrm_resignation', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (222, 80, 1, 'bpmn_app', 'hrm_resignation', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (223, 80, 1, 'bpmn_app', 'hrm_resignation', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (238, 80, 1, 'bpmn_app', 'hrm_talent_pool', 'int', null, 'id', '人才id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (239, 80, 1, 'bpmn_app', 'hrm_talent_pool', 'char', 100, 'applicant_id', '应聘者id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (240, 80, 1, 'bpmn_app', 'hrm_talent_pool', 'char', 20, 'level', '人才等级', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (241, 80, 1, 'bpmn_app', 'hrm_talent_pool', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (242, 80, 1, 'bpmn_app', 'hrm_talent_pool', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (243, 80, 1, 'bpmn_app', 'hrm_talent_pool', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (244, 80, 1, 'bpmn_app', 'hrm_talent_pool', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (224, 80, 1, 'bpmn_app', 'hrm_train', 'int', null, 'id', '培训记录id', 'num', 'Number', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (225, 80, 1, 'bpmn_app', 'hrm_train', 'char', 50, 'name', '培训名称', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (226, 80, 1, 'bpmn_app', 'hrm_train', 'char', 20, 'code', '培训编号', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (227, 80, 1, 'bpmn_app', 'hrm_train', 'datetime', null, 'start_time', '培训开始时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (228, 80, 1, 'bpmn_app', 'hrm_train', 'datetime', null, 'end_time', '培训结束时间', 'date', 'DateTime', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (229, 80, 1, 'bpmn_app', 'hrm_train', 'char', 100, 'dept_org_id', '组织培训的部门id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (230, 80, 1, 'bpmn_app', 'hrm_train', 'char', 100, 'dept_join_id', '参与培训的部门id', 'relation', 'Records', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (231, 80, 1, 'bpmn_app', 'hrm_train', 'int', null, 'budget', '培训预算(元)', 'num', 'Amount', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (232, 80, 1, 'bpmn_app', 'hrm_train', 'varchar', 255, 'goal', '培训宗旨', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (233, 80, 1, 'bpmn_app', 'hrm_train', 'varchar', 255, 'material_path', '培训资料附件地址', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (234, 80, 1, 'bpmn_app', 'hrm_train', 'int', null, 'creator_id', '记录创建人id', 'user', 'Member', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (235, 80, 1, 'bpmn_app', 'hrm_train', 'varchar', 255, 'comment', '备注', 'text', 'Input', 0);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (236, 80, 1, 'bpmn_app', 'hrm_train', 'datetime', null, 'created_at', '数据创建时间', 'date', 'DateTime', 1);
INSERT INTO `bpmn`.datasource_column (id, app_id, schemata_id, table_schema, table_name, data_type, character_maximum_length, column_name, column_comment, show_type, field_type, is_system_field) VALUES (237, 80, 1, 'bpmn_app', 'hrm_train', 'datetime', null, 'updated_at', '数据更新时间', 'date', 'DateTime', 1);

INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (22, 26, 239, 19, 123, 1);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (26, 27, 229, 20, 138, 3);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (27, 27, 230, 20, 138, 3);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (28, 21, 146, 19, 123, 1);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (29, 21, 150, 20, 138, 1);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (31, 22, 179, 19, 123, 3);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (32, 22, 181, 21, 145, 3);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (34, 24, 203, 19, 123, 1);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (35, 24, 204, 20, 138, 1);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (36, 24, 206, 21, 145, 1);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (37, 24, 209, 21, 145, 1);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (39, 23, 190, 21, 145, 1);
INSERT INTO `bpmn`.datasource_column_relation (id, source_table_id, source_column_id, target_table_id, target_column_id, type) VALUES (41, 25, 216, 21, 145, 1);