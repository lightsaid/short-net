package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// Setingenv 从 env map 对象，提取值到 target, target 必须是一个带有 mapstruct tag 结构体指针；
//
// 不支持结构体嵌套
//
// 支持解析类型有：
//
// 支持时间解析：天、时、分、秒(d、h、m、s) 如：3d、12h、7m、10s，但不支持组合 12h7m
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
				// fmt.Println("kind: >>> ", f.Type().Kind())
				// fmt.Println("type: >>> ", f.Type().String())

				// Int, 这里先不设置 int64，int64 最面处理
				if f.Type().Kind() == reflect.Int ||
					f.Type().Kind() == reflect.Int8 ||
					f.Type().Kind() == reflect.Int16 ||
					f.Type().Kind() == reflect.Int32 {
					real, _ := strconv.Atoi(v)
					f.SetInt(int64(real))
				}

				// Uint
				if f.Type().Kind() == reflect.Uint ||
					f.Type().Kind() == reflect.Uint8 ||
					f.Type().Kind() == reflect.Uint16 ||
					f.Type().Kind() == reflect.Uint32 ||
					f.Type().Kind() == reflect.Uint64 {
					real, _ := strconv.Atoi(v)
					f.SetUint(uint64(real))
				}

				// Float
				if f.Type().Kind() == reflect.Float32 || f.Type().Kind() == reflect.Float64 {
					real, _ := strconv.ParseFloat(v, 32)
					f.SetFloat(real)
				}

				// Complex
				if f.Type().Kind() == reflect.Complex64 || f.Type().Kind() == reflect.Complex128 {
					real, _ := strconv.ParseComplex(v, 64)
					f.SetComplex(real)
				}

				// String
				if f.Type().Kind() == reflect.String {
					f.SetString(v)
				}

				// Boolean
				if f.Type().Kind() == reflect.Bool {
					real, _ := strconv.ParseBool(v)
					f.SetBool(real)
				}

				if f.Type().Kind() == reflect.Int64 {
					// 设置 time.Duration
					if f.Type().String() == "time.Duration" {
						// 天、时、分、秒 解析 d,h,m,s
						real, _ := parseTime(v)
						f.SetInt(real)
					} else {
						real, _ := strconv.Atoi(v)
						f.SetInt(int64(real))
					}
				}
			}
		}
	}

	return nil
}

func parseTime(v string) (int64, error) {
	start := time.Now()
	defer slog.Info("loading config parseTime():", time.Since(start))
	// 解析： 3d、12h、7m、10s 时间
	var handler = func(v string, kind string) (int64, error) {
		fmt.Println("parseTime: ", v, kind)
		s := strings.ReplaceAll(v, kind, "")
		t, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return int64(t), nil
	}

	if strings.HasSuffix(v, "d") {
		t, err := handler(v, "d")
		if err != nil {
			return 0, err
		}
		tt := t * 24 * int64(time.Hour)
		return tt, nil
	}

	if strings.HasSuffix(v, "h") {
		t, err := handler(v, "h")
		if err != nil {
			return 0, err
		}
		tt := t * int64(time.Hour)
		return tt, nil
	}

	if strings.HasSuffix(v, "m") {
		t, err := handler(v, "m")
		if err != nil {
			return 0, err
		}
		tt := t * int64(time.Minute)
		return tt, nil
	}

	if strings.HasSuffix(v, "s") {
		t, err := handler(v, "s")
		if err != nil {
			return 0, err
		}
		tt := t * int64(time.Second)
		return tt, nil
	}

	return 0, nil
}
