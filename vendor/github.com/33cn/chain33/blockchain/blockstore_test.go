package blockchain

import (
	"testing"

	"io/ioutil"
	"os"

	dbm "github.com/33cn/chain33/common/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStoreUpgradeMeta(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up
	os.RemoveAll(dir)       //删除已存在目录

	blockStoreDB := dbm.NewDB("blockchain", "leveldb", dir, 100)

	blockStore := NewBlockStore(nil, blockStoreDB, nil)
	require.NotNil(t, blockStore)

	meta, err := blockStore.GetStoreUpgradeMeta()
	require.NoError(t, err)
	require.Equal(t, meta.Version, "0.0.0")

	meta.Version = "1.0.0"
	err = blockStore.SetStoreUpgradeMeta(meta)
	require.NoError(t, err)
	meta, err = blockStore.GetStoreUpgradeMeta()
	require.NoError(t, err)
	require.Equal(t, meta.Version, "1.0.0")
}

func TestSeqSaveAndGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up
	os.RemoveAll(dir)       //删除已存在目录

	blockStoreDB := dbm.NewDB("blockchain", "leveldb", dir, 100)

	blockStore := NewBlockStore(nil, blockStoreDB, nil)
	assert.NotNil(t, blockStore)
	blockStore.saveSequence = true
	blockStore.isParaChain = false

	newBatch := blockStore.NewBatch(true)
	seq, err := blockStore.saveBlockSequence(newBatch, []byte("s0"), 0, 1, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), seq)
	err = newBatch.Write()
	assert.Nil(t, err)

	newBatch = blockStore.NewBatch(true)
	seq, err = blockStore.saveBlockSequence(newBatch, []byte("s1"), 1, 1, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), seq)
	err = newBatch.Write()
	assert.Nil(t, err)

	s, err := blockStore.LoadBlockLastSequence()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), s)

	s2, err := blockStore.GetBlockSequence(s)
	assert.Nil(t, err)
	assert.Equal(t, []byte("s1"), s2.Hash)

	s3, err := blockStore.GetSequenceByHash([]byte("s1"))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), s3)
}

func TestParaSeqSaveAndGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up
	os.RemoveAll(dir)       //删除已存在目录

	blockStoreDB := dbm.NewDB("blockchain", "leveldb", dir, 100)

	blockStore := NewBlockStore(nil, blockStoreDB, nil)
	assert.NotNil(t, blockStore)
	blockStore.saveSequence = true
	blockStore.isParaChain = true

	newBatch := blockStore.NewBatch(true)
	seq, err := blockStore.saveBlockSequence(newBatch, []byte("s0"), 0, 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), seq)
	err = newBatch.Write()
	assert.Nil(t, err)

	newBatch = blockStore.NewBatch(true)
	seq, err = blockStore.saveBlockSequence(newBatch, []byte("s1"), 1, 1, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), seq)
	err = newBatch.Write()
	assert.Nil(t, err)

	s, err := blockStore.LoadBlockLastSequence()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), s)

	s2, err := blockStore.GetBlockSequence(s)
	assert.Nil(t, err)
	assert.Equal(t, []byte("s1"), s2.Hash)

	s3, err := blockStore.GetSequenceByHash([]byte("s1"))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), s3)

	s4, err := blockStore.GetMainSequenceByHash([]byte("s1"))
	assert.Nil(t, err)
	assert.Equal(t, int64(10), s4)

	s5, err := blockStore.LoadBlockLastMainSequence()
	assert.Nil(t, err)
	assert.Equal(t, int64(10), s5)

	s6, err := blockStore.GetBlockByMainSequence(1)
	assert.Nil(t, err)
	assert.Equal(t, []byte("s0"), s6.Hash)

	chain := &BlockChain{
		blockStore: blockStore,
	}
	s7, err := chain.ProcGetMainSeqByHash([]byte("s0"))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), s7)

	_, err = chain.ProcGetMainSeqByHash([]byte("s0-not-exist"))
	assert.NotNil(t, err)
}
