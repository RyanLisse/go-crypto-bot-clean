package crypto

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestKeyRegistry_AddGetActiveGetKeyRetire(t *testing.T) {
	r := NewKeyRegistry()
	key1 := []byte("supersecretkey1")
	key2 := []byte("supersecretkey2")

	// Add first key, set as active
	err := r.AddKey("v1", key1, true)
	assert.NoError(t, err)

	active, err := r.GetActiveKey()
	assert.NoError(t, err)
	assert.Equal(t, "v1", active.Version)
	assert.Equal(t, key1, active.Key)
	assert.Equal(t, "active", active.Status)

	// Add second key, set as active
	err = r.AddKey("v2", key2, true)
	assert.NoError(t, err)
	active, err = r.GetActiveKey()
	assert.NoError(t, err)
	assert.Equal(t, "v2", active.Version)
	assert.Equal(t, key2, active.Key)

	// Get by version
	k1, err := r.GetKey("v1")
	assert.NoError(t, err)
	assert.Equal(t, key1, k1.Key)

	// Retire v1
	err = r.RetireKey("v1")
	assert.NoError(t, err)
	k1, err = r.GetKey("v1")
	assert.NoError(t, err)
	assert.Equal(t, "retired", k1.Status)

	// List keys
	keys := r.ListKeys()
	assert.Len(t, keys, 2)
}
