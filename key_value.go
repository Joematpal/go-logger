package logger

type KV struct {
	key   string
	value interface{}
}

func (k KV) Key() string {
	return k.key
}

func (k KV) Value() interface{} {
	return k.value
}

type Map map[string]interface{}

func (m Map) ToFields() []Field {
	out := []Field{}
	for key, value := range m {
		out = append(out, KV{key, value})
	}
	return out
}
