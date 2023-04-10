package util

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

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
	var path = "../.env"
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
		DBPort     string `mapstruct:"DB_PORT"`
		DBName     string `mapstruct:"DB_NAME"`
		DBPassword string `mapstruct:"DB_PASSWORD"`
	}
	env := readEnv(t)

	var target config
	err := Setingenv(&target, env)
	require.NoError(t, err)

	// fmt.Printf(">>> %v\n", target)

	require.Equal(t, target.DBPort, "3307")
	require.Equal(t, target.DBName, "shortnet")
	require.Equal(t, target.DBPassword, "abc123")
}
