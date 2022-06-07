package core

import (
	"os"
	"time"
	"io/ioutil"
	"strings"
	"path/filepath"
	"github.com/cihub/seelog"
)

type FClearTask struct {
	config *FClearConfigTask
	alioss *AliOssClient
	timeLimit time.Time
}

func RunTask(cfg *FClearConfigTask, ali *AliOssClient) error {
	now := time.Now()
	m, err := time.ParseDuration(cfg.ValidTime)
	if err != nil {
		return err
	}
	task := &FClearTask {
		config: cfg,
		alioss: ali,
		timeLimit: now.Truncate(m),
	}
	seelog.Infof("");
	return task.RunTaskOnDir(cfg.Dir)
}

func ValidExt(ext string, exts []string) bool {
	c := len(exts)
	for i := 0; i < c; i+=1 {
		if strings.EqualFold(exts[i], ext) {
			return true
		}
	}
	return false
}

func (i *FClearTask)Restore(p string) {
	ossp := strings.Replace(strings.Trim(p[len(i.config.Dir):], "/\\"), "\\", "/", -1)
	seelog.Infof("%s => %s", p, ossp)
	err := i.alioss.PutFromFile(ossp, p)
	if err != nil {
		seelog.Error(err)
	} else {
		err = os.Remove(p)
		if err != nil {
			seelog.Error(err)
		}
	}
}

func (i *FClearTask)RunTaskOnDir(dir string) error {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		p := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			i.RunTaskOnDir(p)
		} else {
			ext := filepath.Ext(p)
			if ValidExt(ext, i.config.Exts) {
				fi, err := os.Stat(p)
				if err != nil {
					seelog.Error(err)
				}else {
					mt := fi.ModTime()
					d := i.timeLimit.Sub(mt).Seconds()
					if d > 0 {
						//seelog.Infof("%s size: %d", mt.Format("2006-01-02 15:04:05"), fi.Size())
						i.Restore(p)
					}
				}
			}
		}
	}

	return nil
}