package coreunix

import (
	core "github.com/ipfs/go-ipfs/core"
	cid "gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	dag "gx/ipfs/QmVvNkTCx8V9Zei8xuTYTBdUXmbnDRS4iNuw1SztYyhQwQ/go-merkledag"
	ft "gx/ipfs/QmWE6Ftsk98cG2MTVgH4wJT8VP2nL9TuBkYTrz9GSqcsh5/go-unixfs"
)

func AddMetadataTo(n *core.IpfsNode, skey string, m *ft.Metadata) (string, error) {
	c, err := cid.Decode(skey)
	if err != nil {
		return "", err
	}

	nd, err := n.DAG.Get(n.Context(), c)
	if err != nil {
		return "", err
	}

	mdnode := new(dag.ProtoNode)
	mdata, err := ft.BytesForMetadata(m)
	if err != nil {
		return "", err
	}

	mdnode.SetData(mdata)
	if err := mdnode.AddNodeLink("file", nd); err != nil {
		return "", err
	}

	err = n.DAG.Add(n.Context(), mdnode)
	if err != nil {
		return "", err
	}

	return mdnode.Cid().String(), nil
}

func Metadata(n *core.IpfsNode, skey string) (*ft.Metadata, error) {
	c, err := cid.Decode(skey)
	if err != nil {
		return nil, err
	}

	nd, err := n.DAG.Get(n.Context(), c)
	if err != nil {
		return nil, err
	}

	pbnd, ok := nd.(*dag.ProtoNode)
	if !ok {
		return nil, dag.ErrNotProtobuf
	}

	return ft.MetadataFromBytes(pbnd.Data())
}
