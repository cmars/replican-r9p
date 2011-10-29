
package clnt

import (
	"bytes"
//	"fmt"
	"gob"
	"io"
	"os"
	"path/filepath"
	
	"github.com/cmars/replican-sync/replican/fs"
	
	"github.com/cmars/replican-sync/replican/r9p"
	
	"go9p.googlecode.com/hg/p"
	"go9p.googlecode.com/hg/p/clnt"
)

type RemoteStore struct {
	*clnt.Clnt
	root *fs.Dir
	index *fs.BlockIndex
}

func Connect(addr string) (store *RemoteStore, err os.Error) {
	store = &RemoteStore{}
	user := p.OsUsers.Uid2User(os.Geteuid())
	
	store.Clnt, err = clnt.Mount("tcp", addr, "", user)
	if err != nil {
		return nil, err
	}
	
	return store, nil
}

func (store *RemoteStore) Refresh() (err os.Error) {
	store.root = nil
	store.index = nil
	
	rootFile, err := store.FOpen(filepath.Join("/", r9p.ROOT_FILE), p.OREAD)
	if err != nil {
		return err
	}
	
	buffer := bytes.NewBuffer([]byte{})
	chunk := make([]byte, fs.BLOCKSIZE)
	for {
		n, err := rootFile.Read(chunk)
		if err != nil && err != os.EOF {
			return err
		}
		if n == 0 {
			err = nil
			break
		}
		buffer.Write(chunk[0:n])
	}
	
	decoder := gob.NewDecoder(buffer)
	
	store.root = &fs.Dir{}
	err = decoder.Decode(store.root)
	
	return err
}

func (store *RemoteStore) Root() *fs.Dir {
	return store.root
}

func (store *RemoteStore) Index() *fs.BlockIndex {
	if store.index == nil {
		store.index = fs.IndexBlocks(store.root)
	}
	return store.index
}

func (store *RemoteStore) ReadBlock(strong string) ([]byte, os.Error) {
	path := filepath.Join("/", r9p.STRONG_DIR, strong)
	f, err := store.FOpen(path, p.OREAD)
	if err != nil {
		return nil, err
	}
	
	buf := make([]byte, fs.BLOCKSIZE)
	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	
	return buf, nil
}

func (store *RemoteStore) ReadInto(strong string, from int64, length int64, writer io.Writer) os.Error {
	path := filepath.Join("/", r9p.STRONG_DIR, strong)
	f, err := store.FOpen(path, p.OREAD)
	if err != nil {
		return err
	}
	
	buf := make([]byte, length)
	_, err = f.ReadAt(buf, from)
	if err != nil {
		return err
	}
	writer.Write(buf)
	
	return nil
}



