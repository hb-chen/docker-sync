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

	dst.Id, err = getImageIdWith(src.Repository, src.Digest)
	if err != nil {
		return err
	}

	return nil
}

func getImageIdWith(repo, digest string) (string, error) {
	if len(digest) > 0 {
		repo += "@" + digest
	}
	cmd := exec.Command("docker", "images", "--format", "{{.ID}}", repo)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Debug("\n" + cmd.String())
	err := cmd.Run()
	if err != nil {
		return "", errors.Wrap(err, stderr.String())
	} else {
		log.Debug("\n", stdout.String())
	}

	imageId := ""
	out := strings.Split(stdout.String(), "\n")
	if len(out) > 0 {
		imageId = out[0]
	} else {
		return "", errors.New("image id length 0")
	}

	return imageId, nil
}
