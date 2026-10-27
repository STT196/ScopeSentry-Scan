// config-------------------------------------
// @file      : type.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/9/8 12:32
// -------------------------------------------

package global

// Config 结构体
type Config struct {
	NodeName     string        `yaml:"NodeName"`
	State        int           `yaml:"state"`
	TimeZoneName string        `yaml:"TimeZoneName"`
	Debug        bool          `yaml:"debug"`
	// Proxy is the egress proxy for HTTP-layer scanning tools (http/https/socks5).
	// Empty = direct. Applied process-wide as HTTP_PROXY/HTTPS_PROXY/ALL_PROXY so
	// httpx / nuclei / crawler / dirscan and other HTTP tools honor it. Raw-socket
	// scans (port/DNS) ignore it and still use the node's own IP.
	Proxy   string        `yaml:"proxy"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Redis   RedisConfig   `yaml:"redis"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	URI      string `yaml:"uri"`
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}
