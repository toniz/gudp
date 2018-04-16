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

* noquote:{"table_name":"", "values":""}  
```
not quoted when replacing the parameters.
替换table_name和values这两个参数的时候,不在两边加引号.

eg: 
sql : select * from $table_name$ where uid = $uid$
noquote:{"table_name":""}

table_name = "t_user"
uid = "abc"
**real sql**: select * from t_user where uid = "abc";

```

* noescape: {"condition":""}
```
not escape string when replacing the parameters.
决定加不加转义字符.需要往数据库写入引号的时候,要加上.
```

* check: {"condition": "^.*$"}
```
使用正则表达式校验client传过来的参数是否符合要求.
Use regular expressions to check whether the client parameters match the rule.
eg: "check":   {"id": "^\\d+$"}
The id parameter must be number string.
```

* "db": "db_t_gpsbox_w"
```
对于DB配置里面的db_t_gpsbox_w. 表示这句SQL到这个DB请求。


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

* sharding: {"dbseq": ""}
```
使用client传过来的dbseq值,替换dbname里面对应的值
Use dbseq passed from client to replace the value in dbname

```


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

Put the sql configure in "sqlgroup". It will execute with transaction.
