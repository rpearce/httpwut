// https://www.rfc-editor.org/rfc/rfc9110.html
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var statusMessage, statusLong, statusUrl string

	errorCount := 0
	cmd := os.Args[1]

	flagCmd := flag.NewFlagSet("status", flag.ExitOnError)
	verbosePtr := flagCmd.Bool("v", false, "verbose output")
	flagCmd.Parse(os.Args[2:])

	switch cmd {
	case "100":
		statusMessage = "Continue"
		statusLong = "The initial part of a request has been received and has not yet been rejected by the server. The server intends to send a final response after the request has been fully received and acted upon."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-100-continue"
	case "101":
		statusMessage = "Switching Protocols"
		statusLong = "The server understands and is willing to comply with the client's request, via the Upgrade header field, for a change in the application protocol being used on this connection. The server MUST generate an Upgrade header field in the response that indicates which protocol(s) will be in effect after this response."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-101-switching-protocols"
	case "200":
		statusMessage = "OK"
		statusLong = "The request has succeeded. The content sent in a 200 response depends on the request method."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-200-ok"
	case "201":
		statusMessage = "Created"
		statusLong = "The request has been fulfilled and has resulted in one or more new resources being created."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-201-created"
	case "202":
		statusMessage = "Accepted"
		statusLong = "The request has been accepted for processing, but the processing has not been completed."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-202-accepted"
	case "203":
		statusMessage = "Non-Authoritative Information"
		statusLong = "The request was successful but the enclosed content has been modified from that of the origin server's 200 (OK) response by a transforming proxy."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-203-non-authoritative-infor"
	case "204":
		statusMessage = "No Content"
		statusLong = "The server has successfully fulfilled the request and that there is no additional content to send in the response content."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-204-no-content"
	case "205":
		statusMessage = "Reset Content"
		statusLong = "The server has fulfilled the request and desires that the user agent reset the \"document view\", which caused the request to be sent, to its original state as received from the origin server."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-205-reset-content"
	case "206":
		statusMessage = "Partial Content"
		statusLong = "The server is successfully fulfilling a range request for the target resource by transferring one or more parts of the selected representation."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-206-partial-content"
	case "300":
		statusMessage = "Multiple Choices"
		statusLong = "The target resource has more than one representation, each with its own more specific identifier, and information about the alternatives is being provided so that the user (or user agent) can select a preferred representation by redirecting its request to one or more of those identifiers."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-300-multiple-choices"
	case "301":
		statusMessage = "Moved Permanently"
		statusLong = "The target resource has been assigned a new permanent URI and any future references to this resource ought to use one of the enclosed URIs."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-301-moved-permanently"
	case "302":
		statusMessage = "Found"
		statusLong = "The target resource resides temporarily under a different URI."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-302-found"
	case "303":
		statusMessage = "See Other"
		statusLong = "The server is redirecting the user agent to a different resource, as indicated by a URI in the Location header field, which is intended to provide an indirect response to the original request."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-303-see-other"
	case "304":
		statusMessage = "Not Modified"
		statusLong = "A conditional GET or HEAD request has been received and would have resulted in a 200 (OK) response if it were not for the fact that the condition evaluated to false."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-304-not-modified"
	case "305":
		statusMessage = "Use Proxy"
		statusLong = "(305 was defined in a previous version of this specification and is now deprecated)"
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-305-use-proxy"
	case "306":
		statusMessage = "(Unused)"
		statusLong = "(306 was defined in a previous version of this specification, is no longer used, and the code is reserved.)"
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-306-unused"
	case "307":
		statusMessage = "Temporary Redirect"
		statusLong = "The target resource resides temporarily under a different URI and the user agent MUST NOT change the request method if it performs an automatic redirection to that URI."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-307-temporary-redirect"
	case "308":
		statusMessage = "Permanent Redirect"
		statusLong = "The target resource has been assigned a new permanent URI and any future references to this resource ought to use one of the enclosed URIs."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-308-permanent-redirect"
	case "400":
		statusMessage = "Bad Request"
		statusLong = "The server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request"
	case "401":
		statusMessage = "Not Authorized"
		statusLong = "The request has not been applied because it lacks valid authentication credentials for the target resource."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized"
	case "402":
		statusMessage = "Payment Required"
		statusLong = "(402 is reserved for future use.)"
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-402-payment-required"
	case "403":
		statusMessage = "Forbidden"
		statusLong = "The server understood the request but refuses to fulfill it."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden"
	case "404":
		statusMessage = "Not Found"
		statusLong = "The origin server did not find a current representation for the target resource or is not willing to disclose that one exists."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found"
	case "405":
		statusMessage = "Method Not Allowed"
		statusLong = "That the method received in the request-line is known by the origin server but not supported by the target resource."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-405-method-not-allowed"
	case "406":
		statusMessage = "Not Acceptable"
		statusLong = "The target resource does not have a current representation that would be acceptable to the user agent, according to the proactive negotiation header fields received in the request, and the server is unwilling to supply a default representation."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-406-not-acceptable"
	case "407":
		statusMessage = "Proxy Authentication Required"
		statusLong = "The client needs to authenticate itself in order to use a proxy for this request."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-407-proxy-authentication-re"
	case "408":
		statusMessage = "Request Timeout"
		statusLong = "The server did not receive a complete request message within the time that it was prepared to wait."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-408-request-timeout"
	case "409":
		statusMessage = "Conflict"
		statusLong = "The request could not be completed due to a conflict with the current state of the target resource."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-409-conflict"
	case "410":
		statusMessage = "Gone"
		statusLong = "Access to the target resource is no longer available at the origin server and that this condition is likely to be permanent."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-410-gone"
	case "411":
		statusMessage = "Length Required"
		statusLong = "The server refuses to accept the request without a defined Content-Length."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-411-length-required"
	case "412":
		statusMessage = "Precondition Failed"
		statusLong = "one or more conditions given in the request header fields evaluated to false when tested on the server."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-412-precondition-failed"
	case "413":
		statusMessage = "Content Too Large"
		statusLong = "The server is refusing to process a request because the request content is larger than the server is willing or able to process."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-413-content-too-large"
	case "414":
		statusMessage = "URI Too Long"
		statusLong = "The server is refusing to service the request because the target URI is longer than the server is willing to interpret."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-414-uri-too-long"
	case "415":
		statusMessage = "Unsupported Media Type"
		statusLong = "The origin server is refusing to service the request because the content is in a format not supported by this method on the target resource."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-415-unsupported-media-type"
	case "416":
		statusMessage = "Range Not Satisfiable"
		statusLong = "The set of ranges in the request's Range header field has been rejected either because none of the requested ranges are satisfiable or because the client has requested an excessive number of small or overlapping ranges (a potential denial of service attack)."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-416-range-not-satisfiable"
	case "417":
		statusMessage = "Expectation Failed"
		statusLong = "The expectation given in the request's Expect header field could not be met by at least one of the inbound servers."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-417-expectation-failed"
	case "418":
		statusMessage = "(Unused / I'm a teapot)"
		statusLong = "The server refuses to brew coffee because it is, permanently, a teapot. This error is a reference to Hyper Text Coffee Pot Control Protocol defined in April Fools' jokes in 1998 and 2014."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-418-unused"
	case "421":
		statusMessage = "Misdirected Request"
		statusLong = "The request was directed at a server that is unable or unwilling to produce an authoritative response for the target URI."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request"
	case "422":
		statusMessage = "Unprocessable Content"
		statusLong = "The server understands the content type of the request content (hence a 415 (Unsupported Media Type) status code is inappropriate), and the syntax of the request content is correct, but it was unable to process the contained instructions."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content"
	case "426":
		statusMessage = "Upgrade Required"
		statusLong = "The server refuses to perform the request using the current protocol but might be willing to do so after the client upgrades to a different protocol."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-426-upgrade-required"
	case "500":
		statusMessage = "Internal Server Error"
		statusLong = "The server encountered an unexpected condition that prevented it from fulfilling the request."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error"
	case "501":
		statusMessage = "Not Implemented"
		statusLong = "The server does not support the functionality required to fulfill the request."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-501-not-implemented"
	case "502":
		statusMessage = "Bad Gateway"
		statusLong = "The server, while acting as a gateway or proxy, received an invalid response from an inbound server it accessed while attempting to fulfill the request."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway"
	case "503":
		statusMessage = "Service Unavailable"
		statusLong = "The server is currently unable to handle the request due to a temporary overload or scheduled maintenance, which will likely be alleviated after some delay."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-503-service-unavailable"
	case "504":
		statusMessage = "Gateway Timeout"
		statusLong = "The server, while acting as a gateway or proxy, did not receive a timely response from an upstream server it needed to access in order to complete the request."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-504-gateway-timeout"
	case "505":
		statusMessage = "HTTP Version Not Supported"
		statusLong = "The server does not support, or refuses to support, the major version of HTTP that was used in the request message."
		statusUrl = "https://www.rfc-editor.org/rfc/rfc9110.html#name-505-http-version-not-suppor"
	default:
		statusMessage = "Error: Unknown status code"
		errorCount++
	}

	fmt.Println(statusMessage)

	if *verbosePtr && statusLong != "" {
		fmt.Println(statusLong)
		fmt.Println(statusUrl)
	}

	os.Exit(errorCount)
}
