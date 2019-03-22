package mod

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func NewHttp() *HttpClass {
	return &HttpClass{}
}

type HttpClass struct {
	HttpClient http.Client
}

type HttpClassRet struct {
	Body   string
	Cookie string
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

//合并2个Cookie
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
