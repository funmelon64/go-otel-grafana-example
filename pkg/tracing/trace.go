package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"otel-jaeger-learn/pkg/logging"
)

type Span struct {
	span trace.Span
}

func (s *Span) End() {
	s.span.End()
}

func (s *Span) AddEvent(name string, attrs ...slog.Attr) {
	otelAttrs := convertSlogAttrToOtelAttr(attrs)
	s.span.AddEvent(name, trace.WithAttributes(otelAttrs...))
}

func (s *Span) AddError(msg string, err error, attrs ...slog.Attr) {
	s.span.SetStatus(codes.Error, err.Error())
	otelAttrs := convertSlogAttrToOtelAttr(attrs)
	otelAttrs = append(otelAttrs, attribute.String("error", err.Error()))
	s.span.AddEvent(msg, trace.WithAttributes(otelAttrs...))
}

func convertSlogAttrToOtelAttr(attrs []slog.Attr) []attribute.KeyValue {
	var otelAttrs []attribute.KeyValue

	for _, attr := range attrs {
		otelAttr := attribute.String(attr.Key, attr.String())
		otelAttrs = append(otelAttrs, otelAttr)
	}

	return otelAttrs
}

func TraceError(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetStatus(codes.Error, err.Error())
		otelAttrs := convertSlogAttrToOtelAttr(attrs)
		otelAttrs = append(otelAttrs, attribute.String("error", err.Error()))
		span.AddEvent(msg, trace.WithAttributes(otelAttrs...))
	} else {
		attrs = append(attrs, slog.String("msg", msg))
		logging.ErrorErr("TraceError", err, attrs...)
	}
}

// TraceEvent logs event to span if there is span in context
// otherwise it logs event to default logger
func TraceEvent(ctx context.Context, name string, attrs ...slog.Attr) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		otelAttrs := convertSlogAttrToOtelAttr(attrs)
		span.AddEvent(name, trace.WithAttributes(otelAttrs...))
	} else {
		attrs = append(attrs, slog.String("name", name))
		logging.Info("TraceEvent", attrs...)
	}
}

// NewSpan returns Span that must be closed by Span.End() or memory leaks can occur
func NewSpan(ctx context.Context, name string) (context.Context, Span) {
	newCtx, span := otel.Tracer("").Start(ctx, name)
	return newCtx, Span{span: span}
}

// TraceLogger returns logger with traceID if there is span in context
// otherwise it returns default logger
func TraceLogger(ctx context.Context) *logging.Logger {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return logging.GetDefault()
	}
	return logging.NewWith(slog.String("traceID", span.SpanContext().TraceID().String()))
}
