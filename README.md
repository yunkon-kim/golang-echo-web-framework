# Golang Echo web framework study
Golang관련 Web framework는 많이 있습니다. 그 중에서 간단하게 구현할 수 있는 Echo에 대하여 스터디 하였고, 그에 대한 설명과 소스코드를 공유합니다.

## Getting started with Echo web framework
### 스터디 환경
- Windows 10 x64
- Golang 1.15.3 설치

### 1. Echo package 설치
```
go get -u github.com/labstack/echo
```
결과물을 올리는터라 빠진 패키지가 있을 수 있습니다.

### 2. 디렉토리 구조
스터디하며 개인적으로 설정 및 수정한 디렉토리 구조 입니다. 
```
.
├── assets
│   └── xxx.xx
├── css
│   └── xxx.css
├── js
│   └── xxx.js
├── public
│   ├── index.html
│   └── xxx.html
└── echo-server.go
```
### 3. `echo-server.go`
`main`부터 보시길 권장합니다. `main`위쪽에 있는 구조체와 함수는 페이지를 Rendering하고 할때 필요한 부분 입니다.

전체 소스 보기에 앞서 코드를 간단히 설명 드립니다.
- `e := echo.New()`: 에코 객체 생성
- `e.Static("/", "assets")`: asset 등록, 위 `assets 디렉토리`를 `/`로 인식하도록 합니다. 예를 들어, html 에서 /image/picture.png 등을 설정하실때 사용하시면 됩니다.
- `e.Static("/js", "js")`: js 등록, 위 `js` 디렉토리`를 `/js`로 인식하도록 합니다. 예를 들어, html 에서 /js/xxx.js 등을 설정하실때 사용하시면 됩니다.
- `e.Static("/css", "css")`: css 등록, 위 `css 디렉토리`를 `/css`로 인식하도록 합니다. 예를 들어, html 에서 /css/xxx.css 등을 설정하실때 사용하시면 됩니다.
- `e.File("/", "public/index.html")`: 브라우저에서 `/` 요청시, `public/index.html`을 제공합니다.
` `e.GET("/aaa", Index)`: Render를 통해 `aaa.html` 또는 `hello.html`을 렌더링하여 제공합니다.
- `e.Logger.Fatal(e.Start(":8000"))`: 서버를 8000 포트로 시작하고, Logger에 등록합니다.

Goland에서 개발시, Static 설정을 해줬음에도 IDE에서는 Highlight 표시가 되네요..

**주의!! Working directory 설정 확인하셔요.** 저는 GoLand에서 실행 할 때 `Run Configuration`에서 Working directory 설정 확인 확인하지 않아서 경로 문제로 꽤 오랜시간을 허비했네요 ㅠ

```
package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

// Render renders a template document
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}


func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "aaa.html", nil)
}


func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello.html", map[string]interface{}{
		"name": "Dolly!",
	})
}


func main() {

	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = t

	e.Static("/", "assets")
	e.Static("/js", "js")
	e.Static("/css", "css")
	e.File("/", "public/index.html")

	// Render
	e.GET("/aaa", Index)
	e.GET("/hello", Hello)

	e.Logger.Fatal(e.Start(":8000"))
}
```
