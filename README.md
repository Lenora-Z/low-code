### 拉取mod地址
```bash
go mod init github.com/Lenora-Z/low-code
```
### 生成文档
```bash
#swag命令获取
go get -u github.com/swaggo/swag/cmd/swag
#swagger文档生成
swag init -g ./cmd/main.go
```

### 数据库表结构导出
```shell
mkdir data/sql
cd data/sql
#仅导出表结构
mysqldump -h 192.168.4.108 -P3306 -uroot -p123456 -d bpmn > bpmn.sql
#导出结构&数据
mysqldump -h 192.168.4.108 -P3306 -uroot -p123456 bpmn > bpmn.sql
cat bpmn.sql | mysql -h [hostaddress] -P [port] -u[username] -p[pwd] [database] 
```

### 同步数据库至代码仓库
```bash
#下载依赖包
go get -u -v github.com/xxjwxc/gormt@master
#创建软连接至代码仓库
mkdir gormstruct
ln -s `which gormt` gormstruct
#执行同步
./gormstruct/gormt
```


### 项目配置
字段 | 配置项
--- | ---
base_url | api查看文件地址url
mysql | mysql数据库(项目数据库)
data_db | mysql数据库(业务数据库)
mongo | mongodb数据库
tritium | 氚平台
engine | bpmn引擎
minio | minio客户端
email_sender | 发送邮件邮箱配置

### 技术文档地址
[点击查看技术文档](https://imbni9806o.feishu.cn/docs/doccnVCzviFZC42Q2WAKsxCY6fe)