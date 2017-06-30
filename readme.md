# 新闻聚合网站
## 目的
任选一种技术实现一个新闻资讯收集、整理，进行个性化展示的网站

## 基本功能 1. 定时到各主流新闻门户网站抓取信息  - 用户访问网站时能看到聚合的信息内容  - 实现用户注册、登录功能，用户注册时需要填写必要的信息并验证，如用户名、密码要求在 6 字  节以上，email 的格式验证，并保证用户名和 email 在系统中唯一。  - 用户登录后可以设置感兴趣的新闻资讯栏目，用户访问网站的展示页面会根据用户设置做出相应的调整- 实现一个 Android 或 iphone 客户端软件，功能同网站，但展示界面根据屏幕大小做 UI 的自适应调整，并能实现热点新闻推送- 具体一定的学习能力，能根据用户的使用习惯调整展现的内容

## 过程
### 到各主流新闻门户网站抓取信息#### 信息来源的确定打开网易新闻，今日头条，新浪新闻等网站，查看源代码，可以找到每条新闻的dom  但是用爬虫获取内容的时候并不能找到对应的dom，经过分析后发现大部分页面是用js动态生成的，只有少部分新闻（比如右侧只显示一标题的新闻）是静态的经过一番思考和寻找，找到了Phantom这个工具，它可以模拟浏览器的运行，运行js动态生成页面。尝试抓取了新闻，发现能获取到内容但是速度太慢，需要几十秒，找了半天也没找到让它加速的方法，放弃之（后来想到可能是每次启动都需要初始化，这里可能花了非常多的时间，应该考虑让它启动了之后再不断的模拟打开网页）  打开chrome的调试器，查看Network，仔细测试和查看后发现前端会通过发请求来获取新的新闻，格式是json。这就好办多了，可以直接对这个链接发送请求来获取新闻  网易新闻  
`http://temp.163.com/special/00804KVA/cm_yaowen_03.js?callback=data_callback`新浪新闻  `http://api.roll.news.sina.com.cn/zt_list?channel=news&cat_1=gnxw&cat_2==gdxw1||=gatxw||=zs-pl||=mtjj&level==1||=2&show_ext=1&show_all=1&show_num=22&tag=1&format=json&page=2&callback=newsloadercallback&_=1498536194155`今日头条  `http://www.toutiao.com/api/pc/feed/?category=news_tech&utm_source=toutiao&widen=1&max_behot_time=1498531911&max_behot_time_tmp=1498531911&tadrequire=true&as=A1E5F9A5717D362&cp=59511DA396923E1`
#### 把信息转换成我们想要的格式并存入数据库中以网易新闻为例  
获取到的新闻的格式是这样的  

```go
type Data163 struct {
	Data []News163
}
type Keyword163 struct {
	Akey_link string
	Keyname   string
}
type News163 struct {
	Add1        string
	Add2        string
	Add3        string
	Channelname string
	Commenturl  string
	Digest      string
	Docurl      string
	Imgurl      string
	Keywords    []Keyword163
	Label       string
	Newstype    string
	Pics        []string
	Tienum      int
	Time        string
	Title       string
	Tlastid     string
	Tlink       string
}
```
因为这里面有很多我们不需要的信息，并且我们需要将不同网站的新闻整合到一起，所以我们需要将新闻都转换成相同的格式，初步定成这样

```go
type News struct {
	ID     string `bson:"_id"`
	Title  string
	Docurl string
	Imgurl string
	Time   int64
	Tag    string
	Source string
}
```

为了防止插入重复的新闻，所以我们需要一个标识符`_id`来唯一辨识每条新闻，我的想法是将来源，分类，事件，标题连在一起再进行md5取hash

```go
h := md5.New()
io.WriteString(h, source163+tag+data.Time+data.Title)
id := hex.EncodeToString(h.Sum(nil))
```

光有id还不行，我们还需要知道上次插入的时间，这样才能知道应该抓到什么时间为止。所以爬到数据后，应该先从数据库中获取该来源该分类的最新的一条新闻的时间，设置为lastFecthTime（如果没有新闻就设为0），如果某条新闻的时间lastFecthTime小，就跳过，设置lastFecthTime为此次获取到的新闻的最大时间

爬新闻就先写到这里，有空再写新浪新闻和今日头条的爬虫

#### 用户访问网站时能看到聚合的信息内容
在浏览新闻网页的时候，我们会发现既可以从顶部刷最新的新闻，也可以从底部加载历史新闻，所以我们要针对这两种情况写api

url  
`/api/feed`  
param  
`topTime, bottomTime, show_num, tag, source`  
如果设置了topTime，就返回时间比topTime大的show_num条新闻，如果设置了bottomTime，就返回时间比bottomTime小的show_num条新闻，如果都没设置，就返回最新的show_num条新闻  

#### 实现用户注册、登录功能，用户注册时需要填写必要的信息并验证，如用户名、密码要求在 6 字节以上，email 的格式验证，并保证用户名和 email 在系统中唯一
保证用户名和 email 在系统中唯一的实现方式是为这两个字段添加unique的索引

```js
db.user.ensureIndex({ mail: 1,  },
        {sparse: true, unique: true }
)
db.user.ensureIndex({ nickname: 1 },
        {sparse: true, unique: true }
)
```
用户的数据结构

```go
type User struct {
	Mail     string
	NickName string
	Password string
	Tags     []string
}
```
其他的话就是要在服务器维持一个用户的session，每次服务都需要验证，防止比如在未登录的情况下发出修改tags的请求

#### 用户登录后可以设置感兴趣的新闻资讯栏目，用户访问网站的展示页面会根据用户设置做出相应的调整

对用户的tags字段进行修改，没什么好说的

#### 实现一个 Android 或 iphone 客户端软件，功能同网站，但展示界面根据屏幕大小做 UI 的自适应调整，并能实现热点新闻推送

Android 客户端软件，用webview实现

#### 具有一定的学习能力，能根据用户的使用习惯调整展现的内容

1. 根据用户的关注来展示新闻  
* 根据用户浏览过的新闻，以及对某条新闻的like和dislike以此来推荐  
* 根据用户的位置，年龄，职业等对用户进行分类，推荐这个群体关注的新闻  

#### mangodb脚本文件

```js
db.auth('username','password')
use newsAggregation
db.createUser({
    user: "user",
    pwd: "pwd",
    roles: [{
            "role" : "dbOwner",
            "db" : "newsAggregation"
    }]
})
db.createCollection("news")
db.createCollection("user")
db.user.ensureIndex({ mail: 1,  },
        {sparse: true, unique: true }
)
db.user.ensureIndex({ nickname: 1 },
        {sparse: true, unique: true }
)
db.shutdownServer()
```
