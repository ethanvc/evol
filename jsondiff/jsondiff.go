package jsondiff

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"

	"github.com/ethanvc/evol/base"
)

type JsonDiffer struct {
	MissingEqualZero bool
	NullEqualZero    bool
	changes          []Change
}

func NewJsonDiffer() *JsonDiffer {
	return &JsonDiffer{}
}

func (jd *JsonDiffer) JsonDiffStr(src, dst string) ([]Change, error) {
	return jd.JsonDiff([]byte(src), []byte(dst))
}

func (jd *JsonDiffer) JsonDiff(src, dst []byte) ([]Change, error) {
	var srcAny, dstAny any
	if err := unmarshal(src, &srcAny); err != nil {
		return nil, errors.Join(errors.New("src json unmarshal fail"), err)
	}
	if err := unmarshal(dst, &dstAny); err != nil {
		return nil, errors.Join(errors.New("dst json unmarshal fail"), err)
	}
	jd.changes = nil
	jd.diff(nil, srcAny, dstAny)
	return jd.changes, nil
}

func (jd *JsonDiffer) diff(p []string, src any, dst any) {
	if src == nil && dst == nil {
		return
	}
	if src == nil {
		jd.addChange(ChangeTypeCreate, p, src, dst)
		return
	}
	if dst == nil {
		jd.addChange(ChangeTypeDelete, p, src, dst)
		return
	}
	if reflect.TypeOf(src) != reflect.TypeOf(dst) {
		jd.addChange(ChangeTypeSchema, p, src, dst)
		return
	}
	switch realSrc := src.(type) {
	case map[string]any:
		jd.diffMap(p, realSrc, dst.(map[string]any))
	case string:
		jd.diffStr(p, realSrc, dst.(string))
	case json.Number:
		jd.diffNumber(p, realSrc, dst.(json.Number))
	case []any:
		jd.diffArray(p, realSrc, dst.([]any))
	default:
		panic("type not support")
	}
}

func (jd *JsonDiffer) diffArray(p []string, src, dst []any) {

}

func (jd *JsonDiffer) diffNumber(p []string, src, dst json.Number) {
	if string(src) == string(dst) {
		return
	}
	jd.addChange(ChangeTypeUpdate, p, src, dst)
}

func (jd *JsonDiffer) diffStr(p []string, src, dst string) {
	if src == dst {
		return
	}
	var srcAny, dstAny any
	if err := unmarshal([]byte(src), &srcAny); err == nil {
		if err := unmarshal([]byte(dst), &dstAny); err == nil {
			if base.In(reflect.TypeOf(srcAny).Kind(), reflect.Map, reflect.Array) &&
				base.In(reflect.TypeOf(dstAny).Kind(), reflect.Map, reflect.Array) {
				jd.diff(p, srcAny, dstAny)
				return
			}
		}
	}

	jd.addChange(ChangeTypeUpdate, p, src, dst)
}

func (jd *JsonDiffer) addChange(t ChangeType, p []string, from, to any) {
	change := Change{
		ChangeType: t,
		JsonPath:   slices.Clone(p),
		From:       from,
		To:         to,
	}
	jd.changes = append(jd.changes, change)
}

func (jd *JsonDiffer) diffMap(p []string, src, dst map[string]any) {
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	for k := range dst {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	keys = slices.Compact(keys)
	for _, key := range keys {
		srcVal, ok1 := src[key]
		dstVal, ok2 := dst[key]
		newPath := append(p, key)
		if !ok1 {
			jd.addChange(ChangeTypeCreate, newPath, nil, dstVal)
			continue
		}
		if !ok2 {
			jd.addChange(ChangeTypeDelete, newPath, srcVal, nil)
			continue
		}
		jd.diff(newPath, srcVal, dstVal)
	}
}

type ChangeType int

const (
	ChangeTypeNotSet ChangeType = iota
	ChangeTypeCreate
	ChangeTypeUpdate
	ChangeTypeDelete
	ChangeTypeSchema
)

func (ct ChangeType) String() string {
	switch ct {
	case ChangeTypeNotSet:
		return "ChangeTypeNotSet"
	case ChangeTypeCreate:
		return "ChangeTypeCreate"
	case ChangeTypeUpdate:
		return "ChangeTypeUpdate"
	case ChangeTypeDelete:
		return "ChangeTypeDelete"
	case ChangeTypeSchema:
		return "ChangeTypeSchema"
	default:
		return fmt.Sprintf("ChangeTypeUnknown_%d", ct)
	}
}

func (ct ChangeType) MarshalText() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *ChangeType) UnmarshalText(text []byte) error {
	switch string(text) {
	case ChangeTypeNotSet.String():
		*ct = ChangeTypeNotSet
	case ChangeTypeCreate.String():
		*ct = ChangeTypeCreate
	case ChangeTypeUpdate.String():
		*ct = ChangeTypeUpdate
	case ChangeTypeDelete.String():
		*ct = ChangeTypeDelete
	case ChangeTypeSchema.String():
		*ct = ChangeTypeSchema
	default:
		return errors.New(fmt.Sprintf("unknown change type(%s)", text))
	}
	return nil
}

type Change struct {
	ChangeType ChangeType `json:"change_type"`
	JsonPath   []string   `json:"json_path"`
	From       any        `json:"from"`
	To         any        `json:"to"`
}

func unmarshal(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err := decoder.Decode(v)
	return err
}
