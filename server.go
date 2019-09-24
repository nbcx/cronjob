package main

import (
	"fmt"
	"github.com/nbcx/cronjob/ext"
	"log"
	"net/http"
	"strings"
)

var assign *assignment

func NewServer(conf *Config) ext.ServerInterface {
	assign = NewAssignment(conf)
	return &Server{
		conf: conf,
	}
}

type Server struct {
	conf *Config
}

func (this *Server) Run() {
	conf := this.conf
	addr := fmt.Sprintf("%s:%d", conf.ip, conf.port)

	logger.Info("server run in:", addr)

	http.HandleFunc("/script", script)
	http.HandleFunc("/crontab", crontab)

	err := http.ListenAndServe(addr, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func script(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的

	args := strings.Join(r.PostForm["args"], " ")

	fmt.Fprintf(w, "+ %t", assign.flag)
	key := r.PostFormValue("key")
	value := r.PostFormValue("value")
	if key != "" {
		go assign.run(key, value, args)
	}
}

func crontab(w http.ResponseWriter, r *http.Request) {

}
