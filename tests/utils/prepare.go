package utils

// PrepareNormalTest prepare the env for normal tests, base directory is a temp file
func PrepareNormalTest() (map[string]string, func(map[string]string) error) {
	fileMap := CreateTestConfigFile("copy", "fs", "fs", nil, nil)
	CreateLocalSrcTestRandDirFile(fileMap, 4, 1, MB*2, 2)
	CreateLocalDstDir(fileMap)
	return fileMap, CleanTestTempFile
}

// PrepareDefaultTest prepare the env for default case that the base directory is ~/.qscamel
func PrepareDefaultTest() (map[string]string, func(map[string]string) error) {
	fileMap := CreateTestDefaultFile("copy", "fs", "fs", nil, nil)
	CreateLocalSrcTestRandDirFile(fileMap, 4, 1, MB*2, 2)
	CreateLocalDstDir(fileMap)
	return fileMap, CleanTestTempFile
}

// PrepareDefinedTest will be used to generate user-defined file tree.
func PrepareDefinedTest(filePerDir, dirPerDir int, fileSize int64, depth int) (map[string]string, func(map[string]string) error) {
	fileMap := CreateTestConfigFile("copy", "fs", "fs", nil, nil)
	CreateLocalSrcTestRandDirFile(fileMap, filePerDir, dirPerDir, fileSize, depth)
	CreateLocalDstDir(fileMap)
	return fileMap, CleanTestTempFile
}
