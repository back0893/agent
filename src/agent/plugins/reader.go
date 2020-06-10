package plugins

import (
	"github.com/back0893/goTcp/utils"
	"github.com/toolkits/file"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// 读取目录下 name_%d+ 类似的文件
func ListPlugins(relativePath string) map[string]*Plugin {
	ret := make(map[string]*Plugin)
	dir := filepath.Join(utils.GlobalConfig.GetString("plugin.dir"), relativePath)
	log.Println(dir)
	if !file.IsExist(dir) || file.IsFile(dir) {
		return ret
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("can not list files under", dir)
		return ret
	}

	for _, f := range fs {
		//如果是.开始,说明是隐藏或者特殊直接丢弃
		if strings.Index(f.Name(), ".") == 0 {
			continue
		}
		//继续扫描
		if f.IsDir() {
			tmpRet := ListPlugins(filepath.Join(relativePath, f.Name()))
			for key := range tmpRet {
				ret[key] = tmpRet[key]
			}
			continue
		}

		filename := f.Name()
		arr := strings.Split(filename, "_")
		if len(arr) < 2 {
			continue
		}

		// filename should be: $xx_$interval
		var interval int
		interval, err = strconv.Atoi(arr[0])
		if err != nil {
			continue
		}

		fpath := filepath.Join(relativePath, filename)
		//如果间隔时间为0,意味只执行一次的插件
		plugin := &Plugin{FilePath: fpath, Interval: interval, IsRepeat: true}
		if interval == 0 {
			plugin.IsRepeat = false
		}
		ret[fpath] = plugin
	}

	return ret
}
