package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

type FireConfig struct {
	container framework.IContainer

	env      string
	folder   string
	keyDelim string
	lock     sync.RWMutex // 配置文件读写锁
	envMaps  map[string]string
	confMaps map[string]interface{}
	confRaws map[string][]byte
}

func (conf *FireConfig) loadConfigFile(configFolder string, fileName string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	s := strings.Split(fileName, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		filePath := filepath.Join(configFolder, fileName)
		// read file bytes
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}
		// do replace
		content = replace(content, conf.envMaps)
		// parse yaml
		values := map[string]interface{}{}
		if err = yaml.Unmarshal(content, &values); err != nil {
			return err
		}
		conf.confMaps[name] = values
		conf.confRaws[name] = content

		if name == "app" && conf.container.IsBind(contract.AppKey) {
			if path, ok := values["path"]; ok {
				appService := conf.container.MustMake(contract.AppKey).(contract.App)
				appService.LoadAppConfig(cast.ToStringMapString(path))
			}
		}
	}
	return nil
}

// 删除文件的操作
func (conf *FireConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	s := strings.Split(file, ".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(conf.confRaws, name)
		delete(conf.confMaps, name)
	}
	return nil
}

func NewFireConfig(params ...interface{}) (interface{}, error) {
	if len(params) != 4 {
		return nil, errors.New("NewFireConfigService params error")
	}
	folder := params[0].(string)
	envMaps := params[1].(map[string]string)
	env := params[2].(string)
	container := params[3].(framework.IContainer)
	configFolder := filepath.Join(folder, env)
	// check folder exist
	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		return nil, errors.New("folder " + configFolder + " not exist: " + err.Error())
	}
	fireConf := &FireConfig{
		container: container,
		folder:    folder,
		env:       env,
		keyDelim:  ".",
		lock:      sync.RWMutex{},
		envMaps:   envMaps,
		confMaps:  map[string]interface{}{},
		confRaws:  map[string][]byte{},
	}
	// read all yml/yaml files in folder
	files, err := os.ReadDir(configFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, file := range files {
		err = fireConf.loadConfigFile(configFolder, file.Name())
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// 监控文件夹文件
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(configFolder)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(fmt.Println(err))
			}
		}()
		for {
			select {
			case ev := <-watcher.Events:
				{
					fmt.Println(ev)
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]
					//判断事件发生的类型
					// Create 创建
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
						_ = fireConf.loadConfigFile(folder, fileName)
					}
					// Write 更新
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
						_ = fireConf.loadConfigFile(folder, fileName)
					}
					// Remove 删除
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						_ = fireConf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watcher.Errors:
				{
					log.Println("watcher error : ", err)
					return
				}
			}
		}
	}()
	return fireConf, nil
}

func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}

	return content
}

func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}
		switch next.(type) {
		case map[interface{}]interface{}:
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			return nil
		}
	}
	return nil
}

func (conf *FireConfig) find(key string) interface{} {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

// IsExist check setting is existed
func (conf *FireConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

// Get a new interface
func (conf *FireConfig) Get(key string) interface{} {
	return conf.find(key)
}

// GetBool get bool type
func (conf *FireConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

// GetInt get Int type
func (conf *FireConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// GetFloat64 get float64
func (conf *FireConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *FireConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

// GetString get string typen
func (conf *FireConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *FireConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *FireConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *FireConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *FireConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *FireConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

// Load a config to a struct, val should be a pointer
func (conf *FireConfig) Load(key string, val interface{}) error {
	decode, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}
	return decode.Decode(conf.find(key))
}
