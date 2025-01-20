package cloudian

import (
	"encoding/json"
	"strings"

	"k8s.io/utils/ptr"
)

type Limit struct {
	Soft *int64
	Hard *int64
}

type QoSLimit int

const (
	StorageQuotaKBytes QoSLimit = iota
	StorageQuotaCount
	RequestRate
	DataKBytesIn
	DataKBytesOut
)

var qoSLimitName = map[QoSLimit]string{
	StorageQuotaKBytes: "StorageQuotaKBytes",
	StorageQuotaCount:  "StorageQuotaCount",
	RequestRate:        "RequestRate",
	DataKBytesIn:       "DataKBytesIn",
	DataKBytesOut:      "DataKBytesOut",
}

var qoSLimitJSONName = map[string]QoSLimit{
	"STORAGE_QUOTA_KBYTES": StorageQuotaKBytes,
	"STORAGE_QUOTA_COUNT":  StorageQuotaCount,
	"REQUEST_RATE":         RequestRate,
	"DATAKBYTES_IN":        DataKBytesIn,
	"DATAKBYTES_OUT":       DataKBytesOut,
}

func (ql QoSLimit) String() string {
	return qoSLimitName[ql]
}

type QoS map[QoSLimit]Limit

func (q QoS) QueryParams() map[string]int64 {
	params := make(map[string]int64)
	for ql, l := range q {
		params["wl"+ql.String()] = ptr.Deref(l.Soft, -1)
		params["hl"+ql.String()] = ptr.Deref(l.Hard, -1)
	}
	return params
}

func (q QoS) UnmarshalJSON(b []byte) error {
	var data struct {
		QOSLimitList []struct {
			Type  string `json:"type"`
			Value *int64 `json:"value"`
		} `json:"qosLimitList"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for _, ql := range data.QOSLimitList {
		if before, found := strings.CutSuffix(ql.Type, "_LW"); found {
			l := q[qoSLimitJSONName[before]]
			l.Soft = ql.Value
		}
		if before, found := strings.CutSuffix(ql.Type, "_LH"); found {
			l := q[qoSLimitJSONName[before]]
			l.Hard = ql.Value
		}
	}

	return nil
}
