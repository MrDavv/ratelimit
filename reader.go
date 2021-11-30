// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

package ratelimit

import (
	"io"
	"net/http"

	"github.com/astaxie/beego/context"
)

type reader struct {
	r      io.Reader
	bucket *Bucket
}

// Reader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Reader(r io.Reader, bucket *Bucket) io.Reader {
	return &reader{
		r:      r,
		bucket: bucket,
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}
	r.bucket.Wait(int64(n))
	return n, err
}

type writer struct {
	w      io.Writer
	bucket *Bucket
}

// Writer returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Writer(w io.Writer, bucket *Bucket) io.Writer {
	return &writer{
		w:      w,
		bucket: bucket,
	}
}

func (w *writer) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}

type writer2 struct {
	w      *context.Response
	bucket *Bucket
}

func Writer2(w *context.Response, bucket *Bucket) http.ResponseWriter {
	return &writer2{
		w:      w,
		bucket: bucket,
	}
}

func (w *writer2) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}
func (w *writer2) Header() http.Header {
	return w.w.ResponseWriter.Header()
}
func (w *writer2) WriteHeader(code int) {
	if w.w.Status > 0 {
		//prevent multiple response.WriteHeader calls
		return
	}
	w.w.Status = code
	w.w.Started = true
	w.w.ResponseWriter.WriteHeader(code)
}
