/*
该工具主要是为了方便将数据库结构快速导出到 db/growerlab.sql
避免，修改一个字段或者相关数据库结构操作频繁的手动导出到 db/growerlab.sql，同时避免手动操作可能的失误

使用方法：

$ cd backend
$ go run script/tool/exportddl/main.go

*/

package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bitfield/script"

	"github.com/go-sql-driver/mysql"
	"github.com/growerlab/growerlab/src/common/configurator"
)

var projectPath = filepath.Join(os.Getenv("GOPATH"), "src", "github.com/growerlab/growerlab/src/backend")
var ddlPath = "db/growerlab.sql"

type DBInfo struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

func main() {
	os.Args[0] = projectPath

	info, err := Prepare()
	if err != nil {
		panic(err)
	}

	err = Export(info)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done.")
}

func Export(info *DBInfo) error {
	defaultPath := filepath.Join(projectPath, ddlPath)
	sqlFile, err := os.OpenFile(defaultPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer sqlFile.Close()

	cmd := fmt.Sprintf("docker exec mysql8 mysqldump -u%s -p%s -d %s", info.Username, info.Password, info.DBName)
	script.Exec(cmd).
		ReplaceRegexp(regexp.MustCompile(" AUTO_INCREMENT=[[:digit:]]+"), "").
		ReplaceRegexp(regexp.MustCompile("mysqldump:+.+"), "").
		WithStdout(sqlFile).
		Stdout()
	// Wait()
	return nil
}

func Prepare() (*DBInfo, error) {
	err := configurator.InitConfig()
	if err != nil {
		return nil, err
	}

	dbUrl := configurator.GetConf().DBUrl

	cfg, err := mysql.ParseDSN(dbUrl)
	if err != nil {
		return nil, err
	}

	info := &DBInfo{}
	info.Username = cfg.User
	info.Password = cfg.Passwd
	info.Host, info.Port, _ = net.SplitHostPort(cfg.Addr)
	info.DBName = cfg.DBName
	return info, nil
}
