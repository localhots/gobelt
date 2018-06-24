package log

import "context"

type logContext byte

const ctxFields logContext = iota

// ContextWithFields returns a new context with given fields added.
func ContextWithFields(ctx context.Context, f F) context.Context {
	ctxf, ok := ctx.Value(ctxFields).(F)
	if !ok {
		ctxf = F{}
	}
	for k, v := range f {
		ctxf[k] = v
	}
	return context.WithValue(ctx, ctxFields, ctxf)
}

func contextFields(ctx context.Context) F {
	f, ok := ctx.Value(ctxFields).(F)
	if ok {
		return f
	}
	return F{}
}
