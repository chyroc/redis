package redis

import (
	"bytes"
	"fmt"
	"strconv"
)

// Reply ...
type Reply struct {
	err     error
	null    bool
	str     string
	integer int64

	replys []*Reply
}

// String ...
func (p *Reply) String() string {
	if p.err != nil {
		return fmt.Sprintf("Err: %v", p.err)
	}
	if p.null {
		return "NULL"
	}
	if p.str != "" {
		return fmt.Sprintf("String: %v", p.str)
	}
	if p.integer != 0 {
		return fmt.Sprintf("Integet: %v", p.integer)
	}
	if len(p.replys) > 0 {
		buf := new(bytes.Buffer)
		buf.WriteString("List:")
		for _, v := range p.replys {
			buf.WriteString("  ")
			buf.WriteString(v.String())
		}
		return buf.String()
	}
	return ""
}

// Integer ...
func (p *Reply) int() (int, error) {
	if p.err != nil {
		return 0, p.err
	}
	return int(p.integer), nil // TODO int64?
}

func (p *Reply) string() (NullString, error) {
	if p.err != nil {
		return NullString{}, p.err
	}
	if p.null {
		return NullString{}, nil
	}

	return NullString{String: p.str, Valid: true}, nil
}

func (p *Reply) fixBool() (bool, error) {
	if p.err == nil {
		return p.integer == 1, nil
	}
	return false, p.err
}

func (p *Reply) fixNilInt() (int, error) {
	if p.err != nil {
		return 0, p.err
	} else if p.null {
		return 0, ErrKeyNotExist
	}
	return int(p.integer), nil
}

func (p *Reply) fixFloat64() (float64, error) {
	if p.err != nil {
		return 0, p.err
	}
	return strconv.ParseFloat(p.str, 10)
}

func (p *Reply) fixNullStringSlice() ([]NullString, error) {
	if p.err != nil {
		return nil, p.err
	}

	var ns []NullString
	for _, v := range p.replys {
		if v.err != nil {
			return nil, v.err // TODO 这里真的有error吗
		}
		n, _ := v.string()
		ns = append(ns, n)
	}
	return ns, nil
}

func (p *Reply) fixStringSlice() ([]string, error) {
	if p.err != nil {
		return nil, p.err
	}

	var s []string
	for _, v := range p.replys {
		if v.err != nil {
			return nil, v.err // TODO 真的有吗
		}
		s = append(s, v.str)
	}
	return s, nil
}

func (p *Reply) fixMap() (map[string]string, error) {
	if p.err != nil {
		return nil, p.err
	}

	var s = make(map[string]string)
	for i := 0; i < len(p.replys); i += 2 {
		if p.replys[i].err != nil {
			return nil, p.replys[i].err // TODO 真的有吗
		}
		if p.replys[i+1].err != nil {
			return nil, p.replys[i+1].err // TODO 真的有吗
		}
		s[p.replys[i].str] = p.replys[i+1].str
	}
	return s, nil
}

func (p *Reply) fixGeoLocationSlice() ([]*GeoLocation, error) {
	if p.err != nil {
		return nil, p.err
	}
	var ss []*GeoLocation
	for _, v := range p.replys {
		if v.err != nil {
			return nil, v.err
		}
		if len(v.replys) < 2 {
			return nil, fmt.Errorf("expect 2 string to parse to geo")
		}
		longitude, err := stringToFloat64(v.replys[0].str)
		if err != nil {
			return nil, err
		}
		latitude, err := stringToFloat64(v.replys[1].str)
		if err != nil {
			return nil, err
		}
		ss = append(ss, &GeoLocation{Longitude: longitude, Latitude: latitude})
	}
	return ss, nil
}

func (p *Reply) fixSortedSetSlice() ([]*SortedSet, error) {
	if p.err != nil {
		return nil, p.err
	}
	var ss []*SortedSet
	for _, v := range p.replys {
		if v.err != nil {
			return nil, v.err
		}
		ss = append(ss, &SortedSet{Member: v.str})
	}
	return ss, nil
}

func (p *Reply) fixSortedSetSliceWithScores() ([]*SortedSet, error) {
	if p.err != nil {
		return nil, p.err
	}
	var ss []*SortedSet
	for i := 0; i < len(p.replys); i += 2 {
		if p.replys[i].err != nil {
			return nil, p.replys[i].err // TODO 真的有吗
		}
		if p.replys[i+1].err != nil {
			return nil, p.replys[i+1].err // TODO 真的有吗
		}
		score, err := strconv.Atoi(p.replys[i+1].str)
		if err != nil {
			return nil, err
		}
		ss = append(ss, &SortedSet{Member: p.replys[i].str, Score: score})
	}

	return ss, nil
}

func (p *Reply) fixFloat() (float64, error) {
	if p.err != nil {
		return 0, p.err
	}
	return strconv.ParseFloat(p.str, 64)
}

func errToReply(err error) *Reply {
	if err != nil {
		return &Reply{err: err}
	}
	return nil
}

func bytesToReply(bs []byte) *Reply {
	return &Reply{str: string(bs)}
}
