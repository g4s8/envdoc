package ast

func (r *FieldTypeRefKind) ScanStr(s string) bool {
	for i := 0; i < len(_FieldTypeRefKind_index)-1; i++ {
		from, to := _FieldTypeRefKind_index[i], _FieldTypeRefKind_index[i+1]
		if s == _FieldTypeRefKind_name[from:to] {
			*r = FieldTypeRefKind(i)
			return true
		}
	}
	return false
}
