/**
@Author: wei-g
@Date:   2020/7/23 3:06 下午
@Description:
*/

package mysql

import "testing"

func TestInitDefaultDBFromCfg(t *testing.T) {
	var err error
	cfg := &Config{
		Ip:             "127.0.0.1",
		Port:           3306,
		Database:       "test",
		Charset:        "utf8mb4",
		Username:       "user",
		Password:       "password",
		ParseTime:      true,
		Loc:            "Local",
		SetMaxOpenConn: 100,
		SetMaxIdleConn: 10,
	}
	if _, err = InitDefaultDBFromCfg(cfg); err != nil {
		t.Fatal(err)
	}
	if err = Close(); err != nil {
		t.Fatal(err)
	}
}
