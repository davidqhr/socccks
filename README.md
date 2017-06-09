# Introduce

[![Build Status](https://travis-ci.org/davidqhr/socccks.svg?branch=master)](https://travis-ci.org/davidqhr/socccks)

I create this repo for learning golang, socks5 protocol and shadowsocks(ss).

```
                                                obstacle
                                                   |
                                                   |
                                                   |
         +--------------------+                    |                     +--------------------+
         |   your processes   |                    |                     |  real destination  |
         +--------------------+                    |                     +--------------------+
                   |                               |                               |
                   |                               |                               |
   socks5 protocol |                               |                               |
                   |                               |                               |
                   |                               |                               |
         +--------------------+            Encrypted Data pipe           +--------------------+
         |   socccks-client   | <--------------------------------------> |   socccks-server   |
         +--------------------+             through tcp socket           +--------------------+
                                                   |
----------------------------------------------------------------------------------------------------
                                                   |
              your computer                        |                         remote computer
                                                   |
                                                obstacle		
```

Socccks is a separated socks5 proxy. It allows you to access some resources that are behind a obstacle through a socks5 socket. The Data between socccks-client and socccks-server is encryped(aes-256-cfb) and no features.

## install

- install golang
- `go get github.com/davidqhr/socccks`
- `go install github.com/davidqhr/socccks`

## usage

### server side

```json
# config.json
{
  "address": "0.0.0.0",
  "daemon": true,
  "users": {
    "david": 8112,
    "monika": 8113
  }
}

```

```bash
socccks-server -c config.yml
```

### client side

```bash
socccks-client -s server_ip:server_port -l bindaddress:port -p pass -d
curl --socks5-hostname bindaddress:port https://www.google.com -v

# eg:
socccks-client -s 192.168.1.132:8113 -l localhost:1090 -p david -d
curl --socks5-hostname localhost:1090 https://www.google.com -v
```

### Licence
MIT
