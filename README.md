# go_mongo

Единое подключение к MongoDB на `mongo-driver/v2` для микросервисов SHINPI.

## Установка

```
go get github.com/shinpi-tech/go_mongo
```

## Использование

```go
db, err := mongox.Connect(ctx, mongox.Config{
    Host:                cfg.Mongo.Host,
    Port:                cfg.Mongo.Port,
    Database:            "catalog",
    User:                cfg.Mongo.User,
    Password:            cfg.Mongo.Password,
    AuthSource:          cfg.Mongo.AuthSource,
    DirectConnection:    cfg.Mongo.DirectConnection,
    ObjectIDAsHexString: true,
})
if err != nil {
    log.Fatal(err)
}
defer db.Client().Disconnect(ctx)
```

`ObjectIDAsHexString=true` декодирует `ObjectID` в hex-строки (как в сервисах, работающих со строковыми ID);
сервисы, использующие `bson.ObjectID` напрямую, оставляют `false`.
