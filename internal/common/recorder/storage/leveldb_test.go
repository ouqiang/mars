package storage

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/ouqiang/goutil"
	"github.com/ouqiang/mars/internal/common/recorder"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/ouqiang/mars/internal/common"
)

func TestLevelDB(t *testing.T) {
	queueSize := 10
	queue := common.NewQueue(queueSize)
	randString := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(goutil.RandNumber(10000000, 99999999))
	dbFile := filepath.Join(os.TempDir(), randString)
	db, err := leveldb.OpenFile(dbFile, nil)
	require.NoError(t, err)
	defer os.RemoveAll(dbFile)
	defer db.Close()

	s := NewLevelDB(db, queue)

	tx := recorder.NewTransaction()
	err = s.Put(tx)
	require.NoError(t, err)
	tx, err = s.Get(tx.Id)
	require.NoError(t, err)

	tx, err = s.Get("not found")
	require.Error(t, err)
}
