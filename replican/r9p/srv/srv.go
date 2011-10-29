
package srv

import (
	"bytes"
	"gob"
	"os"
	
	"github.com/cmars/replican-sync/replican/fs"
	"github.com/cmars/replican-sync/replican/r9p"
	
	"go9p.googlecode.com/hg/p"
	"go9p.googlecode.com/hg/p/srv"
)

type StoreSrv struct {
	*srv.Fsrv
	user p.User
	local *fs.LocalStore
	rootGob []byte
}

func Serve(path string) (*StoreSrv, os.Error) {
	store, err := fs.NewLocalStore(path)
	if err != nil {
		return nil, err
	}
	
	return NewStoreSrv(store)
}

func NewStoreSrv(local *fs.LocalStore) (storeSrv *StoreSrv, err os.Error) {
	user := p.OsUsers.Uid2User(os.Geteuid())
	
	root := new(srv.File)
	err = root.Add(nil, "/", user, nil, p.DMDIR|0555, nil)
	
	noderoot := &RootFile{ store: local }
	err = noderoot.Add(root, r9p.ROOT_FILE, user, nil, 0444, noderoot)
	
	strongdir := &StrongDir{ store: local, File: new(srv.File) }
	err = strongdir.Add(root, r9p.STRONG_DIR, user, nil, p.DMDIR|0555, strongdir)
	
	storeSrv = &StoreSrv{ 
		local: local,
		user: user }
	storeSrv.Fsrv = srv.NewFileSrv(root)
	
	return storeSrv, nil
}

type RootFile struct {
	srv.File
	store *fs.LocalStore
	rootGob []byte
}

func (rootFile *RootFile) Read(fid *srv.FFid, buf []byte, offset uint64) (int, os.Error) {
	if rootFile.rootGob == nil {
		buffer := bytes.NewBuffer([]byte{})
		encoder := gob.NewEncoder(buffer)
		encoder.Encode(rootFile.store.Root())
		rootFile.rootGob = buffer.Bytes()
	}
	
	n := len(rootFile.rootGob)
	if offset >= uint64(n) {
		return 0, nil
	}

	b := rootFile.rootGob[int(offset):n]
	n -= int(offset)
	if len(buf) < n {
		n = len(buf)
	}
	
	copy(buf[offset:int(offset)+n], b[offset:])
	return n, nil
}

type StrongDir struct {
	*srv.File
	store *fs.LocalStore	
}

func (strongDir *StrongDir) Find(name string) *srv.File {
	file := strongDir.File.Find(name)
	
	if file == nil {
		strongFile := &StrongFile{ strong: name, store: strongDir.store }
		file = strongFile.File
		user := p.OsUsers.Uid2User(os.Geteuid())
		strongFile.Add(file, name, user, nil, 0444, file)
	}
	
	return file
}

type StrongFile struct {
	*srv.File
	store *fs.LocalStore
	strong string
}

func (strongFile *StrongFile) Read(fid *srv.FFid, buf []byte, offset uint64) (int, os.Error) {
	buffer := bytes.NewBuffer(buf)
	n, err := strongFile.store.ReadInto(
			strongFile.strong, int64(offset), int64(len(buf)), buffer)
	return int(n), err
}



