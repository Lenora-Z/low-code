// Package conf Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 10:56 上午
package conf

type ServerConfig struct {
	Debug            bool         `yaml:"debug"`
	JWTSecret        string       `required:"true" yaml:"jwt_secret"`
	BaseUrl          string       `yaml:"base_url"`
	DbConfig         *DBConfig    `required:"true" yaml:"mysql"`
	BusinessDbConfig *DBConfig    `required:"true" yaml:"mysql_business"`
	MongoConfig      *MongoConfig `required:"true" yaml:"mongo"`
	Engine           *HostConfig  `yaml:"engine"`
	MinioConf        *MinioConf   `required:"true" yaml:"minio"`
	TritiumConfig    *TrimConfig  `yaml:"tritium"`
	SmtpGroupConf    []SmtpConfig `required:"true" yaml:"email_sender"`
	//Consul         *ConsulConfig     `yaml:"consul"`
	//SkyWalking     *SkyWalkingConfig `yaml:"skyWalking"`
}

type DBConfig struct {
	User     string `default:"root" yaml:"user"`
	Password string `default:"" yaml:"password"`
	Name     string `yaml:"ip"`
	Port     uint   `default:"3306" yaml:"port"`
	DbName   string `required:"true" yaml:"db_name"`
	Charset  string `default:"utf8" yaml:"charset"`
	MaxIdle  int    `default:"10" yaml:"max_idle"`
	MaxOpen  int    `default:"50" yaml:"max_open"`
	LogMode  bool   `yaml:"log_mode"`
	Loc      string `required:"true" yaml:"loc"`
}

type MongoConfig struct {
	IP       string `yaml:"ip"`
	User     string `default:"root" yaml:"user"`
	Password string `default:"" yaml:"password"`
	Port     string `default:"27017" yaml:"port"`
	DbName   string `required:"true" yaml:"db_name"`
	MaxIdle  string `default:"1" yaml:"max_idle"`
	MaxOpen  string `default:"10" yaml:"max_open"`
}

type HostConfig struct {
	Api string `yaml:"api"`
}

type TrimConfig struct {
	Api    string `yaml:"api"`
	Switch bool   `yaml:"status"`
}

type MinioConf struct {
	EndpointIP      string `yaml:"endpointIP"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
}

type ConsulConfig struct {
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type SkyWalkingConfig struct {
	OapServer string `yaml:"oap_server"`
}

type SmtpConfig struct {
	Sender   string `yaml:"sender"`
	Password string `yaml:"password"`
	SmtpAddr string `yaml:"smtp_addr"`
	SmtpPort uint64 `yaml:"smtp_port"`
}
