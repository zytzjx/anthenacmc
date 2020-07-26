### redis hash
```
HMSET runoobkey name "redis tutorial" description "redis basic commands for caching" likes 20 visitors 23000
HGetAll runoobkey
HDel runoobkey description
HSETNX runoobkey likes "foo"   #like still is 20
Del runoobkey
```

# Redis Keys
* serialconfig  [HMSET]  
    record cmc server information

* transaction  [HMSET]  
    all info send cmc server