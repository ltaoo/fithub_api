## 本地开发

### 数据库迁移

```bash
# 降级到指定迁移版本
migrate -path ./migrations --database "sqlite3://./myapi.db" goto 8
```

### 手动迁移数据库

## 打包 linux

```bash
CC=x86_64-linux-musl-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o fithub_api -trimpath -ldflags "-extldflags -static" cmd/server/main.go
```

## 常见问题

### error: Dirty database version 2. Fix and force version.