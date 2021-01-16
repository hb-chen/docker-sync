package sync

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/hb-chen/docker-sync/pkg/log"
)

func tag(dst *Image) error {
	cmd := exec.Command("docker", "tag", dst.Id, dst.String())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Debug("\n" + cmd.String())
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, stderr.String(), dst.Id, dst.String())
	}

	return nil
}
