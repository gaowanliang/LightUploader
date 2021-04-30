package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.SimplifiedChinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	dir, err := os.Executable()
	dropErr(err)
	dir = filepath.Dir(dir)
	// fmt.Println(dir)
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}
	for _, fi := range rd {
		if !fi.IsDir() && path.Ext(fi.Name()) == ".toml" {
			bundle.LoadMessageFile(path.Join(dir, fi.Name()))
		}
	}

}

func dropErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func pageDownload(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return ""
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error", err)
		return ""
	}
	return string(body)
}

func GetFileModTime(path string) int64 {
	f, err := os.Open(path)
	dropErr(err)
	defer f.Close()

	fi, err := f.Stat()
	dropErr(err)

	return fi.ModTime().Unix()
}

type Loc struct {
	localize *i18n.Localizer
}

func (loc *Loc) init(locLanguage string) {
	dir, err := os.Executable()
	dropErr(err)
	dir = filepath.Dir(dir)
	//fmt.Println(dir)
	dir = filepath.Join(dir, fmt.Sprintf("%s.toml", locLanguage))
	_, err = os.Stat(dir)
	if err != nil {
		resp, err := http.Get(fmt.Sprintf("https://cdn.jsdelivr.net/gh/gaowanliang/OneDriveUploader/i18n/%s.toml", locLanguage))
		dropErr(err)
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		dropErr(err)
		ioutil.WriteFile(dir, data, 0666)
	} else {
		url := "https://cdn.jsdelivr.net/gh/gaowanliang/OneDriveUploader@latest/i18n/"
		j := pageDownload(url)
		var re = regexp.MustCompile(`(?m)i18n/(.*?)"[\s\S]*?<td class="time">(.*?)</td>`)
		var newLanFileTime int64
		for _, val := range re.FindAllStringSubmatch(j, -1) {
			if fmt.Sprintf("%s.toml", locLanguage) == val[1] {
				t, _ := time.Parse(time.RFC1123, val[2])
				newLanFileTime = t.Unix()
			}

		}
		oldLanFileTime := GetFileModTime(dir)
		if newLanFileTime > oldLanFileTime {
			err = os.RemoveAll(dir)
			dropErr(err)
			resp, err := http.Get(fmt.Sprintf("https://cdn.jsdelivr.net/gh/gaowanliang/OneDriveUploader/i18n/%s.toml", locLanguage))
			dropErr(err)
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			dropErr(err)
			ioutil.WriteFile(dir, data, 0644)
		}
	}
	dir, err = os.Executable()
	// log.Println(dir)
	dropErr(err)
	dir = filepath.Dir(dir)
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}
	for _, fi := range rd {
		if !fi.IsDir() && path.Ext(fi.Name()) == ".toml" {
			// log.Println(path.Join(dir, fi.Name()))
			bundle.LoadMessageFile(path.Join(dir, fi.Name()))
		}
	}
	loc.localize = i18n.NewLocalizer(bundle, locLanguage)
}
func (loc *Loc) print(tag string) string {
	return loc.localize.MustLocalize(&i18n.LocalizeConfig{MessageID: tag})
}
