## 1. 启动数据库

```
cd docker
docker-compose up -d
```

启动数据库后，使用用户root登录数据库，地址为localhost:8000，密码为`docker/docker-compose.yaml`中环境变量`MYSQL_ROOT_PASSWORD`的值。

登录数据库后需要手动创建名为`base-lite`的数据库。

## 2. 启动base-backend-lite

修改`config/base_backend_lite_config.yaml`中数据库`database`中的相关配置，之后使用以下命令启动后端服务。

```
go run main.go
```

## 3. 关闭数据库

```
cd docker
docker-compose down -v
```
