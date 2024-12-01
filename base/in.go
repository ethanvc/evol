package base

func In[T comparable](val T, set1 T, sets ...T) bool {
	if val == set1 {
		return true
	}
	for _, s := range sets {
		if val == s {
			return true
		}
	}
	return false
}
