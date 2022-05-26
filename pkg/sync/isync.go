package sync

type SYNC_STRATEGY int32

const (
	SYNC_STRATEGY_REPLACE = iota
	SYNC_STRATEGY_MERGE   = 1
)

type ISync interface {
	Fetch() (string, error)
}
