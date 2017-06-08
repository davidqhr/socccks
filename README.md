这个项目的初衷是为了了解go，正好最近在看sock5，所以就拿go来实现一个sock5 proxy server

## usage

```bash
ss-server help
ss-server -d -b 211.111.111.111 -p 8118 -key a_password,another_password,test_password
ss-server -d -c config.yml # command params will overwrite yml
ss-server dump-config config.yml #
ss-server status
ss-server stop
```

```bash
ss-local -d -s 211.111.111.111:8118 -p 1080 -key a_password
ss-local -d -c config.yml # command params will overwrite yml
ss-local status
```

## Feature

- Graceful Exit
