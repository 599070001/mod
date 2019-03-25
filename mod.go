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
}

func (*TimeClass) Random() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatFloat(rand.Float64(), 'f', 16, 64)
}

func (*StringsClass) BetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	s := len(start)
	if s > m {
		s = m
	}
	str = string([]byte(str)[s:m])
	return str
}

//httpClass.Get
func (t *HttpClass) Get(url string, header map[string]string) (*HttpClassRet, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for h_key, h_var := range header {
		req.Header.Set(h_key, h_var)
	}
	resp, err := t.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad HTTP Response: %v", resp.Status)
		return nil, fmt.Errorf("error http status code %v", resp.Status)
	}
	//get Set-Cookie
	session := resp.Cookies()
	session_arr := []string{}
	for _, sessionItem := range session {
		session_arr = append(session_arr, fmt.Sprintf("%s=%s", sessionItem.Name, sessionItem.Value))
	}
	session_str := strings.Join(session_arr, "; ")

	ret, _ := ioutil.ReadAll(resp.Body)
	return &HttpClassRet{string(ret), session_str}, nil
}

//httpClass.Post
func (t *HttpClass) Post(url string, body string, header map[string]string) (*HttpClassRet, error) {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))

	if header["content-type"] == "" {
		header["content-type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	}

	for h_key, h_var := range header {
		req.Header.Set(h_key, h_var)
	}
	resp, err := t.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad HTTP Response: %v", resp.Status)
		return nil, fmt.Errorf("error http status code %v", resp.Status)
	}

	//get Set-Cookie
	session := resp.Cookies()
	session_arr := []string{}
	for _, sessionItem := range session {
		session_arr = append(session_arr, fmt.Sprintf("%s=%s", sessionItem.Name, sessionItem.Value))
	}
	session_str := strings.Join(session_arr, "; ")

	ret, _ := ioutil.ReadAll(resp.Body)
	return &HttpClassRet{string(ret), session_str}, nil
}

//httpClass.AddCookie 合并2个Cookie
func (*HttpClass) AddCookie(old string, new string) string {
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
