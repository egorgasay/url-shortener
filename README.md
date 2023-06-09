# go-url-shortener [![autotests](https://github.com/egorgasay/url-shortener/actions/workflows/shortenertest.yml/badge.svg?branch=iter24)](https://github.com/egorgasay/url-shortener/actions/workflows/shortenertest.yml)

### 🔍️ Purpose

Server on Go (Go). Accepts a link to a web resource from the client and, using a text shortening algorithm, shortens it and gives it back.   
The new short link will automatically redirect everyone clients to the original (longer) link.

### 🔴 Endpoints

```http
- Create link 
GET /api/shorten or /:id
- Get all links 
GET /api/user/urls
- Get one link 
POST /
- Ping 
GET /ping
- Get Stats 
GET /api/internal/stats
- Batch create 
POST /api/shorten/batch
- Delete links 
DELETE /api/user/urls
```

### ⚙️ Configuration

#### 🔧 json
```json
{
  "server_address": "localhost:8090",
  "base_url": "http://localhost",
  "enable_https": true,
  "storage": "sqlite3",
  "database_dsn" : "urls_db"
}
```
#### 🚩 flags
```
grpc - ip for gRPC -grpc=host:port
a - ip for REST -a=host
b base url -b=URL
f - path to the file to be used as a database -f=path
stype - storage type (sqlite3, mysql, postgres) -s=storage
d - connection string -d=connection_string
vdb - virtual db name -vdb=qdfh12
s - enable a HTTPS connection -s
c - path to config -c=path/to/conf.json
config - path to config -config=path/to/conf.json
t - trusted subnet -t=192.168.0.0/24
```
