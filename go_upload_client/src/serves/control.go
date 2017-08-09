package serves

import (
	_ "encoding/json"
	"fmt"
	_ "fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	_ "strconv"
	"strings"
	"time"
	_ "uploadTest/src/common"
	"uploadTest/src/upload"

	//"github.com/khlipeng/websocket"
)

var (
	tokenMap map[string]map[string]*upload.TUploadTask
	routers  *TRouters
)

type tokenTaskRequest struct {
	FilePath string `json:"filePath"`
	Token    string `json:"token"`
}

type tokenTaskReturn struct {
	ID           int    `json:"task_id"`
	FileSize     string `json:"size"`
	TransferSize string `json:"transferSize"`
	Rate         string `json:"rate"`
	FileName     string `json:"name"`
	UsedTime     string `json:"time"`
	State        string `json:"state"`
	Progress     string `json:"progress"`
}

type TRouters struct {
	//	handleSoc websocket.Handler
}

func Run() {
	routers = &TRouters{
	//		handleSoc: websocket.Handler(SocketLoop),
	}
	srv := &http.Server{
		Handler: routers,
		Addr:    ":9090",
	}
	tokenMap = make(map[string]map[string]*upload.TUploadTask)
	err := srv.ListenAndServe() //设置监听的端口

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PathSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PathSep+fi.Name())
		}
	}

	return files, nil
}

func uploadFiles(filePath string, taskMap map[string]*upload.TUploadTask) {
	files, err := ListDir(filePath, "mp4")
	if err != nil {
		log.Print(err)
		return
	}
	//PathSep := string(os.PathSeparator)
	for _, fileInfoPath := range files {
		//log.Print(fileInfoPath)
		task := new(upload.TUploadTask)
		task.LocalPath = fileInfoPath
		//index := strings.LastIndex(fileInfoPath, PathSep)
		//name := fileInfoPath[index+1:]
		task.TargetPath = "d:\\123\\" // + GetRandomString(4) + PathSep
		task.TaskID = GetRandomString(16)
		taskMap[fileInfoPath] = task
		task.StartTask()
	}

}

func GetFormatSize(qSize int64) string {
	cKB := int64(1024)
	cMB := int64(1024 * 1024)
	cGB := int64(1024 * 1024 * 1024)

	if qSize >= cGB {
		return strconv.FormatFloat(float64(qSize/cGB), 'f', -1, 64) + " GB"
	} else if qSize >= cMB {
		return strconv.FormatFloat(float64(qSize/cMB), 'f', -1, 64) + " MB"
	} else if qSize >= cKB {
		return strconv.FormatFloat(float64(qSize/cKB), 'f', -1, 64) + " KB"
	} else {
		return strconv.FormatInt(qSize, 10) + " B"
	}
}

func GetTimeString(d time.Duration) string {
	hours := int(d.Hours())
	mins := int(d.Minutes())
	secs := int(d.Seconds())
	msecs := int(d.Nanoseconds() / 1000)
	var timeStr string

	if hours > 0 {
		timeStr = fmt.Sprintf("%d时 %d分 %d秒", hours, mins%60, secs%3600)
	} else if mins > 0 {
		timeStr = fmt.Sprintf("%d分 %d秒", mins, secs%60)
	} else if secs > 0 {
		timeStr = fmt.Sprintf("%d秒 %d毫秒", secs, msecs%1000)
	} else if msecs > 0 {
		timeStr = fmt.Sprintf("%d毫秒", msecs)
	}
	return timeStr
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(16)
}

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
