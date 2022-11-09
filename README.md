# Requests
[![license](http://dmlc.github.io/img/apache2.svg)](https://raw.githubusercontent.com/ahuigo/requests/master/LICENSE)
Requests is an HTTP library like python requests.


# Installation

```
go get -u github.com/ahuigo/requests
```

# Examples
> For more examples, refer to [examples/post](https://github.com/ahuigo/requests/blob/master/examples/post_test.go) and [all examples](https://github.com/ahuigo/requests/tree/master/examples)

## Get
    package main
    import (
        "github.com/ahuigo/requests"
        "fmt"
    )

    func main(){
        var json map[string]interface{}
        params := requests.Params{"name": "ahuigo", "page":"1"}
        resp, err := requests.Get("https://httpbin.org/json", params)
        if err != nil {
            panic(err)
        }else{
            resp.Json(&json)
            for k, v := range json {
                fmt.Println(k, v)
            }
        }
    }


## Post

### Post params

    // Post params
    func TestPostParams(t *testing.T) {
        println("Test POST: post params")
        data := requests.Params{
            "name": "ahuigo",
        }
        resp, err := requests.Post("https://www.httpbin.org/post", data)
        if err == nil {
            fmt.Println(resp.Text())
        }
    }

#### Post application/x-www-form-urlencoded 

    // Post Form UrlEncoded data: application/x-www-form-urlencoded
    func TestPostFormUrlEncode(t *testing.T) {
        resp, err := requests.Post(
            "https://www.httpbin.org/post",
            requests.Datas{
                "name": "ahuigo",
            },
        )
        if err != nil {
            t.Error(err)
        }
        var data = struct {
            Form struct {
                Name string
            }
        }{}
        resp.Json(&data)
        if data.Form.Name != "ahuigo" {
            t.Error("invalid response body:", resp.Text())
        }
    }


### Post multipart/form-data

    // Test POST:  multipart/form-data; boundary=....
    func TestPostFormData(t *testing.T) {
        resp, err := requests.Post(
            "https://www.httpbin.org/post",
            requests.FormData{
                "name": "ahuigo",
            },
        )
        if err != nil {
            t.Error(err)
        }
        var data = struct {
            Form struct {
                Name string
            }
        }{}
        resp.Json(&data)
        if data.Form.Name != "ahuigo" {
            t.Error("invalid response body:", resp.Text())
        }
    }


### Post Json: application/json 
    func TestPostJson(t *testing.T) {
    	// You can also use `json := requests.Jsoni(anyTypeData)`
        json := requests.Json{
            "key": "value",
        }
        resp, err := requests.Post("https://www.httpbin.org/post", json)
        if err == nil {
            fmt.Println(resp.Text())
        }
    }

### Post Raw text/plain
    func TestRawString(t *testing.T) {
        println("Test POST: post data and json")
        rawText := "raw data: Hi, Jack!"
        resp, err := requests.Post(
            "https://www.httpbin.org/post", 
            rawText,
            requests.Header{"Content-Type": "text/plain"},
        )
        if err == nil {
            fmt.Println(resp.Text())
        }
    }

### PostFiles

	path, _ := os.Getwd()
	session := requests.R()

	resp, err := session.SetDebug().Post(
		"https://www.httpbin.org/post",
		requests.Files{
            "file1": path + "/README.md",
            "file2": path + "/version",
        },
	)
	if err == nil {
		fmt.Println(resp.Text())
	}

## Session Support
    // 0. Make a session
	session := r.R()

    // 1. First, set cookies: count=100
	var data struct {
		Cookies struct {
			Count string `json:"count"`
		}
	}
	session.Get("https://httpbin.org/cookies/set?count=100")

	// 2. Second, get cookies
	resp, err := session.Get("https://httpbin.org/cookies")
	if err == nil {
		resp.Json(&data)
        if data.Cookies.Count!="100"{
            t.Fatal("Failed to get valid cookies: "+resp.Text())
        }
	}

Warning: Session is not safe in multi goroutines. You can not do as following:

    // Bad! Do not call session in in multi goroutine!!!!!
    session := requests.R()

    // goroutine 1
    go func(){
       session.Post(url1) 
    }()

    // goroutine 2
    go func(){
       session.Post(url2) 
    }()

## Request Options

### SetTimeout

    session := Requests.R()
    session.SetTimeout(20*time.Second)

### Debug Mode
Refer to https://github.com/ahuigo/requests/blob/master/examples/debug_test.go

	session := requests.R()

	resp, err := session.SetDebug().Post(
		"https://www.httpbin.org/post",
		requests.Files{
            "file1": "/README.md",
            "file2": "/version",
        },
	)

### Set Authentication
    session := requests.R()
    resp,_ := session.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})

### Set Cookie
	cookie1 := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
    session.SetCookie(&cookie1)

### Set header

    func TestGetParamsHeaders(t *testing.T) {
        println("Test Get: custom header and params")
        requests.Get("http://www.zhanluejia.net.cn",
            requests.Header{"Referer": "http://www.jeapedu.com"},
            requests.Params{"page": "1", "size": "20"},
            requests.Params{"name": "ahuio"},
        )
    }

    func TestGetParamsHeaders2(t *testing.T) {
        session := requests.R()
        session.SetHeader("accept-encoding", "gzip, deflate, br")
        session.Run("http://www.zhanluejia.net.cn",
            requests.Params{"page": "1", "size": "20"},
            requests.Params{"name": "ahuio"},
        )
    }

    func TestResponseHeader(t *testing.T) {
        resp, _ := requests.Get("https://www.baidu.com/")
        println(resp.Text())
        println(resp.R.Header.Get("location"))
        println(resp.R.Header.Get("Location"))
    }

two or more headers ...

    headers1 := requests.Header{"Referer": "http://www.jeapedu.com"},
    ....
    resp,_ = session.Get(
        "http://www.zhanluejia.net.cn",
        headers1,
        headers2,
        headers3,
    )


## Response
### Fetch Response Body
https://github.com/ahuigo/requests/blob/master/examples/resp_test.go

    fmt.Println(resp.Text())    // string
    fmt.Println(resp.Body())    //[]byte

### Fetch Response Cookies
https://github.com/ahuigo/requests/blob/master/examples/cookie_test.go

    resp,_ = session.Get("https://www.httpbin.org/get")
    coo := resp.Cookies()
    for _, c:= range coo{
        fmt.Println(c.Name,c.Value)
    }

### Response examples

    func TestResponse(t *testing.T) {
        resp, _ := requests.Get("https://httpbin.org/get")
        fmt.Println("Status Code:", resp.StatusCode())
        fmt.Println("Time:", resp.Time())
        fmt.Println("Size:", resp.Size())
        fmt.Println("Headers:")
        for key, value := range resp.Header() {
            fmt.Println(key, "=", value)
        }
        fmt.Println("Cookies:")
        for i, cookie := range resp.Cookies() {
            fmt.Printf("cookie%d: name:%s value:%s\n", i, cookie.Name, cookie.Value)
        }
    }

## Return Error
    func TestError(t *testing.T) {
        _, err := requests.Get("http://127.0.0.1:12346/connect-refused")
        if err2, ok := err.(*rerrors.Error); !ok {
            t.Fatalf("unexpected error:%+v", err)
        } else {
            switch err2.ErrType {
            case rerrors.NetworkError:
                t.Log(err2.ErrType)
            case rerrors.URLError:
                t.Log(err2.ErrType)
            default:
                t.Log(err2.ErrType)
            }
        }
    }

# Utils

## Build Request

	req, err := r.BuildRequest("post", "http://baidu.com/a/b/c", r.Json{
		"age": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(req.Body)
	expectedBody := `{"age":1}`
	if string(body) != expectedBody {
		t.Fatal("Failed to build request")
	}

## Generate curl shell command

	req, _ := requests.BuildRequest("post", "https://baidu.com/path?q=curl&v=1", requests.Json{
		"age": 1,
	})
	curl := requests.BuildCurlRequest(req)
	if !regexp.MustCompile(`^curl -X POST .+ 'https://baidu.com/path\?q=curl&v=1'`).MatchString(curl) {
		t.Fatal(`bad curl cmd: ` + curl)
	}

## Retry request

	r := requests.R().
		SetRetryCount(maxRetries).
		SetRetryCondition(func(resp *requests.Response, err error) bool {
			if err != nil {
				return true
			}
			var json map[string]interface{}
			resp.Json(&json)
			return json["status"]!="ok"
		})


# Feature Support
  - Set headers, params, body, json, ....
  - Multipart File Uploads
  - R-seesion with Cookie Persistence
  - Proxy
      - Http Proxy
      - [] Socks5 Proxy
  - Authentication
  - JSON
  - Chunked Requests
  - Debug Mode
      - Support Curl Dump
  - SetTimeout
  - Retry
  - [] Support traceInfo(like resty's clientTrace and gout)
  - [] Support http stream and callback(notParseResponse)
  - [] Remove Global Setting
  - [] Support http2
  - [] Support retry
  - [] Support connection pool
  - [] Middleware support?

# Thanks
This project is inspired by [github.com/asmcos/requests](http://github.com/asmcos/requests). 

Great thanks to it :).
