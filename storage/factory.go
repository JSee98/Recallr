package storage

import (
	"github.com/Jsee98/Recallr/storage/dragonfly"
)

func NewDefaultStore(cnfg dragonfly.DragonflyConfig) Store {
	return dragonfly.NewDragonflyStore(&cnfg)
}
