package sync

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

func tag(dst *Image) error {
	cmd := exec.Command("docker", "tag", dst.Id, dst.String())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, stderr.String(), dst.String())
	}

	return nil
}
