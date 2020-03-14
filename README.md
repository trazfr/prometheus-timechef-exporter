# prometheus-timechef-exporter

Prometheus monitoring of one's account on Timechef (https://timechef.elior.com/).

## Purpose

Use this only if you have an account on Timechef and you want to monitor your balance.

## Install

Having a working Golang environment:

```bash
go get github.com/trazfr/prometheus-timechef-exporter
go install github.com/trazfr/prometheus-timechef-exporter
```

## Use

This performs OAuth authentication and queries to Timechef's service to get your balance.  
It exports the following metric:

- `timechef_solde` which is your balance

To use it, just run `prometheus-timechef-exporter config.json`

### Example of configuration file

config.json:

```json
{
    "listen": ":9091",
    "user": "mail@example.com",
    "password": "MyPassword",
    "timeout: 20
}
```

- `user` is your login
- `password` is your password
- `timeout` is `10` seconds by default
- `listen` is `:9091` by default, so you may configure a Prometheus running on the same server:

```yaml
- job_name: prometheus-timechef-exporter
  scrape_timeout: 1m
  scrape_interval: 5m
  static_configs:
  - targets: ['127.0.0.1:9091'] 
```
