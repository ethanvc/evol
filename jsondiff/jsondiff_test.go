package jsondiff

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJsonDiffer_JsonDiff(t *testing.T) {
	jd := NewJsonDiffer()
	changes, err := jd.JsonDiffStr(`{"a":1}`, `{"a":2}`)
	buf, _ := json.Marshal(changes)
	_ = buf
	requireNoChange(t, err, changes)
}

func requireNoChange(t *testing.T, err error, changes []Change) {
	require.NoError(t, err)
	require.Zero(t, len(changes))
}
