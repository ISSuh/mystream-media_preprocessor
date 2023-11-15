package media

type Timestamp struct {
	Pts uint64
	Dts uint64
}

func (t *Timestamp) IsEmpty() bool {
	return (t.Pts == 0) && (t.Dts == 0)
}

func (t *Timestamp) Diff(other Timestamp) uint64 {
	return other.Pts - t.Pts
}
