package spruce

import (
	"log"

	ap "git.yaop.ink/tungyao/awesome-pool"
)

func NewDial() *ap.Pool {
	p, err := ap.NewPool(5, "127.0.0.1:6998")
	if err != nil {
		log.Panic(err)
	}
	return p
}
