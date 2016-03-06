package mojsql

import "database/sql"

func nullint64(d **int64, s sql.NullInt64) {
	if s.Valid {
		*d = new(int64)
		**d = s.Int64()
	} else {
		*d = nil
	}
}

func nullint8(d **int8, s sql.NullInt64) {
	if s.Valid {
		*d = new(int8)
		**d = int8(s.Int64())
	} else {
		*d = nil
	}
}

func tonullint64(d *int64) sql.NullInt64 {
	var n sql.NullInt64
	n.Valid = (d != nil)
	if n.Valid {
		n.Int64 = *d
	}
}

func tonullint8(d *int8) sql.NullInt64 {
	var n sql.NullInt64
	n.Valid = (d != nil)
	if n.Valid {
		n.Int64 = int8(*d)
	}
}

func makebit(b bool) byte {
	if b {
		return '1'
	} else {
		return '0'
	}
}

func tobit(b []bool) []byte {
	a := make([]byte, 0, len(b))
	for i := 0; i < len(b); i++ {
		a = append(a, makebit(b[i]))
	}
	return a
}

func makebool(b byte) bool {
	switch b {
	case '1':
		return true
	case '0':
		return false
	default:
		panic(b)
	}
}

func tobool(b []byte) []bool {
	a := make([]bool, 0, len(b))
	for i := 0; i < len(b); i++ {
		a = append(a, makebool(b[i]))
	}
	return a
}
