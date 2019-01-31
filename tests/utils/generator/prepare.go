package generator

// PrepareNormalTest prepare the env for normal tests, base directory is a temp file
func PrepareNormalTest(p bool) (*map[string]string, func(*map[string]string) error, error) {
	fileMap, err := CreateTestConfigFile("copy", "fs", "fs", nil, nil, p)
	if err != nil {
		return nil, nil, err
	}
	err = CreateLocalSrcTestRandDirFile(fileMap, 4, 1, MB*2, 2)
	if err != nil {
		return nil, nil, err
	}
	err = CreateLocalDstDir(fileMap)
	if err != nil {
		return nil, nil, err
	}
	return fileMap, CleanTestTempFile, nil
}

// PrepareDefaultTest prepare the env for default case that the base directory is ~/.qscamel
func PrepareDefaultTest(p bool) (*map[string]string, func(*map[string]string) error, error) {
	fileMap, err := CreateTestConfigFile("copy", "fs", "fs", nil, nil, p)
	if err != nil {
		return nil, nil, err
	}
	err = CreateLocalSrcTestRandDirFile(fileMap, 4, 1, MB*2, 2)
	if err != nil {
		return nil, nil, err
	}
	err = CreateLocalDstDir(fileMap)
	if err != nil {
		return nil, nil, err
	}
	return fileMap, CleanTestTempFile, nil
}

// PrepareDefinedTest will be used to generate user-defined file tree.
func PrepareDefinedTest(filePerDir, dirPerDir int, fileSize int64, depth int, p bool) (
	*map[string]string, func(*map[string]string) error, error) {
	fileMap, err := CreateTestConfigFile("copy", "fs", "fs", nil, nil, p)
	if err != nil {
		return nil, nil, err
	}
	err = CreateLocalSrcTestRandDirFile(fileMap, filePerDir, dirPerDir, fileSize, depth)
	if err != nil {
		return nil, nil, err
	}
	err = CreateLocalDstDir(fileMap)
	if err != nil {
		return nil, nil, err
	}
	return fileMap, CleanTestTempFile, nil
}
