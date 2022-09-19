package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/browser"
)

// =============================================================================
// Usage and help

func usage() {
	fmt.Println("Get HTTP status code information")
	fmt.Printf("usage: %s [100..511] [-v|--verbose] [-c|--cats] [-d|--dogs]\n", os.Args[0])
}

const helpUsage = "Print usage information"

func isAskingForHelp() bool {
	var help bool

	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage)
	flag.Parse()

	return help
}

// =============================================================================

type Status struct {
	code        string
	title       string
	description string
	url         string
}

type StatusCommand struct {
	fs *flag.FlagSet

	verbose bool
	cats    bool
	dogs    bool
}

var statuses = map[string]Status{
	"100": {
		code:        "100",
		title:       "Continue",
		description: "The initial part of a request has been received and has not yet been rejected by the server. The server intends to send a final response after the request has been fully received and acted upon.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-100-continue",
	},
	"101": {
		code:        "101",
		title:       "Switching Protocols",
		description: "The server understands and is willing to comply with the client's request, via the Upgrade header field, for a change in the application protocol being used on this connection. The server MUST generate an Upgrade header field in the response that indicates which protocol(s) will be in effect after this response.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-101-switching-protocols",
	},
	"103": {
		code:        "103",
		title:       "Early Hints",
		description: "Allows user-agents to perform some operations, such as to speculatively load resources that are likely to be used by the document, before the navigation request is fully handled by the server and a response code is served.",
		url:         "https://html.spec.whatwg.org/multipage/semantics.html#early-hints",
	},
	"200": {
		code:        "200",
		title:       "OK",
		description: "The request has succeeded. The content sent in a 200 response depends on the request method.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-200-ok",
	},
	"201": {
		code:        "201",
		title:       "Created",
		description: "The request has been fulfilled and has resulted in one or more new resources being created.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-201-created",
	},
	"202": {
		code:        "202",
		title:       "Accepted",
		description: "The request has been accepted for processing, but the processing has not been completed.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-202-accepted",
	},
	"203": {
		code:        "203",
		title:       "Non-Authoritative Information",
		description: "The request was successful but the enclosed content has been modified from that of the origin server's 200 (OK) response by a transforming proxy.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-203-non-authoritative-infor",
	},
	"204": {
		code:        "204",
		title:       "No Content",
		description: "The server has successfully fulfilled the request and that there is no additional content to send in the response content.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-204-no-content",
	},
	"205": {
		code:        "205",
		title:       "Reset Content",
		description: "The server has fulfilled the request and desires that the user agent reset the \"document view\", which caused the request to be sent, to its original state as received from the origin server.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-205-reset-content",
	},
	"206": {
		code:        "206",
		title:       "Partial Content",
		description: "The server is successfully fulfilling a range request for the target resource by transferring one or more parts of the selected representation.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-206-partial-content",
	},
	"300": {
		code:        "300",
		title:       "Multiple Choices",
		description: "The target resource has more than one representation, each with its own more specific identifier, and information about the alternatives is being provided so that the user (or user agent) can select a preferred representation by redirecting its request to one or more of those identifiers.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-300-multiple-choices",
	},
	"301": {
		code:        "301",
		title:       "Moved Permanently",
		description: "The target resource has been assigned a new permanent URI and any future references to this resource ought to use one of the enclosed URIs.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-301-moved-permanently",
	},
	"302": {
		code:        "302",
		title:       "Found",
		description: "The target resource resides temporarily under a different URI.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-302-found",
	},
	"303": {
		code:        "303",
		title:       "See Other",
		description: "The server is redirecting the user agent to a different resource, as indicated by a URI in the Location header field, which is intended to provide an indirect response to the original request.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-303-see-other",
	},
	"304": {
		code:        "304",
		title:       "Not Modified",
		description: "A conditional GET or HEAD request has been received and would have resulted in a 200 (OK) response if it were not for the fact that the condition evaluated to false.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-304-not-modified",
	},
	"305": {
		code:        "305",
		title:       "Use Proxy",
		description: "(305 was defined in a previous version of this specification and is now deprecated)",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-305-use-proxy",
	},
	"306": {
		code:        "306",
		title:       "(Unused)",
		description: "(306 was defined in a previous version of this specification, is no longer used, and the code is reserved.)",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-306-unused",
	},
	"307": {
		code:        "307",
		title:       "Temporary Redirect",
		description: "The target resource resides temporarily under a different URI and the user agent MUST NOT change the request method if it performs an automatic redirection to that URI.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-307-temporary-redirect",
	},
	"308": {
		code:        "308",
		title:       "Permanent Redirect",
		description: "The target resource has been assigned a new permanent URI and any future references to this resource ought to use one of the enclosed URIs.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-308-permanent-redirect",
	},
	"400": {
		code:        "400",
		title:       "Bad Request",
		description: "The server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing).",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request",
	},
	"401": {
		code:        "401",
		title:       "Not Authorized",
		description: "The request has not been applied because it lacks valid authentication credentials for the target resource.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized",
	},
	"402": {
		code:        "402",
		title:       "Payment Required",
		description: "(402 is reserved for future use.)",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-402-payment-required",
	},
	"403": {
		code:        "403",
		title:       "Forbidden",
		description: "The server understood the request but refuses to fulfill it.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden",
	},
	"404": {
		code:        "404",
		title:       "Not Found",
		description: "The origin server did not find a current representation for the target resource or is not willing to disclose that one exists.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found",
	},
	"405": {
		code:        "405",
		title:       "Method Not Allowed",
		description: "That the method received in the request-line is known by the origin server but not supported by the target resource.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-405-method-not-allowed",
	},
	"406": {
		code:        "406",
		title:       "Not Acceptable",
		description: "The target resource does not have a current representation that would be acceptable to the user agent, according to the proactive negotiation header fields received in the request, and the server is unwilling to supply a default representation.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-406-not-acceptable",
	},
	"407": {
		code:        "407",
		title:       "Proxy Authentication Required",
		description: "The client needs to authenticate itself in order to use a proxy for this request.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-407-proxy-authentication-re",
	},
	"408": {
		code:        "408",
		title:       "Request Timeout",
		description: "The server did not receive a complete request message within the time that it was prepared to wait.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-408-request-timeout",
	},
	"409": {
		code:        "409",
		title:       "Conflict",
		description: "The request could not be completed due to a conflict with the current state of the target resource.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-409-conflict",
	},
	"410": {
		code:        "410",
		title:       "Gone",
		description: "Access to the target resource is no longer available at the origin server and that this condition is likely to be permanent.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-410-gone",
	},
	"411": {
		code:        "411",
		title:       "Length Required",
		description: "The server refuses to accept the request without a defined Content-Length.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-411-length-required",
	},
	"412": {
		code:        "412",
		title:       "Precondition Failed",
		description: "one or more conditions given in the request header fields evaluated to false when tested on the server.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-412-precondition-failed",
	},
	"413": {
		code:        "413",
		title:       "Content Too Large",
		description: "The server is refusing to process a request because the request content is larger than the server is willing or able to process.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-413-content-too-large",
	},
	"414": {
		code:        "414",
		title:       "URI Too Long",
		description: "The server is refusing to service the request because the target URI is longer than the server is willing to interpret.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-414-uri-too-long",
	},
	"415": {
		code:        "415",
		title:       "Unsupported Media Type",
		description: "The origin server is refusing to service the request because the content is in a format not supported by this method on the target resource.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-415-unsupported-media-type",
	},
	"416": {
		code:        "416",
		title:       "Range Not Satisfiable",
		description: "The set of ranges in the request's Range header field has been rejected either because none of the requested ranges are satisfiable or because the client has requested an excessive number of small or overlapping ranges (a potential denial of service attack).",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-416-range-not-satisfiable",
	},
	"417": {
		code:        "417",
		title:       "Expectation Failed",
		description: "The expectation given in the request's Expect header field could not be met by at least one of the inbound servers.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-417-expectation-failed",
	},
	"418": {
		code:        "418",
		title:       "(Unused / I'm a teapot)",
		description: "The server refuses to brew coffee because it is, permanently, a teapot. This error is a reference to Hyper Text Coffee Pot Control Protocol defined in April Fools' jokes in 1998 and 2014.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-418-unused",
	},
	"421": {
		code:        "421",
		title:       "Misdirected Request",
		description: "The request was directed at a server that is unable or unwilling to produce an authoritative response for the target URI.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request",
	},
	"422": {
		code:        "422",
		title:       "Unprocessable Content",
		description: "The server understands the content type of the request content (hence a 415 (Unsupported Media Type) status code is inappropriate), and the syntax of the request content is correct, but it was unable to process the contained instructions.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content",
	},
	"426": {
		code:        "426",
		title:       "Upgrade Required",
		description: "The server refuses to perform the request using the current protocol but might be willing to do so after the client upgrades to a different protocol.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-426-upgrade-required",
	},
	"429": {
		code:        "429",
		title:       "Too Many Requests",
		description: "The user has sent too many requests in a given amount of time (\"rate limiting\").",
		url:         "https://www.rfc-editor.org/rfc/rfc6585#section-4",
	},
	"431": {
		code:        "431",
		title:       "Request Header Fields Too Large",
		description: "The server is unwilling to process the request because its header fields are too large.",
		url:         "https://www.rfc-editor.org/rfc/rfc6585#section-5",
	},
	"451": {
		code:        "451",
		title:       "Unavailable For Legal Reasons",
		description: "The server is denying access to the resource as a consequence of a legal demand.",
		url:         "https://httpwg.org/specs/rfc7725.html#n-451-unavailable-for-legal-reasons",
	},
	"500": {
		code:        "500",
		title:       "Internal Server Error",
		description: "The server encountered an unexpected condition that prevented it from fulfilling the request.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error",
	},
	"501": {
		code:        "501",
		title:       "Not Implemented",
		description: "The server does not support the functionality required to fulfill the request.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-501-not-implemented",
	},
	"502": {
		code:        "502",
		title:       "Bad Gateway",
		description: "The server, while acting as a gateway or proxy, received an invalid response from an inbound server it accessed while attempting to fulfill the request.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway",
	},
	"503": {
		code:        "503",
		title:       "Service Unavailable",
		description: "The server is currently unable to handle the request due to a temporary overload or scheduled maintenance, which will likely be alleviated after some delay.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-503-service-unavailable",
	},
	"504": {
		code:        "504",
		title:       "Gateway Timeout",
		description: "The server, while acting as a gateway or proxy, did not receive a timely response from an upstream server it needed to access in order to complete the request.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-504-gateway-timeout",
	},
	"505": {
		code:        "505",
		title:       "HTTP Version Not Supported",
		description: "The server does not support, or refuses to support, the major version of HTTP that was used in the request message.",
		url:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-505-http-version-not-suppor",
	},
	"506": {
		code:        "506",
		title:       "Variant Also Negotiates",
		description: "The server has an internal configuration error: the chosen variant resource is configured to engage in transparent content negotiation itself, and is therefore not a proper end point in the negotiation process.",
		url:         "https://www.rfc-editor.org/rfc/rfc2295#section-8.1",
	},
	"507": {
		code:        "507",
		title:       "Insufficient Storage",
		description: "The method could not be performed on the resource because the server is unable to store the representation needed to successfully complete the request.",
		url:         "https://www.rfc-editor.org/rfc/rfc4918#section-11.5",
	},
	"508": {
		code:        "508",
		title:       "Loop Detected",
		description: "The server terminated an operation because it encountered an infinite loop while processing a request with \"Depth: infinity\".",
		url:         "https://www.rfc-editor.org/rfc/rfc5842#section-7.2",
	},
	"510": {
		code:        "510",
		title:       "Not Extended",
		description: "The policy for accessing the resource has not been met in the request.",
		url:         "https://www.rfc-editor.org/rfc/rfc2774#section-7",
	},
	"511": {
		code:        "511",
		title:       "Network Authentication Required",
		description: "The client needs to authenticate to gain network access.",
		url:         "https://www.rfc-editor.org/rfc/rfc6585#section-6",
	},
}

const helpVerbose = "Print status code description and URL"
const helpCats = "Open HTTP Status Cats webpage for status"
const helpDogs = "Open HTTP Status Dogs webpage for status"
const httpStatusCatsUrl = "https://http.cat/"
const httpStatusDogsUrl = "https://httpstatusdogs.com/"

func NewStatusCommand(code string) *StatusCommand {
	sc := &StatusCommand{
		fs: flag.NewFlagSet(code, flag.ContinueOnError),
	}

	sc.fs.BoolVar(&sc.verbose, "verbose", false, helpVerbose)
	sc.fs.BoolVar(&sc.verbose, "v", false, helpVerbose)
	sc.fs.BoolVar(&sc.cats, "cats", false, helpCats)
	sc.fs.BoolVar(&sc.cats, "c", false, helpCats)
	sc.fs.BoolVar(&sc.dogs, "dogs", false, helpDogs)
	sc.fs.BoolVar(&sc.dogs, "d", false, helpDogs)

	return sc
}

func (s *StatusCommand) Init(args []string) error {
	return s.fs.Parse(args)
}

func (s *StatusCommand) Run() error {
	status := statuses[s.Name()]

	fmt.Println(status.code + " - " + status.title)

	if s.verbose {
		fmt.Println(status.description)
		fmt.Println(status.url)
	}

	if s.cats {
		browser.OpenURL(fmt.Sprint(httpStatusCatsUrl, status.code))
	}

	if s.dogs {
		browser.OpenURL(fmt.Sprint(httpStatusDogsUrl, status.code))
	}

	return nil
}

func (s *StatusCommand) Name() string {
	return s.fs.Name()
}

// =============================================================================

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("httpwut: Please provide an HTTP status code")
	}

	flag.Usage = usage

	if isAskingForHelp() {
		usage()
		return nil
	}

	subcommand := os.Args[1]

	if status, exists := statuses[subcommand]; exists {
		cmd := NewStatusCommand(status.code)
		cmd.Init(os.Args[2:])
		return cmd.Run()
	}

	return fmt.Errorf("httpwut: '%s' might not be an HTTP status code", subcommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
