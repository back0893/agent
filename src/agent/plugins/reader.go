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
	if relativePath == "" {
		return ret
	}

	dir := filepath.Join(utils.GlobalConfig.GetString("plugin.dir"), relativePath)

	if !file.IsExist(dir) || file.IsFile(dir) {
		return ret
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("can not list files under", dir)
		return ret
	}

	for _, f := range fs {
		//继续扫描
		if f.IsDir() {
			tmpRet := ListPlugins(filepath.Join(relativePath, f.Name()))
			for key := range tmpRet {
				ret[key] = tmpRet[key]
			}
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
		plugin := &Plugin{FilePath: fpath, Interval: interval}
		ret[fpath] = plugin
	}

	return ret
}
