## 1. Start the database

``
cd docker
docker-compose up -d
```

After starting the database, log in to the database using the user root at localhost:8000 and the password as the value of the environment variable `MYSQL_ROOT_PASSWORD` in `docker/docker-compose.yaml`.

You need to manually create a database named `base-lite` after logging in to the database.

## 2. Starting base-backend-lite

Modify the configuration of the database `database` in `config/base_backend_lite_config.yaml`.

After that use the following command to start the backend service.

``
go run main.go
```

## 3. Close the database

```
cd docker
docker-compose down -v
```
