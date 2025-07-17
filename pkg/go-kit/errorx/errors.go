package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/1612421/cinema-booking/pkg/go-kit/log"
)

type IsErrorDetail interface {
	isErrorDetail()
}

type ErrorWrapper struct {
	ErrorBody ErrorBody `json:"error"`
}

type ErrorBody struct {
	Error   string       `json:"error"`
	Code    int32        `json:"code"`
	Details *ErrorDetail `json:"-"`
}

type ErrorDetail struct {
	ErrorInfo           *ErrorInfo           `json:"error_info,omitempty"`
	LocalizedMessage    *LocalizedMessage    `json:"localized_message,omitempty"`
	BadRequest          *BadRequest          `json:"bad_request,omitempty"`
	PreconditionFailure *PreconditionFailure `json:"precondition_failure,omitempty"`
	ResourceInfo        *ResourceInfo        `json:"resource_info,omitempty"`
	QuotaFailure        *QuotaFailure        `json:"quota_failure,omitempty"`
	DebugInfo           *DebugInfo           `json:"debug_info,omitempty"`
	Help                *Help                `json:"help,omitempty"`
}

func (e *ErrorWrapper) Error() string {
	return e.ErrorBody.Error
}

func New(code int32, message string, details ...IsErrorDetail) *ErrorWrapper {
	errResult := &ErrorWrapper{
		ErrorBody: ErrorBody{
			Error: message,
			Code:  code,
		},
	}
	var err error
	errResult.ErrorBody.Details, err = transformErrDetails(details...)
	if err != nil {
		errResult.ErrorBody.Details = nil
	}
	return errResult
}

func ParseValidateDetails(err error) *BadRequest {
	var validateErr validator.ValidationErrors
	ok := errors.As(err, &validateErr)
	if !ok || validateErr == nil {
		return nil
	}

	result := &BadRequest{}
	for _, e := range validateErr {
		description := fmt.Sprintf("Key: '%s' failed on the '%s' tag.", e.Field(), e.Tag())
		if len(e.Param()) > 0 {
			description += " Accepted values: [" + strings.Join(strings.Split(e.Param(), " "), ", ") + "]"
		}
		result.FieldViolations = append(
			result.FieldViolations, &BadRequestFieldViolation{
				Field:       e.Field(),
				Description: description,
			},
		)
	}
	return result
}

type ErrorInfo struct {
	Reason   string         `json:"reason,omitempty"`
	Domain   string         `json:"domain,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (e *ErrorInfo) isErrorDetail() {}

type LocalizedMessage struct {
	Locale  string `json:"locale,omitempty"`
	Message string `json:"message,omitempty"`
}

func (l *LocalizedMessage) isErrorDetail() {}

type BadRequest struct {
	FieldViolations []*BadRequestFieldViolation `json:"field_violations,omitempty"`
}

func (b *BadRequest) isErrorDetail() {}

type BadRequestFieldViolation struct {
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
}

type PreconditionFailure struct {
	Violations []*PreconditionFailureViolation `json:"violations,omitempty"`
}

func (p *PreconditionFailure) isErrorDetail() {}

type PreconditionFailureViolation struct {
	Type        string `json:"type,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *PreconditionFailureViolation) isErrorDetail() {}

type ResourceInfo struct {
	ResourceType string `protobuf:"bytes,1,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	ResourceName string `protobuf:"bytes,2,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	Owner        string `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	Description  string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
}

func (r *ResourceInfo) isErrorDetail() {}

type QuotaFailure struct {
	Violations []*QuotaFailureViolation `protobuf:"bytes,1,rep,name=violations,proto3" json:"violations,omitempty"`
}

func (q *QuotaFailure) isErrorDetail() {}

type QuotaFailureViolation struct {
	Subject     string `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

type DebugInfo struct {
	StackEntries []string `protobuf:"bytes,1,rep,name=stack_entries,json=stackEntries,proto3" json:"stack_entries,omitempty"`
	Detail       string   `protobuf:"bytes,2,opt,name=detail,proto3" json:"detail,omitempty"`
}

type Help struct {
	Links []*HelpLink `json:"links,omitempty"`
}

func (h Help) isErrorDetail() {}

type HelpLink struct {
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
}

func (d *DebugInfo) isErrorDetail() {}

var (
	_ IsErrorDetail = (*ErrorInfo)(nil)
	_ IsErrorDetail = (*LocalizedMessage)(nil)
	_ IsErrorDetail = (*BadRequest)(nil)
	_ IsErrorDetail = (*PreconditionFailure)(nil)
	_ IsErrorDetail = (*ResourceInfo)(nil)
	_ IsErrorDetail = (*QuotaFailure)(nil)
	_ IsErrorDetail = (*DebugInfo)(nil)
	_ IsErrorDetail = (*Help)(nil)
)

func GetGRPCCode(err error) int32 {
	s, ok := status.FromError(err)
	if !ok {
		return int32(codes.Unknown)
	}
	return int32(s.Code()) //nolint:gosec
}

func GetHTTPCode(err error) int {
	var (
		errWrapper *ErrorWrapper
		statusCode = http.StatusInternalServerError
	)
	ok := errors.As(err, &errWrapper)
	if ok {
		if errWrapper.ErrorBody.Code >= http.StatusOK {
			statusCode = int(errWrapper.ErrorBody.Code)
		} else {
			statusCode = runtime.HTTPStatusFromCode(codes.Code(errWrapper.ErrorBody.Code)) //nolint:gosec
		}
	}
	return statusCode
}

// AttachContextError attaches an error to the current context. The error is pushed to a list of errors.
func AttachContextError(ctx *gin.Context, logger log.Logger, err error) {
	// This function return error wrapper so that not need to log error again. This log for avoiding lint error
	ctxErr := ctx.Error(err)
	if ctxErr != nil {
		logger.Debug("Attach error into context", zap.Error(err))
	}
}

func convertErrDetails(errDetailItem IsErrorDetail, errDetail *ErrorDetail) {
	switch errDetailItem.(type) {
	case *ErrorInfo:
		errDetail.ErrorInfo = errDetailItem.(*ErrorInfo)
	case *LocalizedMessage:
		errDetail.LocalizedMessage = errDetailItem.(*LocalizedMessage)
	case *BadRequest:
		errDetail.BadRequest = errDetailItem.(*BadRequest)
	case *PreconditionFailure:
		errDetail.PreconditionFailure = errDetailItem.(*PreconditionFailure)
	case *ResourceInfo:
		errDetail.ResourceInfo = errDetailItem.(*ResourceInfo)
	case *QuotaFailure:
		errDetail.QuotaFailure = errDetailItem.(*QuotaFailure)
	case *DebugInfo:
		errDetail.DebugInfo = errDetailItem.(*DebugInfo)
	case *Help:
		errDetail.Help = errDetailItem.(*Help)
	}
}

func transformErrDetails(details ...IsErrorDetail) (*ErrorDetail, error) {
	result := &ErrorDetail{}

	for _, detail := range details {
		convertErrDetails(detail, result)
	}

	return result, nil
}

var (
	ErrorInternal         = New(int32(codes.Internal), "Internal server error")
	ErrorIsNotInWhitelist = New(int32(codes.PermissionDenied), "Not in whitelist")
)
