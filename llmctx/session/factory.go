package session

import (
	"github.com/JSee98/Recallr/constants"
	"github.com/JSee98/Recallr/session/dragonfly"
	"github.com/JSee98/Recallr/storage"
)

func NewSessionManager(storageLayer storage.Store) SessionManager {
	if storageLayer == nil {
		panic("SessionManager requires a valid storage layer")
	}

	if storageLayer.Type() == constants.DragonFlyType {
		sessionStore := dragonfly.NewDragonflySessionStore(storageLayer)
		return dragonfly.NewSessionManager(sessionStore)
	}
	panic("Unsupported storage type for SessionManager: " + storageLayer.Type())
}
