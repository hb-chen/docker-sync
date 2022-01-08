package sync

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/hb-chen/docker-sync/pkg/log"
)

func pull(src *Image, dst *Image) error {
	var err error
	cmd := exec.Command("docker", "pull", src.String())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Debug("\n" + cmd.String())
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, stderr.String())
	} else {
		log.Debug("\n", stdout.String())
	}

	out := strings.Split(stdout.String(), "\n")
	for _, line := range out {
		if strings.HasPrefix(line, "Digest: ") {
			src.Digest = line[len("Digest: "):]
		}
	}

	return nil
}
