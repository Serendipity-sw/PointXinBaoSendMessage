package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
)

var (
	rootPrefix                           string // 场景导航的route prefix
	configFn                                    = flag.String("config", "./config.json", "config file path")
	debugFlag                                   = flag.Bool("d", false, "debug mode")
	rt                                   *gin.Engine
addflow string //点信网获取手机号码对应连接接口地址
	smspassword string //短信发送密码
	sendMessageUrl string //短信发送接口地址
	doublekill string //追销点信网获取手机号码对应连接接口地址
	notepadProcessDir string //监控目标文件夹地址
)

func serverRun(cfn string, debug bool) {
	config.ReadCfg(cfn)
	logInit(debug)

	// 读取配置文件参数
	// js文件夹名称
	port := config.GetIntDefault("port", 8000)
	rootPrefix = strings.TrimSpace(config.GetStringDefault("rootprefix", ""))
	addflow = strings.TrimSpace(config.GetString("addflow"))
	smspassword = strings.TrimSpace(config.GetString("smspassword"))
	sendMessageUrl = strings.TrimSpace(config.GetString("sendMessageUrl"))
	doublekill = strings.TrimSpace(config.GetString("doublekill"))
	notepadProcessDir = strings.TrimSpace(config.GetString("notepadProcessDir"))

	if len(rootPrefix) != 0 {
		if !strings.HasPrefix(rootPrefix, "/") {
			rootPrefix = "/" + rootPrefix
		}
		if strings.HasSuffix(rootPrefix, "/") {
			rootPrefix = rootPrefix[0 : len(rootPrefix)-1]
		}
	}

	// 初始化
	deferinit.InitAll()
	glog.Info("init all module successfully.\n")

	// 设置多cpu运行
	runtime.GOMAXPROCS(runtime.NumCPU())

	deferinit.RunRoutines()
	glog.Info("run routines successfully.\n")

	// 注册term console的命令
	registerCommands()
	// gin的工作模式
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	rt = gin.Default()
	//loadTemplates(rt, false)
	router(rt)
	go rt.Run(fmt.Sprintf(":%d", port))
}

// 结束进程
func serverExit() {
	// 结束所有go routine
	deferinit.StopRoutines()
	glog.Info("stop routine successfully.\n")

	deferinit.FiniAll()
	glog.Info("fini all modules successfully.\n")
}

func main() {
	//判断进程是否存在
	if checkPid() {
		return
	}

	flag.Parse()

	serverRun(*configFn, *debugFlag)

	c := make(chan os.Signal, 1)
	writePid()
	// 信号处理
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	// 等待信号
	<-c

	serverExit()
	rmPidFile()
	glog.Close()
	os.Exit(0)
}

func router(r *gin.Engine) {

	g := &r.RouterGroup
	if rootPrefix != "" {
		g = r.Group(rootPrefix)
	}

	{
		g.GET("/", func(c *gin.Context) { c.String(200, "ok") })
	}
}

// 报告用户请求相关信息
func userReqInfo(req *http.Request) (info string) {
	info += fmt.Sprintf("ipaddr: %s user-agent: %s referer: %s",
		req.RemoteAddr, req.UserAgent(), req.Referer())
	return info
}
