package retry

import (
	"context"
	"time"
)

// Retryer handles the retry logic
type Retryer struct {
	Ctx   context.Context
	F     Func
	Max   int
	ran   int
	Delay *time.Duration
}

// Retrier is a type alias for Retryer
type Retrier = Retryer

// Run runs the function to be retried. If the function returns an error, it will be called again
// unless it has reached the max amount of retries, or the context has expired
func (r Retryer) Run() (interface{}, error) {
	res, err := r.F()
	if err != nil {
		if ctxErr := r.Ctx.Err(); ctxErr != nil {
			return res, err
		}
		if r.ran == r.Max {
			return res, err
		}
		if r.Delay != nil {
			time.Sleep(*r.Delay)
		}
		r.ran++
		return r.Run()
	}
	return res, err
}

// Func is a retryable function
type Func func() (interface{}, error)

// Once retries F once if it returns an error
func Once(f Func) (interface{}, error) {
	return OnceCtx(context.Background(), f)
}

// OnceCtx retries function F once, if it returns an error.
// Returns if the context is done.
func OnceCtx(ctx context.Context, f Func) (interface{}, error) {
	r := Retryer{F: f, Max: 1, Ctx: ctx}
	return r.Run()
}

// OnceDelayed retries F once, after the delay if it returns an error.
func OnceDelayed(f Func, delay time.Duration) (interface{}, error) {
	return OnceDelayedCtx(context.Background(), f, delay)
}

// OnceDelayedCtx retries function F once, after the delay if it returns an error.
// Returns if the context is done.
func OnceDelayedCtx(ctx context.Context, f Func, delay time.Duration) (interface{}, error) {
	r := Retryer{F: f, Max: 1, Ctx: ctx, Delay: &delay}
	return r.Run()
}

// Times retries function F, x amount of times if it returns an error
func Times(amount int, f Func) (interface{}, error) {
	return TimesCtx(context.Background(), amount, f)
}

// TimesCtx retries function F, x amount of times if it returns an error.
// Returns if the context is done.
func TimesCtx(ctx context.Context, amount int, f Func) (interface{}, error) {
	r := Retryer{F: f, Max: amount, Ctx: ctx}
	return r.Run()
}

// TimesDelayed retries function F, x amount of times with a delay, if it returns an error
func TimesDelayed(amount int, f Func, delay time.Duration) (interface{}, error) {
	return TimesDelayedCtx(context.Background(), amount, f, delay)
}

// TimesDelayedCtx retries function f an x amount of times with a delay, if it returns an error.
// Returns if the context is done.
func TimesDelayedCtx(ctx context.Context, amount int, f Func, delay time.Duration) (interface{}, error) {
	r := Retryer{F: f, Max: amount, Ctx: ctx, Delay: &delay}
	return r.Run()
}
