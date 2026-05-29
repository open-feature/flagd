package bloburi

import (
	"regexp"
	"testing"
)

var schemeRegexes = map[string]*regexp.Regexp{
	"s3":     regexp.MustCompile("^s3://.+?/"),
	"gs":     regexp.MustCompile("^gs://.+?/"),
	"azblob": regexp.MustCompile("^azblob://.+?/"),
}

func TestSplit(t *testing.T) {
	tests := map[string]struct {
		uri    string
		scheme string
		bucket string
		object string
	}{
		"s3 simple": {
			uri:    "s3://my-bucket/flags.json",
			scheme: "s3",
			bucket: "s3://my-bucket/",
			object: "flags.json",
		},
		"s3 with query (use_path_style)": {
			uri:    "s3://my-bucket/example_flags.json?use_path_style=true&region=garage&endpoint=http://127.0.0.1:3900",
			scheme: "s3",
			bucket: "s3://my-bucket?use_path_style=true&region=garage&endpoint=http://127.0.0.1:3900",
			object: "example_flags.json",
		},
		"gs simple": {
			uri:    "gs://my-bucket/path/to/object",
			scheme: "gs",
			bucket: "gs://my-bucket/",
			object: "path/to/object",
		},
		"azblob simple": {
			uri:    "azblob://my-bucket/flags.yaml",
			scheme: "azblob",
			bucket: "azblob://my-bucket/",
			object: "flags.yaml",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			bucket, object := Split(tt.uri, schemeRegexes[tt.scheme])
			if bucket != tt.bucket {
				t.Errorf("Split bucket = %q, want %q", bucket, tt.bucket)
			}
			if object != tt.object {
				t.Errorf("Split object = %q, want %q", object, tt.object)
			}
		})
	}
}

// TestSplitJoinInverse ensures that Join is the inverse of Split
func TestSplitJoinInverse(t *testing.T) {
	uris := []string{
		"s3://b/o",
		"s3://b/path/to/object",
		"s3://b/o?use_path_style=true",
		"s3://b/o?use_path_style=true&region=garage&endpoint=http://127.0.0.1:3900",
		"s3://b/path/to/object?a=1&b=2",
		"gs://b/o",
		"gs://b/path/to/object",
		"azblob://b/o",
		"azblob://b/flags.yaml",
	}
	for _, uri := range uris {
		t.Run(uri, func(t *testing.T) {
			for _, reg := range schemeRegexes {
				if !reg.MatchString(uri) {
					continue
				}
				bucket, object := Split(uri, reg)
				if got := Join(bucket, object); got != uri {
					t.Errorf("Join(Split(%q)) = %q, want %q (bucket=%q object=%q)",
						uri, got, uri, bucket, object)
				}
			}
		})
	}
}
