package spruce

func EntrySet(key, value []byte, ti int) []byte {
	out := make([]byte, 11)
	out[0] = 1
	out[1] = 2
	tm := ParsingExpirationDate(ti).([]byte)
	out[2] = tm[0]
	out[3] = tm[1]
	for _, v := range key {
		out = append(out, v)
	}
	out = append(out, 0)
	for _, v := range value {
		out = append(out, v)
	}
	return out
}
func EntryGet(key []byte) []byte {
	out := make([]byte, 11)
	out[0] = 2
	out[1] = 2
	for _, v := range key {
		out = append(out, v)
	}
	return out
}
func EntryHashSet(key, value []byte, ti int) []byte {
	out := make([]byte, 11)
	out[0] = 1
	out[1] = 2
	tm := ParsingExpirationDate(ti).([]byte)
	out[2] = tm[0]
	out[3] = tm[1]
	for _, v := range MD5(key) {
		out = append(out, v)
	}
	out = append(out, 0)
	for _, v := range value {
		out = append(out, v)
	}
	return out
}
func EntryHashGet(key []byte) []byte {
	out := make([]byte, 11)
	out[0] = 2
	out[1] = 2
	for _, v := range MD5(key) {
		out = append(out, v)
	}
	return out
}
func EntryDelete(key []byte) []byte {
	out := make([]byte, 11)
	out[0] = 0
	out[1] = 2
	for _, v := range key {
		out = append(out, v)
	}
	return out
}
func EntryHashDelete(key []byte) []byte {
	out := make([]byte, 11)
	out[0] = 0
	out[1] = 2
	for _, v := range MD5(key) {
		out = append(out, v)
	}
	return out
}
