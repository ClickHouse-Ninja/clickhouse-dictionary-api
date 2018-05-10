[WIP] The simple API that provides an HTTP source for the ClickHouse external dictionaries.

Supported databases: 

* PostgreSQL 
* MySQL  

Usage: 

```sh
$ go build 

$ ./clickhouse-dictionary-api -h

Usage of ./clickhouse-dictionary-api:
  -addr string
         (default ":8080")
  -c string
         (default "config.yaml")

curl http://127.0.0.1:8080/dictionary/countries

{"id":316,"name":"RU"}
{"id":316,"name":"CY"}

```