package middleware

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	DataDogTraceID = "DatadogTraceId"
	DataDogSpanID  = "DatadogSpanId"
)

// DataDogTrace Middleware returns middleware that will trace incoming requests .
func DataDogTrace(service string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			resource := c.Path()

			opts := []ddtrace.StartSpanOption{
				tracer.ServiceName(service),
				tracer.ResourceName(resource),
				tracer.SpanType(ext.SpanTypeWeb),
				tracer.Tag(ext.HTTPMethod, request.Method),
				tracer.Tag(ext.HTTPURL, request.URL.Path),
				tracer.Tag(ext.Environment, os.Getenv("STAGE")),
			}

			if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(request.Header)); err == nil {
				opts = append(opts, tracer.ChildOf(spanctx))
			}

			span, ctx := tracer.StartSpanFromContext(request.Context(), "http.request", opts...)

			defer span.Finish()

			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))
			c.Set(DataDogTraceID, span.Context().TraceID())
			c.Set(DataDogSpanID, span.Context().SpanID())

			// serve the request to the next middleware
			span.SetTag(ext.HTTPCode, strconv.Itoa(c.Response().Status))

			err := next(c)
			if err != nil {
				span.SetTag(ext.Error, err)
			} else if c.Response().Status >= 400 {
				span.SetTag(ext.Error, http.StatusText(c.Response().Status))
			}

			return err
		}
	}
}
