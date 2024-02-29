package util

func ValidateCommitId(commitId string) bool {
	return len(commitId) == 64
}

func ValidateRef(ref string) bool {
	return len(ref) <= 64 && len(ref) > 0
}
