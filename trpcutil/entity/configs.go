// Package entity 公共的结构体
package entity

// DBInfo 配置
type DBInfo struct {
	DSN            string `yaml:"dsn"`
	MaxIdleConnNum int    `yaml:"maxIdleConnNum"`
	MaxConnNum     int    `yaml:"maxConnNum"`
}

// TrpcInfo 配置
type TrpcInfo struct {
	Namespace string `yaml:"namespace"`
	Env       string `yaml:"env"`
	Target    string `yaml:"target"`
	Timeout   int32  `yaml:"timeout"`
}

// TafInfo 配置
type TafInfo struct {
	Namespace string `yaml:"namespace"`
	Env       string `yaml:"env"`
	Target    string `yaml:"target"`
	Timeout   int32  `yaml:"timeout"`
	TafObj    string `yaml:"tafObj"`
}

// RedisInfo redis配置
type RedisInfo struct {
	Address string `yaml:"address"`
	Passwd  string `yaml:"passwd"`
	DB      int    `yaml:"db"`
}

// MongodbInfo mongo配置
type MongodbInfo struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}
