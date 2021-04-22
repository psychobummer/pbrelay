package keystore

import (
	"fmt"
	"time"
)

// An Artifact is what will be stored by a `Backend`
type Artifact interface {
	CreatedAt() time.Time
	Data() []byte
}

// A KeyArtifact is an artifact with some expectations around
// what will be stored
type KeyArtifact struct {
	key       []byte
	createdAt time.Time
}

// NewKeyArtifact returns a KeyArtifact initailized with the supplied data
func NewKeyArtifact(data []byte) (KeyArtifact, error) {
	if len(data) != 32 {
		return KeyArtifact{}, fmt.Errorf("this doesnt look like a key")
	}
	artifact := KeyArtifact{
		key:       data,
		createdAt: time.Now(),
	}
	return artifact, nil
}

// CreatedAt returns the creation time of this key
func (k KeyArtifact) CreatedAt() time.Time {
	return k.createdAt
}

// Data returns the data for this key
func (k KeyArtifact) Data() []byte {
	return k.key
}

// Backend defines the contract a storage backend must satisfy.
type Backend interface {
	Delete(key string) error
	Get(key string) (Artifact, error)
	HasKey(key string) bool
	Set(key string, artifact Artifact) error
}
