package env

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/YunzeGao/fire/framework/contract"
)

type FireEnv struct {
	folder string
	maps   map[string]string
}

func (fire FireEnv) AppEnv() string {
	return fire.Get("APP_ENV")
}

func (fire FireEnv) IsExist(key string) bool {
	_, ok := fire.maps[key]
	return ok
}

func (fire FireEnv) Get(key string) string {
	if val, ok := fire.maps[key]; ok {
		return val
	}
	return ""
}

func (fire FireEnv) All() map[string]string {
	return fire.maps
}

// NewFireEnv 需要一个参数: .env文件所在的目录
// example: NewFireEnv("/env/folder/") 会读取文件: /env/folder/.env
// .env的文件格式 FOO_ENV=BAR
func NewFireEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewFireEnv param error")
	}
	// 读取folder文件
	folder := params[0].(string)
	filePath := path.Join(folder, ".env")

	// 实例化
	env := &FireEnv{
		folder: folder,
		// 实例化环境变量，APP_ENV默认设置为开发环境
		maps: map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	// 读取文件配置
	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		rd := bufio.NewReader(file)
		for {
			line, _, readErr := rd.ReadLine()
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				err = readErr
				break
			}
			s := bytes.SplitN(line, []byte{'='}, 2)
			if len(s) != 2 {
				continue
			}
			env.maps[string(s[0])] = string(s[1])
		}
	}
	// 运行是覆盖读取配置
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			env.maps[pair[0]] = pair[1]
		}
	}
	if err != nil {
		log.Println("[NewFireEnv Error]\t", err)
	}
	return env, nil
}
