package requests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
)

func TestAddTextReq(t *testing.T) {
	type args struct {
		data storage.TextData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusOK},
			want: false, wantErr: true},
		{name: "Ok",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusAccepted},
			want: true, wantErr: false},
		{name: "Witherr2",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := AddTextReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTextReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddTextReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTextReq(t *testing.T) {
	type args struct {
		data storage.TextData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusOK},
			want: false, wantErr: true},
		{name: "Ok",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusAccepted},
			want: true, wantErr: false},
		{name: "Witherr2",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := UpdateTextReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTextReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateTextReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteTextReq(t *testing.T) {
	type args struct {
		data storage.TextData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusOK},
			want: true, wantErr: false},
		{name: "Ok",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusAccepted},
			want: false, wantErr: true},
		{name: "Witherr2",
			args: args{data: storage.TextData{Text: "Text", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := DeleteTextReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTextReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteTextReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncTextReq(t *testing.T) {
	tests := []struct {
		name     string
		want     bool
		wantErr  bool
		resp     int
		respData []storage.TextData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "NoData1",
			resp: http.StatusAccepted,
			want: false, wantErr: true},
		{name: "WithData",
			resp:     http.StatusAccepted,
			respData: []storage.TextData{{ID: 1, Text: "111", Comment: "333"}},
			want:     true, wantErr: false},
		{name: "NoData2",
			resp: http.StatusOK,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1)
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(tt.resp)
				if len(tt.respData) != 0 {
					mock.EXPECT().AddData(tt.respData[0]).Times(1)
					mock.EXPECT().CreateNewRecord(gomock.Any()).Times(1)
					if !tt.wantErr {
						mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(tt.respData, nil)
					} else {
						mock.EXPECT().GetAllRecords(gomock.Any()).Times(1)
					}
					data, _ := json.Marshal(tt.respData)
					w.Write(data)
				} else {
					w.Write([]byte(`{"message": "Hello, World!"}`))
				}
			}))
			storage.ServURL = mockServer.URL

			storage.TLiteS = mock
			got, err := SyncTextReq()
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncTextReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SyncTextReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_textSliceToMap(t *testing.T) {
	type args struct {
		slice []storage.TextData
	}
	tests := []struct {
		name  string
		args  args
		want  map[int]storage.TextData
		want1 map[int]int
	}{
		{name: "11", args: args{slice: []storage.TextData{{ID: 1, Text: "1111"}}}, want: map[int]storage.TextData{1: {ID: 1, Text: "1111"}}, want1: map[int]int{1: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := textSliceToMap(tt.args.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("textSliceToMap() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("textSliceToMap() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_compareUnorderedTextSlices(t *testing.T) {
	type args struct {
		s []storage.TextData
		c []storage.TextData
	}
	tests := []struct {
		name string
		args args
		want []storage.TextData
	}{
		{name: "1", args: args{s: []storage.TextData{1: {ID: 1, Text: "1111"}},
			c: []storage.TextData{1: {ID: 1, Text: "1111"}, 2: {ID: 2, Text: "2222"}}},
			want: []storage.TextData{0: {ID: 2, Text: "2222"}}},
		{name: "2", args: args{s: []storage.TextData{1: {ID: 1, Text: "1111"}},
			c: []storage.TextData{1: {ID: 1, Text: "1111"}}},
			want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareUnorderedTextSlices(tt.args.s, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compareUnorderedTextSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}
