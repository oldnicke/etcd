package httpproxy

import (
	"net/url"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestNewDirectorScheme(t *testing.T) {
	tests := []struct {
		urls []string
		want []string
	}{
		{
			urls: []string{"http://192.0.2.8:4002", "http://example.com:8080"},
			want: []string{"http://192.0.2.8:4002", "http://example.com:8080"},
		},
		{
			urls: []string{"https://192.0.2.8:4002", "https://example.com:8080"},
			want: []string{"https://192.0.2.8:4002", "https://example.com:8080"},
		},

		// accept urls without a port
		{
			urls: []string{"http://192.0.2.8"},
			want: []string{"http://192.0.2.8"},
		},

		// accept urls even if they are garbage
		{
			urls: []string{"http://."},
			want: []string{"http://."},
		},
	}

	for i, tt := range tests {
		uf := func() []string {
			return tt.urls
		}
		got := newDirector(uf, time.Minute, time.Minute)

		var gep []string
		for _, ep := range got.ep {
			gep = append(gep, ep.URL.String())
		}
		sort.Strings(tt.want)
		sort.Strings(gep)
		if !reflect.DeepEqual(tt.want, gep) {
			t.Errorf("#%d: want endpoints = %#v, got = %#v", i, tt.want, gep)
		}
	}
}

func TestDirectorEndpointsFiltering(t *testing.T) {
	d := director{
		ep: []*endpoint{
			{
				URL:       url.URL{Scheme: "http", Host: "192.0.2.5:5050"},
				Available: false,
			},
			{
				URL:       url.URL{Scheme: "http", Host: "192.0.2.4:4000"},
				Available: true,
			},
		},
	}

	got := d.endpoints()
	want := []*endpoint{
		{
			URL:       url.URL{Scheme: "http", Host: "192.0.2.4:4000"},
			Available: true,
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("directed to incorrect endpoint: want = %#v, got = %#v", want, got)
	}
}
