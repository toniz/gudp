# SQL File Parameter Instructions

## Example:
```
"ACCOUNT_t_user_select_by_uids" :
{
    "sql" : "SELECT user_id, user_name, type FROM t_user WHERE user_id in ($condition$)",
    "noquote": {"condition":""},
    "noescape":{"condition":""},
    "check": {"condition": "^.*$"},
    "db" : "db_t_gpsbox_w"
}
```

### 1. noquote:{"table_name":"", "values":""}  
```
not quoted when replacing the parameters.
替换table_name和values这两个参数的时候,不在两边加引号.
```

eg: 
```
sql configure: 
"example1" :
{
    "sql": "select * from $table_name$ where uid = $uid$"
     noquote:{"table_name":""}
}
```

``` 
client pass value: 
req := pb.Query{
    Engine: "mysql",
    Ident:  "example1",
    Params: map[string]string{
        "table_name": "t_user",
        "uid": "abc",
    },  
} 
``` 

**real sql**: select * from t_user where uid = "abc";


### 2. noescape: {"condition":""}
```
not escape string when replacing the parameters.
决定加不加转义字符.需要往数据库写入引号的时候,要加上.
```

### 3. check: {"condition": "^.*$"}
```
使用正则表达式校验client传过来的参数是否符合要求.
Use regular expressions to check whether the client parameters match the rule.
eg: "check":   {"id": "^\\d+$"}
The id parameter must be number string.
```

### 4. "db": "db_t_gpsbox_w"
```
db_t_gpsbox_w对应的数据库配置在DB配置文件中.
Connect to this database: db_t_gpsbox_w
db_t_gpsbox_w is definded in DBConfigure.
```


## Sharding Examle:
```
"ACCOUNT_t_user_sharding":
{   
    "sql" : "SELECT user_id, user_name, type FROM t_user;",
    "sharding": {"dbseq": ""},
    "db" : "db_t_account_$dbseq$"
} 
```

### 1. sharding: {"dbseq": ""}
```
使用client传过来的dbseq值,替换dbname里面的“$dbseq$”。
Replace the value 'dbseq' in dbname. 
```

eg: Mysql Sharding Example   
[MYSQL数据库分片实现](mysql_db_sharding.md)。 


## Trancation Example
```
"ACCOUNT_t_user_insert_transaction" : 
{   
"sqlgroup": 
    [   
        {   
            "sql" : "INSERT INTO t_user(user_id, user_name, type) VALUES($id$, $name$ ,$type$);",
            "noquote": {"id":""},
            "check":   {"id": "^\\d+$"},
            "db" : "db_t_account_w"
        },  
        {   
            "sql" : "INSERT INTO t_images(id, name, image) VALUES($id$, $name$ ,$image$);",
            "noquote": {"id":""},
            "check":   {"id": "^\\d+$"},
            "db" : "db_t_image_w"
        }   
    ]   
}
```

### 1. Put the sql configure in "sqlgroup". It will execute with transaction.   

eg: 
Mysql Multi DB Transcation  
[MYSQL多数据库事务实现](mysql_multi_db_transaction.md)。 


