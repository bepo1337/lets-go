package models

import (
	"letsgo.bepo1337/internal/assert"
	"testing"
)

type testCaseData struct {
	name   string
	userID int
	want   bool
}

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	tests := testCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDb(t)
			userModel := &UserModel{db}
			exists, err := userModel.Exists(tt.userID)
			assert.NilError(t, err)
			assert.Equal(t, exists, tt.want)
		})
	}

}

func testCases() []testCaseData {
	return []testCaseData{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
	}
}
