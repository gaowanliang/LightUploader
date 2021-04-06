# OneDriveUploader

萌咖大佬写了一个 [非常好的版本](https://github.com/MoeClub/OneList/tree/master/OneDriveUploader) ，可惜并没有开源，而且已经好久都没有更新了。这个项目作为从 [DownloadBot](https://github.com/gaowanliang/DownloadBot) 中独立出来的一个简易上传工具，使得上传到OneDrive更加方便。


- 支持 国际版, 个人版(家庭版).
- 支持上传文件和文件夹到指定目录,并保持上传前的目录结构.
- 支持命令参数使用, 方便外部程序调用.
- 支持自定义上传分块大小.
- 支持多线程上传(多文件同时上传).
- 支持根据文件大小动态调整重试次数.
- 支持通过Telegram Bot实时监控上传进度，方便使用全自动下载脚本时对上传的实时监控

## 授权
### 通过下面URL登录 (右键新标签打开)
#### 国际版, 个人版(家庭版)
[https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=ad5e65fd-856d-4356-aefc-537a9700c137&response_type=code&redirect_uri=http://localhost/onedrive-login&response_mode=query&scope=offline_access%20User.Read%20Files.ReadWrite.All](https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=ad5e65fd-856d-4356-aefc-537a9700c137&response_type=code&redirect_uri=http://localhost/onedrive-login&response_mode=query&scope=offline_access%20User.Read%20Files.ReadWrite.All)


### 初始化配置文件
```bash
# 国际版
OneDriveUploader -a "url"
# 国际版，并使用中文语言包
OneDriveUploader -a "url" -l zh-CN
# 个人版(家庭版)
OneDriveUploader -a "url" -v 1
# 中国版(世纪互联) 目前设计中，暂不可用
OneDriveUploader -a "url" -v 2

# 在浏览器地址栏中获取以 http://loaclhost 开头的整个url内容
# 将获取的完整url内容替换命令中的 url 三个字母
# 每次产生的 url 只能用一次, 重试请重新获取 url
# 此操作将会自动初始化的配置文件
```

## 使用
```c
Usage of OneDriveUploader:
  -a string
        // 初始化授权
        Setup and Init auth.json.
  -b string
        // 自定义上传分块大小, 可以提高网络吞吐量, 受限于磁盘性能和网络速度.
  -c string
        // 配置文件路径

  -r string
        // 上传到网盘中的某个目录, 默认: 根目录
  -l string
        // 软件语言
  -f string
        // *必要参数, 要上传的文件或文件夹
  -t string
        // 线程数, 同时上传文件的个数. 默认: 3
  -to int
        //单个数据包超时时间，默认为60s
  -tgbot string
        //使用Telegram机器人实时监控上传，此处需填写机器人的access token，形如123456789:xxxxxxxxx，输入时需使用双引号包裹。当写入内容为“1”时，使用配置文件中的BotKey和UserID作为载入项
  -uid string
        // 使用Telegram机器人实时监控上传，此处需填写接收人的userID，形如123456789
  -v int
        // 选择版本，其中0为国际版，1为个人版(家庭版)，默认为0
```

## 配置
```jsonc
{
    // 授权令牌
    "RefreshToken": "1234564567890ABCDEF",
    // 最大线程数.(同时上传文件的数量)
    "ThreadNum": 2,
    // 最大上传分块大小.(每次上传文件的最大分块大小,网络不好建议调低. 单位:MB)
    "BlockSize": 10,
    // 最大单文件大小.(目前: 个人版(家庭版)单文件限制为100GB; 其他版本单文件限制为15GB,微软将逐步更新为100GB. 单位:GB)
    "SigleFile": 100,
    // 如果是中国版(世纪互联), 此项应为 true.
    "MainLand": false,
    //软件语言
    "Language": "zh-CN",
    //超时时间
    "TimeOut": 60,
    //Telegram Bot的key
    "BotKey": "",
    //Telegram 用户ID
    "UserID": ""
}
```
注意，当同时使用配置文件，并且配置文件中的参数不为默认值时，默认将使用配置文件中的参数，当在命令行中输入参数时，会使用当前在命令行中输入的参数，并改变配置文件中的参数。即命令行中的参数优先级高于配置文件中的参数，并会改变配置文件中的参数.

## 示例
```bash
# 一些示例:

# 将同目录下的 mm00.jpg 文件上传到 OneDrive 网盘根目录
OneDriveUploader -c xxx.json -f "mm00.jpg"

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录
OneDriveUploader -c xxx.json -f "Download" 

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘Test目录中
OneDriveUploader -c xxx.json -f "Download" -r "Test"

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程
OneDriveUploader -c xxx.json -t 10 -f "Download" 

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程，同时设置超时时间为30秒
OneDriveUploader -c xxx.json -t 10 -f "Download" -to 30

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程，同时使用 Telegram Bot 实时监控上传进度
OneDriveUploader -c xxx.json -t 10 -f "Download" -tgbot "123456:xxxxxxxx" -uid 123456789

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程，同时使用配置文件中的 Telegram Bot 参数载入程序实时监控上传进度（前提是配置文件中含有Telegram Bot 的参数）
OneDriveUploader -c xxx.json -t 10 -f "Download" -tgbot "1"

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 15 线程, 并设置分块大小为 20M
OneDriveUploader -c xxx.json -t 15 -b 20 -f "Download" 

```

## 注意
当上传未出现问题，返回0，可作为上传是否失败的凭证