package jsondiff

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
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

func Test_ChangeType(t *testing.T) {
	buf, err := json.Marshal(ChangeTypeCreate)
	require.NoError(t, err)
	require.Equal(t, `"ChangeTypeCreate"`, string(buf))
	ct := ChangeTypeNotSet
	err = json.Unmarshal([]byte(`"ChangeTypeCreate"`), &ct)
	require.NoError(t, err)
	require.Equal(t, ChangeTypeCreate, ct)
}
