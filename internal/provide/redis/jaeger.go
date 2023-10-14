package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go/ext"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TracingHook struct {
	tracer trace.Tracer
}

var _ redis.Hook = TracingHook{}

// NewHook creates a new go-redis hook instance and that will collect spans using the provided tracer.
func NewHook(tracer trace.Tracer) redis.Hook {
	return &TracingHook{
		tracer: tracer,
	}
}

func (hook TracingHook) createSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}
	return hook.tracer.Start(ctx, spanName)
}

func (hook TracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx, span := hook.createSpan(ctx, cmd.FullName())
	if span == nil {
		return ctx, nil
	}
	span.SetAttributes([]attribute.KeyValue{
		attribute.String("cache.type", "redis"),
		attribute.String("cache.command", cmd.String()),
	}...)
	return ctx, nil
}

func (hook TracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}
	defer span.End()

	if err := cmd.Err(); err != nil {
		recordError(ctx, "cache.error", span, err)
	} else {
		span.SetAttributes([]attribute.KeyValue{
			attribute.String("cache.result", cmd.(*redis.Cmd).String()),
		}...)
	}
	return nil
}

func (hook TracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	ctx, span := hook.createSpan(ctx, "pipeline")
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

func (hook TracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
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
				attribute.String("cache.result", cmd.(*redis.Cmd).String()),
			}...)
		}
	}
	return nil
}

func recordError(ctx context.Context, errorTag string, span trace.Span, err error) {
	if err != redis.Nil {
		span.SetAttributes([]attribute.KeyValue{
			attribute.Bool(string(ext.Error), true),
			attribute.String(errorTag, err.Error()),
		}...)
	}
}
