package util

//maps val from [ istart , iend ] to [ dstart , dend ]
func Map(val, iStart, iEnd, DStart, Dend float64) float64 {
	return DStart + ((val-iStart)/(iEnd-iStart))*(Dend-DStart)
}

func Max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func Min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
