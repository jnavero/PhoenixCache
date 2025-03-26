package utils

import (
	"bytes"
	"compress/gzip"
	"io"
)

func CompressData(data []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, _ = w.Write(data)
	w.Close()
	return b.Bytes()
}

func DecompressData(compressed []byte) ([]byte, error) {
	b := bytes.NewBuffer(compressed)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	return decompressed, err
}

// truncateString acorta un string a `n` caracteres
func TruncateString(s string, n int) string {
	if len(s) > n {
		return s[:n-3] + "..."
	}
	return s
}
