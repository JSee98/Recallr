package storage

import (
	"github.com/JSee98/Recallr/storage/dragonfly"
)

func NewDefaultStore(cnfg dragonfly.DragonflyConfig) Store {
	return dragonfly.NewDragonflyStore(&cnfg)
}
