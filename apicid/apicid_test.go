package apicid

import (
	"encoding/json"
	"testing"

	cid "github.com/ipfs/go-cid"
)

func TestJson(t *testing.T) {
	cid, _ := cid.Decode("zb2rhak9iRgDiik36KQBRr2qiCJHdyBH7YxFmw7FTdM6zo31m")
	hash := FromCid(cid)
	data, err := json.Marshal(hash)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"zb2rhak9iRgDiik36KQBRr2qiCJHdyBH7YxFmw7FTdM6zo31m"` {
		t.Fatalf("json string incorrect: %s\n", data)
	}
	var hash2 Hash
	err = json.Unmarshal(data, &hash2)
	if err != nil {
		t.Fatal(err)
	}
	if hash != hash2 {
		t.Fatal("round trip failed")
	}
}

func TestJsonMap(t *testing.T) {
	cid1, _ := cid.Decode("zb2rhak9iRgDiik36KQBRr2qiCJHdyBH7YxFmw7FTdM6zo31m")
	cid2, _ := cid.Decode("QmRJggJREPCt7waGQKMXymrXRvrvsSiiPjgFbLK9isuM8K")
	hash1 := FromCid(cid1)
	hash2 := FromCid(cid2)
	m := map[Hash]string{hash1: "a value", hash2: "something else"}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	m2 := map[Hash]string{}
	err = json.Unmarshal(data, &m2)
	if err != nil {
		t.Fatal(err)
	}
	if len(m2) != 2 || m[hash1] != m2[hash1] || m[hash2] != m2[hash2] {
		t.Fatal("round trip failed")
	}
}
