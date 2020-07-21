![](https://www.futuredial.com/wp-content/uploads/2020/02/futuredial-logo-color.svg)

### Discraptioncd 
    This is Athena project background function.

### Technology stack
| Syntax      | Description |
| :----------- | :----------- |
| Logger      | [lfshook](https://github.com/rifflocklfshook)|
|Redis|![](https://redis.io/images/redis-white.png)|


```
go get github.com/go-redis/redis/v8
go get github.com/go-resty/resty/v2 
go get github.com/lestrrat-go/file-rotatelogs 
go get github.com/rifflock/lfshook
go get github.com/sirupsen/logrus 
go get gopkg.in/yaml.v2 
```


$$ J_\alpha(x) = \sum_{m=0}^\infty \frac{(-1)^m}{m! \Gamma (m + \alpha + 1)} {\left({ \frac{x}{2} }\right)}^{2m + \alpha} \text {，行内公式示例  } $$


this is athena function library. any other function only main.go.
uuid get config
Results[0] save redis DB, Key is **serialconfig**, value is hash.
```
{
   "ok":1,
   "results":[
      {
         "companyid":55,
         "staticfileserver":"http://cmc-dl.futuredial.com/",
         "siteid":"73",
         "installitunes":"True",
         "adminconsoleserver":"http://cmc.futuredial.com/",
         "pname":"CMC GreenT V5",
         "_id":"6c87ceb9-25a3-4e09-b81b-fb0a57b64d42",
         "solutionid":"2",
         "webserviceserver":"http://cmc.futuredial.com/ws/",
         "productid":"2"
      }
   ],
   "id":34787
}

```
### CommandLine
```
anthenacmc -udid="6c87ceb9-25a3-4e09-b81b-fb0a57b64d42"
anthenacmc -login -username=qa -password=qa
```