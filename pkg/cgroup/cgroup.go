package cgroup

import (
	"fmt"
	"io/ioutil"
	"os"
)

const CGroupPath = "/sys/fs/cgroup/memory"
const CGroupMemStatsPath = CGroupPath + "/memory.stat"

func CGroupsEnabled() (bool, error) {
	fileInfo, err := os.Stat(CGroupPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Entry does not exist - no cgroups
			return false, nil
		}

		return false, fmt.Errorf("error statting dir %s: %w", CGroupPath, err)
	}

	if !fileInfo.IsDir() {
		// Entry exists but is not a dir
		return false, nil
	}

	return true, nil
}

func ReadCGroupStats() (string, error) {
	memstat, err := ioutil.ReadFile(CGroupMemStatsPath)
	if err != nil {
		return "", err
	}

	return string(memstat), nil
}
