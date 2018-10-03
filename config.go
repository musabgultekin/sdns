package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jonboulle/clockwork"
	"github.com/semihalev/log"
)

// BuildVersion returns the build version of sdns, this should be incremented every new release
var BuildVersion = "2.0.2"

// ConfigVersion returns the version of sdns, this should be incremented every time the config changes so sdns presents a warning
var ConfigVersion = "2.0.2"

type config struct {
	Version        string
	Sources        []string
	SourceDirs     []string
	Log            string
	LogLevel       string
	Bind           string
	API            string
	Nullroute      string
	Nullroutev6    string
	OutboundIP     string
	Interval       int
	Timeout        int
	ConnectTimeout int
	Expire         uint32
	Maxcount       int
	Maxdepth       int
	Blocklist      []string
	Whitelist      []string
}

var defaultConfig = `# version this config was generated from
version = "%s"

# list of sources to pull blocklists from, stores them in ./sources
sources = [
"http://mirror1.malwaredomains.com/files/justdomains",
"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
"http://sysctl.org/cameleon/hosts",
"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist",
"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt",
"https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt",
"http://hosts-file.net/ad_servers.txt",
"https://raw.githubusercontent.com/quidsup/notrack/master/trackers.txt"
]

# list of locations to recursively read blocklists from (warning, every file found is assumed to be a hosts-file or domain list)
sourcedirs = [
"sources"
]

# what kind of information should be logged, Log verbosity level [crit,error,warn,info,debug]
loglevel = "info"

# address to bind to for the DNS server
bind = "0.0.0.0:53"

# outbound ip address
#outboundip = ""

# address to bind to for the API server
api = "127.0.0.1:8080"

# ipv4 address to forward blocked queries to
nullroute = "0.0.0.0"

# ipv6 address to forward blocked queries to
nullroutev6 = "0:0:0:0:0:0:0:0"

# concurrency interval for lookups in miliseconds
interval = 200

# query timeout for dns lookups in seconds
timeout = 5

# connect timeout for dns lookups in seconds
connecttimeout = 2

# cache entry lifespan in seconds
expire = 600

# cache capacity, 0 for infinite
maxcount = 0

# maximum recursion depth for nameservers
maxdepth = 30

# manual blocklist entries
blocklist = []

# manual whitelist entries
whitelist = [
	"getsentry.com",
	"www.getsentry.com"
]
`

// WallClock is the wall clock
var WallClock = clockwork.NewRealClock()

// Config is the global configuration
var Config config

// LoadConfig loads the given config file
func LoadConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := generateConfig(path); err != nil {
			return err
		}
	}

	if _, err := toml.DecodeFile(path, &Config); err != nil {
		return fmt.Errorf("could not load config: %s", err)
	}

	if Config.Version != ConfigVersion {
		if Config.Version == "" {
			Config.Version = "none"
		}

		log.Warn("Config file sdns.toml is out of date!")
	}

	return nil
}

func generateConfig(path string) error {
	output, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not generate config: %s", err)
	}
	defer output.Close()

	r := strings.NewReader(fmt.Sprintf(defaultConfig, ConfigVersion))
	if _, err := io.Copy(output, r); err != nil {
		return fmt.Errorf("could not copy default config: %s", err)
	}

	if abs, err := filepath.Abs(path); err == nil {
		log.Info("Default config file generated", "config", abs)
	}

	return nil
}
