package common

const (
	PublicName = "sta音乐推荐系统"
	AnyUser    = "anyUser"
)

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}
