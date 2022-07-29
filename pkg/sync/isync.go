package sync

import "context"

/*
ISync implementations watch for changes in the flag source
(HTTP backend, local file, s3 bucket), and fetch the latest values.
*/
type ISync interface {
	Fetch(ctx context.Context) (string, error)
	// Notify implementor should signal its readiness on the ready chan
	Notify(ctx context.Context, ready chan<- struct{}, c chan<- INotify)
}
