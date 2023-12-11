package strategy

type IStrategy interface {
	Evict() (uint64, bool)
}
