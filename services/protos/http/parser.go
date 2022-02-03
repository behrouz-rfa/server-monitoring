package http

import (
	"bytes"
	"errors"
	"net/http"
	"server-monitoring/services/caputure/utils"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/rs/zerolog/log"
)

const MAX_RECV_SIZE = 10 * 1024

type parseState int

const (
	stateStart parseState = iota
	stateHeaders
	stateBody
	stateBodyChunkedStart
	stateBodyChunked
	stateBodyChunkedWaitFinalCRLF
	stateCompleted
)

type parser struct {
	showBody     bool
	parseOffset  int
	headerOffset int
	bodyReceived int
	state        parseState
}

type parserStream struct {
	data []byte
}

type httpParser struct {
	requestParser  *parser
	responseParser *parser
	requestStream  *parserStream
	responseStream *parserStream

	Result *Http
}

var (
	transferEncodingChunked = "chunked"

	constCRLF = []byte("\r\n")

	constClose       = []byte("close")
	constKeepAlive   = []byte("keep-alive")
	constHTTPVersion = []byte("HTTP/")

	nameContentLength    = []byte("content-length")
	nameContentType      = []byte("content-type")
	nameTransferEncoding = []byte("transfer-encoding")
	nameContentEncoding  = []byte("content-encoding")
	nameConnection       = []byte("connection")
	nameHost             = []byte("host")
	nameLocation         = []byte("location")
	nameReferer          = []byte("referer")
	nameUserAgent        = []byte("user-agent")

	includeRequestBodyFor  []string = []string{}
	includeResponseBodyFor []string = []string{
		"text/html",
		"text/plain",
		"application/javascript",
		"application/json",
	}

	httpValidMethods = map[string]bool{
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodPatch:   true,
		http.MethodDelete:  true,
		http.MethodConnect: true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}
)

func newParser() *httpParser {
	return &httpParser{
		requestParser:  &parser{},
		responseParser: &parser{},
		requestStream:  &parserStream{},
		responseStream: &parserStream{},
		Result: &Http{
			Request:  &HttpMessage{isRequest: true, Ts: time.Now()},
			Response: &HttpMessage{},
		},
	}
}

func (t *httpParser) parseRequest(request []byte) (bool, bool) {
	extraMsgSize := t.bufferData(t.Result.Request, t.requestStream, request)
	return t.requestParser.doParse(t.Result.Request, t.requestStream, extraMsgSize)
}

func (t *httpParser) parseResponse(response []byte) (bool, bool) {
	extraMsgSize := t.bufferData(t.Result.Response, t.responseStream, response)
	return t.responseParser.doParse(t.Result.Response, t.responseStream, extraMsgSize)
}

func (t *httpParser) bufferData(m *HttpMessage, stream *parserStream, data []byte) int {
	extraMsgSize := 0
	if m.TotalReceivedSize < MAX_RECV_SIZE {
		stream.data = append(stream.data, data...)
	} else {
		extraMsgSize = len(data)
		m.TotslExtraMsgSize += len(data)
	}
	m.TotalReceivedSize += len(data)
	return extraMsgSize
}

func (parser *parser) doParse(m *HttpMessage, s *parserStream, extraMsgSize int) (bool, bool) {
	if extraMsgSize > 0 {
		// A packet of extraMsgSize size was seen, but we don't have
		// its actual bytes. This is only usable in the `stateBody` state.
		if parser.state != stateBody {
			return false, false
		}

		return parser.eatBody(m, s, extraMsgSize)
	}

	for parser.parseOffset < len(s.data) {
		switch parser.state {
		case stateStart:
			if cont, ok, complete := parser.parseHTTPLine(m, s); !cont {
				return ok, complete
			}
		case stateHeaders:
			if cont, ok, complete := parser.parseHeaders(m, s); !cont {
				return ok, complete
			}
		case stateBody:
			return parser.parseBody(m, s)
		case stateBodyChunkedStart:
			if cont, ok, complete := parser.parseBodyChunkedStart(m, s); !cont {
				return ok, complete
			}
		case stateBodyChunked:
			if cont, ok, complete := parser.parseBodyChunked(m, s); !cont {
				return ok, complete
			}
		case stateBodyChunkedWaitFinalCRLF:
			return parser.parseBodyChunkedWaitFinalCRLF(m, s)
		case stateCompleted:
			return true, true
		}
	}

	return true, false
}

func (parser *parser) parseHTTPLine(m *HttpMessage, s *parserStream) (cont, ok, complete bool) {
	i := bytes.Index(s.data[parser.parseOffset:], []byte("\r\n"))
	if i == -1 {
		return false, false, false
	}

	// Very basic tests on the first line. Just to check that
	// we have what looks as an HTTP message
	var version []byte
	var err error
	fline := s.data[parser.parseOffset:i]
	if len(fline) < 9 {
		log.Debug().Msg("First line too small")
		return false, false, false
	}

	if bytes.Equal(fline[0:5], constHTTPVersion) {
		//RESPONSE
		version = fline[5:8]
		m.StatusCode, m.StatusPhrase, err = parseResponseStatus(fline[9:])
		if err != nil {
			log.Warn().Msgf("Failed to understand HTTP response status: %s", fline[9:])
			return false, false, false
		}

		log.Debug().Msgf("HTTP status_code=%d, status_phrase=%s", m.StatusCode, "-")
	} else {
		// REQUEST
		afterMethodIdx := bytes.IndexFunc(fline, unicode.IsSpace)
		afterRequestURIIdx := bytes.LastIndexFunc(fline, unicode.IsSpace)

		// Make sure we have the VERB + URI + HTTP_VERSION
		if afterMethodIdx == -1 || afterRequestURIIdx == -1 || afterMethodIdx == afterRequestURIIdx {
			log.Debug().Msg("Couldn't understand HTTP request")
			return false, false, false
		}

		m.Method = string(fline[:afterMethodIdx])
		if _, found := httpValidMethods[m.Method]; !found {
			log.Debug().Msg("Couldn't understand HTTP request")
			return false, false, false
		}
		m.RequestURI = string(fline[afterMethodIdx+1 : afterRequestURIIdx])

		versionIdx := afterRequestURIIdx + len(constHTTPVersion) + 1
		if len(fline) > versionIdx && bytes.Equal(fline[afterRequestURIIdx+1:versionIdx], constHTTPVersion) {
			version = fline[versionIdx:]
		} else {
			log.Debug().Msg("Couldn't understand HTTP version")
			return false, false, false
		}
	}

	m.Version.major, m.Version.minor, err = parseVersion(version)
	if err != nil {
		log.Debug().Msgf("Failed to understand HTTP version: %v", version)
		m.Version.major = 1
		m.Version.minor = 0
	}
	log.Debug().Msgf("HTTP version %d.%d", m.Version.major, m.Version.minor)

	// ok so far
	parser.parseOffset = i + 2
	parser.headerOffset = parser.parseOffset
	parser.state = stateHeaders

	return true, true, true
}

func parseResponseStatus(s []byte) (int, string, error) {
	log.Trace().Msgf("parseResponseStatus: %s", s)

	var phrase []byte
	p := bytes.IndexByte(s, ' ')
	if p == -1 {
		p = len(s)
	} else {
		phrase = s[p+1:]
	}
	statusCode := utils.MustParseInt(string(s[0:p]))
	return statusCode, string(phrase), nil
}

func parseVersion(s []byte) (uint8, uint8, error) {
	if len(s) < 3 {
		return 0, 0, errors.New("Invalid version")
	}

	major := s[0] - '0'
	minor := s[2] - '0'
	if major > 1 || minor > 2 {
		return 0, 0, errors.New("unsupported version")
	}
	return uint8(major), uint8(minor), nil
}

func (parser *parser) parseHeaders(m *HttpMessage, s *parserStream) (cont, ok, complete bool) {
	if len(s.data)-parser.parseOffset >= 2 &&
		bytes.Equal(s.data[parser.parseOffset:parser.parseOffset+2], []byte("\r\n")) {
		// EOH
		m.size = uint64(parser.parseOffset + 2)
		m.RawHeaders = s.data[:m.size]
		s.data = s.data[m.size:] // split only body data
		parser.parseOffset = 0   // reset to body offset

		if !m.isRequest && ((100 <= m.StatusCode && m.StatusCode < 200) || m.StatusCode == 204 || m.StatusCode == 304) {
			//response with a 1xx, 204 , or 304 status  code is always terminated
			// by the first empty line after the  header fields
			log.Debug().Msgf("Terminate response, status code %d", m.StatusCode)
			return false, true, true
		}

		if !parser.showBody {
			log.Debug().Msg("Ignore parse body")
			parser.state = stateCompleted
			return false, true, true
		}

		// fmt.Printf("isRequest: %v ContentType: %s\n", m.isRequest, m.ContentType)
		if m.isRequest {
			m.sendBody = parser.shouldIncludeInBody(m.ContentType, includeRequestBodyFor)
		} else {
			m.sendBody = parser.shouldIncludeInBody(m.ContentType, includeResponseBodyFor)
		}
		m.saveBody = m.sendBody || (m.ContentLength > 0 && strings.Contains(m.ContentType, "urlencoded")) || (m.ContentLength > 0 && strings.Contains(m.ContentType, "application/json"))
		if m.isChunked {
			// support for HTTP/1.1 Chunked transfer
			// Transfer-Encoding overrides the Content-Length
			log.Debug().Msg("Read chunked body")

			parser.state = stateBodyChunkedStart
			return true, true, true
		}

		if m.ContentLength == 0 && (m.isRequest || m.HasContentLength) {
			log.Debug().Msg("Empty content length, ignore body")

			// Ignore body for request that contains a message body but not a Content-Length
			parser.state = stateCompleted
			return false, true, true
		}

		log.Debug().Msg("Read body")
		parser.state = stateBody
	} else {
		ok, hfcomplete, offset := parser.parseHeader(m, s.data[parser.parseOffset:])
		if !ok {
			return false, false, false
		}
		if !hfcomplete {
			return false, true, false
		}
		parser.parseOffset += offset
	}
	return true, true, true
}

func (parser *parser) parseHeader(m *HttpMessage, data []byte) (bool, bool, int) {
	if m.Headers == nil {
		m.Headers = make(map[string]string)
	}
	i := bytes.Index(data, []byte(":"))
	if i == -1 {
		// Expected \":\" in headers. Assuming incomplete"
		return true, false, 0
	}

	// skip folding line
	for p := i + 1; p < len(data); {
		q := bytes.Index(data[p:], constCRLF)
		if q == -1 {
			// Assuming incomplete
			return true, false, 0
		}
		p += q
		if len(data) > p && (data[p+1] == ' ' || data[p+1] == '\t') {
			p = p + 2
		} else {
			headerName := bytes.ToLower(data[:i])
			headerVal := trim(data[i+1 : p])
			log.Debug().Msgf("Header: '%s' Value: '%s'", data[:i], headerVal)

			// Headers we need for parsing. Make sure we always
			// capture their value
			if bytes.Equal(headerName, nameContentLength) {
				m.ContentLength = utils.MustParseInt(string(headerVal))
				m.HasContentLength = m.ContentLength > 0
			} else if bytes.Equal(headerName, nameContentType) {
				m.ContentType = string(headerVal)
			} else if bytes.Equal(headerName, nameLocation) {
				m.Location = string(headerVal)
			} else if bytes.Equal(headerName, nameTransferEncoding) {
				encodings := parseCommaSeparatedList(headerVal)
				// 'chunked' can only appear at the end
				if n := len(encodings); n > 0 && encodings[n-1] == transferEncodingChunked {
					m.isChunked = true
					encodings = encodings[:n-1]
				}
				if len(encodings) > 0 {
					// Append at the end of encodings. If a content-encoding
					// header is also present, it was applied by sender before
					// transfer-encoding.
					m.encodings = append(m.encodings, encodings...)
				}
			} else if bytes.Equal(headerName, nameContentEncoding) {
				encodings := parseCommaSeparatedList(headerVal)
				// Append at the beginning of m.encodings, as Content-Encoding
				// is supposed to be applied before Transfer-Encoding.
				m.encodings = append(encodings, m.encodings...)
			} else if bytes.Equal(headerName, nameConnection) {
				m.Connection = string(headerVal)
			} else if bytes.Equal(headerName, nameHost) {
				m.Host = string(headerVal)
			} else if bytes.Equal(headerName, nameReferer) {
				m.Referer = string(headerVal)
			} else if bytes.Equal(headerName, nameUserAgent) {
				m.UserAgent = string(headerVal)
			}

			return true, true, p + 2
		}
	}

	return true, false, len(data)
}

func parseCommaSeparatedList(s []byte) (list []string) {
	values := bytes.Split(s, []byte(","))
	list = make([]string, len(values))
	for idx := range values {
		list[idx] = string(bytes.ToLower(bytes.Trim(values[idx], " ")))
	}
	return list
}

func (parser *parser) parseBody(m *HttpMessage, s *parserStream) (ok, complete bool) {
	nbytes := len(s.data)
	if !m.HasContentLength && (bytes.Equal([]byte(m.Connection), constClose) ||
		(isVersion(m.Version, 1, 0) && !bytes.Equal([]byte(m.Connection), constKeepAlive))) {
		m.size += uint64(nbytes)
		parser.bodyReceived += nbytes
		m.ContentLength += nbytes

		// HTTP/1.0 no content length. Add until the end of the connection
		log.Debug().Msgf("http conn close, received %d", len(s.data))
		if m.saveBody {
			m.Body = append(m.Body, s.data...)
		}
		s.data = nil
		return true, false
	} else if nbytes >= m.ContentLength-parser.bodyReceived {
		wanted := m.ContentLength - parser.bodyReceived
		if m.saveBody {
			m.Body = append(m.Body, s.data[:wanted]...)
		}
		parser.bodyReceived = m.ContentLength
		m.size += uint64(wanted)
		s.data = s.data[wanted:]
		return true, true
	} else {
		if m.saveBody {
			m.Body = append(m.Body, s.data...)
		}
		s.data = nil
		parser.bodyReceived += nbytes
		m.size += uint64(nbytes)
		log.Trace().Msgf("bodyReceived: %d", parser.bodyReceived)
		return true, false
	}
}

// eatBody acts as if size bytes were received, without having access to
// those bytes.
func (parser *parser) eatBody(m *HttpMessage, s *parserStream, nbytes int) (ok, complete bool) {
	log.Debug().Msg("eatBody body")
	if !m.HasContentLength && (bytes.Equal([]byte(m.Connection), constClose) ||
		(isVersion(m.Version, 1, 0) && !bytes.Equal([]byte(m.Connection), constKeepAlive))) {

		// HTTP/1.0 no content length. Add until the end of the connection
		log.Debug().Msgf("http conn close, received %d", nbytes)
		m.size += uint64(nbytes)
		parser.bodyReceived += nbytes
		m.ContentLength += nbytes
		return true, false
	} else if nbytes >= m.ContentLength-parser.bodyReceived {
		wanted := m.ContentLength - parser.bodyReceived
		parser.bodyReceived = m.ContentLength
		m.size += uint64(wanted)
		return true, true
	} else {
		parser.bodyReceived += nbytes
		m.size += uint64(nbytes)
		log.Debug().Msgf("bodyReceived: %d", parser.bodyReceived)
		return true, false
	}
}

func (parser *parser) parseBodyChunkedStart(m *HttpMessage, s *parserStream) (cont, ok, complete bool) {
	// read hexa length
	i := bytes.Index(s.data, constCRLF)
	if i == -1 {
		return false, true, false
	}
	line := string(s.data[:i])
	chunkLength, err := strconv.ParseInt(line, 16, 32)
	if err != nil {
		log.Warn().Msg("Failed to understand chunked body start line")
		return false, false, false
	}
	m.chunkedLength = int(chunkLength)

	s.data = s.data[i+2:] //+ \r\n
	m.size += uint64(i + 2)

	if m.chunkedLength == 0 {
		if len(s.data) < 2 {
			parser.state = stateBodyChunkedWaitFinalCRLF
			return false, true, false
		}
		m.size += 2
		if s.data[0] != '\r' || s.data[1] != '\n' {
			log.Warn().Msg("Expected CRLF sequence at end of message")
			return false, false, false
		}
		s.data = s.data[2:]
		parser.state = stateCompleted
		return false, true, true
	}
	parser.bodyReceived = 0
	parser.state = stateBodyChunked

	return true, true, false
}

func (parser *parser) parseBodyChunked(m *HttpMessage, s *parserStream) (cont, ok, complete bool) {
	wanted := m.chunkedLength - parser.bodyReceived
	if len(s.data) >= wanted+2 /*\r\n*/ {
		// Received more data than expected
		if m.saveBody {
			m.Body = append(m.Body, s.data[:wanted]...)
		}
		m.size += uint64(wanted + 2)
		s.data = s.data[wanted+2:]
		m.ContentLength += m.chunkedLength
		parser.state = stateBodyChunkedStart
		return true, true, false
	}

	if len(s.data) >= wanted {
		// we need need to wait for the +2, else we can crash on next call
		return false, true, false
	}

	// Received less data than expected
	if m.saveBody {
		m.Body = append(m.Body, s.data...)
	}
	parser.bodyReceived += len(s.data)
	m.size += uint64(len(s.data))
	s.data = nil
	return false, true, false
}

func (parser *parser) parseBodyChunkedWaitFinalCRLF(m *HttpMessage, s *parserStream) (ok, complete bool) {
	if len(s.data) < 2 {
		return true, false
	}

	m.size += 2
	if s.data[0] != '\r' || s.data[1] != '\n' {
		log.Warn().Msg("Expected CRLF sequence at end of message")
		return false, false
	}

	s.data = s.data[2:]
	parser.state = stateCompleted
	return true, true
}

func (parser *parser) shouldIncludeInBody(contenttype string, capturedContentTypes []string) bool {
	for _, include := range capturedContentTypes {
		if strings.Contains(contenttype, include) {
			log.Debug().Msgf("Should Include Body = true Content-Type %s include_body %s",
				contenttype, include)
			return true
		}
	}
	log.Debug().Msgf("Should Include Body = false Content-Type %s", contenttype)
	return false
}

func isVersion(v version, major, minor uint8) bool {
	return v.major == major && v.minor == minor
}

func trim(buf []byte) []byte {
	return trimLeft(trimRight(buf))
}

func trimLeft(buf []byte) []byte {
	for i, b := range buf {
		if b != ' ' && b != '\t' {
			return buf[i:]
		}
	}
	return nil
}

func trimRight(buf []byte) []byte {
	for i := len(buf) - 1; i > 0; i-- {
		b := buf[i]
		if b != ' ' && b != '\t' {
			return buf[:i+1]
		}
	}
	return nil
}
