package sync

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/hb-chen/docker-sync/pkg/log"
)

func push(dst *Image) error {
	str := dst.String()
	idx := strings.LastIndex(str, "@")
	if idx > 0 {
		str = str[:idx]
	}
	cmd := exec.Command("docker", "push", str)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Debug("\n" + cmd.String())
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, stderr.String())
	} else {
		log.Debug("\n", stdout.String())
	}

	out := strings.Split(stdout.String(), "\n")
	for _, line := range out {
		split := strings.Split(line, " ")
		if len(split) > 2 && split[1] == "digest:" && (dst.Tag == "" || dst.Tag == split[0][:len(split[0])-1]) {
			dst.Digest = split[2]
			break
		}
	}

	return nil
}
