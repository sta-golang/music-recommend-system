package plugin

func Normalization(x, max, min float64) float64 {
	if max == min {
		if x == min {
			return min
		}
		return (x - min) / 1.0
	}
	return (x - min) / (max - min)
}
