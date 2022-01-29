package http

import (
	"fmt"
	"sync"
	"time"
)

var parserCache sync.Map

type Http struct {
	Request  *HttpMessage
	Response *HttpMessage
}

func (t *Http) GetSearchKey() string {
	return fmt.Sprintf("//%s%s", t.Request.Host, t.Request.RequestURI)
}

// Http Message
type HttpMessage struct {
	Ts               time.Time
	Version          version
	HasContentLength bool
	Connection       string
	chunkedLength    int

	isRequest bool

	//Request Info
	RequestURI   string
	Method       string
	StatusCode   int
	StatusPhrase string

	// Http Headers
	ContentLength int
	ContentType   string
	Host          string
	Referer       string
	UserAgent     string
	Location      string
	encodings     []string
	isChunked     bool
	Headers       map[string]string
	size          uint64

	RawHeaders []byte

	// sendBody determines if the body must be sent along with the event
	// because the content-type is included in the send_body_for setting.
	sendBody bool
	// saveBody determines if the body must be saved. It is set when sendBody
	// is true or when the body type is form-urlencoded.
	saveBody          bool
	Body              []byte
	TotalReceivedSize int
	TotslExtraMsgSize int
}

type version struct {
	major uint8
	minor uint8
}

func (v version) String() string {
	if v.major == 1 && v.minor == 1 {
		return "1.1"
	}
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

func ParseRequest(id string, request []byte, showBody bool) (*Http, bool) {
	p := getParser(id)
	p.requestParser.showBody = showBody
	p.responseParser.showBody = showBody
	ok, complete := p.parseRequest(request)
	updateParser(id, p)
	if ok {
		return p.Result, complete
	} else {
		return nil, complete
	}
}

func ParseReponse(id string, request []byte, showBody bool) (*Http, bool) {

	p := getParser(id)
	p.requestParser.showBody = showBody
	p.responseParser.showBody = showBody
	ok, complete := p.parseResponse(request)
	if complete {
		deleteParser(id)
	} else {
		updateParser(id, p)
	}
	if ok {
		return p.Result, complete
	} else {
		return nil, complete
	}
}

func GetResult(id string) *Http {
	if v, ok := parserCache.Load(id); ok {
		result := v.(*httpParser).Result
		if result.Request.RequestURI != "" {
			return result
		}
	}
	return nil
}

func getParser(id string) *httpParser {
	if v, ok := parserCache.Load(id); ok {
		return v.(*httpParser)
	}

	return newParser()
}

func updateParser(id string, state *httpParser) {
	parserCache.Store(id, state)
}

func deleteParser(id string) {
	parserCache.Delete(id)
}

func Close(id string) {
	deleteParser(id)
}
