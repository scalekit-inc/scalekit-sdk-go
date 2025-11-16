package scalekit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/errdetails"
)

// gRPC to HTTP status mapping
var grpcToHTTP = map[connect.Code]int{
	connect.CodeInvalidArgument:    http.StatusBadRequest,
	connect.CodeFailedPrecondition: http.StatusBadRequest,
	connect.CodeOutOfRange:         http.StatusBadRequest,
	connect.CodeUnauthenticated:    http.StatusUnauthorized,
	connect.CodePermissionDenied:   http.StatusForbidden,
	connect.CodeNotFound:           http.StatusNotFound,
	connect.CodeAlreadyExists:      http.StatusConflict,
	connect.CodeAborted:            http.StatusConflict,
	connect.CodeResourceExhausted:  http.StatusTooManyRequests,
	connect.CodeCanceled:           499, // Client Closed Request
	connect.CodeDataLoss:           http.StatusInternalServerError,
	connect.CodeUnknown:            http.StatusInternalServerError,
	connect.CodeInternal:            http.StatusInternalServerError,
	connect.CodeUnimplemented:       http.StatusNotImplemented,
	connect.CodeUnavailable:         http.StatusServiceUnavailable,
	connect.CodeDeadlineExceeded:   http.StatusGatewayTimeout,
}

// HTTP to gRPC status mapping
var httpToGRPC = map[int]connect.Code{
	http.StatusBadRequest:          connect.CodeInvalidArgument,
	http.StatusUnauthorized:        connect.CodeUnauthenticated,
	http.StatusForbidden:           connect.CodePermissionDenied,
	http.StatusNotFound:            connect.CodeNotFound,
	http.StatusConflict:            connect.CodeAlreadyExists,
	http.StatusTooManyRequests:     connect.CodeResourceExhausted,
	http.StatusInternalServerError: connect.CodeInternal,
	http.StatusNotImplemented:      connect.CodeUnimplemented,
	http.StatusServiceUnavailable: connect.CodeUnavailable,
	http.StatusGatewayTimeout:      connect.CodeDeadlineExceeded,
}

// HTTP status constants
const (
	HTTPStatusOK                  = http.StatusOK
	HTTPStatusBadRequest         = http.StatusBadRequest
	HTTPStatusUnauthorized       = http.StatusUnauthorized
	HTTPStatusForbidden          = http.StatusForbidden
	HTTPStatusNotFound           = http.StatusNotFound
	HTTPStatusConflict           = http.StatusConflict
	HTTPStatusTooManyRequests    = http.StatusTooManyRequests
	HTTPStatusInternalServerError = http.StatusInternalServerError
	HTTPStatusNotImplemented      = http.StatusNotImplemented
	HTTPStatusServiceUnavailable  = http.StatusServiceUnavailable
	HTTPStatusGatewayTimeout      = http.StatusGatewayTimeout
)

// Error types
var (
	ErrRefreshTokenRequired  = errors.New("refresh token is required")
	ErrTokenExpired          = errors.New("token has expired")
	ErrInvalidExpClaimFormat = errors.New("invalid exp claim format")
	ErrAuthRequestIdRequired = errors.New("authRequestId is required")
)

// ScalekitException is the base exception class for all scalekit exceptions
type ScalekitException struct {
	message string
	cause   error
}

func (e *ScalekitException) Error() string {
	if e.message != "" {
		return e.message
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return "Unknown error"
}

func (e *ScalekitException) Unwrap() error {
	return e.cause
}

// WebhookVerificationError is raised for webhook verification failure
type WebhookVerificationError struct {
	ScalekitException
}

func NewWebhookVerificationError(err error) *WebhookVerificationError {
	return &WebhookVerificationError{
		ScalekitException: ScalekitException{
			message: fmt.Sprintf("Webhook verification failed: %v", err),
			cause:   err,
		},
	}
}

// ScalekitValidateTokenFailureException is raised for token validation failure
type ScalekitValidateTokenFailureException struct {
	ScalekitException
}

func NewScalekitValidateTokenFailureException(err error) *ScalekitValidateTokenFailureException {
	return &ScalekitValidateTokenFailureException{
		ScalekitException: ScalekitException{
			message: fmt.Sprintf("Token validation failed: %v", err),
			cause:   err,
		},
	}
}

// ScalekitServerException is the base class for all scalekit server exceptions
type ScalekitServerException struct {
	ScalekitException
	grpcStatus    connect.Code
	httpStatus    int
	errorCode     string
	errDetails    interface{}
	message       string
	unpackedDetails []*errdetails.ErrorInfo
}

func (e *ScalekitServerException) GRPCStatus() connect.Code {
	return e.grpcStatus
}

func (e *ScalekitServerException) HTTPStatus() int {
	return e.httpStatus
}

func (e *ScalekitServerException) ErrorCode() string {
	return e.errorCode
}

func (e *ScalekitServerException) ErrDetails() interface{} {
	return e.errDetails
}

func (e *ScalekitServerException) Message() string {
	return e.message
}

func (e *ScalekitServerException) UnpackedDetails() []*errdetails.ErrorInfo {
	return e.unpackedDetails
}

func (e *ScalekitServerException) Error() string {
	border := strings.Repeat("=", 40)
	
	if len(e.unpackedDetails) > 0 {
		detailsJSON, err := json.MarshalIndent(e.unpackedDetails, "", "  ")
		detailsStr := string(detailsJSON)
		if err != nil {
			detailsStr = fmt.Sprintf("%v", e.unpackedDetails)
		}
		
		// Format the JSON string for better readability
		if strings.HasPrefix(detailsStr, "[") && strings.Contains(detailsStr, "\n") {
			detailsStr = "[\n" + detailsStr[1:]
		}
		
		return fmt.Sprintf("\n%s\nError Code: %s\nGRPC: (%s: %d)\nHTTP: (%s: %d)\nError Details:\n%s: %s\n%s\n",
			border,
			e.errorCode,
			e.getGRPCStatusName(),
			int(e.grpcStatus),
			e.getHTTPStatusName(),
			e.httpStatus,
			e.message,
			detailsStr,
			border,
		)
	}
	
	return fmt.Sprintf("\n%s\nError Code: %s\nGRPC: (%s: %d)\nHTTP: (%s: %d)\nError Details: %v\n%s\n",
		border,
		e.errorCode,
		e.getGRPCStatusName(),
		int(e.grpcStatus),
		e.getHTTPStatusName(),
		e.httpStatus,
		e.errDetails,
		border,
	)
}

func (e *ScalekitServerException) getGRPCStatusName() string {
	switch e.grpcStatus {
	case connect.CodeInvalidArgument:
		return "INVALID_ARGUMENT"
	case connect.CodeFailedPrecondition:
		return "FAILED_PRECONDITION"
	case connect.CodeOutOfRange:
		return "OUT_OF_RANGE"
	case connect.CodeUnauthenticated:
		return "UNAUTHENTICATED"
	case connect.CodePermissionDenied:
		return "PERMISSION_DENIED"
	case connect.CodeNotFound:
		return "NOT_FOUND"
	case connect.CodeAlreadyExists:
		return "ALREADY_EXISTS"
	case connect.CodeAborted:
		return "ABORTED"
	case connect.CodeResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case connect.CodeCanceled:
		return "CANCELED"
	case connect.CodeDataLoss:
		return "DATA_LOSS"
	case connect.CodeUnknown:
		return "UNKNOWN"
	case connect.CodeInternal:
		return "INTERNAL"
	case connect.CodeUnimplemented:
		return "UNIMPLEMENTED"
	case connect.CodeUnavailable:
		return "UNAVAILABLE"
	case connect.CodeDeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	default:
		return "UNKNOWN"
	}
}

func (e *ScalekitServerException) getHTTPStatusName() string {
	switch e.httpStatus {
	case http.StatusOK:
		return "OK"
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusTooManyRequests:
		return "TOO_MANY_REQUESTS"
	case http.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	case http.StatusNotImplemented:
		return "NOT_IMPLEMENTED"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	case http.StatusGatewayTimeout:
		return "GATEWAY_TIMEOUT"
	default:
		return "INTERNAL_SERVER_ERROR"
	}
}

// newScalekitServerException creates a new ScalekitServerException from a connect.Error
func newScalekitServerException(err error) *ScalekitServerException {
	exception := &ScalekitServerException{
		unpackedDetails: []*errdetails.ErrorInfo{},
	}
	
	if connectErr, ok := err.(*connect.Error); ok {
		// Handle gRPC ConnectError
		exception.grpcStatus = connectErr.Code()
		exception.httpStatus = grpcToHTTP[connectErr.Code()]
		if exception.httpStatus == 0 {
			exception.httpStatus = http.StatusInternalServerError
		}
		exception.message = connectErr.Message()
		if exception.message == "" {
			exception.message = "An error occurred"
		}
		exception.cause = err
		
		// Extract error details
		for _, detail := range connectErr.Details() {
			msg, detailErr := detail.Value()
			if detailErr != nil {
				continue
			}
			
			if info, ok := msg.(*errdetails.ErrorInfo); ok {
				exception.unpackedDetails = append(exception.unpackedDetails, info)
				if exception.errorCode == "" && info.ErrorCode != "" {
					exception.errorCode = info.ErrorCode
				}
			}
		}
		
		exception.errDetails = exception.unpackedDetails
		if exception.errorCode == "" {
			exception.errorCode = exception.getGRPCStatusName()
		}
	} else if httpErr, ok := err.(*httpError); ok {
		// Handle HTTP errors
		exception.httpStatus = httpErr.StatusCode
		if exception.httpStatus == 0 {
			exception.httpStatus = http.StatusInternalServerError
		}
		exception.grpcStatus = httpToGRPC[exception.httpStatus]
		if exception.grpcStatus == 0 {
			exception.grpcStatus = connect.CodeUnknown
		}
		exception.message = httpErr.Error()
		if exception.message == "" {
			exception.message = "An HTTP error occurred"
		}
		exception.cause = err
		exception.errorCode = exception.getHTTPStatusName()
		exception.errDetails = httpErr.Error()
	} else {
		// Handle generic errors
		exception.grpcStatus = connect.CodeUnknown
		exception.httpStatus = http.StatusInternalServerError
		exception.message = err.Error()
		if exception.message == "" {
			exception.message = "An unknown error occurred"
		}
		exception.cause = err
		exception.errorCode = "UNKNOWN"
		exception.errDetails = err.Error()
	}
	
	return exception
}

// Promote converts a connect.Error or HTTP error to a specific ScalekitServerException type
func PromoteError(err error) error {
	if err == nil {
		return nil
	}
	
	baseException := newScalekitServerException(err)
	grpcStatus := baseException.grpcStatus
	
	switch grpcStatus {
	case connect.CodeInvalidArgument, connect.CodeFailedPrecondition, connect.CodeOutOfRange:
		return &ScalekitBadRequestException{*baseException}
	case connect.CodeUnauthenticated:
		return &ScalekitUnauthorizedException{*baseException}
	case connect.CodePermissionDenied:
		return &ScalekitForbiddenException{*baseException}
	case connect.CodeNotFound:
		return &ScalekitNotFoundException{*baseException}
	case connect.CodeAlreadyExists, connect.CodeAborted:
		return &ScalekitConflictException{*baseException}
	case connect.CodeResourceExhausted:
		return &ScalekitTooManyRequestsException{*baseException}
	case connect.CodeCanceled:
		return &ScalekitCancelledException{*baseException}
	case connect.CodeDataLoss, connect.CodeUnknown, connect.CodeInternal:
		return &ScalekitInternalServerException{*baseException}
	case connect.CodeUnimplemented:
		return &ScalekitNotImplementedException{*baseException}
	case connect.CodeUnavailable:
		return &ScalekitServiceUnavailableException{*baseException}
	case connect.CodeDeadlineExceeded:
		return &ScalekitGatewayTimeoutException{*baseException}
	default:
		return &ScalekitUnknownException{*baseException}
	}
}

// Specific exception types

// ScalekitBadRequestException is raised for bad requests
type ScalekitBadRequestException struct {
	ScalekitServerException
}

// ScalekitUnauthorizedException is raised for unauthorized access
type ScalekitUnauthorizedException struct {
	ScalekitServerException
}

// ScalekitForbiddenException is raised for forbidden access
type ScalekitForbiddenException struct {
	ScalekitServerException
}

// ScalekitNotFoundException is raised when a resource is not found
type ScalekitNotFoundException struct {
	ScalekitServerException
}

// ScalekitConflictException is raised for conflicts, such as duplicate resources
type ScalekitConflictException struct {
	ScalekitServerException
}

// ScalekitTooManyRequestsException is raised when too many requests are made in a short time
type ScalekitTooManyRequestsException struct {
	ScalekitServerException
}

// ScalekitInternalServerException is raised for internal server errors
type ScalekitInternalServerException struct {
	ScalekitServerException
}

// ScalekitNotImplementedException is raised when a feature is not implemented
type ScalekitNotImplementedException struct {
	ScalekitServerException
}

// ScalekitServiceUnavailableException is raised when the service is unavailable
type ScalekitServiceUnavailableException struct {
	ScalekitServerException
}

// ScalekitGatewayTimeoutException is raised when a gateway timeout occurs
type ScalekitGatewayTimeoutException struct {
	ScalekitServerException
}

// ScalekitCancelledException is raised when an operation is cancelled
type ScalekitCancelledException struct {
	ScalekitServerException
}

// ScalekitUnknownException is raised for unknown errors
type ScalekitUnknownException struct {
	ScalekitServerException
}
