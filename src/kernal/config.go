package kernal

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"time"
)

const (
    config_prefix       = "/config/"
    global_config_name  = "__global__"
)

type global_config struct {
    Env         string
    Log_path    string
    Jaeger_url  string
}

func _watch_config(sname string, c interface{}, init_chann *chan int) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ETCD_SERVER},
		DialTimeout: 2 * time.Second,
	})
	ErrorPanic(err)
	defer cli.Close()

    ctx := context.TODO()

	// 先获取一次解析配置
	resp, err := cli.Get(ctx, config_prefix + sname)
	if err != nil {
		panic(err)
	}
	for _, ev := range resp.Kvs {
		// Debug(ctx, "%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		// Debug(ctx, "ev.Value: %s", ev.Value)
		/*err = */json.Unmarshal(ev.Value, c)
		/*if err != nil {
			Error(ctx, "[_watch_config] err: %+v", err)
		}*/
		*init_chann <- 0
	}

	// 监控key的变化
	wcli := clientv3.NewWatcher(cli)

	rch := wcli.Watch(ctx, config_prefix + sname, clientv3.WithPrefix())

	for wresp := range rch {
		for _, ev := range wresp.Events {
			// Debug(ctx, "%s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
			/*err = */json.Unmarshal(ev.Kv.Value, c)
			/*if err != nil {
				Error(ctx, "[_watch_config] err: %+v", err)
			}*/
		}
	}
}

func WatchConfig(sname string, c interface{}) {
	init_chan := make(chan int)
	go _watch_config(sname, c, &init_chan)

	<-init_chan
}
