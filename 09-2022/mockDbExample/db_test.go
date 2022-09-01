package db

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))

	if v := GetFromDB(m, "Tom"); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, nil)

	if v := GetFromDB(m, "Tom"); v != 100 {
		t.Fatal("expected -1, but got", v)
	}
}
