package gecko

import (
	"errors"
	"fmt"
	"net/http"
)

func FatalErr(err error) {
	if err != nil {
		panic(err)
	}
}

func FatalFmt(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

// Errors
var (
	ErrBadRequest                    = NewErr(http.StatusBadRequest)                    // HTTP 400 Bad Request
	ErrUnauthorized                  = NewErr(http.StatusUnauthorized)                  // HTTP 401 Unauthorized
	ErrPaymentRequired               = NewErr(http.StatusPaymentRequired)               // HTTP 402 Payment Required
	ErrForbidden                     = NewErr(http.StatusForbidden)                     // HTTP 403 Forbidden
	ErrNotFound                      = NewErr(http.StatusNotFound)                      // HTTP 404 Not Found
	ErrMethodNotAllowed              = NewErr(http.StatusMethodNotAllowed)              // HTTP 405 Method Not Allowed
	ErrNotAcceptable                 = NewErr(http.StatusNotAcceptable)                 // HTTP 406 Not Acceptable
	ErrProxyAuthRequired             = NewErr(http.StatusProxyAuthRequired)             // HTTP 407 Proxy AuthRequired
	ErrRequestTimeout                = NewErr(http.StatusRequestTimeout)                // HTTP 408 Request Timeout
	ErrConflict                      = NewErr(http.StatusConflict)                      // HTTP 409 Conflict
	ErrGone                          = NewErr(http.StatusGone)                          // HTTP 410 Gone
	ErrLengthRequired                = NewErr(http.StatusLengthRequired)                // HTTP 411 Length Required
	ErrPreconditionFailed            = NewErr(http.StatusPreconditionFailed)            // HTTP 412 Precondition Failed
	ErrStatusRequestEntityTooLarge   = NewErr(http.StatusRequestEntityTooLarge)         // HTTP 413 Payload Too Large
	ErrRequestURITooLong             = NewErr(http.StatusRequestURITooLong)             // HTTP 414 URI Too Long
	ErrUnsupportedMediaType          = NewErr(http.StatusUnsupportedMediaType)          // HTTP 415 Unsupported Media Type
	ErrRequestedRangeNotSatisfiable  = NewErr(http.StatusRequestedRangeNotSatisfiable)  // HTTP 416 Range Not Satisfiable
	ErrExpectationFailed             = NewErr(http.StatusExpectationFailed)             // HTTP 417 Expectation Failed
	ErrTeapot                        = NewErr(http.StatusTeapot)                        // HTTP 418 I'm a teapot
	ErrMisdirectedRequest            = NewErr(http.StatusMisdirectedRequest)            // HTTP 421 Misdirected Request
	ErrUnprocessableEntity           = NewErr(http.StatusUnprocessableEntity)           // HTTP 422 Unprocessable Entity
	ErrLocked                        = NewErr(http.StatusLocked)                        // HTTP 423 Locked
	ErrFailedDependency              = NewErr(http.StatusFailedDependency)              // HTTP 424 Failed Dependency
	ErrTooEarly                      = NewErr(http.StatusTooEarly)                      // HTTP 425 Too Early
	ErrUpgradeRequired               = NewErr(http.StatusUpgradeRequired)               // HTTP 426 Upgrade Required
	ErrPreconditionRequired          = NewErr(http.StatusPreconditionRequired)          // HTTP 428 Precondition Required
	ErrTooManyRequests               = NewErr(http.StatusTooManyRequests)               // HTTP 429 Too Many Requests
	ErrRequestHeaderFieldsTooLarge   = NewErr(http.StatusRequestHeaderFieldsTooLarge)   // HTTP 431 Request Header Fields Too Large
	ErrUnavailableForLegalReasons    = NewErr(http.StatusUnavailableForLegalReasons)    // HTTP 451 Unavailable For Legal Reasons
	ErrInternalServerError           = NewErr(http.StatusInternalServerError)           // HTTP 500 Internal Server Error
	ErrNotImplemented                = NewErr(http.StatusNotImplemented)                // HTTP 501 Not Implemented
	ErrBadGateway                    = NewErr(http.StatusBadGateway)                    // HTTP 502 Bad Gateway
	ErrServiceUnavailable            = NewErr(http.StatusServiceUnavailable)            // HTTP 503 Service Unavailable
	ErrGatewayTimeout                = NewErr(http.StatusGatewayTimeout)                // HTTP 504 Gateway Timeout
	ErrHTTPVersionNotSupported       = NewErr(http.StatusHTTPVersionNotSupported)       // HTTP 505 HTTP Version Not Supported
	ErrVariantAlsoNegotiates         = NewErr(http.StatusVariantAlsoNegotiates)         // HTTP 506 Variant Also Negotiates
	ErrInsufficientStorage           = NewErr(http.StatusInsufficientStorage)           // HTTP 507 Insufficient Storage
	ErrLoopDetected                  = NewErr(http.StatusLoopDetected)                  // HTTP 508 Loop Detected
	ErrNotExtended                   = NewErr(http.StatusNotExtended)                   // HTTP 510 Not Extended
	ErrNetworkAuthenticationRequired = NewErr(http.StatusNetworkAuthenticationRequired) // HTTP 511 Network Authentication Required

	ErrValidatorNotRegistered = errors.New("validator not registered")
	ErrRendererNotRegistered  = errors.New("renderer not registered")
	ErrInvalidRedirectCode    = errors.New("invalid redirect status code")
	ErrCookieNotFound         = errors.New("cookie not found")
	ErrInvalidCertOrKeyType   = errors.New("invalid cert or key type, must be string or []byte")
	ErrInvalidListenerNetwork = errors.New("invalid listener network")
)
