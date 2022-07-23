package sync

import "context"

/*
ISync implementations watch for changes in the flag source
(HTTP backend, local file, s3 bucket), and fetch the latest values.
*/
type ISync interface {
	Fetch(ctx context.Context) (string, error)
	Notify(ctx context.Context, c chan<- INotify)
}
