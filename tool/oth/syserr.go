package oth

import "log"

type SysErr struct {
	Code    int
	Message string
}

func (s SysErr) Error() string {
	return s.Message
}

func ErrHandle(err error) error {
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
