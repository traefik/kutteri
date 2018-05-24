package locker

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeLocker_GetLastTime(t *testing.T) {
	dir, clean, err := createTempDir()
	if err != nil {
		return
	}
	defer clean()

	recorder := TimeLocker{
		FilePath: dir + "/foo",
		HourBack: 1,
	}

	lastTime, err := recorder.GetLastTime()
	require.NoError(t, err)

	assert.NotEmpty(t, lastTime)
}

func TestTimeLocker_SaveLastTime(t *testing.T) {
	dir, clean, err := createTempDir()
	if err != nil {
		return
	}
	defer clean()

	recorder := TimeLocker{
		FilePath: dir + "/foo",
		HourBack: 1,
	}

	lastTime, err := recorder.SaveLastTime()
	require.NoError(t, err)

	assert.NotEmpty(t, lastTime)
}

func createTempDir() (string, func(), error) {
	dir, err := ioutil.TempDir("", "kutteri")
	if err != nil {
		return "", func() {}, err
	}

	return dir, func() {
		errRemove := os.RemoveAll(dir)
		if errRemove != nil {
			log.Println(errRemove)
		}
	}, nil
}
