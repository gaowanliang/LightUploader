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

func ApplyForNewPass(url string, ms int) string {
	return httpLocal.NewPassCheck(url, ms)
}

func Upload(infoPath string, filePath string, threads int, sendMsg func() (func(text string), string, string), locText func(text string) string) {

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
	restoreSrvc := upload.GetRestoreService(http.DefaultClient)

	//Get the list of files that needs to be restore with the actual backed up path. 获取需要使用实际备份路径还原的文件列表。
	fileInfoToUpload, err := fileutil.GetAllUploadItemsFrmSource(filePath)
	if err != nil {
		log.Fatalf("Failed to Load Files from source :%v", err)
	}

	//Call restore process based on alternate or original location 基于备用或原始位置调用还原过程
	/*if restoreOption == "alt" {
		restoreToAltLoc(restoreSrvc, fileInfoToUpload)
	} else {
		restore(restoreSrvc, fileInfoToUpload, threads)
	}*/
	restore(restoreSrvc, fileInfoToUpload, threads, sendMsg, locText, infoPath)
	err = os.Chdir(oldDir)
	if err != nil {
		log.Panic(err)
	}
}
func changeBlockSize(MB int) {
	fileutil.SetDefaultChunkSize(MB)
}

//Restore to original location
func restore(restoreSrvc *upload.RestoreService, filesToRestore map[string]fileutil.FileInfo, threads int, sendMsg func() (func(text string), string, string), locText func(text string) string, infoPath string) {
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
			tip := "`" + filePath + "`" + "开始上传至OneDrive"
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
			restoreSrvc.SimpleUploadToOriginalLoc(userID, bearerToken, "rename", filePath, fileInfo, temp, locText, username)

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
func restoreToAltLoc(restoreSrvc *upload.RestoreService, filesToRestore map[string]fileutil.FileInfo, sendMsg func() func(text string), locText func(text string) string, infoPath string) {
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
			temp(filePath + "开始上传至OneDrive")
			us := ""
			userID, bearerToken := httpLocal.GetMyIDAndBearer(infoPath)
			//username := strings.ReplaceAll(filepath.Base(infoPath), ".json", "")
			restoreSrvc.SimpleUploadToAlternateLoc(userID, bearerToken, "rename", rootFilePath, fileItem, temp, locText, us)

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
		log.Panicf("Telegram Send Error:%s", description)
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
				log.Panicf("Telegram Send Error:%s", description)
			}
		}

	}
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
	// StringVar用指定的名称、控制台参数项目、默认值、使用信息注册一个string类型flag，并将flag的值保存到p指向的变量
	flag.StringVar(&codeURL, "a", "", "通过登录 https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=ad5e65fd-856d-4356-aefc-537a9700c137&response_type=code&redirect_uri=http://localhost/onedrive-login&response_mode=query&scope=offline_access%20User.Read%20Files.ReadWrite.All 后跳转的网址(输入网址时需要使用双引号包裹网址)")
	flag.IntVar(&ms, "v", 0, "选择版本，其中0为国际版，1为个人版(家庭版)，默认为0")
	flag.StringVar(&configFile, "c", "", "登录文件位置")
	flag.StringVar(&folder, "f", "", "欲上传的文件/文件夹")
	flag.IntVar(&thread, "t", 3, "线程数，默认为3")
	flag.IntVar(&block, "b", 10, "自定义上传分块大小, 可以提高网络吞吐量, 受限于磁盘性能和网络速度，默认为 10 (单位MB)")
	flag.StringVar(&botKey, "tgbot", "", "使用Telegram机器人实时监控上传，此处需填写机器人的access token，形如123456789:xxxxxxxxx，输入时需使用双引号包裹")
	flag.StringVar(&iuserID, "uid", "", "使用Telegram机器人实时监控上传，此处需填写接收人的userID，形如123456789")

	// 从arguments中解析注册的flag。必须在所有flag都注册好而未访问其值时执行。未注册却使用flag -help时，会返回ErrHelp。

	flag.Parse()
	if configFile != "" && folder != "" {
		fileutil.SetDefaultChunkSize(block)
		startTime := time.Now().Unix()
		writer := uilive.New()
		writer.Start()
		var sendMsg func(string)
		if botKey != "" && iuserID != "" {
			sendMsg = botSend(botKey, iuserID, fmt.Sprintf("`%s` 开始上传", folder))
		}

		Upload(strings.ReplaceAll(configFile, "\\", "/"), strings.ReplaceAll(folder, "\\", "/"), thread, func() (func(text string), string, string) {
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
			return oldFunc(text)
		})
		_, _ = fmt.Fprintf(writer, "%s上传完成，耗时 %d 秒\n", folder, time.Now().Unix()-startTime)
		if botKey != "" && iuserID != "" {
			sendMsg(fmt.Sprintf("`%s`上传完成，耗时 %d 秒\n", folder, time.Now().Unix()-startTime))
		}
	} else {
		if codeURL == "" {
			flag.PrintDefaults()
		} else {
			log.Printf("注册成功，已在运行目录下新建登录文件%s", ApplyForNewPass(codeURL, ms))
		}
	}

	// 打印
	//fmt.Printf("username=%v password=%v host=%v port=%v", username, password, host, port)
}

/*
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
*/
