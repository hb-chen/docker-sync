package sync

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hb-chen/docker-sync/config"
	"github.com/hb-chen/docker-sync/pkg/log"
	"gopkg.in/yaml.v2"
)

func Run(conf []*config.Sync) {
	synced := make([]*config.Sync, 0, len(conf))
	for _, s := range conf {

		src := NewImageWith(s.Src)
		dst := NewImageWith(s.Dst)

		if err := pull(src, dst); err != nil {
			log.Error(err)
			continue
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

	storeSynced(synced)
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
