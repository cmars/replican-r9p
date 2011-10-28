
# replican-r9p - Remote filesystem synchronization support #

r9p adds remote filesystem synchronization support to replican-sync with 
the Plan9 filesystem protocol.

9P is used to implement fs.BlockStore, which is a suitable match & patch 
endpoint for replican-sync.

r9p is not serving the remote files directly over 9P. Instead, the server 
indexes locally, and serves up the metadata and raw block access in a custom 
filesystem structure.

This pre-indexing and fetching of only the delta follows a similar 
synchronization pattern as zsync (pull, rather than push).



