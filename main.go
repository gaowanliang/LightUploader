package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"main/api/restore/upload"
	"main/fileutil"
	httpLocal "main/graph/net/http"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gosuri/uilive"
)

var loc Loc

func ApplyForNewPass(url string, ms int) string {
	return httpLocal.NewPassCheck(url, ms)
}

func Upload(infoPath string, filePath string, targetFolder string, threads int, sendMsg func() (func(text string), string, string), locText func(text string) string) {

	programPath, err := filepath.Abs(filepath.Dir(infoPath))
	if err != nil {
		log.Panic(err)
	}
	infoPath = filepath.Base(infoPath)
	infoPath = filepath.Join(programPath, infoPath)

	// restoreOption := "orig"
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
	for filePath, fileInfo := range filesToRestore {
		wg.Add(1)
		pool <- struct{}{}
		go func(filePath string, fileInfo fileutil.FileInfo) {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			temps, botKey, iUserID := sendMsg()
			var iSendMsg func(string)
			tip := "`" + filePath + "`" + loc.print("startToUpload1")
			if botKey != "" && iUserID != "" {
				iSendMsg = botSend(botKey, iUserID, tip)
			}
			temp := func(text string) {
				temps(text)
				if botKey != "" && iUserID != "" {
					iSendMsg(text)
				}
			}
			temp(tip)
			userID, bearerToken := httpLocal.GetMyIDAndBearer(infoPath)
			username := strings.ReplaceAll(filepath.Base(infoPath), ".json", "")
			restoreSrvc.SimpleUploadToOriginalLoc(userID, bearerToken, "rename", targetFolder, filePath, fileInfo, temp, locText, username)

			//printResp(resp)
			defer temp("close")
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
			userID, bearerToken := httpLocal.GetMyIDAndBearer(infoPath)
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
	var message_id = int64(0)
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&parse_mode=MarkdownV2&text=%s", botKey, iuserID, url.QueryEscape(initText)))
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(string(body))
	ok, _ := jsonparser.GetBoolean(body, "ok")
	if ok {
		message_id, _ = jsonparser.GetInt(body, "result", "message_id")
	} else {
		description, _ := jsonparser.GetString(body, "description")
		log.Panicf(loc.print("telegramSendError"), description)
	}
	return func(text string) {
		if text == "close" {
			resp, err = http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/deleteMessage?chat_id=%s&message_id=%d", botKey, iuserID, message_id))
			if err != nil {
				log.Panic(err)
			}
			defer resp.Body.Close()
			return
		}
		resp, err = http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText?chat_id=%s&parse_mode=MarkdownV2&message_id=%d&text=%s", botKey, iuserID, message_id, url.QueryEscape(text)))
		if err != nil {
			log.Panic(err)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
		}
		//fmt.Println(string(body))
		ok, _ = jsonparser.GetBoolean(body, "ok")
		if !ok {
			description, _ := jsonparser.GetString(body, "description")
			if !strings.Contains(string(body), "message is not modified") {
				log.Panicf(loc.print("telegramSendError"), description)
			}
		}

	}
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func main() {
	var codeURL string
	var configFile string
	var folder string
	var thread int
	var botKey string
	var iuserID string
	var ms int
	var block int
	var targetFolder string
	var timeOut int
	var lang string
	// StringVar用指定的名称、控制台参数项目、默认值、使用信息注册一个string类型flag，并将flag的值保存到p指向的变量

	flag.StringVar(&codeURL, "a", "", "Jump to the https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=ad5e65fd-856d-4356-aefc-537a9700c137&response_type=code&redirect_uri=http://localhost/onedrive-login&response_mode=query&scope=offline_access%20User.Read%20Files.ReadWrite.All after logging in to the website (you need to use double quotation marks to wrap the URL when entering the URL)")
	flag.IntVar(&ms, "v", 0, "Select the version, where 0 is the business version and 1 is the personal version (home version). The default is 0")
	flag.StringVar(&configFile, "c", "", "Login config file location")
	flag.StringVar(&folder, "f", "", "Files / folders to upload")
	flag.IntVar(&thread, "t", 3, "The number of threads")
	flag.IntVar(&block, "b", 10, "User defined upload block size can improve network throughput. Limited by disk performance and network speed, the default is 10 (unit: MB)")
	flag.StringVar(&botKey, "tgbot", "", "Use the telegram robot to monitor the upload in real time. Here you need to fill in the robot's access token, such as 123456789:XXXXXXXX, and use double quotation marks when entering")
	flag.StringVar(&iuserID, "uid", "", "Use the telegram robot to monitor the upload in real time. Fill in the user ID of the receiver, such as 123456789")
	flag.StringVar(&targetFolder, "r", "", "Set the directory you want to upload to onedrive")
	flag.IntVar(&timeOut, "to", 60, "When uploading, the timeout of each block is 60s by default")
	flag.StringVar(&lang, "l", "en", "Set the software language, English by default")
	// 从arguments中解析注册的flag。必须在所有flag都注册好而未访问其值时执行。未注册却使用flag -help时，会返回ErrHelp。
	flag.Parse()
	loc = Loc{}
	loc.init(lang)

	if configFile != "" && folder != "" {
		fileutil.SetDefaultChunkSize(block)
		fileutil.SetTimeOut(timeOut)
		startTime := time.Now().Unix()
		writer := uilive.New()
		writer.Start()
		size, err := DirSize(folder)
		if err != nil {
			log.Panic(err)
		}

		_, _ = fmt.Fprintf(writer, loc.print("startToUpload"), folder, fileutil.Byte2Readable(float64(size)))

		var sendMsg func(string)
		if botKey != "" && iuserID != "" {
			sendMsg = botSend(botKey, iuserID, fmt.Sprintf(loc.print("startToUpload"), folder, fileutil.Byte2Readable(float64(size))))
		}

		Upload(strings.ReplaceAll(configFile, "\\", "/"), strings.ReplaceAll(folder, "\\", "/"), targetFolder, thread, func() (func(text string), string, string) {
			if botKey != "" && iuserID != "" {
				return func(text string) {
					_, _ = fmt.Fprintf(writer, "%s\n", text)
				}, botKey, iuserID
			} else {
				return func(text string) {
					_, _ = fmt.Fprintf(writer, "%s\n", text)
				}, "", ""
			}

		}, func(text string) string {
			return loc.print(text)
		})
		cost := time.Now().Unix() - startTime
		speed := fileutil.Byte2Readable(float64(size) / float64(cost))
		_, _ = fmt.Fprintf(writer, loc.print("completeUpload"), folder, cost, speed)
		if botKey != "" && iuserID != "" {
			sendMsg(fmt.Sprintf(loc.print("completeUpload"), folder, cost, speed))
		}
	} else {
		if codeURL == "" {
			flag.PrintDefaults()
		} else {
			log.Printf(loc.print("configCreateSuccess"), ApplyForNewPass(codeURL, ms))
		}
	}

	// 打印
	//fmt.Printf("username=%v password=%v host=%v port=%v", username, password, host, port)
}

/*
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o OneDriveUploader .
/usr/local/bin
*/
