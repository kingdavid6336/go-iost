package snapshot

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/core/block"
	"github.com/iost-official/go-iost/db"
	"github.com/iost-official/go-iost/db/kv/leveldb"
	"github.com/iost-official/go-iost/ilog"
)

// Save the function for saving block's head from snapshot.
func Save(db db.MVCCDB, blk *block.Block) error {
	bhJSON, err := json.Marshal(blk.Head)
	if err != nil {
		return fmt.Errorf("json fail: %v", err)
	}
	err = db.Put("snapshot", "blockHead", string(bhJSON))
	if err != nil {
		return fmt.Errorf("state db put fail: %v", err)
	}
	return nil
}

// Load the function for loading block's head from snapshot.
func Load(db db.MVCCDB) (*block.Block, error) {
	bhJSON, err := db.Get("snapshot", "blockHead")
	if err != nil {
		return nil, fmt.Errorf("get current block head from state db failed. err: %v", err)
	}
	bh := &block.BlockHead{}
	err = json.Unmarshal([]byte(bhJSON), bh)
	if err != nil {
		return nil, fmt.Errorf("block head decode failed. err: %v", err)
	}

	blk := &block.Block{Head: bh}
	return blk, blk.CalculateHeadHash()
}

// ToSnapshot the function for saving db to snapshot.
func ToSnapshot(conf *common.Config) error {
	var src string
	src = conf.DB.LdbPath + "StateDB"
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	file, err := os.Create(conf.DB.LdbPath + "Snapshot.tar.gz")
	if err != nil {
		ilog.Fatal(err)
	}
	defer file.Close()

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		f.Close()
		return nil
	})
	return nil
}

// FromSnapshot the function for loading db from snapshot.
func FromSnapshot() error {
	return nil
}

// ToFile the function for saving db to File.
func ToFile(conf *common.Config) error {
	var src string
	src = conf.DB.LdbPath + "StateDB"

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}
	db, err := leveldb.NewDB(src)
	defer db.Close()

	file, err := os.Create(conf.DB.LdbPath + "Snapshot.iost")
	if err != nil {
		ilog.Fatal(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	iter := db.NewIteratorByPrefix([]byte("")).(*leveldb.Iter)
	for iter.Next() {
		k := string(iter.Key())
		v := string(iter.Value())
		writer.WriteString(k + "\n")
		writer.WriteString(v + "\n")
	}

	return nil
}

// FromFile the function for loading db from File.
func FromFile() error {
	return nil
}
