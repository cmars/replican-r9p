
package clnt

import (
	"bytes"
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
	c *clnt.Clnt
	user p.User
	root *fs.Dir
	index *fs.BlockIndex
}

func Connect(addr string) (store *RemoteStore, err os.Error) {
	store = &RemoteStore{}
	store.user = p.OsUsers.Uid2User(os.Geteuid())
	
	store.c, err = clnt.Mount("tcp", addr, "", store.user)
	if err != nil {
		return nil, err
	}
	
	return store, nil
}

func (store *RemoteStore) Refresh() (err os.Error) {
	store.root = nil
	store.index = nil
	
	info, err := store.c.FStat(r9p.ROOT_FILE)
	if err != nil {
		return err
	}
	
	rootFile, err := store.c.FOpen(r9p.ROOT_FILE, p.OREAD)
	if err != nil {
		return err
	}
	
	raw := make([]byte, info.Length)
	_, err = rootFile.Read(raw)
	if err != nil {
		return err
	}
	
	buffer := bytes.NewBuffer(raw)
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
	path := filepath.Join(r9p.STRONG_DIR, strong)
	f, err := store.c.FOpen(path, p.OREAD)
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
	path := filepath.Join(r9p.STRONG_DIR, strong)
	f, err := store.c.FOpen(path, p.OREAD)
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



