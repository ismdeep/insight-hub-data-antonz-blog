package core

import (
	"bytes"
	_ "embed"
	"os"
	"testing"
	"time"
)

//go:embed data.example.txt
var exampleData []byte

func TestStore_Write(t *testing.T) {
	s := NewStore(os.Stdout)

	if err := s.Load(bytes.NewReader(exampleData)); err != nil {
		t.Errorf("failed to load data: %s", err)
		t.FailNow()
	}

	if err := s.Save(Record{
		Source:      "demo-source",
		Link:        "https://ismdeep.github.io/hi.txt",
		Title:       "Hello | World",
		Author:      "x98",
		Content:     "<p>Hi, there.<p>",
		PublishedAt: time.Now(),
	}); err != nil {
		t.Errorf("failed to run write: %v", err)
		t.FailNow()
	}
}

func TestLinkIsTidy(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				link: "https://ismdeep.github.io/",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				link: "http://127.0.0.1/hello/world/",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				link: "https://ismdeep.github.io//hi",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LinkIsTidy(tt.args.link); got != tt.want {
				t.Errorf("LinkIsTidy() = %v, want %v", got, tt.want)
			}
		})
	}
}
