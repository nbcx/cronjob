package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nbcx/cronjob/ext"
	"log"
	"os"
	"strings"
	"unicode"
)

//var logger *log.Logger
var logger ext.LoggerInterface

type Config struct {
	protocol string
	ip       string
	port     int
	command  string
	path     string
	data     map[string]interface{} /*创建指令集合 */
	jobs     map[string][]string    /*创建指令集合 */
}

// test data
var (
	h bool
	s string
)

func NewConfigObject() *Config {
	conf := &Config{
		jobs: make(map[string][]string),
		data: make(map[string]interface{}),
	}
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&conf.ip, "ip", "127.0.0.1", "info")
	flag.IntVar(&conf.port, "port", 4567, "info")
	flag.StringVar(&conf.protocol, "protocol", "http", "info")
	flag.StringVar(&conf.command, "command", "http", "info")
	flag.StringVar(&conf.path, "config", "gox.ini", "info")
	// 注意 `signal`。默认是 -s string，有了 `signal` 之后，变为 -s signal
	flag.StringVar(&s, "s", "", "send `signal` to a master process: stop, quit, reopen, reload")

	// 改变默认的 Usage，flag包中的Usage 其实是一个函数类型。这里是覆盖默认函数实现，具体见后面Usage部分的分析
	flag.Usage = usage
	flag.Parse()
	if h {
		flag.Usage()
	}
	logger = NewLogger(conf)

	conf.job()
	conf.init("cronjob.ini")

	return conf
}

func (conf *Config) job() error {
	file, err := os.Open("job.sh")
	if err != nil {
		log.Printf("Cannot open text file: %s, err: [%v]", "job.sh", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // or

		strRune := []rune(line)
		//fmt.Println("first is ", string(nameRune[:1]), "\n")
		if string(strRune[:1]) == "#" {
			continue
		}

		//line := scanner.Bytes()
		args := strings.FieldsFunc(line, unicode.IsSpace) //strings.Split(line, " ")

		if len(args) < 4 {
			continue
		}

		conf.jobs[args[0]] = args
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: %s, err: [%v]", "job.sh", err)
		return err
	}
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `cronjob version: ejob/1.10.0
Usage: cronjob [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

Options:
`)
	flag.PrintDefaults()
}

// 根据文件名，段名，键名获取ini的值
func (conf *Config) init(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	// 当前读取的段的名字
	var sectionName string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linestr := scanner.Text() // or

		// 切掉行的左右两边的空白字符
		linestr = strings.TrimSpace(linestr)
		// 忽略空行
		if linestr == "" {
			continue
		}
		// 忽略注释
		if linestr[0] == ';' {
			continue
		}
		// 行首和尾巴分别是方括号的，说明是段标记的起止符
		if linestr[0] == '[' && linestr[len(linestr)-1] == ']' {
			// 将段名取出
			sectionName = linestr[1 : len(linestr)-1]
			conf.data[sectionName] = make(map[string]interface{})
			continue
		}

		// 切开等号分割的键值对
		pair := strings.Split(linestr, "=")
		// 保证切开只有1个等号分割的键值情况
		if len(pair) != 2 {
			continue
		}

		// 去掉键的多余空白字符
		key := strings.TrimSpace(pair[0])

		if sectionName == "" {
			conf.data[key] = strings.TrimSpace(pair[1])
		} else {
			conf.data[sectionName].(map[string]interface{})[key] = strings.TrimSpace(pair[1])
		}
	}
	return nil
}

func (conf *Config) Get(key string) interface{} {
	m, ok := conf.data[key]
	if !ok {
		return nil
	}
	return m
}

func (conf *Config) Set(key string, value interface{}) {
	conf.data[key] = value
}

func (conf *Config) GetString(key string, defaultValue ...string) string {
	value := conf.Get(key)
	if value != nil {
		return value.(string)
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (conf *Config) GetInt(key string, defaultValue ...int64) int64 {
	value := conf.Get(key)
	if value != nil {
		return value.(int64)
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (conf *Config) GetSection(section string, key string) interface{} {
	m, ok := conf.data[section].(map[string]interface{})
	if !ok {
		return nil
	}
	value, ok := m[key]
	if !ok {
		return nil
	}
	return value
}

func (conf *Config) GetSectionString(section string, key string, defaultValue ...string) string {
	value := conf.GetSection(section, key)
	if value != nil {
		return value.(string)
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (conf *Config) GetSectionInt(section string, key string, defaultValue ...int64) int64 {
	value := conf.GetSection(section, key)
	if value != nil {
		return value.(int64)
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}
