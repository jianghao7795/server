package utils

import "os/exec"

// NetWorkStatus check network status by curl 测试网址是否联网
func NetWorkStatus(path string) bool {
	cmd := exec.Command("curl", path)

	err := cmd.Run()
	return err == nil
}
