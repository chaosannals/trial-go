package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"strings"
	"github.com/cihub/seelog"
	"github.com/chaosannals/fclear/core"
)

var CFG core.FClearConfig

func init() {
	root, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}

	// 日志初始化
	logCfgPath := filepath.Join(root, "seelog.xml")
	log, err := seelog.LoggerFromConfigAsFile(logCfgPath)
	if err != nil {
		fmt.Printf("seelog.xml 失败: %v\n", err)
		return
	}
	seelog.ReplaceLogger(log)

	// 配置初始化
	cfgPath := filepath.Join(root, "config.json")
	f,err := os.Open(cfgPath)
	if err != nil {
		seelog.Error(err)
		return
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	err = json.Unmarshal(content, &CFG)
	if err != nil {
		seelog.Error(err)
		return
	}
	seelog.Info("init final.")
}

func main() {
	seelog.Info("start.")
	wkdir, err := os.Getwd()
	if err != nil {
		seelog.Error(err)
	}
	tc := len(CFG.Tasks)
	seelog.Infof("tasks count: %d", tc)
	for i := 0; i < tc; i+=1 {
		task := &CFG.Tasks[i]
		dir := strings.Replace(task.Dir, "${wkdir}", wkdir, -1)
		dir, err = filepath.Abs(filepath.Clean(dir))
		if err != nil {
			seelog.Error(err)
		}
		task.Dir = dir
		seelog.Infof("clear: %s", dir)
		ali, err := core.MatchAliOssClient(CFG.Alioss, task.Alioss)
		if err != nil {
			seelog.Error(err)
		}
		core.RunTask(task, ali)
	}
}