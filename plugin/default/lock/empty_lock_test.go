package lock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyLocker(t *testing.T) {
	locker := &emptyLocker{}

	res, err := locker.Lock(context.Background(), "", 0)
	assert.NoError(t, err)
	assert.Equal(t, res, "")

	locker.Unlock(context.Background(), "", "")

	err = locker.Close()
	assert.NoError(t, err)
}
