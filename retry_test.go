package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/semlette/retry"
)

func TestOnce_ReturnsIfSuccessful(t *testing.T) {
	t.Parallel()
	var times = 0
	f := func() (interface{}, error) {
		times++
		return "value", nil
	}
	value, err := retry.Once(f)
	if value != "value" {
		t.Error("value is not the value returned by the function")
	}
	if err != nil {
		t.Errorf("err should be nil: %s", err.Error())
	}
	if times != 1 {
		t.Errorf("F() was not called exactly once, instead: %d", times)
	}
}

func TestOnce_RunsAgainIfError(t *testing.T) {
	t.Parallel()
	var times = 0
	e := errors.New("error")
	f := func() (interface{}, error) {
		times++
		return "error", e
	}
	_, err := retry.Once(f)
	if err != e {
		t.Errorf("the returned error was not the one returned by F(), instead: %s", err.Error())
	}
	if times != 2 {
		t.Errorf("F() was not called exactly 2 times, instead: %d", times)
	}
}

func TestOnceCtx_ReturnsIfSuccessful(t *testing.T) {
	t.Parallel()
	var times = 0
	f := func() (interface{}, error) {
		times++
		return "value", nil
	}
	value, err := retry.OnceCtx(context.Background(), f)
	if value != "value" {
		t.Error("value is not the value returned by the function")
	}
	if err != nil {
		t.Errorf("err should be nil: %s", err.Error())
	}
	if times != 1 {
		t.Errorf("F() was not called exactly once, instead: %d", times)
	}
}

func TestOnceCtx_RunsAgainIfError(t *testing.T) {
	t.Parallel()
	var times = 0
	e := errors.New("error")
	f := func() (interface{}, error) {
		times++
		return "error", e
	}
	_, err := retry.OnceCtx(context.Background(), f)
	if err != e {
		t.Errorf("the returned error was not the one returned by F(), instead: %s", err.Error())
	}
	if times != 2 {
		t.Errorf("F() was not called exactly 2 times, instead: %d", times)
	}
}

func TestOnceCtx_StopsIfContextCancelled(t *testing.T) {
	t.Parallel()
	var times = 0
	ctx, cancel := context.WithCancel(context.Background())
	e := errors.New("error")
	f := func() (interface{}, error) {
		times++
		cancel()
		return "error", e
	}
	_, err := retry.OnceCtx(ctx, f)
	if err != e {
		t.Errorf("the returned error was not the one returned by F(), instead: %s", err.Error())
	}
	if times != 1 {
		t.Errorf("F() was not called exactly once, instead: %d", times)
	}
}

func TestOnceDelayed(t *testing.T) {
	t.Parallel()
	var times = 0
	e := errors.New("error")
	f := func() (interface{}, error) {
		times++
		return "error", e
	}
	start := time.Now()
	value, err := retry.OnceDelayed(f, 1*time.Second)
	if value != "error" {
		t.Errorf("OnceDelayed() did not return the value, instead: %s", value)
	}
	if err != e {
		t.Errorf("OnceDelayed() did not return the error, instead: %s", err.Error())
	}
	duration := time.Since(start)
	if duration.Seconds() < 1 {
		t.Errorf("OnceDelayed() did not retry after the given delay")
	}
	if times != 2 {
		t.Errorf("F() was not called exactly 2 times, instead: %d", times)
	}
}

func TestTimes_ReturnsIfSuccessful(t *testing.T) {
	t.Parallel()
	var times = 0
	f := func() (interface{}, error) {
		times++
		return "value", nil
	}
	value, err := retry.Times(3, f)
	if value != "value" {
		t.Errorf("Times() did not return the value")
	}
	if err != nil {
		t.Errorf("Times() return an error when there wasn't any: %s", err.Error())
	}
	if times != 1 {
		t.Errorf("f() did not run exactly once, instead: %d", times)
	}
}

// testing is really boring
