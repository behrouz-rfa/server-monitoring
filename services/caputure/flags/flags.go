package flags

type options struct {
	Debug         bool
	Verbose       bool
	InterfaceName string
	Ip            string
	Port          int
	Protos        string // choose protocol to display
	Body          bool   // show request / response body
	Raw           bool   // show raw request / response headers
	MoreDetail    bool
	Keyword       string
}

var Options options
