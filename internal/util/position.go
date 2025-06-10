package util

// NormPos normalizes a position within a list of given length.
func NormPos(pos, length int) int {
	for pos < 0 {
		pos += length
	}
	if length > 0 {
		pos %= length
	}
	return pos
}
