package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"main/fileutil"
	httpLocal "main/graph/net/http"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

const (
	simpleUploadPath = "/users/%s/drive/root:/%s:/content"
)

func GetRestoreService(c *http.Client) *RestoreService {
	return &RestoreService{
		httpLocal.NewOneDriveClient(c, false),
	}
}

// RestoreService ItemService manages the communication with Item related API endpoints
type RestoreService struct {
	*httpLocal.OneDrive
}

func (rs *RestoreService) GetDriveItem(userId string, bearerToken string, targetFolder string) (map[string]bool, error) {
	uploadPath := fmt.Sprintf("/users/%s/drive/root:/%s:/children", userId, targetFolder)
	req, err := rs.NewRequest("GET", uploadPath, getSimpleUploadHeader(bearerToken), nil)
	if err != nil {
		return nil, err
	}
	resp, err := rs.Do(req)
	if err != nil {
		//log.Panicln(err)
		return nil, err
	}
	//s, _ := ioutil.ReadAll(resp.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	a := make(map[string]bool, 0)

	jsonparser.ArrayEach(buf.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		i, _ := jsonparser.GetString(value, "name")
		//fmt.Println(i)
		a[i] = true
	}, "value")

	return a, nil
}

// SimpleUploadToOriginalLoc allows you to provide the contents of a new file or update the
// contents of an existing file in a single API call. This method only supports
// files up to 4MB in size. For larger files use ResembleUpload().
//SimpleUploadToOriginalLoc 允许您在单个API调用中提供新文件的内容或更新现有文件的内容。此方法只支持4MB大小的文件。
//对于较大的文件，请使用ResembleUpload()。
//@userId will be extracted as sent from the restore input xml
//@bearerToken will be extracted as sent from the restore input xml
//@filePath will be extracted from the file hierarchy the needs to be restored
//@fileInfo it is the file info struct that contains the actual file reference and the size_type
func (rs *RestoreService) SimpleUploadToOriginalLoc(userId string, bearerToken string, conflictOption string, targetFolder string, filePath string, fileInfo fileutil.FileInfo, sendMsg func(text string), locText func(text string) string, username string) interface{} {
	if fileInfo.SizeType == fileutil.SizeTypeLarge {
		//For Large file type use resemble onedrive upload API
		//log.Printf("Processing Large File: %s", filePath)
		sendMsg(fmt.Sprintf(locText("oneDriveBigFile"), filePath, username))
		return rs.recoverableUpload(userId, bearerToken, conflictOption, targetFolder, filePath, fileInfo, sendMsg, locText, username)
	} else {
		//log.Printf("Processing Small File: %s", filePath)
		sendMsg(fmt.Sprintf(locText("oneDriveSmallFile"), filePath, username))
		targetPath := strings.ReplaceAll(filepath.Join(targetFolder, filePath), "\\", "/")
		startTime := time.Now().Unix()
		_size, err := fileutil.GetFileSize(filePath)
		uploadPath := fmt.Sprintf(simpleUploadPath, userId, targetPath)
		req, err := rs.NewRequest("PUT", uploadPath, getSimpleUploadHeader(bearerToken), fileInfo.FileData)
		if err != nil {
			log.Panicf(locText("failToStore"), err)
		}
		//Handle query parameter for conflict resolution 冲突解决的句柄查询参数
		//The different values for @microsoft.graph.conflictBehavior= rename|replace|fail
		q := url.Values{}
		q.Add("@microsoft.graph.conflictBehavior", conflictOption)
		req.URL.RawQuery = q.Encode()

		//Execute the request
		var resp *http.Response
		for errCount := 1; errCount < 10; errCount++ {
			resp, err = rs.Do(req)
			if err != nil {
				sendMsg(fmt.Sprintf(locText("failToLink"), username, filePath, errCount))
			} else {
				break
			}

		}

		if err != nil {
			sendMsg(fmt.Sprintf(locText("filenameFail"), filePath))
			return nil
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		//Convert to simple map
		respMap := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&respMap)
		if err != nil {
			log.Panicf(locText("failToStore"), err)
		}
		timeUnix := time.Now().UnixNano()
		sendMsg("close=" + fmt.Sprintf(fmt.Sprintf(locText("completeUpload"), filePath, time.Now().Unix()-startTime, byte2Readable(float64(_size)/float64(time.Now().UnixNano()-timeUnix)*float64(1000000000)))))
		return respMap
	}

}

// SimpleUploadToAlternateLoc allows you to provide the contents of a new file or update the
// contents of an existing file in a single API call. This method only supports
// files up to 4MB in size. For larger files use ResumableUpload().
//@userId will be extracted as sent from the restore input xml
//@filePath will be extracted from the file hierarchy the needs to be restored
//@fileInfo it is the file info struct that contains the actual file reference and the size_type
func (rs *RestoreService) SimpleUploadToAlternateLoc(altUserId string, bearerToken string, targetFolder string, conflictOption string, filePath string, fileInfo fileutil.FileInfo, sendMsg func(text string), locText func(text string) string, username string) interface{} {
	if fileInfo.SizeType == fileutil.SizeTypeLarge {
		//For Large file type use resemble onedrive upload API
		return rs.recoverableUpload(altUserId, bearerToken, conflictOption, targetFolder, filePath, fileInfo, sendMsg, locText, username)
	} else {

		uploadPath := fmt.Sprintf(simpleUploadPath, altUserId, filePath)
		req, err := rs.NewRequest("PUT", uploadPath, getSimpleUploadHeader(bearerToken), fileInfo.FileData)
		if err != nil {
			log.Panicf("Failed to Restore :%v", err)
		}
		//Handle query parameter for conflict resolution
		//The different values for @microsoft.graph.conflictBehavior= rename|replace|fail
		q := url.Values{}
		q.Add("@microsoft.graph.conflictBehavior", conflictOption)
		req.URL.RawQuery = q.Encode()

		//Execute the request
		resp, err := rs.Do(req)
		if err != nil {
			log.Panicf("Failed to Restore :%v", err)
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		//Convert to simple map
		respMap := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&respMap)
		if err != nil {
			log.Panicf("Failed to Restore :%v", err)
		}
		return respMap
	}
}

//Get response as string
func readRespAsString(resp *http.Response) string {
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return bodyString
	}
	return ""
}

func getSimpleUploadHeader(accessToken string) map[string]string {
	//As a work around for now, ultimately this will be recived as a part of restore xml
	bearerToken := fmt.Sprintf("bearer %s", accessToken)
	return map[string]string{
		"Content-Type":  "application/octet-stream",
		"Authorization": bearerToken,
	}
}
