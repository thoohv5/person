package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go/ext"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracingHook struct {
	tracer trace.Tracer
}

var _ redis.Hook = tracingHook{}

// NewTracingHook creates a new go-redis hook instance and that will collect spans using the provided tracer.
func NewTracingHook() redis.Hook {
	return &tracingHook{
		tracer: otel.Tracer("github.com/go-redis/redis"),
	}
}

func (h tracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}
	ctx, span := h.tracer.Start(ctx, cmd.FullName())
	if span == nil {
		return ctx, nil
	}
	span.SetAttributes([]attribute.KeyValue{
		attribute.String("cache.type", "redis"),
		attribute.String("cache.command", cmd.String()),
	}...)
	return ctx, nil
}

func (h tracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}
	defer span.End()

	if err := cmd.Err(); err != nil {
		recordError(ctx, "cache.error", span, err)
	} else {
		span.SetAttributes([]attribute.KeyValue{
			attribute.String("cache.result", cmd.String()),
		}...)
	}
	return nil
}

func (h tracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}
	ctx, span := h.tracer.Start(ctx, "pipeline")
	if span == nil {
		return ctx, nil
	}

	for _, cmd := range cmds {
		span.SetAttributes([]attribute.KeyValue{
			attribute.String("cache.type", "redis"),
			attribute.Int("cache.redis.num_cmd", len(cmds)),
			attribute.String("cache.command", cmd.String()),
		}...)
	}

	return ctx, nil
}

func (h tracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}
	defer span.End()

	for i, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			recordError(ctx, "cache.error"+strconv.Itoa(i), span, err)
		} else {
			span.SetAttributes([]attribute.KeyValue{
				attribute.String("cache.result", cmd.String()),
			}...)
		}
	}
	return nil
}

func recordError(_ context.Context, errorTag string, span trace.Span, err error) {
	if err != redis.Nil {
		span.SetAttributes([]attribute.KeyValue{
			attribute.Bool(string(ext.Error), true),
			attribute.String(errorTag, err.Error()),
		}...)
	}
}
