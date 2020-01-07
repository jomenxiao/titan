package main

import (
	"bytes"
	"encoding/hex"

	"github.com/pingcap/tidb/util/codec"
)

func decodeKey(key string) (dkey []byte, err error) {
	sk, err := hex.DecodeString(key)
	if err != nil {
		return sk, err
	}
	_, dkey, err = codec.DecodeBytes(sk, nil)
	if err != nil {
		return dkey, err
	}
	return dkey, nil
}

func getSize(namespace string, regions []*RegionInfo) (size int64, regionIdList []uint64, err error) {
	size = 0
	regionIdList = []uint64{}
	var startKey, endKey []byte
	for _, region := range regions {
		if region.StartKey != "" {
			startKey, err = decodeKey(region.StartKey)
			if err != nil {
				return size, regionIdList, err
			}
		}
		if region.EndKey != "" {
			endKey, err = decodeKey(region.EndKey)
			if err != nil {
				return size, regionIdList, err
			}
		}
		// fmt.Printf("start key %s   end key %s \n", startKey, endKey)

		if bytes.HasPrefix(startKey, []byte(namespace)) || bytes.HasPrefix(endKey, []byte(namespace)) {
			size += region.ApproximateSize
			regionIdList = append(regionIdList, region.ID)
		}
	}
	return size, regionIdList, nil
}
