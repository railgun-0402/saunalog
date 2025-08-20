package domain

import (
	"testing"
	"time"
)

func validExperienceParams() ExperienceLog {
	return ExperienceLog{
		ID:              "e-123",
		UserID:          "u-1",
		SaunaFacilityID: "s-1",
		Date:            time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		CongestionLevel: 3,
		CostPerformance: 4,
		TotonoiLevel:    5,
		Comment:         "great",
	}
}

func TestNewExperienceLog_Success(t *testing.T) {
	in := validExperienceParams()
	before := time.Now()
	log, err := NewExperienceLog(in)
	after := time.Now()
	if err != nil {
		t.Fatalf("NewExperienceLog() error = %v, want nil", err)
	}
	if log == nil {
		t.Fatalf("NewExperienceLog() returned nil log")
	}
	if log.Date != in.Date {
		t.Errorf("Date mismatch got %v want %v", log.Date, in.Date)
	}
	if log.CongestionLevel != in.CongestionLevel || log.CostPerformance != in.CostPerformance || log.TotonoiLevel != in.TotonoiLevel {
		t.Errorf("rating fields mismatch")
	}
	if log.CreatedAt.Before(before) || log.CreatedAt.After(after) {
		t.Errorf("CreatedAt out of range")
	}
}

func TestNewExperienceLog_InvalidRating(t *testing.T) {
	cases := []struct {
		name   string
		modify func(*ExperienceLog)
		want   string
	}{
		{
			name:   "invalid totonoi",
			modify: func(e *ExperienceLog) { e.TotonoiLevel = 6 },
			want:   "整い度は1〜5で指定してください",
		},
		{
			name:   "invalid congestion",
			modify: func(e *ExperienceLog) { e.CongestionLevel = 0 },
			want:   "混雑度は1〜5で指定してください",
		},
		{
			name:   "invalid cost performance",
			modify: func(e *ExperienceLog) { e.CostPerformance = 10 },
			want:   "コスパは1〜5で指定してください",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := validExperienceParams()
			tc.modify(&in)
			log, err := NewExperienceLog(in)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if log != nil {
				t.Fatalf("expected nil log, got %#v", log)
			}
			if err.Error() != tc.want {
				t.Errorf("error message mismatch got %q want %q", err.Error(), tc.want)
			}
		})
	}
}
