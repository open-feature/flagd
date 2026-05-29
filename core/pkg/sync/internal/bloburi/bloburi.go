// Package bloburi splits and rejoins blob sync URIs (s3://, gs://, azblob://).
package bloburi

import (
	"regexp"
	"strings"
)

// Split parses "scheme://bucket/key?opt=1" into bucket URL ("scheme://bucket?opt=1")
// and object key ("key"). The query moves to the bucket URL because gocloud blob
// drivers read driver options (e.g. s3blob use_path_style, region) from there.
// schemeRegex must match through the first "/" after the scheme (e.g. "^s3://.+?/").
func Split(uri string, schemeRegex *regexp.Regexp) (bucket, object string) {
	raw, query, hasQuery := strings.Cut(uri, "?")
	bucket = schemeRegex.FindString(raw)
	object = schemeRegex.ReplaceAllString(raw, "")
	if hasQuery && query != "" {
		bucket = strings.TrimSuffix(bucket, "/") + "?" + query
	}
	return bucket, object
}

// Join is the inverse of Split. The reconstructed URI must match what was
// registered in the store (see flagd#1971).
func Join(bucket, object string) string {
	i := strings.Index(bucket, "?")
	if i < 0 {
		return bucket + object
	}
	base := strings.TrimSuffix(bucket[:i], "/") + "/"
	return base + object + bucket[i:]
}
