package status

// statuses maps each known HTTP status code to its title, description, and
// reference URL. Lookup and All fill in Code from the map key.
var statuses = map[int]Status{
	100: {
		Title:       "Continue",
		Description: "The initial part of a request has been received and has not yet been rejected by the server. The server intends to send a final response after the request has been fully received and acted upon.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-100-continue",
	},
	101: {
		Title:       "Switching Protocols",
		Description: "The server understands and is willing to comply with the client's request, via the Upgrade header field, for a change in the application protocol being used on this connection. The server MUST generate an Upgrade header field in the response that indicates which protocol(s) will be in effect after this response.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-101-switching-protocols",
	},
	102: {
		Title:       "Processing",
		Description: "An interim response used to inform the client that the server has accepted the complete request but has not yet completed it. Deprecated by RFC 4918.",
		URL:         "https://www.rfc-editor.org/rfc/rfc2518#section-10.1",
	},
	103: {
		Title:       "Early Hints",
		Description: "Allows user-agents to perform some operations, such as to speculatively load resources that are likely to be used by the document, before the navigation request is fully handled by the server and a response code is served.",
		URL:         "https://www.rfc-editor.org/rfc/rfc8297",
	},
	200: {
		Title:       "OK",
		Description: "The request has succeeded. The content sent in a 200 response depends on the request method.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-200-ok",
	},
	201: {
		Title:       "Created",
		Description: "The request has been fulfilled and has resulted in one or more new resources being created.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-201-created",
	},
	202: {
		Title:       "Accepted",
		Description: "The request has been accepted for processing, but the processing has not been completed.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-202-accepted",
	},
	203: {
		Title:       "Non-Authoritative Information",
		Description: "The request was successful but the enclosed content has been modified from that of the origin server's 200 (OK) response by a transforming proxy.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-203-non-authoritative-infor",
	},
	204: {
		Title:       "No Content",
		Description: "The server has successfully fulfilled the request and there is no additional content to send in the response content.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-204-no-content",
	},
	205: {
		Title:       "Reset Content",
		Description: "The server has fulfilled the request and desires that the user agent reset the \"document view\", which caused the request to be sent, to its original state as received from the origin server.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-205-reset-content",
	},
	206: {
		Title:       "Partial Content",
		Description: "The server is successfully fulfilling a range request for the target resource by transferring one or more parts of the selected representation.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-206-partial-content",
	},
	207: {
		Title:       "Multi-Status",
		Description: "The message body that follows is by default an XML message and can contain a number of separate response codes, depending on how many sub-requests were made.",
		URL:         "https://www.rfc-editor.org/rfc/rfc4918#section-11.1",
	},
	208: {
		Title:       "Already Reported",
		Description: "The members of a DAV binding have already been enumerated in a preceding part of the multistatus response, and are not being included again.",
		URL:         "https://www.rfc-editor.org/rfc/rfc5842#section-7.1",
	},
	226: {
		Title:       "IM Used",
		Description: "The server has fulfilled a GET request for the resource, and the response is a representation of the result of one or more instance-manipulations applied to the current instance.",
		URL:         "https://www.rfc-editor.org/rfc/rfc3229#section-10.4.1",
	},
	300: {
		Title:       "Multiple Choices",
		Description: "The target resource has more than one representation, each with its own more specific identifier, and information about the alternatives is being provided so that the user (or user agent) can select a preferred representation by redirecting its request to one or more of those identifiers.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-300-multiple-choices",
	},
	301: {
		Title:       "Moved Permanently",
		Description: "The target resource has been assigned a new permanent URI and any future references to this resource ought to use one of the enclosed URIs.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-301-moved-permanently",
	},
	302: {
		Title:       "Found",
		Description: "The target resource resides temporarily under a different URI.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-302-found",
	},
	303: {
		Title:       "See Other",
		Description: "The server is redirecting the user agent to a different resource, as indicated by a URI in the Location header field, which is intended to provide an indirect response to the original request.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-303-see-other",
	},
	304: {
		Title:       "Not Modified",
		Description: "A conditional GET or HEAD request has been received and would have resulted in a 200 (OK) response if it were not for the fact that the condition evaluated to false.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-304-not-modified",
	},
	305: {
		Title:       "Use Proxy",
		Description: "305 was defined in a previous version of HTTP/1.1 and is now deprecated.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-305-use-proxy",
	},
	306: {
		Title:       "(Unused)",
		Description: "306 was defined in a previous version of HTTP/1.1, is no longer used, and the code is reserved.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-306-unused",
	},
	307: {
		Title:       "Temporary Redirect",
		Description: "The target resource resides temporarily under a different URI and the user agent MUST NOT change the request method if it performs an automatic redirection to that URI.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-307-temporary-redirect",
	},
	308: {
		Title:       "Permanent Redirect",
		Description: "The target resource has been assigned a new permanent URI and any future references to this resource ought to use one of the enclosed URIs.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-308-permanent-redirect",
	},
	400: {
		Title:       "Bad Request",
		Description: "The server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing).",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request",
	},
	401: {
		Title:       "Unauthorized",
		Description: "The request has not been applied because it lacks valid authentication credentials for the target resource.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized",
	},
	402: {
		Title:       "Payment Required",
		Description: "(402 is reserved for future use.)",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-402-payment-required",
	},
	403: {
		Title:       "Forbidden",
		Description: "The server understood the request but refuses to fulfill it.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden",
	},
	404: {
		Title:       "Not Found",
		Description: "The origin server did not find a current representation for the target resource or is not willing to disclose that one exists.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found",
	},
	405: {
		Title:       "Method Not Allowed",
		Description: "The method received in the request-line is known by the origin server but not supported by the target resource.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-405-method-not-allowed",
	},
	406: {
		Title:       "Not Acceptable",
		Description: "The target resource does not have a current representation that would be acceptable to the user agent, according to the proactive negotiation header fields received in the request, and the server is unwilling to supply a default representation.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-406-not-acceptable",
	},
	407: {
		Title:       "Proxy Authentication Required",
		Description: "The client needs to authenticate itself in order to use a proxy for this request.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-407-proxy-authentication-re",
	},
	408: {
		Title:       "Request Timeout",
		Description: "The server did not receive a complete request message within the time that it was prepared to wait.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-408-request-timeout",
	},
	409: {
		Title:       "Conflict",
		Description: "The request could not be completed due to a conflict with the current state of the target resource.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-409-conflict",
	},
	410: {
		Title:       "Gone",
		Description: "Access to the target resource is no longer available at the origin server and this condition is likely to be permanent.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-410-gone",
	},
	411: {
		Title:       "Length Required",
		Description: "The server refuses to accept the request without a defined Content-Length.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-411-length-required",
	},
	412: {
		Title:       "Precondition Failed",
		Description: "One or more conditions given in the request header fields evaluated to false when tested on the server.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-412-precondition-failed",
	},
	413: {
		Title:       "Content Too Large",
		Description: "The server is refusing to process a request because the request content is larger than the server is willing or able to process.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-413-content-too-large",
	},
	414: {
		Title:       "URI Too Long",
		Description: "The server is refusing to service the request because the target URI is longer than the server is willing to interpret.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-414-uri-too-long",
	},
	415: {
		Title:       "Unsupported Media Type",
		Description: "The origin server is refusing to service the request because the content is in a format not supported by this method on the target resource.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-415-unsupported-media-type",
	},
	416: {
		Title:       "Range Not Satisfiable",
		Description: "The set of ranges in the request's Range header field has been rejected either because none of the requested ranges are satisfiable or because the client has requested an excessive number of small or overlapping ranges (a potential denial of service attack).",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-416-range-not-satisfiable",
	},
	417: {
		Title:       "Expectation Failed",
		Description: "The expectation given in the request's Expect header field could not be met by at least one of the inbound servers.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-417-expectation-failed",
	},
	418: {
		Title:       "I'm a teapot",
		Description: "The server refuses to brew coffee because it is, permanently, a teapot. This is a reference to the Hyper Text Coffee Pot Control Protocol from the 1998 and 2014 April Fools' RFCs. RFC 9110 lists 418 as unused.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-418-unused",
	},
	421: {
		Title:       "Misdirected Request",
		Description: "The request was directed at a server that is unable or unwilling to produce an authoritative response for the target URI.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request",
	},
	422: {
		Title:       "Unprocessable Content",
		Description: "The server understands the content type of the request content (hence a 415 (Unsupported Media Type) status code is inappropriate), and the syntax of the request content is correct, but it was unable to process the contained instructions.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content",
	},
	423: {
		Title:       "Locked",
		Description: "The source or destination resource of a method is locked.",
		URL:         "https://www.rfc-editor.org/rfc/rfc4918#section-11.3",
	},
	424: {
		Title:       "Failed Dependency",
		Description: "The method could not be performed on the resource because the requested action depended on another action and that action failed.",
		URL:         "https://www.rfc-editor.org/rfc/rfc4918#section-11.4",
	},
	425: {
		Title:       "Too Early",
		Description: "The server is unwilling to risk processing a request that might be replayed.",
		URL:         "https://www.rfc-editor.org/rfc/rfc8470#section-5.2",
	},
	426: {
		Title:       "Upgrade Required",
		Description: "The server refuses to perform the request using the current protocol but might be willing to do so after the client upgrades to a different protocol.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-426-upgrade-required",
	},
	428: {
		Title:       "Precondition Required",
		Description: "The origin server requires the request to be conditional.",
		URL:         "https://www.rfc-editor.org/rfc/rfc6585#section-3",
	},
	429: {
		Title:       "Too Many Requests",
		Description: "The user has sent too many requests in a given amount of time (\"rate limiting\").",
		URL:         "https://www.rfc-editor.org/rfc/rfc6585#section-4",
	},
	431: {
		Title:       "Request Header Fields Too Large",
		Description: "The server is unwilling to process the request because its header fields are too large.",
		URL:         "https://www.rfc-editor.org/rfc/rfc6585#section-5",
	},
	451: {
		Title:       "Unavailable For Legal Reasons",
		Description: "The server is denying access to the resource as a consequence of a legal demand.",
		URL:         "https://www.rfc-editor.org/rfc/rfc7725#section-3",
	},
	500: {
		Title:       "Internal Server Error",
		Description: "The server encountered an unexpected condition that prevented it from fulfilling the request.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error",
	},
	501: {
		Title:       "Not Implemented",
		Description: "The server does not support the functionality required to fulfill the request.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-501-not-implemented",
	},
	502: {
		Title:       "Bad Gateway",
		Description: "The server, while acting as a gateway or proxy, received an invalid response from an inbound server it accessed while attempting to fulfill the request.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway",
	},
	503: {
		Title:       "Service Unavailable",
		Description: "The server is currently unable to handle the request due to a temporary overload or scheduled maintenance, which will likely be alleviated after some delay.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-503-service-unavailable",
	},
	504: {
		Title:       "Gateway Timeout",
		Description: "The server, while acting as a gateway or proxy, did not receive a timely response from an upstream server it needed to access in order to complete the request.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-504-gateway-timeout",
	},
	505: {
		Title:       "HTTP Version Not Supported",
		Description: "The server does not support, or refuses to support, the major version of HTTP that was used in the request message.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-505-http-version-not-suppor",
	},
	506: {
		Title:       "Variant Also Negotiates",
		Description: "The server has an internal configuration error: the chosen variant resource is configured to engage in transparent content negotiation itself, and is therefore not a proper end point in the negotiation process.",
		URL:         "https://www.rfc-editor.org/rfc/rfc2295#section-8.1",
	},
	507: {
		Title:       "Insufficient Storage",
		Description: "The method could not be performed on the resource because the server is unable to store the representation needed to successfully complete the request.",
		URL:         "https://www.rfc-editor.org/rfc/rfc4918#section-11.5",
	},
	508: {
		Title:       "Loop Detected",
		Description: "The server terminated an operation because it encountered an infinite loop while processing a request with \"Depth: infinity\".",
		URL:         "https://www.rfc-editor.org/rfc/rfc5842#section-7.2",
	},
	510: {
		Title:       "Not Extended",
		Description: "The policy for accessing the resource has not been met in the request. This code is obsolete: RFC 2774 has been moved to Historic status.",
		URL:         "https://www.rfc-editor.org/rfc/rfc2774#section-7",
	},
	511: {
		Title:       "Network Authentication Required",
		Description: "The client needs to authenticate to gain network access.",
		URL:         "https://www.rfc-editor.org/rfc/rfc6585#section-6",
	},
}
