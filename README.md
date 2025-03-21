# logtail
日志收集工具

### 执行方式

``` shell
# 指定配置目录，重新递归读取目录中的yaml配置文件。如果新增配置文件需要重启
# 需要管理员权限，默认工作目录在/var/lib/logtail
$ sudo logtail


```

### 注册服务

``` shell
sudo cp logtail.service /etc/systemd/system/logtail.service
sudo systemctl daemon-reload
sudo systemctl start logtail
sudo systemctl stop logtail
sudo systemctl restart logtail
sudo systemctl status logtail
```

### yaml配置格式


``` yaml
name: dati
files:
    -
        name: /mnt/games/tools/nginx/html/slog/dati/{date}/c_client.log
        topic: dati:c_client
    -
        name: /mnt/games/tools/nginx/html/slog/dati/{date}/c_real.log
        topic: dati:c_real
    -
        name: /mnt/games/tools/nginx/html/slog/dati/{date}/s_login.log
        topic: dati:s_login
    -
        name: /mnt/games/tools/nginx/html/slog/dati/{date}/s_logout.log
        topic: dati:s_logout
    -
        name: /mnt/games/tools/nginx/html/slog/dati/{date}/s_register.log
        topic: dati:s_register
    -
        name: /mnt/games/tools/nginx/html/slog/dati/{date}/s_resource.log
        topic: dati:s_resource

# out:
#     kafak:
#         hosts: ["kafka01:8101"]
#         sasl: plain # plain or scram
#         username:
#         password:

```