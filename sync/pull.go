package sync

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/hb-chen/docker-sync/pkg/log"
	"github.com/pkg/errors"
)

func pull(src *Image, dst *Image, opts ...string) error {
	var err error
	args := append([]string{"pull", src.String()}, opts...)
	cmd := exec.Command("docker", args...)

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
