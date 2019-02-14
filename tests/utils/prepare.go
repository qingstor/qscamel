package utils

import "testing"

// PrepareNormalTest prepare the env for normal tests, base directory is a temp file
func PrepareNormalTest(t *testing.T) (*map[string]string, func(*map[string]string) error) {
	fileMap := CreateTestConfigFile(t, "copy", "fs", "fs", nil, nil)
	CreateLocalSrcTestRandDirFile(t, fileMap, 4, 1, MB*2, 2)
	CreateLocalDstDir(t, fileMap)
	return fileMap, CleanTestTempFile
}

// PrepareDefaultTest prepare the env for default case that the base directory is ~/.qscamel
func PrepareDefaultTest(t *testing.T) (*map[string]string, func(*map[string]string) error) {
	fileMap := CreateTestDefaultFile(t, "copy", "fs", "fs", nil, nil)
	CreateLocalSrcTestRandDirFile(t, fileMap, 4, 1, MB*2, 2)
	CreateLocalDstDir(t, fileMap)
	return fileMap, CleanTestTempFile
}

// PrepareDefinedTest will be used to generate user-defined file tree.
func PrepareDefinedTest(t *testing.T, filePerDir, dirPerDir int, fileSize int64, depth int) (*map[string]string, func(*map[string]string) error) {
	fileMap := CreateTestConfigFile(t, "copy", "fs", "fs", nil, nil)
	CreateLocalSrcTestRandDirFile(t, fileMap, filePerDir, dirPerDir, fileSize, depth)
	CreateLocalDstDir(t, fileMap)
	return fileMap, CleanTestTempFile
}
