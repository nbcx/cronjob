package main

import (
	"github.com/nbcx/cronjob/ext"
	"plugin"
)

//加载插件
func loadOn(path string) (plugin.Symbol, error) {
	has, err := ext.FileExists(path)

	if !has {
		return nil, err
	}

	//打开动态库
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	//接口验证
	f, err := p.Lookup("New")
	if err != nil {
		return nil, err
	}

	return f, err
}
