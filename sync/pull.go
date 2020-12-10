package sync

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/hb-chen/docker-sync/pkg/log"
	"github.com/pkg/errors"
)

func pull(src *Image, dst *Image) error {
	cmd := exec.Command("docker", "pull", src.String())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, stderr.String())
	} else {
		log.Debug("\n"+cmd.String()+"\n\n", stdout.String())
	}

	out := strings.Split(stdout.String(), "\n")
	for _, line := range out {
		if strings.HasPrefix(line, "Digest: ") {
			src.Digest = line[len("Digest: "):]
		}
	}

	dst.Id, err = getImageIdWith(src.Repository, src.Tag)
	if err != nil {
		return err
	}

	return nil
}

func getImageIdWith(repo, tag string) (string, error) {
	if len(tag) > 0 {
		repo += ":" + tag
	}
	cmd := exec.Command("docker", "images", "--format", "{{.ID}}", repo)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.Wrap(err, stderr.String())
	} else {
		log.Debug("\n"+cmd.String()+"\n\n", stdout.String())
	}

	imageId := ""
	out := strings.Split(stdout.String(), "\n")
	if len(out) > 0 {
		imageId = out[0]
	}

	return imageId, nil
}
