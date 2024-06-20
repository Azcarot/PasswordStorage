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

func TestAddFileReq(t *testing.T) {
	type args struct {
		data storage.FileData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusOK},
			want: false, wantErr: true},
		{name: "Ok",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusAccepted},
			want: true, wantErr: false},
		{name: "Witherr2",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := AddFileReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFileReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddFileReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateFileReq(t *testing.T) {
	type args struct {
		data storage.FileData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusOK},
			want: false, wantErr: true},
		{name: "Ok",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusAccepted},
			want: true, wantErr: false},
		{name: "Witherr2",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := UpdateFileReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateFileReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateFileReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteFileReq(t *testing.T) {
	type args struct {
		data storage.FileData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusOK},
			want: true, wantErr: false},
		{name: "Ok",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusAccepted},
			want: false, wantErr: true},
		{name: "Witherr2",
			args: args{data: storage.FileData{FileName: "Name", Path: "Path", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := DeleteFileReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFileReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteFileReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncFileReq(t *testing.T) {
	tests := []struct {
		name     string
		want     bool
		wantErr  bool
		resp     int
		respData []storage.FileData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "NoData1",
			resp: http.StatusAccepted,
			want: false, wantErr: true},
		{name: "WithData",
			resp:     http.StatusAccepted,
			respData: []storage.FileData{{ID: 1, FileName: "111", Path: "333"}},
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

			storage.FLiteS = mock
			got, err := SyncFileReq()
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncCardReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SyncCardReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileSliceToMap(t *testing.T) {
	type args struct {
		slice []storage.FileData
	}
	tests := []struct {
		name  string
		args  args
		want  map[int]storage.FileData
		want1 map[int]int
	}{
		{name: "11", args: args{slice: []storage.FileData{{ID: 1, FileName: "1111"}}}, want: map[int]storage.FileData{1: {ID: 1, FileName: "1111"}}, want1: map[int]int{1: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := fileSliceToMap(tt.args.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileSliceToMap() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("fileSliceToMap() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_compareUnorderedFileSlices(t *testing.T) {
	type args struct {
		s []storage.FileData
		c []storage.FileData
	}
	tests := []struct {
		name string
		args args
		want []storage.FileData
	}{
		{name: "1", args: args{s: []storage.FileData{1: {ID: 1, FileName: "1111"}},
			c: []storage.FileData{1: {ID: 1, FileName: "1111"}, 2: {ID: 2, FileName: "2222"}}},
			want: []storage.FileData{0: {ID: 2, FileName: "2222"}}},
		{name: "2", args: args{s: []storage.FileData{1: {ID: 1, FileName: "1111"}},
			c: []storage.FileData{1: {ID: 1, FileName: "1111"}}},
			want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareUnorderedFileSlices(tt.args.s, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compareUnorderedfileSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}
