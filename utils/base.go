package utils

// Abstract
type Base struct {
	ID          int64  `db:"id" json:"id"`
	CreatedTime Time   `db:"created_time" json:"created_time"`
	UpdatedTime Time   `db:"updated_time" json:"updated_time"`
	UpdatedSeq  int64  `db:"updated_seq" json:"updated_seq"` // 乐观锁标记
	Region      string `db:"region" json:"region"`
}

// Abstract
type EnableBase struct {
	Base
	Removed bool `db:"removed" json:"removed"` // 软删标记
}
