package jaeger

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func NewSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, spanName)
}

func RecordError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func AddAttributes(span trace.Span, key, value string) {
	attrs := []attribute.KeyValue{attribute.String(key, value)}
	span.SetAttributes(attrs...)
}

func AddSqlEvent(span trace.Span, sql string, values ...interface{}) {
	span.AddEvent("sql_query", trace.WithAttributes(attribute.String("text", sql), attribute.String("values", fmt.Sprint(values))))
}
