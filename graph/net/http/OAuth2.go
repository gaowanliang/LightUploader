package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type Certificate struct {
	Drive        string      `json:"Drive"`
	RefreshToken string      `json:"RefreshToken"`
	ThreadNum    int         `json:"ThreadNum"`
	BlockSize    int         `json:"BlockSize"`
	MainLand     bool        `json:"MainLand"`
	Language     string      `json:"Language"`
	TimeOut      int         `json:"TimeOut"`
	BotKey       string      `json:"BotKey"`
	UserID       string      `json:"UserID"`
	Other        interface{} `json:"Other"`
}

var GraphURL = "https://graph.microsoft.com/v1.0/me/"
var TokenURL = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
var Host = "login.microsoftonline.com"
var isCN = false

func ChangeCNURL() {
	GraphURL = "https://microsoftgraph.chinacloudapi.cn/v1.0/me/"
	TokenURL = "https://login.chinacloudapi.cn/common/oauth2/v2.0/token"
	Host = "login.chinacloudapi.cn"
	isCN = true
	BaseURL = "https://microsoftgraph.chinacloudapi.cn/v1.0"
}
func NewPassCheck(oauth2URL string, ms int, lang string) string {
	Bearer := getAccessToken(oauth2URL, ms, lang)

	url := GraphURL
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+Bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var mail string
	if ms == 1 {
		mail, err = jsonparser.GetString(body, "userPrincipalName")
	} else {
		mail, err = jsonparser.GetString(body, "mail")
	}

	if err != nil {
		mail, err = jsonparser.GetString(body, "userPrincipalName")
		if err != nil {
			log.Println(string(body))
			log.Panicln(err)
		}
	}
	err = os.Rename("./amazing.json", "./"+mail+".json")
	if err != nil {
		log.Panic(err)
	}

	return "./" + mail + ".json"
}

// GetMyIDAndBearer is get microsoft ID and access Certificate
func GetMyIDAndBearer(infoPath string, Thread int, BlockSize int, Language string, TimeOut int, BotKey string, UserID string) (string, string) {
	MyID := ""
	Bearer := ""
	_, err := os.Stat(infoPath)
	Bearer = refreshAccessToken(infoPath, Thread, BlockSize, Language, TimeOut, BotKey, UserID)
	url := GraphURL
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+Bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	MyID, err = jsonparser.GetString(body, "id")
	if err != nil {
		log.Println(string(body))
		log.Panicln(err)
	}

	//os.Rename("info.json", mail+".json")
	// log.Println(MyID)

	return MyID, Bearer
}

func getAccessToken(oauth2URL string, ms int, lang string) string {
	var re *regexp.Regexp
	if ms == 1 {
		re = regexp.MustCompile(`(?m)code=(.*?)$`)
	} else {
		re = regexp.MustCompile(`(?m)code=(.*?)&`)
	}
	var str = oauth2URL
	/*log.Printf(
		`%s https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=ad5e65fd-856d-4356-aefc-537a9700c137&response_type=code&redirect_uri=http://localhost/onedrive-login&response_mode=query&scope=offline_access%%20User.Read%%20Files.ReadWrite.All`,
		"*请打开下面的网址，登录OneDrive账户后，将跳转后的网址复制后，发送给本Bot*\n注意：本程序不会涉及您的隐私信息，请放心使用，后续会提供更换上传API的方法",
	)*/

	code := re.FindStringSubmatch(str)[1]
	//fmt.Println(code)
	url := TokenURL
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("client_id=ad5e65fd-856d-4356-aefc-537a9700c137&scope=offline_access%%20User.Read%%20Files.ReadWrite.All&code=%s&redirect_uri=http://localhost/onedrive-login&grant_type=authorization_code", code)))
	if isCN {
		req, err = http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("client_id=4fbf37cf-dc83-4b60-b6c1-6230546e247b&code=%s&redirect_uri=http%%3A%%2F%%2Flocalhost%%2Fonedrive-login&grant_type=authorization_code&client_secret=y-L73QIBxO_UmJvOVw8YMlX~8B_h4D6zzT", code)))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Host = Host
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	accessToken, err := jsonparser.GetString(body, "access_token")
	if err != nil {
		log.Println(string(body))
		log.Println(code)
		errorMsg, _ := jsonparser.GetString(body, "error_description")
		log.Panicln(errorMsg)
	}
	//log.Println(accessToken)
	refreshToken, err := jsonparser.GetString(body, "refresh_token")
	if err != nil {
		log.Println(string(body))
		log.Panicln(err)
	}
	//log.Println(refreshToken)

	info := Certificate{
		Drive:        "OneDrive",
		RefreshToken: refreshToken,
		ThreadNum:    3,
		BlockSize:    10,
		MainLand:     isCN,
		Language:     lang,
		TimeOut:      60,
		BotKey:       "",
		UserID:       "",
	}
	// 创建文件
	filePtr, err := os.Create("./amazing.json")
	if err != nil {
		log.Panicln(err.Error())
		return ""
	}
	defer filePtr.Close()
	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(info)
	if err != nil {
		log.Panicln(err.Error())
	}
	return accessToken
}

func refreshAccessToken(path string, Thread int, BlockSize int, Language string, TimeOut int, BotKey string, UserID string) string {
	filePtr, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
		return ""
	}
	defer filePtr.Close()
	var info Certificate
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		log.Panicln(err.Error())
	}
	url := TokenURL
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("client_id=ad5e65fd-856d-4356-aefc-537a9700c137&scope=offline_access%%20User.Read%%20Files.ReadWrite.All&refresh_token=%s&redirect_uri=http://localhost/onedrive-login&grant_type=refresh_token", info.RefreshToken)))
	if isCN {
		req, err = http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("client_id=4fbf37cf-dc83-4b60-b6c1-6230546e247b&scope=offline_access%%20User.Read%%20Files.ReadWrite.All&refresh_token=%s&redirect_uri=http://localhost/onedrive-login&grant_type=refresh_token&client_secret=y-L73QIBxO_UmJvOVw8YMlX~8B_h4D6zzT", info.RefreshToken)))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Host = Host
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	accessToken, err := jsonparser.GetString(body, "access_token")
	if err != nil {
		log.Panicln(err)
	}
	//log.Println(accessToken)
	refreshToken, err := jsonparser.GetString(body, "refresh_token")
	if err != nil {
		log.Panicln(err)
	}
	// log.Println(refreshToken)

	info = Certificate{
		Drive:        "OneDrive",
		RefreshToken: refreshToken,
		ThreadNum:    Thread,
		BlockSize:    BlockSize,
		MainLand:     isCN,
		Language:     Language,
		TimeOut:      TimeOut,
		BotKey:       BotKey,
		UserID:       UserID,
	}
	// 创建文件
	filePtr, err = os.Create(path)
	if err != nil {
		log.Panicln(err.Error())
		return ""
	}
	defer filePtr.Close()
	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(info)
	if err != nil {
		log.Panicln(err.Error())
	}
	return accessToken
}
