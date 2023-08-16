package template

var ConfigTemplate = `name: "%s"
mode: "%s"
version: "0.0.1"
port: %s

# 日志配置
LogC:
  level: "info"
  filename: "%s.log"
  max_size: 200
  max_age: 30
  max_backups: 7

# 微信配置
WechatC:
  app_id: ""
  secret: ""
  mch_id: ""
  app_key: ""
  notify_url: ""
  prefix: ""

# 数据库配置
MysqlC:
  host: ""
  port:
  user: ""
  password: ""
  dbname: ""
  max_open: 200
  max_idle: 50

# redis配置
RedisC:
  host: 127.0.0.1
  port: 6379
  password:
  db: 0
  pool_size: 100
  prefix: ""

# 七牛配置
QiNiuC:
  secret_key: ""
  access_key: ""
  pic_domain: ""
  bucket: ""
`
