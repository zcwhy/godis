package parser

import (
	"godis/lib/util"
	"testing"
)

func TestParseOne(t *testing.T) {
	replies := []Reply{
		MakeIntReply(1),
		MakeStatusReply("OK"),
		MakeErrReply("ERR unknown"),
		MakeBulkReply([]byte("a\r\nb")), // test binary safe
		MakeNullBulkReply(),
		MakeMultiBulkReply([][]byte{
			[]byte("a"),
			[]byte("\r\n"),
		}),
		MakeEmptyMultiBulkReply(),
	}
	for _, re := range replies {
		result, err := ParseOne(re.ToBytes())
		if err != nil {
			t.Error(err)
			continue
		}
		if !util.BytesEquals(result.ToBytes(), re.ToBytes()) {
			t.Error("parse failed: " + string(re.ToBytes()))
		}
	}
}
