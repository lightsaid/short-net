package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/exp/slog"
)

type Config struct {
	DBPort int `mapstruct:"DB_PORT"`
}

// LoadingEnv 从配置文件加载配置，默认 .env 文件
func Loadingenv(paths ...string) (map[string]string, error) {
	var filename = ".env"
	var configs = map[string]string{}

	if len(paths) > 0 {
		filename = paths[0]
	}

	file, err := os.Open(filename)
	if err != nil {
		return configs, err
	}

	// 将字符串里的空格、'、" 替换成空字符正则
	var regx = regexp.MustCompile(`[\s*|\'*|\"*]`)

	// 创建一个扫描器
	scan := bufio.NewScanner(file)

	// 扫描文档
	for scan.Scan() {
		lineText := string(regx.ReplaceAll([]byte(scan.Text()), []byte("")))

		// 检查是否是注释
		if strings.HasPrefix(lineText, "#") {
			continue
		}

		// 从 # 分割, 防止添加了注释进来，分割取第一个
		part := strings.Split(lineText, "#")[0]
		key, value, found := strings.Cut(part, "=")
		if found {
			slog.Info("read config: ", key, value)
			configs[key] = value
			err := os.Setenv(key, value)
			if err != nil {
				slog.Error("loadenv.go os.Setenv: "+err.Error(), key, value)
			}
		}
	}

	return configs, nil
}

// Setingenv 从 env map 对象，提取值到 target, target 必须是一个带有 mapstruct tag 结构体指针
func Setingenv(target interface{}, env map[string]string) error {
	if reflect.TypeOf(target).Elem().Kind() != reflect.Struct {
		return errors.New("targe must struct pointer")
	}

	value := reflect.ValueOf(target).Elem()

	for i := 0; i < value.NumField(); i++ {
		// 获取 target 结构体 field reflect.StructField
		field := value.Type().Field(i)

		// 获取 mapstruct tag 值
		tag := field.Tag.Get("mapstruct")

		v, exists := env[tag]
		if exists {
			// 根据字段名获取 reflect.Value
			f := value.FieldByName(field.Name)
			if f.CanSet() {
				fmt.Println("v: ", v)
				f.SetString(v)
			}
		}
	}

	return nil
}
