[WIP] The simple API that provides an HTTP source for the ClickHouse external dictionaries.

Supported databases: 

* PostgreSQL 
* MySQL  

Supported output formats:

* TabSeparated
* JSONEachRow

Usage: 

```sh
$ go build 

$ ./clickhouse-dictionary-api -h

Usage of ./clickhouse-dictionary-api:
  -c string
         (default "config.yaml")

curl http://127.0.0.1:8080/dictionary/countries?format=JSONEachRow

{"id":316,"name":"RU"}
{"id":316,"name":"CY"}

```