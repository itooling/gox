package oth

import "github.com/jinzhu/copier"

type Option copier.Option

func Copy(dst, src any) error {
	return copier.Copy(dst, src)
}

func CopyWithOption(dst, src any, opt Option) error {
	return copier.CopyWithOption(dst, src, copier.Option(opt))
}
