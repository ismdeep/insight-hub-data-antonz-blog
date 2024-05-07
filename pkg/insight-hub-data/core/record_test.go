package core

import (
	"reflect"
	"testing"
	"time"
)

func TestRecordUnmarshal(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    Record
		wantErr bool
	}{
		{
			name: "",
			args: args{
				line: "607c0d932d085f76c670cc2c0c065b9e982452d09107ecd348f8e522f2c11fc3|1715087479560048951|demo-source|https%3A%2F%2Fismdeep.github.io%2Fhi.txt|Hello+%7C+World|x98|%3Cp%3EHi%2C+there.%3Cp%3E",
			},
			want: Record{
				Source:      "demo-source",
				Link:        "https://ismdeep.github.io/hi.txt",
				Title:       "Hello | World",
				Author:      "x98",
				Content:     "<p>Hi, there.<p>",
				PublishedAt: time.Unix(1715087479, 560048951),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RecordUnmarshal(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecordUnmarshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
