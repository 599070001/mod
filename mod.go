package mod

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewHttp() *HttpClass {
	return &HttpClass{}
}

func NewTime() *TimeClass {
	return &TimeClass{}
}
func NewStrings() *StringsClass {
	return &StringsClass{}
}

func NewFile() *FileClass {
	return &FileClass{}
}

func NewGoroutinePool(wokerNum int) *GoroutinePool {
	return &GoroutinePool{
		workerNum:      wokerNum,
		TaskChannel:    make(chan Task),
		AddTaskChannel: make(chan Task),
	}
}

//IO类操作
type FileClass struct {
}

//时间类
type TimeClass struct {
}

//字符串处理类
type StringsClass struct {
}

//http处理类
type HttpClass struct {
	HttpClient http.Client
}

type Task struct {
	F func(map[string]string)
	P map[string]string
}

//type GoroutinePool
type GoroutinePool struct {
	workerNum      int
	TaskChannel    chan Task
	AddTaskChannel chan Task
}

type HttpClassRet struct {
	Body   string
	Cookie string
	Header http.Header
}

type TimerRet struct {
	StopCannel chan interface{}
}

func (t *Task) Exec() {
	t.F(t.P)
}

func (pool *GoroutinePool) Run() {

	for i := 0; i < pool.workerNum; i++ {
		go pool.worker()
	}

	for task := range pool.AddTaskChannel {
		pool.TaskChannel <- task
	}

}

func (pool *GoroutinePool) worker() {
	for task := range pool.TaskChannel {
		task.Exec()
	}

}

func (*FileClass) WriteString(path string, data string) error {
	return ioutil.WriteFile(path, []byte(data), os.ModePerm)
}

func (*FileClass) ReadString(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func (*FileClass) AppendString(path string, data string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	n, err := f.Write([]byte(data))
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

//获取随机16位小数 Math.random
func (*TimeClass) Random() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatFloat(rand.Float64(), 'f', 16, 64)
}

//定时器
func (*TimeClass) Timer(t time.Duration, f func()) *TimerRet {
	timerRun := time.NewTicker(t)
	closeChannel := make(chan interface{}, 0)
	go timerF(timerRun, closeChannel, f)
	return &TimerRet{closeChannel}
}

//定时器body
func timerF(t *time.Ticker, closeChannel chan interface{}, f func()) {
	for {
		select {
		case <-t.C:
			f()
		case <-closeChannel:
			t.Stop()
			return
		}
	}
}

//时间戳 10,13
func (*TimeClass) TimeStamp(size int) string {
	timestamp := ""
	switch size {
	case 13:
		timestamp = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	case 10:
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}
	return timestamp
}

//截取中间字符串
func (*StringsClass) BetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start)
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

//获取随机 Math.random
func (*StringsClass) RandomInt(s, e int) string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(e-s) + s)
}

//字符串相似度
func (*StringsClass) SimilarText(first, second string, percent *float64) int {
	var similarText func(string, string, int, int) int
	similarText = func(str1, str2 string, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		// Find the longest segment of the same section in two strings
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for l := 0; (i+l < len1) && (j+l < len2) && (str1[i+l] == str2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(str1, str2, pos1, pos2)
			}
			if (pos1+max < len1) && (pos2+max < len2) {
				s1 := []byte(str1)
				s2 := []byte(str2)
				sum += similarText(string(s1[pos1+max:]), string(s2[pos2+max:]), len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}
	sim := similarText(first, second, l1, l2)
	if percent != nil {
		*percent = float64(sim*200) / float64(l1+l2)
	}
	return sim
}

//关键词组过滤，检查是否存在
func (*StringsClass) FitterKeyWords(input string, words []string) bool {

	for _, value := range words {
		if strings.Contains(input, value) {
			return true
		}
	}
	return false
}

//httpClass.Get
func (t *HttpClass) Get(url string, header map[string]string) (*HttpClassRet, error) {
	req, _ := http.NewRequest("GET", url, nil)
	header = t.initHttpRequst(header)
	for h_key, h_var := range header {
		req.Header.Set(h_key, h_var)
	}
	//t.HttpClient.Timeout = time.Second * 15
	resp, err := t.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return &HttpClassRet{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		fmt.Printf("Bad HTTP Response: %v", resp.Status)
		return &HttpClassRet{}, fmt.Errorf("error http status code %v", resp.Status)
	}
	//get Set-Cookie
	session := resp.Cookies()
	session_arr := []string{}
	for _, sessionItem := range session {
		session_arr = append(session_arr, fmt.Sprintf("%s=%s", sessionItem.Name, sessionItem.Value))
	}
	session_str := strings.Join(session_arr, "; ")

	ret, _ := ioutil.ReadAll(resp.Body)
	return &HttpClassRet{string(ret), session_str, resp.Header}, nil
}

//httpClass.Post
func (t *HttpClass) Post(url string, body string, header map[string]string) (*HttpClassRet, error) {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	header = t.initHttpRequst(header)
	for h_key, h_var := range header {
		req.Header.Set(h_key, h_var)
	}
	//t.HttpClient.Timeout = time.Second * 15
	resp, err := t.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return &HttpClassRet{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		fmt.Printf("Bad HTTP Response: %v", resp.Status)
		return &HttpClassRet{}, fmt.Errorf("error http status code %v", resp.Status)
	}

	//get Set-Cookie
	session := resp.Cookies()
	session_arr := []string{}
	for _, sessionItem := range session {
		session_arr = append(session_arr, fmt.Sprintf("%s=%s", sessionItem.Name, sessionItem.Value))
	}
	session_str := strings.Join(session_arr, "; ")

	ret, _ := ioutil.ReadAll(resp.Body)
	return &HttpClassRet{string(ret), session_str, resp.Header}, nil
}

//初始化http请求header
func (t *HttpClass) initHttpRequst(header map[string]string) map[string]string {
	if header == nil {
		header = map[string]string{}
	}
	if header["content-type"] == "" {
		header["content-type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	}

	if header["user-agent"] == "" {
		header["user-agent"] = "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16"
	}

	t.HttpClient.Timeout = time.Second * 30
	t.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	//忽略https 证书验证
	t.HttpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return header
}

//httpClass.AddCookie 合并2个Cookie
func (t *HttpClass) AddCookie(old string, new string) string {
	cookie := map[string]string{}
	old_arr := strings.Split(old, ";")
	new_arr := strings.Split(new, ";")
	cookie_str := []string{}

	for _, item := range old_arr {
		temp := strings.Split(item, "=")
		if len(temp) == 2 {
			cookie[strings.Trim(temp[0], " ")] = strings.Trim(temp[1], " ")
		}
	}
	for _, item := range new_arr {
		temp := strings.Split(item, "=")
		if len(temp) == 2 {
			cookie[strings.Trim(temp[0], " ")] = strings.Trim(temp[1], " ")
		}
	}
	for c_key, c_val := range cookie {
		cookie_str = append(cookie_str, fmt.Sprintf("%s=%s", c_key, c_val))
	}
	return strings.Join(cookie_str, "; ")
}

func CheckError(e error) {
	if e != nil {
		fmt.Printf("【软件错误】%s", e.Error())
		os.Exit(1)
	}
}

func Info(o interface{}) {
	fmt.Printf("%+v\n", o)
}

func RunPath() string {
	currentPath, ok := os.Getwd()
	if ok != nil {
		panic(ok)
	}
	return currentPath + "/"
}
