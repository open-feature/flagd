package trigger

type Trigger interface {
	Setup() error
	Update() error
}
