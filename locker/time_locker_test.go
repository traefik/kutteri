package locker

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeLocker_GetLastTime(t *testing.T) {
	dir := t.TempDir()

	recorder := TimeLocker{
		FilePath: filepath.Join(dir, "foo"),
		HourBack: 1,
	}

	lastTime, err := recorder.GetLastTime()
	require.NoError(t, err)

	assert.NotEmpty(t, lastTime)
}

func TestTimeLocker_SaveLastTime(t *testing.T) {
	dir := t.TempDir()

	recorder := TimeLocker{
		FilePath: filepath.Join(dir, "foo"),
		HourBack: 1,
	}

	lastTime, err := recorder.SaveLastTime()
	require.NoError(t, err)

	assert.NotEmpty(t, lastTime)
}
