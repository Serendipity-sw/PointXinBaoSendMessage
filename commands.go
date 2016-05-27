package main

import (
	//"errors"
	"fmt"
	//"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	//"sync/atomic"
	//"time"

	"github.com/smtc/glog"
)

const (
	_version = "v1.4 rev-20151203"
)
var uniqStatics bool = false // 是否去重统计

func trimHost(ohost string) string {
	host := strings.ToLower(strings.TrimSpace(ohost))
	if strings.HasPrefix(host, "http://") {
		host = host[7:]
	}
	return host
}

// 转换
// 2015-09-07 10.10.133.99转换为8610010133099
func inet_addr(s string) (msisdn uint64, err error) {
	s = strings.TrimSpace(s)
	items := strings.Split(s, ".")
	if len(items) != 4 {
		return 0, fmt.Errorf("ip address should be 4 items after split with '.'.")
	}
	a, err := strconv.Atoi(items[0])
	if err != nil {
		return 0, err
	}
	b, _ := strconv.Atoi(items[1])
	if err != nil {
		return 0, err
	}
	c, _ := strconv.Atoi(items[2])
	if err != nil {
		return 0, err
	}
	d, _ := strconv.Atoi(items[3])
	if err != nil {
		return 0, err
	}
	//msisdn = uint64(a + (b << 8) + (c << 16) + (d << 24))
	msisdn = uint64(8600000000000 + a*1000000000 + b*1000000 + c*1000 + d)

	return
}

// 输入k1=v1 k2=v2 k3=v3...
func getCmdArgs(args []string, key string) string {
	m := map[string]string{}
	for _, a := range args {
		entries := strings.Split(a, "=")
		if len(entries) != 2 {
			glog.Error("argument %s format error, should be a=b\n", a)
			continue
		}
		m[entries[0]] = entries[1]
	}
	return m[key]
}

// 注册命令
func registerCommands() {
	RegisterTermCmd("help", 2, 1, false, func(argv []string) (string, error) {
		return consoleUsage(), nil
	})

	RegisterTermCmd("version", 2, 1, false, func(argv []string) (string, error) {
		return _version, nil
	})

	RegisterTermCmd("uniqStatics", 2, 1, false, func(argv []string) (string, error) {
		argc := len(argv)
		if argc == 1 {
			return fmt.Sprint(uniqStatics), nil
		}

		res := ""
		switch strings.ToLower(argv[1]) {
		case "true":
			fallthrough
		case "1":
			uniqStatics = true
			res = "set uniqStatics to true."

		case "false":
			fallthrough
		case "0":
			uniqStatics = false
			res = "set uniqStatics to false."

		default:
			res = "invalid parms"
		}

		return res, nil
	})

	RegisterTermCmd("ps", 1, 1, false, func(argv []string) (string, error) {
		var stackBuf []byte = make([]byte, 500000)
		runtime.Stack(stackBuf, true)
		return strings.Replace(string(stackBuf), "\n", "\r\n", -1), nil
	})

}

// 打印命令帮助
func consoleUsage() string {
	return strings.Replace(`
ps:                      -- show stack
help:                    -- show this usage
version:                 -- show version
exit:                    -- exit console
uniqStatics:             -- show/set uniqStatics param

	`, "\n", "\r\n", -1)
}
