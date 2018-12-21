package kernal

import (
    "time"
    "context"
    "encoding/json"
	"github.com/coreos/etcd/clientv3"
)

func _watch_config(ctx context.Context, sname string, c interface{}, init_chann *chan int) {
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{ETCD_SERVER},
        DialTimeout: 2 * time.Second,
    })
    ErrorPanic(err)
    defer cli.Close()

    // 先获取一次解析配置
    resp, err := cli.Get(ctx, "/config/" + sname)
    if err != nil {
        panic(err)
    }
    for _, ev := range resp.Kvs {
        // Debug(ctx, "%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
        Debug(ctx, "ev.Value: %s", ev.Value)
        err = json.Unmarshal(ev.Value, c)
        if err != nil {
            Error(ctx, "[_watch_config] err: %+v", err)
        }
        *init_chann <- 0
    }

    // 监控key的变化
    wcli := clientv3.NewWatcher(cli)

    rch := wcli.Watch(ctx, "/config/" + sname, clientv3.WithPrefix())

    for wresp := range rch {
        for _, ev := range wresp.Events {
            Debug(ctx, "%s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
            err = json.Unmarshal(ev.Kv.Value, c)
            if err != nil {
                Error(ctx, "[_watch_config] err: %+v", err)
            }
        }
    }
}

func WatchConfig(ctx context.Context, sname string, c interface{}) {
    init_chan := make(chan int)
    go _watch_config(ctx, sname, c, &init_chan)

    <-init_chan
}
