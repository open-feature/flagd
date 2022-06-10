package sync

/*
ISync implementations watch for changes in the flag source (HTTP backend, local file, s3 bucket), and fetch the latest values.
*/
type ISync interface {
	Fetch() (string, error)
	Notify(chan<- INotify)
}
