package util

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	var regx = regexp.MustCompile(`[\s*|\'*|\"*]`)
	// newStr := string(regx.ReplaceAll([]byte(`   # #  serverport  ''= '20000''' # ppp= "ddd" # ddddd"`), []byte("")))
	newStr := string(regx.ReplaceAll([]byte(` "  port = ' '"`), []byte("")))
	log.Printf("%q\n", newStr)
	if strings.HasPrefix(newStr, "#") {
		log.Println("Yes")
	}
	// 从 # 分割
	parts := strings.Split(newStr, "#")
	fmt.Println(parts)
	fmt.Printf("%q\n", parts[0])

	key, value, found := strings.Cut(parts[0], "=")
	fmt.Printf("%q - %q - %t\n", key, value, found)

}

func readEnv(t *testing.T) map[string]string {
	var path = "./test.env"
	env, err := Loadingenv(path)
	require.NoError(t, err)
	for k, v := range env {
		require.NotEmpty(t, k, v)
		// fmt.Printf("key=%q value=%q\n", k, v)
	}

	return env
}

func TestLoadingenv(t *testing.T) {
	readEnv(t)
}

func TestSetingenv(t *testing.T) {
	type config struct {
		DBPort int           `mapstruct:"DB_PORT"`
		DBName string        `mapstruct:"DB_NAME"`
		IsDev  bool          `mapstruct:"IS_DEV"`
		Pi     float32       `mapstruct:"PI"`
		Day    time.Duration `mapstruct:"Day"`
		Hour   time.Duration `mapstruct:"Hour"`
		Minute time.Duration `mapstruct:"Minute"`
		Second time.Duration `mapstruct:"Second"`
		/*
			Day=3d
			Hour=10h
			Minute=15m
			Second=30s
		*/
	}
	env := readEnv(t)

	var target config
	err := Setingenv(&target, env)
	require.NoError(t, err)

	require.Equal(t, target.DBPort, 3307)
	require.Equal(t, target.DBName, "shortnet")
	require.Equal(t, target.IsDev, true)
	require.Equal(t, target.Pi, float32(3.14))
	require.Equal(t, target.Day, 3*24*time.Hour)
	require.Equal(t, target.Hour, 10*time.Hour)
	require.Equal(t, target.Minute, 15*time.Minute)
	require.Equal(t, target.Second, 30*time.Second)
}
