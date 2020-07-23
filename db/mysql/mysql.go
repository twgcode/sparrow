/**
@Author: wei-g
@Date:   2020/7/23 2:30 下午
@Description: mysql 工具
*/

package mysql

import (
	"fmt"
	"net/url"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB        *sqlx.DB
	once      sync.Once
	onceClose sync.Once
)

// ConfigMysql mysql 有关配置
type ConfigMysql struct {
	Ip             string `mapstructure:"ip" json:"ip"`
	Port           int    `mapstructure:"port" json:"port"`
	Database       string `mapstructure:"database" json:"database"`
	Charset        string `mapstructure:"charset" json:"charset"`
	Username       string `mapstructure:"username" json:"username"`
	Password       string `mapstructure:"password" json:"password"`
	ParseTime      bool   `mapstructure:"parse_time" json:"parse_time"`
	Loc            string `mapstructure:"loc" json:"loc"`
	SetMaxOpenConn int    `mapstructure:"set_max_open_conn" json:"set_max_open_conn"`
	SetMaxIdleConn int    `mapstructure:"set_max_idle_conn" json:"set_max_idle_conn"`
}

// 初始化默认的 mysql 连接池
func InitDefaultDB(username, password, ip string, port int, database, charset string, parseTime bool, loc string, setMaxOpenConn, setMaxIdleConn int) (*sqlx.DB, error) {
	var err error
	once.Do(func() {
		DB, err = NewDB(username, password, ip, port, database, charset, parseTime, loc, setMaxOpenConn, setMaxIdleConn)
	})
	return DB, err
}

// 初始化默认的 mysql 连接池 从 ConfigMysql 获取有关配置
func InitDefaultDBFromCfg(cfg *ConfigMysql) (*sqlx.DB, error) {
	var err error
	once.Do(func() {
		DB, err = NewDBFromCfg(cfg)
	})
	return DB, err
}

// Close 关闭数据库
func Close() (err error) {
	onceClose.Do(func() {
		if DB != nil {
			err = DB.Close()
		}
	})
	return
}

// NewDBFromCfg 构建一个 mysql 连接池 从 ConfigMysql 获取有关配置
func NewDBFromCfg(cfg *ConfigMysql) (*sqlx.DB, error) {
	return newDB(cfg.Username, cfg.Password, cfg.Ip, cfg.Port, cfg.Database, cfg.Charset, cfg.ParseTime, cfg.Loc, cfg.SetMaxOpenConn, cfg.SetMaxIdleConn)
}

// NewDB  构建一个 mysql 连接池
func NewDB(username, password, ip string, port int, database, charset string, parseTime bool, loc string, setMaxOpenConn, setMaxIdleConn int) (*sqlx.DB, error) {
	return newDB(username, password, ip, port, database, charset, parseTime, loc, setMaxOpenConn, setMaxIdleConn)

}

// newDB 构建一个 mysql 连接池
func newDB(username, password, ip string, port int, database, charset string, parseTime bool, loc string, setMaxOpenConn, setMaxIdleConn int) (db *sqlx.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s", username, password, ip, port, database, charset, parseTime, url.QueryEscape(loc))
	// 连接数据库
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(setMaxOpenConn)
	db.SetMaxIdleConns(setMaxIdleConn)
	return
}
