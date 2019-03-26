package mod

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
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

type HttpClassRet struct {
	Body   string
	Cookie string
	Header http.Header
}

type TimerRet struct {
	StopCannel chan interface{}
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

//httpClass.Get
func (t *HttpClass) Get(url string, header map[string]string) (*HttpClassRet, error) {
	req, _ := http.NewRequest("GET", url, nil)
	header = t.initHttpRequst(header)
	for h_key, h_var := range header {
		req.Header.Set(h_key, h_var)
	}
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
	t.HttpClient.Timeout = time.Second * 5
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
	if header["content-type"] == "" {
		header["content-type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	}

	if header["user-agent"] == "" {
		header["user-agent"] = "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16"
	}

	t.HttpClient.Timeout = time.Second * 20
	t.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
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
