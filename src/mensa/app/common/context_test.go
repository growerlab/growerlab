package common

import (
	"net/url"
	"reflect"
	"testing"
)

func TestBuildContextFromHTTP(t *testing.T) {
	type args struct {
		uri *url.URL
	}

	uri, _ := url.Parse("https://github.com/growerlab/growerlab/src/mensa.git")
	uriAuth, _ := url.Parse("https://moli:pwd@github.com/growerlab/growerlab/src/mensa.git")

	tests := []struct {
		name string
		args args
		want *Context
	}{
		{
			name: "normal",
			args: args{
				uri: uri,
			},
			want: &Context{
				Type:       ProtTypeHTTP,
				RawURL:     "https://github.com/growerlab/growerlab/src/mensa.git",
				RequestURL: uri,
				RepoOwner:  "growerlab",
				RepoName:   "mensa",
			},
		},
		{
			name: "auth",
			args: args{
				uri: uriAuth,
			},
			want: &Context{
				Type:       ProtTypeHTTP,
				RawURL:     "https://moli:pwd@github.com/growerlab/growerlab/src/mensa.git",
				RequestURL: uriAuth,
				RepoOwner:  "growerlab",
				RepoName:   "mensa",
				Operator: &Operator{
					HttpUser: url.UserPassword("moli", "pwd"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildContextFromHTTP(tt.args.uri)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nBuildContextFromHTTP() = \n%+v,\n want \n%+v", got, tt.want)
			}
		})
	}
}
