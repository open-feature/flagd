package service

import (
	"context"

	"github.com/open-feature/flagd/pkg/eval"
)

type IServiceConfiguration interface{}

/*
IService implementations define handlers for a particular transport, which call the IEvaluator implementation.
*/
type IService interface {
	Serve(ctx context.Context, eval eval.IEvaluator) error
	Notify()
}
