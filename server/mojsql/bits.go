package mojsql

import "database/sql"

func fourbyte(s [4]int8) [4]byte {
	return [4]byte{byte(s[0]), byte(s[1]), byte(s[2]), byte(s[3])}
}

func fourint8(s [4]byte) [4]int8 {
	return [4]int8{int8(s[0]), int8(s[1]), int8(s[2]), int8(s[3])}
}

func bas3(s []bool) [3]bool {
	return [3]bool{s[0], s[1], s[2]}
}

func bas6(s []bool) [6]bool {
	return [6]bool{s[0], s[1], s[2], s[3], s[4], s[5]}
}

func yas4(s []byte) [4]byte {
	return [4]byte{s[0], s[1], s[2], s[3]}
}

func nullint64(d **int64, s sql.NullInt64) {
	if s.Valid {
		*d = new(int64)
		**d = s.Int64
		return
	}
	*d = nil
	return
}

func nullint8(d **int8, s sql.NullInt64) {
	if s.Valid {
		*d = new(int8)
		**d = int8(s.Int64)
		return
	}
	*d = nil
	return
}

func tonullint64(d *int64) sql.NullInt64 {
	var n sql.NullInt64
	n.Valid = (d != nil)
	if !n.Valid {
		return n
	}
	n.Int64 = *d
	return n
}

func tonullint8(d *int8) sql.NullInt64 {
	var n sql.NullInt64
	n.Valid = (d != nil)
	if n.Valid {
		return n
	}
	n.Int64 = int64(*d)
	return n
}

func makebit(b bool) byte {
	if b {
		return '1'
	}
	return '0'
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
