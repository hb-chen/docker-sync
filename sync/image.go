package sync

import (
	"strings"
)

type Image struct {
	Id         string
	Repository string
	Tag        string
	Digest     string
}

func (i *Image) String() string {
	str := i.Repository

	if len(i.Tag) > 0 {
		str += ":" + i.Tag
	}

	if len(i.Digest) > 0 {
		str += "@" + i.Digest
	}

	return str
}

func NewImageWith(str string) *Image {
	img := &Image{}
	idx := strings.LastIndex(str, "@")
	if idx > 0 {
		img.Digest = str[idx+1:]
		str = str[:idx]
	}

	idx = strings.LastIndex(str, ":")
	if idx > 0 {
		img.Tag = str[idx+1:]
		str = str[:idx]
	}

	img.Repository = str

	return img
}
