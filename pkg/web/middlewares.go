package web

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Iiqbal2000/bareknews"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Middleware func(Handler) Handler

func SetMiddlewares(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}
	return handler
}

func ContentTypeJSON() Middleware {
	m := func(next Handler) Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("content-type", "application/json;charset=utf-8")
			return next(ctx, w, r)
		}
		return h
	}
	return m
}

func CORS() Middleware {
	m := func(next Handler) Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			return next(ctx, w, r)
		}

		return h
	}

	return m
}

func Errors(log *zap.SugaredLogger) Middleware {
	return func(handler Handler) Handler {
		
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var errResp ErrorResponse
			var status int

			if err := handler(ctx, w, r); err != nil {
				switch errors.Cause(err).(type) {
				case validation.Errors, validation.Error:
					status = http.StatusBadRequest
					// If the error is validation.Errors we want to return it
					// as [fieldName]: error message.
					if ve, ok := err.(validation.Errors); ok {
						errResp = ErrorResponse{
							Error:  "invalid data",
							Fields: make(map[string]interface{}),
						}

						for k, v := range ve {
							errResp.Fields[k] = v.Error()
						}

					} else {
						errResp = ErrorResponse{
							Error: errors.Cause(err).Error(),
						}
					}
				case *RequestError:
					reqErr := GetRequestError(err)
					status = reqErr.Status
					errResp = ErrorResponse{
						Error: reqErr.Error(),
					}
				default:
					log.Errorw(fmt.Sprint(err.Error()))
					
					status = http.StatusInternalServerError
					errResp = ErrorResponse{
						Error: bareknews.ErrInternalServer.Error(),
					}
				}

				if err := Respond(w, errResp, status); err != nil {
					return err
				}
			}

			return nil
		}

		return h
	}
}

func Panics() Middleware {
	return func(next Handler) Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rc := recover(); rc != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rc, string(trace))
				}
			}()
			return next(ctx, w, r)
		}
		return h
	}
}