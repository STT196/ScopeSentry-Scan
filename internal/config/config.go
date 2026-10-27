// config-------------------------------------
// @file      : config.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/9/6 22:01
// -------------------------------------------

package config

import (
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// LoadConfig 读取配置文件并解析
func LoadConfig() error {
	// 尝试打开配置文件
	if _, err := os.Stat(global.ConfigPath); err == nil {
		global.FirstRun = false
		// 配置文件存在，从文件读取
		if err := utils.Tools.ReadYAMLFile(global.ConfigPath, &global.AppConfig); err != nil {
			return err
		}
	} else {
		global.FirstRun = true
		// 配置文件不存在，从环境变量读取
		global.AppConfig = global.Config{
			NodeName:     getEnv("NodeName", ""),
			TimeZoneName: getEnv("TimeZoneName", "Asia/Shanghai"),
			Proxy:        getEnv("PROXY", ""),
			MongoDB: global.MongoDBConfig{
				URI:      getEnv("MONGODB_URI", ""),
				IP:       getEnv("MONGODB_IP", ""),
				Port:     getEnv("MONGODB_PORT", "27017"),
				User:     getEnv("MONGODB_USER", ""),
				Password: getEnv("MONGODB_PASSWORD", ""),
				Database: getEnv("MONGODB_DATABASE", ""),
			},
			Redis: global.RedisConfig{
				URI:      getEnv("REDIS_URL", ""),
				IP:       getEnv("REDIS_IP", ""),
				Port:     getEnv("REDIS_PORT", "6379"),
				Password: getEnv("REDIS_PASSWORD", ""),
			},
		}
		debug := getEnv("DEBUG", "false")
		if debug == "true" {
			global.AppConfig.Debug = true
		} else {
			global.AppConfig.Debug = false
		}
		nodeName := global.AppConfig.NodeName
		if nodeName == "" {
			hostname, err := os.Hostname()
			if err != nil {
				fmt.Println("Error:", err)
				hostname = utils.Tools.GenerateRandomString(6)
			}
			global.AppConfig.NodeName = hostname + "-" + utils.Tools.GenerateRandomString(6)
		}
		global.AppConfig.State = 1
		// 创建配置文件
		if err := createConfigFile(global.ConfigPath, global.AppConfig); err != nil {
			return err
		}
		if global.AppConfig.MongoDB.URI == "" && global.AppConfig.MongoDB.IP == "" {
			return fmt.Errorf("missing required MongoDB IP or URI configuration")
		}

	}
	return nil
}

// getEnv 从环境变量中读取配置值，如果环境变量未设置，则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func createConfigFile(configFile string, config global.Config) error {
	if err := utils.Tools.WriteYAMLFile(configFile, config); err != nil {
		return err
	}
	log.Printf("Configuration file created: %s", configFile)
	return nil
}

func CreateDir() {
	dirs := []string{
		global.ConfigDir,
		filepath.Join(global.ConfigDir, "dir"),
		filepath.Join(global.ConfigDir, "subdomain"),
		global.DictPath,
		global.ExtDir,
		filepath.Join(global.ExtDir, "rad"),
		global.PocDir,
		filepath.Join(global.AbsolutePath, "data"),
		global.PluginDir,
		global.TmpDir,
	}
	for _, scanMopule := range global.ScanModule {
		dirs = append(dirs, filepath.Join(global.PluginDir, scanMopule))
	}

	for _, dir := range dirs {
		err := utils.Tools.EnsureDir(dir)
		if err != nil {
			log.Fatalf("%s create error: %v", dir, err)
		}
	}

}

func InitFilterUrlRe() {
	disallowedRegex := `(?i)\.(png|apng|bmp|gif|ico|cur|jpg|jpeg|jfif|pjp|pjpeg|svg|tif|tiff|webp|xbm|3gp|aac|flac|mpg|mpeg|mp3|mp4|m4a|m4v|m4p|oga|ogg|ogv|mov|wav|webm|eot|woff|woff2|ttf|otf)(?:\?|#|$)`
	global.DisallowedURLFilters = append(global.DisallowedURLFilters, regexp.MustCompile(disallowedRegex))
}

// proxyEnv returns the process env vars to set when an egress proxy is
// configured. Empty proxy returns nil (direct connection). localhost is excluded
// from proxying so the node's own DB/Redis connections are never routed out.
func proxyEnv(proxy string) map[string]string {
	proxy = strings.TrimSpace(proxy)
	if proxy == "" {
		return nil
	}
	return map[string]string{
		"HTTP_PROXY":  proxy,
		"HTTPS_PROXY": proxy,
		"ALL_PROXY":   proxy,
		"http_proxy":  proxy,
		"https_proxy": proxy,
		"all_proxy":   proxy,
		"NO_PROXY":    "localhost,127.0.0.1",
		"no_proxy":    "localhost,127.0.0.1",
	}
}

// applyProxyConfig sets proxy env vars from the loaded config so HTTP-layer
// scanning tools route through the configured egress.
func applyProxyConfig() {
	for key, value := range proxyEnv(global.AppConfig.Proxy) {
		_ = os.Setenv(key, value)
	}
}

// ApplyTaskProxy applies a per-task egress proxy supplied by the orchestrator.
// Empty proxy leaves the node default untouched. External tools spawned by the
// node inherit these vars per task; in-process HTTP clients may cache the proxy
// on first use, so strict per-task isolation across concurrent tasks on one node
// is not guaranteed — prefer one node per egress in that case.
func ApplyTaskProxy(proxy string) {
	for key, value := range proxyEnv(proxy) {
		_ = os.Setenv(key, value)
	}
}

func Initialize() {
	global.VERSION = "1.9.1"
	fmt.Printf("version %v\n", global.VERSION)
	global.AbsolutePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	global.ConfigDir = filepath.Join(global.AbsolutePath, "config")
	global.TmpDir = filepath.Join(global.AbsolutePath, "tmp")
	global.ConfigPath = filepath.Join(global.ConfigDir, "config.yaml")
	ModulesConfigPath = filepath.Join(global.ConfigDir, "modules.yaml")
	global.DictPath = filepath.Join(global.AbsolutePath, "dictionaries")
	global.ExtDir = filepath.Join(global.AbsolutePath, "ext")
	global.PocDir = filepath.Join(global.AbsolutePath, "poc")
	global.PluginDir = filepath.Join(global.AbsolutePath, "plugin")
	CreateDir()
	err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	applyProxyConfig()

	// 初始化子域名接管规则
	err = json.Unmarshal(global.TakeoverFinger, &global.SubdomainTakerFingers)
	if err != nil {
        log.Fatalf("failed to initialize subdomain takeover rules: %v", err)
	}
	InitFilterUrlRe()
	global.CustomMapParameter = sync.Map{}
	global.TmpCustomMapParameter = sync.Map{}
}
