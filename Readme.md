[简体中文](Readme-zh-CN.md)
# LightUploader

MoeClub wrote a [very good version](https://github.com/MoeClub/OneList/tree/master/OneDriveUploader), but unfortunately it's not open source and hasn't been updated in a while. This project is a simple upload tool separate from [DownloadBot](https://github.com/gaowanliang/DownloadBot), designed to be a lightweight way to quickly upload data to various network drives on all platforms.

## Features

- Supports OneDrive Business, Personal (Home) versions, 21vianet (CN) version, Google Drive.
- Support for uploading files and folders to specified directories, keeping the directory structure as it was before the upload.
- Supports the use of command parameters for external applications.
- Support for customising the upload chunk size.
- Supports multi-threaded uploads (multiple files at the same time).
- Support for dynamically adjusting the number of retries according to the file size.
- Supports skipping the existing files with the same name in the OneDrive.
- Support for real-time monitoring of upload progress via Telegram Bot, for easy monitoring of uploads when using fully automated download scripts.


## Authorize
See [wiki](https://github.com/gaowanliang/LightUploader/wiki) for details



### Initialization profile
```bash
# OneDrive Business
LightUploader -a "url"

# OneDrive Business, and use Chinese language pack
LightUploader -a "url" -l zh-CN

# OneDrive Personal (Home)
LightUploader -a "url" -v 1

# OneDrive 21vianet (CN) version, and use Chinese language pack
LightUploader -a "url" -v 2 -l zh-CN

# Get the entire url in the browser address bar starting with http://loaclhost
# Replace the full url with the three letters of the "url" in the command
# Each url generated can only be used once, try again to retrieve the url
# This action will automatically initialise the configuration file

# Google Drive
LightUploader -v 3

```

## Use
```c
Usage of LightUploader:
  -a string
        //Setup and Init auth.json.
  -b string
        //Set block size. [Unit: M; 5<=b<=60;] (default "10")
  -c string
        //Config file.

  -r string
        //Upload to reomte path.
  -l string
        // Software language
  -f string
        // *Necessary parameters, file or folder to upload
  -t string
        // Number of threads, number of files uploaded at the same time. Default: 3
  -to int
        // The timeout time of a single packet is 60s by default
  -tgbot string
        //Use the Telegram bot to monitor uploads in real time, here you need to fill in the access token of the bot, e.g. 123456789:xxxxxxxxxx, use double quotes to wrap it
  -uid string
        // Use the Telegram bot to monitor uploads in real time, here you need to fill in the recipient's userID, shaped like 123456789
  -m int
        // Select the mode, 0 is to replace the file with the same name in cloud drive, 1 is to skip, the default is 0
  -v int
        // Select the version, where 0 is the OneDrive Business version and 1 is the OneDrive Personal (Home) version, 2 is OneDrive 21vianet (CN) version, 3 is Google Drive, the default is 0
```

## Config
```jsonc
{
    "Drive":"OneDrive",
    // Authorisation tokens
    "RefreshToken": "1234564567890ABCDEF",
    // Maximum number of threads. (Number of simultaneous file uploads)
    "ThreadNum": 2,
    // Maximum upload chunk size. (The maximum chunk size for each file upload, it is recommended to reduce it if the network is not good. Unit: MB)
    "BlockSize": 10,
    // Maximum single file size. (Currently: 100GB single file limit for Personal Edition (Home Edition); 15GB single file limit for other editions, Microsoft will update to 100GB over time. Unit: GB)
    "SigleFile": 100,
    // If this is the Chinese version (CenturyLink), this should be true.
    "MainLand": false,
    //Software language
    "Language": "zh-CN",
    //timeout
    "TimeOut": 60,
    //Telegram Bot key
    "BotKey": "",
    //Telegram User ID
    "UserID": ""
}
```
Note that when a configuration file is used at the same time and the parameters in the configuration file are not the default values, the parameters in the configuration file will be used by default, and when parameters are entered on the command line, the parameters currently entered on the command line will be used and the parameters in the configuration file will be changed. That is, the parameters in the command line take precedence over the parameters in the configuration file and will change the parameters in the configuration file.

## Examples
```bash
# Some examples:

# Upload the mm00.jpg file from the same directory to the root of the OneDrive
LightUploader -c xxx.json -f "mm00.jpg"

# Upload the Download folder from the same directory to the root of the OneDrive
LightUploader -c xxx.json -f "Download" 

# Upload the Download folder from the same directory to the Test directory of the OneDrive
LightUploader -c xxx.json -f "Download" -r "Test"

# Upload the Download folder from the same directory to the root of the OneDrive, using 10 threads
LightUploader -c xxx.json -t 10 -f "Download" 

# Upload the download folder in the same directory to the root directory of onedrive, use 10 threads, and skip the file with the same name
LightUploader -c xxx.json -t 10 -f "Download" -m 1

# Upload the download folder in the same directory to the root directory of onedrive, use 10 threads, and set the timeout to 30 seconds
LightUploader -c xxx.json -t 10 -f "Download" -to 30

# Upload the Download folder from the same directory to the root of the OneDrive, using 10 threads, while using Telegram Bot to monitor the progress of the upload in real time
LightUploader -c xxx.json -t 10 -f "Download" -tgbot "123456:xxxxxxxx" -uid 123456789

# Upload the download folder in the same directory to the root directory of onedrive network disk, use 10 threads, and use the telegram BOT parameter loader in the configuration file to monitor the upload progress in real time (provided that the configuration file contains the parameters of telegram BOT)
LightUploader -c xxx.json -t 10 -f "Download" -tgbot "1"


# Upload the Download folder from the same directory to the root of the OneDrive, using 15 threads, and setting the chunk size to 20M
LightUploader -c xxx.json -t 15 -b 20 -f "Download" 
```

## Note

Returns 0 when there is no problem with the upload, which can be used as evidence of whether the upload has failed or not