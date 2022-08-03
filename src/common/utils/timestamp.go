package utils

// AlignTimestamp 对齐时间戳
func AlignTimestamp(ts int64, cycle int32) int64 {
	if cycle < 1 {
		return ts
	}
	return ts - (ts % int64(cycle))
}
