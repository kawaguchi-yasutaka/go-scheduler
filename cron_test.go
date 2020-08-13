package scheduler

import (
	"reflect"
	"testing"
)

func TestValidateMinuteFormat(t *testing.T) {
	cases := []struct{
		caseName string
		in string
		want bool
	}{
		{
			caseName: "*を受け取ったらtrue",
			in: "*",
			want: true,
		},
		{
			caseName: "数字の文字列を受け取ったらtrue",
			in: "0",
			want: true,
		},
		{
			caseName: "カンマ（,）で区切られて複数の数字の文字列を受け取ったらtrue",
			in: "0,59",
			want: true,
		},
		{
			caseName: "ハイフンで指定された数字の文字列を受け取ったらtrue",
			in: "0-59",
			want: true,
		},
		{
			caseName: "ハイフンで指定された数字と複数の文字列が複合されているものを受け取ったらtrue",
			in: "1,0-59,2-30,2,2",
			want: true,
		},
		{
			caseName: "*/[0-59]",
			in: "*/59",
			want: true,
		},
		{
			caseName: "[0-59]-[0-59]/[0-59]",
			in: "0-59/59",
			want: true,
		},
	}

	for _,tt := range cases {
		result := ValidateMinuteFormat(tt.in)
		if  result != tt.want {
			t.Errorf("caseName %v, want: %v result: %v",tt.caseName,tt.want,result)
		}
	}
}



func TestParseMinute(t *testing.T) {
	cases := []struct{
		caseName string
		in string
		want cronField
		wantErr error
	}{
		{
			caseName: "*を受け取ったら、int64のスライスで返す",
			in: "*",
			want: cronField{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,
				22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,
				44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59},
			wantErr: nil,
		},
		{
			caseName: "数字の文字列を受け取ったら、int64のスライスで返す",
			in: "0",
			want: cronField{0},
			wantErr: nil,
		},
		{
			caseName: "複数の数字の文字列を受け取ったら、int64のスライスで返す",
			in: "0,1",
			want: cronField{0,1},
			wantErr: nil,
		},
		{
			caseName: "ハイフンで区切られた範囲を受け取ったら、int64のスライスで返す",
			in: "0-10",
			want: cronField{0,1,2,3,4,5,6,7,8,9,10},
			wantErr: nil,
		},
		{
			caseName: "ハイフンで指定された数字と複数の文字列が複合されているものを受け取ったらtrue",
			in: "1,2-30,3",
			want: cronField{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,
				17,18,19,20,21,22,23,24,25,26,27,28,29,30},
			wantErr: nil,
		},
		{
			caseName: "*/5",
			in: "*/5",
			want: cronField{0,5,10,15,20,25,30,35,40,45,50,55},
			wantErr: nil,
		},
		{
			caseName: "0-40/5",
			in: "0-40/5",
			want: cronField{0,5,10,15,20,25,30,35,40},
			wantErr: nil,
		},
	}

	for _,tt := range cases {
		err,result := ParseMinute(tt.in)
		if err != tt.wantErr {
			t.Errorf("\ncaseName %v\n wanterr: %v\n resut err: %v\n",tt.caseName,tt.wantErr,err)
		}
		if !reflect.DeepEqual(result, tt.want) {
			t.Errorf("\ncaseName %v\n want: %v\n result: %v\n",tt.caseName,tt.want,result)
		}
	}
}

