package configure
import(
	"log"
	"os"
	"bufio"
	"io"
	"strings"
)

const(
	DEFAULT_CONFIG_PATH = "./config.ini"
)

func InitConfig(configPath string) (config *Config){
	if strings.TrimSpace(configPath) == "" {
		configPath = DEFAULT_CONFIG_PATH
	}
	str := readConfig(configPath)
	config = parseConfigStr(str)
	return
}

func readConfig(configPath string) (configStr string){
	f, err := os.Open(configPath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	buf := bufio.NewReader(f)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			} else if len(line) == 0 {
				break
			}
		}
		if len(line) > 0 {
			configStr += line + "\n"
		}
	}
	return
}

func parseConfigStr(configStr string) (config *Config) {
	config = &Config{}
	config.Init()

	lines := strings.Split(configStr, "\n")
	var e *Entity

	for _, line := range lines {
		if len(line) > 0 {
			if line[0] == '#' {
				// 注释
				continue
			}
			if line[0] == '[' && line[len(line)-1] == ']' {
				name := parseConfigEName(line)
				e = &Entity{name, make(map[string]string)}
				config.AddEntity(e)
			} else {
				key,value := parseConfigLine(line)
				if e == nil {
					e = config.GetGloablEntity()
				}
				e.AddAttr(key, value)
			}
		}

	}
	return
}

func parseConfigEName(line string) (name string) {
	name = strings.Trim(line, " []")
	return
}

func parseConfigLine(line string) (key string, value string) {
	kv := strings.Split(line, "=")
	if len(kv) != 2 {
		log.Fatal("配置文件格式错误")
	}

	key = strings.TrimSpace(kv[0])
	value = strings.TrimSpace(kv[1])
	return
}
