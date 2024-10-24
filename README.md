# Vault - Debichыы

### Поднять Redis:

Для этого Вам необходимы: 

1) Иметь Dockerfile

```Dockerfile
FROM redis:6.2

COPY --chown=redis:redis ./cert_redis/redis-server.crt /etc/redis/redis-server.crt
COPY --chown=redis:redis ./cert_redis/redis-server.key /etc/redis/redis-server.key
COPY --chown=redis:redis ./cert_redis/ca-redis.crt /etc/redis/ca.crt

COPY redis.conf /usr/local/etc/redis/redis.conf

EXPOSE 6379 6380

CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]
```

2) Иметь сгенерированные сертификаты. Перейдите в папку `./cert_redis`

Создаём корневой сертификат, которым будем подписывать (время жизни побольше)

```bash
openssl genpkey -algorithm RSA -out ca-redis.key
openssl req -x509 -new -nodes -key ca-redis.key -days 3650 -out ca-redis.crt -subj "/CN=MyRedisCA"
```

Сертификаты сервера и клиента подписываем на год. 

> Создаём и подписываем сертификаты сервера:

```bash
openssl genpkey -algorithm RSA -out redis-server.key
openssl req -new -key redis-server.key -out redis-server.csr -subj "/CN=redis-server"
openssl x509 -req -in redis-server.csr -CA ca-redis.crt -CAkey ca-redis.key -CAcreateserial -out redis-server.crt -days 365
```

> А также создаём и подписываем сертификаты клиента:

```bash
openssl genpkey -algorithm RSA -out redis-client.key
openssl req -new -key redis-client.key -out redis-client.csr -subj "/CN=redis-client"
openssl x509 -req -in redis-client.csr -CA ca-redis.crt -CAkey ca-redis.key -CAcreateserial -out redis-client.crt -days 365
```

3) Иметь `redis.conf`

```conf
tls-port 6380  # avail tls port
port 0  # reject 6379

tls-cert-file /etc/redis/redis-server.crt
tls-key-file /etc/redis/redis-server.key
tls-ca-cert-file /etc/redis/ca.crt

tls-auth-clients yes  # req cert
```

4) Собрать образ

```bash
docker build -t redis-tls .
```

5) Запустить контейнер

```bash
docker run -d --name redis-tls --network docker_network -p 6380:6380 redis-tls
```

6) Проверьте состояние контейнера 

```bash
docker ps | grep redis-tls
```

Вы должны увидеть:
```
cf54398cc549   redis-tls             "docker-entrypoint.s…"   4 minutes ago   Up 4 minutes   6379/tcp, 0.0.0.0:6380->6380/tcp, :::6380->6380/tcp             redis-tls
```

7) Перейдите в папку `./cert_redis` и попробуйте подключиться к БД

Я устанавливал для проверки redis-tools (5:7.0.15-1build2)

```bash
redis-cli --tls   --cert redis-client.crt   --key redis-client.key   --cacert ca-redis.crt   -h 127.0.0.1 -p 6380
```
