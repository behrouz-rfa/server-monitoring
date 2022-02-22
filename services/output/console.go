package output

import (
	"bytes"
	"fmt"
	"server-monitoring/domain/requests"
	"server-monitoring/services/color"
	"server-monitoring/services/protos"
	"server-monitoring/services/protos/http"
	"server-monitoring/services/protos/tls"
	"strings"
	"time"
)

type Console struct {
	hasShowHeader  bool
	showBody       bool
	showRaw        bool
	showMoreDetail bool
}

type consoleContent struct {
	srcAddr string
	srcPort int
	dstAddr string
	dstPort int

	method        string
	ts            time.Time
	statusCode    int
	contentLength int
	url           string
	UserAgent     string
	Body          []byte
	Response      []byte
}

//prin on conlos filter the on base console content
func (t *Console) Print(data protos.Protos, srcAddr string, srcPort int, dstAddr string, dstPort int) {
	t.printHeaderDescription()
	switch v := data.(type) {
	case *http.Http:
		url := fmt.Sprintf("http://%s%s", v.Request.Host, v.Request.RequestURI)
		if strings.Contains(v.Response.StatusPhrase, "Switching Protocols") {
			url = fmt.Sprintf("ws://%s%s", v.Request.Host, v.Request.RequestURI)
		}
		if (v.Response.StatusCode == 301 || v.Response.StatusCode == 302) && v.Response.Location != "" {
			url = fmt.Sprintf("%s -> %s", url, v.Response.Location)
		}

		content := consoleContent{
			srcAddr:       srcAddr,
			srcPort:       srcPort,
			dstAddr:       dstAddr,
			dstPort:       dstPort,
			method:        v.Request.Method,
			ts:            v.Request.Ts,
			statusCode:    v.Response.StatusCode,
			contentLength: v.Response.ContentLength,
			url:           url,
			UserAgent:     v.Request.UserAgent,
			Body:          v.Request.Body,
			Response:      v.Response.Body,
		}
		t.printContent(&content)
		t.printHttpRaw(v)
	case *tls.ClientHello:
		content := consoleContent{
			srcAddr:       srcAddr,
			srcPort:       srcPort,
			dstAddr:       dstAddr,
			dstPort:       dstPort,
			method:        v.Method,
			ts:            v.Ts,
			statusCode:    v.StatusCode,
			contentLength: v.ContentLength,
			url:           v.RequestURI,
		}
		t.printContent(&content)
	}
}

func (t *Console) printContent(content *consoleContent) {
	request := requests.Request{
		SrcAddr:       content.srcAddr,
		SrcPort:       content.srcPort,
		DstAddr:       content.dstAddr,
		DstPort:       content.dstPort,
		Method:        content.method,
		Ts:            content.ts,
		StatusCode:    content.statusCode,
		ContentLength: content.contentLength,
		Url:           content.url,
		UserAgent:     content.UserAgent,
		Body:          content.Body,
		Response:      content.Response,
	}
	request.InsertConsoleLog()
	//database.InsertConsoleLog(request)

	fmt.Println("body: ", string(content.Body))
	fmt.Println("response: ", string(content.Response))
	if t.showMoreDetail {
		color.Printf("%-23s %-42s %-5d %-7s %-5s %s\n",
			color.MethodColor(content.method),
			content.ts.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%s:%d->%s:%d", content.srcAddr, content.srcPort, content.dstAddr, content.dstPort),
			content.statusCode,
			t.humanizeSize(int64(content.contentLength)),
			content.method,
			content.url)
	} else {
		color.Printf("%-23s %-5d %-7s %-5s %s\n",
			color.MethodColor(content.method),
			content.ts.Format("2006-01-02 15:04:05"),
			content.statusCode,
			t.humanizeSize(int64(content.contentLength)),
			content.method,
			content.url)
	}
}

func (t *Console) printHeaderDescription() {
	if t.hasShowHeader {
		return
	}

	fmt.Printf("%-23s %-5s %-7s %-10s %s\n",
		"time",
		"status",
		"length",
		"method",
		"url")
	t.hasShowHeader = true
}

func (t *Console) printHttpRaw(v *http.Http) {
	if !t.showRaw && !t.showBody {
		return
	}

	requestHeaders := ""
	responseHeaders := ""
	if t.showRaw {
		requestHeaders = fmt.Sprintf("%v", string(v.Request.RawHeaders))
		responseHeaders = fmt.Sprintf("%v", string(v.Response.RawHeaders))
	}

	requestBody := ""
	responseBody := ""
	if t.showBody {
		if bytes.Equal(v.Request.Body, []byte("")) {
			requestBody = "***** REQUEST BODY EMPTY *****\n"
		} else {
			if v.Request.TotslExtraMsgSize > 0 {
				requestBody = fmt.Sprintf("%v\n***** REQUEST BODY TOO LARGE (over limit of %d) *****\n", string(v.Request.Body), v.Request.TotslExtraMsgSize)
			} else {
				requestBody = fmt.Sprintf("%v\n", string(v.Request.Body))
			}
		}
		if bytes.Equal(v.Response.Body, []byte("")) {
			responseBody = "***** RESPONSE BODY EMPTY *****\n"
		} else {
			if v.Response.TotslExtraMsgSize > 0 {
				responseBody = fmt.Sprintf("%v\n***** RESPONSE BODY TOO LARGE (over limit of %d) *****\n", string(v.Response.Body), v.Response.TotslExtraMsgSize)
			} else {
				responseBody = fmt.Sprintf("%v\n", string(v.Response.Body))
			}
		}
	}

	fmt.Printf("==============================================================================\n"+
		"~~~ REQUEST ~~~\n\n%s", requestHeaders)
	color.PrintlnResponse(strings.TrimSpace(requestBody))
	fmt.Println("\n------------------------------------------------------------------------------")

	fmt.Printf("~~~ RESPONSE ~~~\n\n%s", responseHeaders)
	color.PrintlnResponse(strings.TrimSpace(responseBody))
	fmt.Println("==============================================================================")
}

func (t *Console) humanizeSize(length int64) string {
	if length >= 1024*1024 {
		return fmt.Sprintf("%dMB", length/(1024*1024))
	} else if length >= 1024 {
		return fmt.Sprintf("%dKB", length/1024)
	} else {
		return fmt.Sprintf("%dB", length)
	}
}
