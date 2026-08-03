package interfaceweather

import "context"

type Usage interface {
	Reserve(context.Context) (allowed bool, count int64, err error)
}
