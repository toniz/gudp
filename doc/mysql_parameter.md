# SQL File Parameter Instructions

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

* sharding: {"dbseq": ""}
```
使用client传过来的dbseq值,替换dbname里面对应的值
Use dbseq passed from client to replace the value in dbname

```



