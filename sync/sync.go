package sync

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/hb-chen/docker-sync/config"
	"github.com/hb-chen/docker-sync/pkg/log"
)

func Run(conf *config.Config) {
	sync := conf.Sync
	synced := make([]*config.Sync, 0, len(sync))
	for _, s := range sync {

		src := NewImageWith(s.Src)
		dst := NewImageWith(s.Dst)

		if conf.Pull {
			if err := pull(src, dst); err != nil {
				log.Error(err)
				continue
			}
		}

		if conf.Push {
			var err error
			dst.Id, err = getImageIdWith(src.Repository, src.Digest)
			if err != nil {
				log.Error(err)
			}

			if err := tag(dst); err != nil {
				log.Error(err)
				continue
			}

			if err := push(dst); err != nil {
				log.Error(err)
				continue
			}

			synced = append(synced, &config.Sync{
				Src: src.String(),
				Dst: dst.String(),
			})
		}

	}

	if conf.Push {
		storeSynced(synced)
	}
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

func storeSynced(synced []*config.Sync) error {
	bytes, err := yaml.Marshal(synced)
	if err != nil {
		return err
	}

	// 输出文件夹不存在时创建
	filePath := "synced.yml"
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// 写入配置文件
	err = ioutil.WriteFile(filePath, bytes, 0755)
	if err != nil {
		log.Debugw("write file error", "err", err)
		return err
	}

	return nil
}
