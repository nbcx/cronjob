package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ServerInterface interface {
	run()
}

func NewServer(conf *Config) ServerInterface {
	return &Server{
		conf: conf,
	}
}

type Server struct {
	conf *Config
}

func (this *Server) run() {
	this.http()
}

//启动webserver
func (this *Server) http() {

	conf := this.conf
	addr := fmt.Sprintf("%s:%d", conf.ip, conf.port)

	assignment := NewAssignment(conf)
	logger.Info("server run in:", addr)

	http.HandleFunc("/script", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm() //解析参数，默认是不会解析的

		args := strings.Join(r.PostForm["args"], " ")

		fmt.Fprintf(w, "+ %t", assignment.flag)
		cmd := r.PostFormValue("cmd")
		if cmd != "" {
			go assignment.run(cmd, args)
		}
	})

	err := http.ListenAndServe(addr, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
