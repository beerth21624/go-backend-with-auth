package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"venturex-backend/internal/pkg/validator"
)

type ResponseStatus struct {
	Code        int    `json:"code"`
	Header      string `json:"header"`
	Description string `json:"description"`
}

type ApiResponse[T any] struct {
	Status ResponseStatus `json:"status"`
	Data   T              `json:"data,omitempty"`
}

type ResponseContext interface {
	JSON(code int, obj interface{})
}

func ResponseOk[T any](c ResponseContext, data T) {
	c.JSON(http.StatusOK, ApiResponse[T]{
		Status: ResponseStatus{
			Code:        http.StatusOK,
			Header:      "Success",
			Description: "Request completed successfully",
		},
		Data: data,
	})
}

func ResponseCreated[T any](c ResponseContext, data T) {
	c.JSON(http.StatusCreated, ApiResponse[T]{
		Status: ResponseStatus{
			Code:        http.StatusCreated,
			Header:      "Created",
			Description: "Resource created successfully",
		},
		Data: data,
	})
}

func ResponseNoContent(c ResponseContext) {
	c.JSON(http.StatusNoContent, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusNoContent,
			Header:      "No Content",
			Description: "Request completed successfully with no content",
		},
	})
}

func ResponseBadRequest(c ResponseContext, description string) {
	c.JSON(http.StatusBadRequest, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusBadRequest,
			Header:      "Bad Request",
			Description: description,
		},
	})
}

func ResponseUnauthorized(c ResponseContext, description string) {
	c.JSON(http.StatusUnauthorized, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusUnauthorized,
			Header:      "Unauthorized",
			Description: description,
		},
	})
}

func ResponseForbidden(c ResponseContext, description string) {
	c.JSON(http.StatusForbidden, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusForbidden,
			Header:      "Forbidden",
			Description: description,
		},
	})
}

func ResponseNotFound(c ResponseContext, description string) {
	c.JSON(http.StatusNotFound, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusNotFound,
			Header:      "Not Found",
			Description: description,
		},
	})
}

func ResponseConflict(c ResponseContext, description string) {
	c.JSON(http.StatusConflict, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusConflict,
			Header:      "Conflict",
			Description: description,
		},
	})
}

func ResponseUnprocessableEntity(c ResponseContext, description string) {
	c.JSON(http.StatusUnprocessableEntity, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusUnprocessableEntity,
			Header:      "Unprocessable Entity",
			Description: description,
		},
	})
}

func ResponseInternalError(c ResponseContext, description string) {
	c.JSON(http.StatusInternalServerError, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusInternalServerError,
			Header:      "Internal Server Error",
			Description: description,
		},
	})
}

func ResponseServiceUnavailable(c ResponseContext, description string) {
	c.JSON(http.StatusServiceUnavailable, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusServiceUnavailable,
			Header:      "Service Unavailable",
			Description: description,
		},
	})
}

func ResponseBusinessError(c ResponseContext, description string) {
	c.JSON(http.StatusUnprocessableEntity, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusUnprocessableEntity,
			Header:      "Business Logic Error",
			Description: description,
		},
	})
}

func ResponseValidationError(c ResponseContext, description string) {
	c.JSON(http.StatusBadRequest, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        http.StatusBadRequest,
			Header:      "Validation Error",
			Description: description,
		},
	})
}

type PaginatedData[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func ResponsePaginated[T any](c ResponseContext, items []T, pagination Pagination) {
	data := PaginatedData[T]{
		Items:      items,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedData[T]]{
		Status: ResponseStatus{
			Code:        http.StatusOK,
			Header:      "Success",
			Description: "Paginated data retrieved successfully",
		},
		Data: data,
	})
}

func NewPagination(page, limit int, total int64) Pagination {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	ResponseOk(c, data)
}

func AbortWithError(c *gin.Context, err error) {
	handleError(c, err)
	c.Abort()
}

func handleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, ApiResponse[any]{
			Status: ResponseStatus{
				Code:        appErr.HTTPStatus,
				Header:      string(appErr.Code),
				Description: appErr.Message,
			},
			Data: appErr.Details,
		})
		return
	}

	mapped := MapDomainError(err)

	c.JSON(mapped.HTTPStatus, ApiResponse[any]{
		Status: ResponseStatus{
			Code:        mapped.HTTPStatus,
			Header:      string(mapped.Code),
			Description: mapped.Message,
		},
		Data: mapped.Details,
	})
}

func BindAndValidate[T any](c *gin.Context, obj *T, rules ...validator.FieldRule[T]) *AppError {
	if err := c.ShouldBindJSON(obj); err != nil {
		return NewBadRequestError("Invalid request format: " + err.Error())
	}

	if len(rules) > 0 {
		result := validator.CollectFieldErrors(*obj, rules...)
		if !result.IsValid() {
			return NewValidationError(result.Errors())
		}
	}

	return nil
}

func ValidateStruct[T any](obj T, rules ...validator.FieldRule[T]) *AppError {
	result := validator.CollectFieldErrors(obj, rules...)
	if !result.IsValid() {
		return NewValidationError(result.Errors())
	}
	return nil
}

func BindAndValidateLegacy(c *gin.Context, obj interface{}) *AppError {
	if err := c.ShouldBindJSON(obj); err != nil {
		return NewBadRequestError("Invalid request format")
	}
	return nil
}

func GetPagination(c *gin.Context) (int, int, *AppError) {
	page := 1
	limit := 10

	if p := c.DefaultQuery("page", "1"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.DefaultQuery("limit", "10"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return page, limit, nil
}

func GetUUIDParam(c *gin.Context, param string) (string, *AppError) {
	value := c.Param(param)
	if value == "" {
		return "", NewBadRequestError("Missing required parameter: " + param)
	}
	return value, nil
}
