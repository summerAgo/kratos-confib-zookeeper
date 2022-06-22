
# kratos-contrib-zookeeper

在kratos框架中使用zookeeper读取配置

```go
zksource, err := zkconf.New(conn, zkconf.WithPath(sc.Serverpath))
if err != nil {
   panic(err)
}

c := config.New(config.WithSource(zksource), config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
    return yaml.Unmarshal(kv.Value, v)
}))

defer c.Close()

if err = c.Load(); err != nil {
   panic(err)
}
if err = c.Scan(&bc); err != nil {
   panic(err)
}
```