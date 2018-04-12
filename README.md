# GRPC Unified Data Proxy(GUDP)
GUDP is a Unified Data Proxy using GRPC. Currently support mysql and redis(redis is an EXPERIMENTAL Feature.). Later will add more data source support. Developers only need to use GRPC to invoke this service without concern for database connection and query details. Meet most data service requirements, and pursue simplicity, lightweight and scalability. 

## GUDP features:
* Unify data access. business layer
* Connections convergence. Several thousands of data request connections in the business layer. If you directly connect to the database, it will cause a waste of database resources. After GUDP, the number of links to the database can be optimally utilized.
* Work hard to make it more scalable and support more data sources.
### MYSQL:
* Easy to use, modify the query logic does not need to change the business code, just change the GUDP SQL configuration statement. If you want to switching the database, just change the GUDP DB configuration.
* Security, business layer will not be exposed to the data source. Through the replacement of parameters to achieve data query, a good defense SQL injection.
* Supports DB fragmentation.
* Support multi-database transaction commit.
### REDIS:
* Easy to use. 

___

# 简介
GUDP是一个统一访问代理．目前支持mysql. Redis是实验特性,仅实现部分. 后面会添加更多数据源支持.
开发人员只需要使用GRPC调用本服务,无需关心数据库连接和查询细节.满足大部分的数据服务需求,并且追求简单,轻量化以及扩展性.使用GO语言编码.
GUDP特性:
* 使得数据访问统一化
* 链接收敛.业务层几千个数据请求链接,如果直接访问数据库,会造成数据库资源浪费的问题.通过GUDP后,到数据库的链接个数能达到最优利用.
* 努力使其具备更好的扩展性,支持更多的数据源.

MYSQL功能比较完善，目前有如下特性：
* 简单易用,修改查询逻辑不需要改业务代码或者切换数据库,只需要改动GUDP的SQL配置语句.
* 读写分离
* 安全,业务层不会接触到数据源.通过参数替换达到数据查询效果，屏蔽了注入SQL的途径.
* 支持数据库分片.
* 支持多数据库事务提交.

REDIS是实验性功能:
* 简单易用,配置好配置文件就可访问redis.业务层不无关系链接细节.


## Example:
Mysql Read Write Spilting  
[MYSQL读写分离实现](doc/mysql_read_write_splitting.md)。 

Mysql Sharding Example  
[MYSQL数据库分片实现](doc/mysql_db_sharding.md)。 

Mysql Multi DB Transcation  
[MYSQL多数据库事务实现](doc/mysql_multi_db_transaction.md)。 


### Mysql:
[Mysql Test Client](client/mysqlcli.go)  
[mysql db configure](conf/mysql/db)  
[mysql sql configure](conf/mysql/sql)  
[testdata](doc/mysqldata.sql)  

```
go get -u github.com/toniz/gudp
export GRPC_GO_LOG_SEVERITY_LEVEL=INFO
nohup gudp &

cd $GOPATH/src/github.com/toniz/gudp/client
go build mysqlcli.go 
./mysqlcli
```

### Redis:
[Redis Test Client](client/rediscli.go)  
[Redis Srv Configure](conf/redis/srv)  

```
cd $GOPATH/src/github.com/toniz/gudp/client
go build rediscli.go 
./rediscli

```


