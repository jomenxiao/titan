package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pingcap/kvproto/pkg/metapb"
	"github.com/pingcap/kvproto/pkg/pdpb"
	"github.com/pkg/errors"
)

const (
	regionListApi string = "/pd/api/v1/regions"
)

type RegionInfo struct {
	ID          uint64              `json:"id"`
	StartKey    string              `json:"start_key"`
	EndKey      string              `json:"end_key"`
	RegionEpoch *metapb.RegionEpoch `json:"epoch,omitempty"`
	Peers       []*metapb.Peer      `json:"peers,omitempty"`

	Leader          *metapb.Peer      `json:"leader,omitempty"`
	DownPeers       []*pdpb.PeerStats `json:"down_peers,omitempty"`
	PendingPeers    []*metapb.Peer    `json:"pending_peers,omitempty"`
	WrittenBytes    uint64            `json:"written_bytes"`
	ReadBytes       uint64            `json:"read_bytes"`
	WrittenKeys     uint64            `json:"written_keys"`
	ReadKeys        uint64            `json:"read_keys"`
	ApproximateSize int64             `json:"approximate_size"`
	ApproximateKeys int64             `json:"approximate_keys"`
}

type RegionsInfo struct {
	Count   int           `json:"count"`
	Regions []*RegionInfo `json:"regions"`
}

func httpGet(url string) (regionsInfo *RegionsInfo, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http get code is %s, not 200", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("read data error %v", err)
	}
	defer resp.Body.Close()
	regionsInfo = &RegionsInfo{}
	if err := json.Unmarshal(body, regionsInfo); err != nil {
		return nil, errors.Errorf("data unmarshal error %v", err)
	}

	return regionsInfo, nil
}
