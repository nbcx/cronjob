package ext

type Config struct {
	Id   string
	Cmd  string
	Args string
	Ct   string
}

type ConfigInterface interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	GetString(key string, defaultValue ...string) string
	GetSectionString(section string, key string, defaultValue ...string) string
}
