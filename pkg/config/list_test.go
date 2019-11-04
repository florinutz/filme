package config

import (
	"testing"
)

func TestGet1337xListUrlFromArgs(t *testing.T) {
	type args struct {
		args       []string
		defaultUrl string
	}
	tests := []struct {
		name        string
		args        args
		wantListUrl string
	}{
		{
			name: "search", args: struct {
				args       []string
				defaultUrl string
			}{
				args:       []string{"mama", "are", "mere"},
				defaultUrl: "http://something.com",
			},
			wantListUrl: "https://1337x.to/srch?search=mama+are+mere",
		},
		{
			name: "url", args: struct {
				args       []string
				defaultUrl string
			}{
				args:       []string{"https://1337x.to/smth"},
				defaultUrl: "http://something.com",
			},
			wantListUrl: "https://1337x.to/smth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotListUrl := Get1337xListUrlFromArgs(tt.args.args); gotListUrl != tt.wantListUrl {
				t.Errorf("Get1337xListUrlFromArgs() = %v, want %v", gotListUrl, tt.wantListUrl)
			}
		})
	}
}

func TestGetGoogleUrlFromArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		results int
		lang    string
		want    string
		wantErr bool
	}{
		{
			name:    "search",
			args:    []string{"mama", "are", "mere"},
			results: 30,
			lang:    "lang_en",
			want:    "https://www.google.com/search?num=30&lr=lang_en&q=imdb|rottentomatoes|metacritic|netflix|wikipedia+mama+are+mere",
			wantErr: false,
		},
		{
			name:    "not enough params",
			args:    []string{},
			results: 30,
			lang:    "lang_en",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGoogleUrlFromArgs(tt.args, tt.results, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGoogleUrlFromArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGoogleUrlFromArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetUrl(t *testing.T) {
	type args struct {
		args        []string
		domain      string
		templateUrl string
		defaultUrl  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "match",
			args: args{
				args:        []string{"something.com mama are mere"},
				domain:      `something.com`,
				templateUrl: "http://something.com/?s=%s",
				defaultUrl:  "",
			},
			want: "http://something.com/?s=something.com+mama+are+mere",
		},
		{
			name: "nomatch",
			args: args{
				args:        []string{"mama are mere"},
				domain:      `caca.com`,
				templateUrl: "http://something.com/?s=%s",
				defaultUrl:  "",
			},
			want: "http://something.com/?s=mama+are+mere",
		},
		{
			name: "fallback to default",
			args: args{
				args:        nil,
				domain:      `caca.com`,
				templateUrl: "http://something.com/?s=%s",
				defaultUrl:  "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUrl(tt.args.args, tt.args.domain, tt.args.templateUrl, tt.args.defaultUrl); got != tt.want {
				t.Errorf("getUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
