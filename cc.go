package closureCompiler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type CCEval struct {
	input, useragent string
	api              string
	mp               map[string]string
	urlVal           url.Values
	cc               *ClosureCompiler
	httpClient       *http.Client
}

func (cce *CCEval) resetMap() {
	for k := range cce.mp {
		delete(cce.mp, k)
	}
}

type ClosureCompiler struct {
	CompiledCode   string    `json:"compiledCode"`
	Errors         []CCError `json:"errors"`
	Error          string    `json:"error"`
	Line           string    `json:"line"`
	OutputFilePath string    `json:"outputFilePath"`
}

type CCError struct {
	Type   string `json:"type"`
	LineNo int    `json:"lineno"`
	CharNo int    `json:"charno"`
	Error  string `json:"error"`
	Line   string `json:"line"`
}

func (cce *CCEval) setupFormData() {
	cce.resetMap()
	cce.mp["output_format"] = `json`
	cce.mp["compilation_level"] = `SIMPLE_OPTIMIZATIONS`
	cce.mp["warning_level"] = `default`
	cce.mp["output_file_name"] = `default.js`
	cce.mp["js_code"] = cce.input
	data := url.Values{}
	for k, v := range cce.mp {
		data.Add(k, v)
	}
	cce.urlVal = data
}

func (cce *CCEval) setupHeaders(r *http.Request) {
	cce.resetMap()
	cce.mp["origin"] = `https://closure-compiler.appspot.com`
	cce.mp["user-agent"] = cce.useragent
	cce.mp["content-type"] = "application/x-www-form-urlencoded;charset=UTF-8"
	cce.mp["accept-language"] = "en-US,en;q=0.8"

	for k, v := range cce.mp {
		r.Header.Add(k, v)
	}
}

// Run example can be seen at NewCCEval
func (cce *CCEval) Run() (string, error) {
	/*ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(5*time.Second, func() {
		cancel()
	})*/
	cce.setupFormData()
	encoded := cce.urlVal.Encode()
	r, err := http.NewRequest("POST", cce.api, bytes.NewBufferString(encoded))
	if err != nil {
		return "", err
	}
	//r = r.WithContext(ctx)
	cce.setupHeaders(r)
	resp, err := cce.httpClient.Do(r)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var b []byte
	b, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(b, cce.cc)
	if err != nil {
		return "", err
	}
	if er := cce.anyErrors(); er == "" {
		return cce.cc.CompiledCode, nil
	}
	return "", errors.New(cce.anyErrors())
}

func (cce *CCEval) anyErrors() string {
	var allErr string
	if cce.cc.CompiledCode == "" && cce.cc.Errors != nil {
		for _, v := range cce.cc.Errors {
			allErr += fmt.Sprintf("\n Compilation error at lineNo %v charNo %v :: %q ", v.LineNo, v.CharNo, v.Line)
			allErr += fmt.Sprintf("\n Erro Message > %v", v.Error)
		}
	}
	return allErr
}

/*NewCCEval example
jsCode = getSomeRawJSCode()
cce = closureCompiler.NewCCEval(jsCode, *ua)
		minified, err = cce.Run() // minified is the result from closureCompiler
*/
func NewCCEval(input, useragent string) *CCEval {
	api := `https://httpbin.org/post`
	api = `https://closure-compiler.appspot.com/compile?output_info=warnings&output_info=errors&output_info=statistics&output_info=compiled_code`
	ua := `Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 OPR/43.0.2442.1144`
	if useragent != "" {
		ua = useragent
	}

	cce := &CCEval{
		cc:        &ClosureCompiler{},
		input:     input,
		api:       api,
		mp:        make(map[string]string),
		useragent: ua,
		httpClient: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}

	return cce
}
