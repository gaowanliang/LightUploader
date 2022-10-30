package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"main/api/restore/upload"
	"main/fileutil"
	"main/googledrive"
	httpLocal "main/graph/net/http"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gosuri/uilive"
)

var loc Loc

func ApplyForNewPass(url string, ms int) string {
	if ms == 2 {
		httpLocal.ChangeCNURL()
	}
	return httpLocal.NewPassCheck(url, ms, lang)
}

func Upload(infoPath string, filePath string, targetFolder string, threads int, sendMsg func() (func(text string), string, string), locText func(text string) string) {

	programPath, err := filepath.Abs(filepath.Dir(infoPath))
	if err != nil {
		log.Panic(err)
	}
	infoPath = filepath.Base(infoPath)
	infoPath = filepath.Join(programPath, infoPath)

	// restoreOption := "orig"
	pathLastChar := filePath[len(filePath)-1]
	if pathLastChar == '/' || pathLastChar == '\\' { // 当最后一位是/时，可能会出现找不到文件的情况，这里进行首先处理
		filePath = filePath[:len(filePath)-1]
	}

	oldDir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	err = os.Chdir(filepath.Dir(filePath))
	if err != nil {
		log.Panic(err)
	}
	filePath = path.Base(filePath)
	//Initialize the upload restore service
	restoreSrvc := upload.GetRestoreService(&http.Client{Timeout: time.Duration(fileutil.GetTimeOut()) * time.Second})

	//Get the list of files that needs to be restore with the actual backed up path. 获取需要使用实际备份路径还原的文件列表。
	fileInfoToUpload, err := fileutil.GetAllUploadItemsFrmSource(filePath)
	if err != nil {
		log.Fatalf(loc.print("failToLoadFiles"), err)
	}

	//Call restore process based on alternate or original location 基于备用或原始位置调用还原过程
	/*if restoreOption == "alt" {
		restoreToAltLoc(restoreSrvc, fileInfoToUpload)
	} else {
		restore(restoreSrvc, fileInfoToUpload, threads)
	}*/

	restore(restoreSrvc, fileInfoToUpload, targetFolder, threads, sendMsg, locText, infoPath)
	err = os.Chdir(oldDir)
	if err != nil {
		log.Panic(err)
	}
}
func changeBlockSize(MB int) {
	fileutil.SetDefaultChunkSize(MB)
}

//Restore to original location
func restore(restoreSrvc *upload.RestoreService, filesToRestore map[string]fileutil.FileInfo, targetFolder string, threads int, sendMsg func() (func(text string), string, string), locText func(text string) string, infoPath string) {
	var wg sync.WaitGroup
	pool := make(chan struct{}, threads)
	checkPath := make(map[string]bool, 0)
	pathFiles := make(map[string]map[string]bool, 0)

	filePaths := make([]string, 0, len(filesToRestore))
	for k := range filesToRestore {
		filePaths = append(filePaths, k)
	}
	sort.Sort(sort.StringSlice(filePaths)) //将文件按照文件名顺序排序
	for _, filePath := range filePaths {
		wg.Add(1)
		pool <- struct{}{}
		fileInfo := filesToRestore[filePath]
		paths, fileName := filepath.Split(filepath.Join(targetFolder, filePath))
		if mode == 1 {
			if paths == "" {
				paths = "/"
			}
			paths = strings.ReplaceAll(paths, "\\", "/")
			if paths[len(paths)-1] == '/' {
				paths = paths[:len(paths)-1]
			}
			if _, ok := checkPath[paths]; !ok {
				userID, bearerToken := httpLocal.GetMyIDAndBearer(infoPath, thread, block, lang, timeOut, botKey, _UserID)
				files, _ := restoreSrvc.GetDriveItem(userID, bearerToken, paths)
				checkPath[paths] = true
				pathFiles[paths] = files
			}
			// log.Println(checkPath, paths, pathFiles, fileName, filePath)
		}

		go func(filePath string, fileInfo fileutil.FileInfo) {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			temps, _botKey, iUserID := sendMsg()
			var iSendMsg func(string)
			tip := "`" + filePath + "`" + loc.print("startToUpload1")
			if _botKey != "" && iUserID != "" {
				iSendMsg = botSend(_botKey, iUserID, tip)
			}
			temp := func(text string) {
				temps(text)
				if _botKey != "" && iUserID != "" {
					iSendMsg(text)
				}
			}
			if _, ok := pathFiles[paths][fileName]; !ok || mode == 0 {
				temp(tip)
				userID, bearerToken := httpLocal.GetMyIDAndBearer(infoPath, thread, block, lang, timeOut, _botKey, _UserID)
				username := strings.ReplaceAll(filepath.Base(infoPath), ".json", "")
				restoreSrvc.SimpleUploadToOriginalLoc(userID, bearerToken, "replace", targetFolder, filePath, fileInfo, temp, locText, username)
			} else {
				tip = filePath + "已存在，自动跳过"
				if _botKey != "" && iUserID != "" {
					iSendMsg = botSend(_botKey, iUserID, tip)
				}
				temp(tip)
				time.Sleep(time.Second * 3)
			}
			defer temp("close=")
		}(filePath, fileInfo)
	}
	wg.Wait()

}

func printResp(resp interface{}) {
	switch resp.(type) {
	case map[string]interface{}:
		fmt.Printf("\n%+v\n", resp)
		break
	case []map[string]interface{}:
		for _, rs := range resp.([]map[string]interface{}) {
			fmt.Printf("\n%+v\n", rs)
		}
	}
}

//Restore to Alternate location 还原到备用位置
func restoreToAltLoc(restoreSrvc *upload.RestoreService, filesToRestore map[string]fileutil.FileInfo, targetFolder string, sendMsg func() func(text string), locText func(text string) string, infoPath string) {
	rootFolder := fileutil.GetAlternateRootFolder()
	var wg sync.WaitGroup
	pool := make(chan struct{}, 10)
	for filePath, fileItem := range filesToRestore {
		rootFilePath := fmt.Sprintf("%s/%s", rootFolder, filePath)
		wg.Add(1)
		pool <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			temp := sendMsg()
			temp(filePath + loc.print("startToUpload1"))
			us := ""
			userID, bearerToken := httpLocal.GetMyIDAndBearer(infoPath, thread, block, lang, timeOut, botKey, _UserID)
			//username := strings.ReplaceAll(filepath.Base(infoPath), ".json", "")
			restoreSrvc.SimpleUploadToAlternateLoc(userID, bearerToken, "rename", targetFolder, rootFilePath, fileItem, temp, locText, us)

		}()
		wg.Wait()
		// fmt.Println(respStr)
	}
}

func oldFunc(a string) string {
	return a
}

func botSend(botKey string, iuserID string, initText string) func(string) {
	var messageId = int64(0)
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&parse_mode=MarkdownV2&text=%s", botKey, iuserID, url.QueryEscape(initText)))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(string(body))
	ok, _ := jsonparser.GetBoolean(body, "ok")
	if ok {
		messageId, _ = jsonparser.GetInt(body, "result", "message_id")
	} else {
		description, _ := jsonparser.GetString(body, "description")
		log.Println(loc.print("telegramSendError"), description)
	}
	return func(text string) {
		if text[:5] == "close" && text[5] != '|' {
			// msg 头部的 close 用作输出时定位，带有 close 在输出时不会被刷新走
			// close= 表示文件传输结束，此时会同步删除tg发出的消息
			// close| 则不会删除消息
			resp, err = http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/deleteMessage?chat_id=%s&message_id=%d", botKey, iuserID, messageId))
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()
			return
		}
		resp, err = http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText?chat_id=%s&parse_mode=MarkdownV2&message_id=%d&text=%s", botKey, iuserID, messageId, url.QueryEscape(text)))
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		//fmt.Println(string(body))
		ok, _ = jsonparser.GetBoolean(body, "ok")
		if !ok {
			description, _ := jsonparser.GetString(body, "description")
			if !strings.Contains(string(body), "message is not modified") && !strings.Contains(string(body), "Too Many Requests") {
				log.Println(loc.print("telegramSendError"), description)
			}
		}

	}
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {

		if info != nil && !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

var timeOut int
var lang string
var block int
var botKey string
var _UserID string
var thread int
var mode int

func main() {
	var codeURL string
	var configFile string
	var folder string
	var ms int
	var targetFolder string

	// StringVar用指定的名称、控制台参数项目、默认值、使用信息注册一个string类型flag，并将flag的值保存到p指向的变量

	flag.StringVar(&codeURL, "a", "", "Please refer to \"https://github.com/gaowanliang/LightUploader/wiki\"")
	flag.IntVar(&ms, "v", 0, "Select the version, where 0 is the business version, 1 is the personal version (home version), 2 is 21vianet (CN) version, 3 is Google Drive. The default is 0")
	flag.StringVar(&configFile, "c", "", "Authorize config file location")
	flag.StringVar(&folder, "f", "", "Files / folders to upload")
	flag.IntVar(&thread, "t", 3, "The number of threads")
	flag.IntVar(&block, "b", 10, "User defined upload block size can improve network throughput. Limited by disk performance and network speed, the default is 10 (unit: MB)")
	flag.StringVar(&botKey, "tgbot", "", "Use the telegram robot to monitor the upload in real time. Here you need to fill in the robot's access token, such as 123456789:XXXXXXXX, and use double quotation marks when entering")
	flag.StringVar(&_UserID, "uid", "", "Use the telegram robot to monitor the upload in real time. Fill in the user ID of the receiver, such as 123456789")
	flag.StringVar(&targetFolder, "r", "", "Set the directory you want to upload to onedrive")
	flag.IntVar(&timeOut, "to", 60, "When uploading, the timeout of each block is 60s by default")
	flag.IntVar(&mode, "m", 0, "Select the mode, 0 is to replace the file with the same name in onedrive, 1 is to skip, the default is 0")
	flag.StringVar(&lang, "l", "en", "Set the software language, English by default")
	// 从arguments中解析注册的flag。必须在所有flag都注册好而未访问其值时执行。未注册却使用flag -help时，会返回ErrHelp。
	flag.Parse()

	loc = Loc{}
	if configFile != "" && folder != "" {
		filePtr, err := os.Open(configFile)
		if err != nil {
			log.Panicln(err)
		}
		defer filePtr.Close()
		var info httpLocal.Certificate
		// 创建json解码器
		decoder := json.NewDecoder(filePtr)
		err = decoder.Decode(&info)
		if err != nil {
			log.Panicln(err.Error())
		}

		if info.MainLand {
			httpLocal.ChangeCNURL()
		}

		if info.Language != "en" && lang == "en" {
			loc.init(info.Language)
			lang = info.Language
		} else {
			loc.init(lang)
		}

		if info.BlockSize != 10 && block == 10 {
			fileutil.SetDefaultChunkSize(info.BlockSize)
		} else {
			fileutil.SetDefaultChunkSize(block)
		}

		if info.TimeOut != 60 && timeOut == 60 {
			fileutil.SetTimeOut(info.TimeOut)
			timeOut = info.TimeOut
		} else {
			fileutil.SetTimeOut(timeOut)
		}

		if info.BotKey != "" && info.UserID != "" && botKey == "1" {
			botKey = info.BotKey
			_UserID = info.UserID
		}

		startTime := time.Now().Unix()
		writer := uilive.New()
		writer.Start()
		size, err := DirSize(folder)
		if err != nil {
			log.Panic(err)
		}

		_, _ = fmt.Fprintf(writer, loc.print("startToUpload"), folder, fileutil.Byte2Readable(float64(size)))

		updateOutput := func(text string) {
			if text[:5] != "close" {
				_, _ = fmt.Fprintf(writer.Newline(), "%s\n", text)
			} else {
				_, _ = fmt.Fprintf(writer.Bypass(), "%s\n", text[6:])
			}
		}

		var sendMsg func(string)
		if botKey != "" && _UserID != "" {
			sendMsg = botSend(botKey, _UserID, fmt.Sprintf(loc.print("startToUpload"), folder, fileutil.Byte2Readable(float64(size))))
		}
		switch info.Drive {
		case "OneDrive":
			Upload(strings.ReplaceAll(configFile, "\\", "/"), strings.ReplaceAll(folder, "\\", "/"), targetFolder, thread, func() (func(text string), string, string) {
				if botKey != "" && _UserID != "" {
					return updateOutput, botKey, _UserID
				} else {
					return updateOutput, "", ""
				}
			}, func(text string) string {
				return loc.print(text)
			})
		case "GoogleDrive":
			googledrive.Upload(strings.ReplaceAll(configFile, "\\", "/"), strings.ReplaceAll(folder, "\\", "/"), func() (func(text string), string, string, func(string, string, string) func(string)) {
				if botKey != "" && _UserID != "" {
					return updateOutput, botKey, _UserID, botSend
				} else {
					return updateOutput, "", "", nil
				}
			}, func(text string) string {
				return loc.print(text)
			}, thread, block, lang, timeOut, botKey, _UserID)

		}

		cost := time.Now().Unix() - startTime
		speed := fileutil.Byte2Readable(float64(size) / float64(cost))
		_, _ = fmt.Fprintf(writer.Bypass(), loc.print("completeUpload"), folder, cost, speed)
		if botKey != "" && _UserID != "" {
			log.Printf(fmt.Sprintf(loc.print("completeUpload"), folder, cost, speed))
			sendMsg(fmt.Sprintf(loc.print("completeUpload"), folder, cost, speed))
		}
	} else {
		loc.init(lang)
		if codeURL == "" {
			if ms != 3 {
				flag.PrintDefaults()
			} else {
				log.Printf(loc.print("googleDriveGetAccess"), googledrive.GetURL())
				inputReader := bufio.NewReader(os.Stdin)
				code, err := inputReader.ReadString('\n')
				if err != nil {
					fmt.Println("There ware errors reading,exiting program.")
					return
				}
				mail := googledrive.CreateNewInfo(code, lang)
				log.Println(loc.print("googleDriveOAuthFileCreateSuccess") + mail)
			}

		} else {
			log.Printf(loc.print("configCreateSuccess"), ApplyForNewPass(codeURL, ms))
		}
	}

	// 打印
	//fmt.Printf("username=%v password=%v host=%v port=%v", username, password, host, port)
}

/*
SET CGO_ENABLED=0
$env:GOOS="linux"
SET GOARCH=amd64
go build -o LightUploader .
/usr/local/bin

set HTTPS_PROXY=http://127.0.0.1:2334
*/
