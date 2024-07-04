package util

func MergeMap(v1, v2 map[string]string) map[string]string {
	ret := make(map[string]string, len(v1)+len(v2))
	for k, v := range v1 {
		ret[k] = v
	}
	for k, v := range v2 {
		ret[k] = v
	}
	return ret
}
