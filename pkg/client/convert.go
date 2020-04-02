package client

import (
	"encoding/binary"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
)

func Merge(payload []byte, timestamp uint64) (merged []byte) {
	ts := make([]byte, 8)
	binary.LittleEndian.PutUint64(ts, timestamp)
	merged = append(ts, payload[:]...)
	return merged
}

func Split(merged []byte) (payload []byte, timestamp uint64) {
	if len(merged) < 9 {
		return payload, 0
	}
	payload = merged[8:]
	ts := merged[0:8]
	timestamp = binary.LittleEndian.Uint64(ts)
	return payload, timestamp
}

func NewSKV(kv *schema.KeyValue) *schema.StructuredKeyValue {
	return &schema.StructuredKeyValue{
		Key: kv.Key,
		Value: &schema.Content{
			Timestamp: uint64(time.Now().Unix()),
			Payload:   kv.Value,
		},
	}
}

func NewSKV2(key []byte, value []byte) *schema.StructuredKeyValue {
	return &schema.StructuredKeyValue{
		Key: key,
		Value: &schema.Content{
			Timestamp: uint64(time.Now().Unix()),
			Payload:   value,
		},
	}
}

func ToSItem(item *schema.Item) *schema.StructuredItem {
	payload, ts := Split(item.Value)
	return &schema.StructuredItem{
		Index: item.Index,
		Key:   item.Key,
		Value: &schema.Content{
			Payload:   payload,
			Timestamp: ts,
		},
	}
}

func ToKV(skv *schema.StructuredKeyValue) *schema.KeyValue {
	return &schema.KeyValue{
		Key:   skv.Key,
		Value: Merge(skv.Value.Payload, skv.Value.Timestamp),
	}
}
