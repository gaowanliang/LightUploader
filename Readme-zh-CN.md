# LightUploader

萌咖大佬写了一个 [非常好的版本](https://github.com/MoeClub/OneList/tree/master/OneDriveUploader) ，可惜并没有开源，而且已经好久都没有更新了。这个项目作为从 [DownloadBot](https://github.com/gaowanliang/DownloadBot) 中独立出来的一个简易上传工具，旨在用更轻量化的方式让在各种平台都能快速的向各个网络硬盘上传数据。

- 支持 OneDrive 国际版, 个人版(家庭版)，世纪互联，Google Drive.
- 支持上传文件和文件夹到指定目录,并保持上传前的目录结构.
- 支持命令参数使用, 方便外部程序调用.
- 支持自定义上传分块大小.
- 支持多线程上传(多文件同时上传).
- 支持根据文件大小动态调整重试次数.
- 支持跳过网盘中已存在的同名文件.
- 支持通过Telegram Bot实时监控上传进度，方便使用全自动下载脚本时对上传的实时监控

## 授权
详见[wiki](https://github.com/gaowanliang/LightUploader/wiki)


### 初始化配置文件
```bash
# OneDrive 国际版
LightUploader -a "url"
# OneDrive 国际版，并使用中文语言包
LightUploader -a "url" -l zh-CN
# OneDrive 个人版(家庭版)
LightUploader -a "url" -v 1
# OneDrive 中国版(世纪互联)，并使用中文语言包
LightUploader -a "url" -v 2 -l zh-CN

# 在浏览器地址栏中获取以 http://loaclhost 开头的整个url内容
# 将获取的完整url内容替换命令中的 url 三个字母
# 每次产生的 url 只能用一次, 重试请重新获取 url
# 此操作将会自动初始化的配置文件

# Google Drive
LightUploader -v 3
```

## 使用
```c
Usage of LightUploader:
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
  -m int
        // 选择模式，0为替换网盘中同名文件，1为跳过，默认为0
  -v int
        // 选择版本，其中0为 OneDrive 国际版，1为 OneDrive 个人版(家庭版)，2为 OneDrive 世纪互联，3为Google Drive，默认为0
```

## 配置
```jsonc
{
    // 网盘类型
    "Drive":"OneDrive",
    // 授权令牌
    "RefreshToken": "1234564567890ABCDEF",
    // 最大线程数.(同时上传文件的数量)
    "ThreadNum": 2,
    // 最大上传分块大小.(每次上传文件的最大分块大小,网络不好建议调低. 单位:MB)
    "BlockSize": 10,
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
LightUploader -c xxx.json -f "mm00.jpg"

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录
LightUploader -c xxx.json -f "Download" 

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘Test目录中
LightUploader -c xxx.json -f "Download" -r "Test"

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程
LightUploader -c xxx.json -t 10 -f "Download" 

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程，并跳过同名文件
LightUploader -c xxx.json -t 10 -f "Download" -m 1

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程，同时设置超时时间为30秒
LightUploader -c xxx.json -t 10 -f "Download" -to 30

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程，同时使用 Telegram Bot 实时监控上传进度
LightUploader -c xxx.json -t 10 -f "Download" -tgbot "123456:xxxxxxxx" -uid 123456789

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 10 线程，同时使用配置文件中的 Telegram Bot 参数载入程序实时监控上传进度（前提是配置文件中含有Telegram Bot 的参数）
LightUploader -c xxx.json -t 10 -f "Download" -tgbot "1"

# 将同目录下的 Download 文件夹上传到 OneDrive 网盘根目录中, 使用 15 线程, 并设置分块大小为 20M
LightUploader -c xxx.json -t 15 -b 20 -f "Download" 

```

## 注意
当上传未出现问题，返回0，可作为上传是否失败的凭证