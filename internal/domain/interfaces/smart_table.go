package interfaces

import (
	"belcamp/internal/domain/valueobject"
)

// SmartTableProvider is an interface that entities can implement to provide their own smart table configuration
type SmartTableProvider interface {
	GetSmartTableConfig() valueobject.SmartTableConfig
}
