package main

import (
	"sync"
	"os"
	"github.com/howeyc/fsnotify"
	"github.com/guotie/deferinit"
	"bufio"
	"strings"
	"io"
	"github.com/smtc/glog"
	"sync/atomic"
)
type sendMessageMap struct  {
	mob string//手机号码
	content string //短息内容
	price int32 //价格
	state int32 // 1成功 2url获取失败 3短信发送失败
	zx int32//是否为追销 1追销 0正常
	province string//省份
}

type sendMessageList struct  {
	sync.RWMutex
	sendMessageTable map[string][]sendMessageMap
}
var (
	//sendMessageList=sendMessageList{
	//	sendMessageTable:map[string][]sendMessageMap{},
	//}
)
func init() {
	deferinit.AddRoutine(notepadProcess)
}
type counter struct {
	val int32
}
func (c *counter) increment() {
	atomic.AddInt32(&c.val, 1)
}

/**
记事本处理文件
如当前目录不存在则不做任何处理,该方法直接不做任何处理
创建人:邵炜
创建时间:2016年4月12日11:37:41
*/
func notepadProcess(ch chan struct{}, wg *sync.WaitGroup) {
	fi, err := os.Stat(notepadProcessDir)
	if err != nil {
		glog.Error("notepadProcess: file data is error! err: %s \n", err.Error())
		return
	}
	if !fi.IsDir() {
		glog.Error("notepadProcess: message file name :%s is not defind! \n", notepadProcessDir)
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Error("notepadProcess: fsnotify newWatcher is error! err: %s \n", err.Error())
		return
	}
	var  modifyReceived counter
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				glog.Info("messageCenterFtp: fsnotify watcher fileName: %s is change!  ev: %v \n", ev.Name, ev)
					if ev.IsModify() {
						modifyReceived.increment()
						if modifyReceived.val % 2 == 0 {
							go func(filePath string) {
								readFileMobs(ev.Name)
							}(ev.Name)
						}
					}
			case err := <-watcher.Error:
				glog.Error("messageCenterFtp: fsnotify watcher is error! err: %s \n", err.Error())
			}
		}
		done <- true
	}()
	err = watcher.WatchFlags(notepadProcessDir,fsnotify.FSN_MODIFY)
	if err != nil {
		glog.Error("messageCenterFtp watch error. messageCenterDir: %s  err: %s \n", notepadProcessDir, err.Error())
	}

	// Hang so program doesn't exit
	<-ch

	/* ... do stuff ... */
	watcher.Close()
	wg.Done()
}

func readFileMobs(filePath string) error{
	fs,err:=os.Open(filePath)
	defer func() {
		err=fs.Close()
		if err != nil {
			glog.Error("readFileMobs close error! filePath: %s error: %s \n",filePath,err.Error())
		}
	}()
	if err != nil {
		glog.Error("readFileMobs open file error ! file: %s err: %s \n",filePath,err.Error())
		return err
	}
	buf := bufio.NewReader(fs)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if len(line) != 0 {
			go func(fileContent string) {
				var contents []string
				contents=strings.Split(line,",")
				if len(contents)>=5 {
					messageProcess(&contents,filePath)
				}
			}(line)
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			glog.Error("readFileMobs read error! file: %s err: %s \n",filePath,err.Error())
			return err
		}
	}
}


/**
用户短信处理
创建人:邵炜
创建时间:2016年5月27日15:03:46
输入参数:用户短信列表数组
 */
func messageProcess(contents *[]string,filePath string) {
	content:=*contents
model,err:=getMyUserJumpUrl(content[0],content[1],content[2],content[3])
	if err != nil {
		return
	}
	if model.State == "success" {
		sendMessage(content[0],strings.Replace(content[4],"+链接",model.Msg,-1),content[3],filePath)
	}
}