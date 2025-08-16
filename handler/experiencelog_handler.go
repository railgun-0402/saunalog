package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	domain "saunalog/domain/experience_log"
	facility "saunalog/domain/facility"
	"saunalog/usecase"
	"time"
)

type ExperienceLogHandler struct {
	Usecase usecase.ExperienceLogUseCase
}

type ExperienceLogRequest struct {
	UserID          string `json:"user_id"`
	FacilityID      string `json:"facility_id"`
	Date            string `json:"date"` // e.g. "2025-08-02"
	CongestionLevel int    `json:"congestion_level"`
	TotonoiLevel    int    `json:"totonoi_level"`
	CostPerformance int    `json:"cost_performance"`
	ServiceQuality  int    `json:"service_quality"`
	Comment         string `json:"comment"`
}

type ExperienceLogResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (h *ExperienceLogHandler) CreateExperienceLog(c echo.Context) error {
	var req ExperienceLogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "日付フォーマットが不正です")
	}

	logID := domain.ExperienceID(uuid.New().String())
	userID := req.UserID
	facilityID := facility.SaunaFacilityID(req.FacilityID)

	log, err := domain.NewExperienceLog(domain.ExperienceLog{
		ID:              logID,
		UserID:          userID,
		SaunaFacilityID: facilityID,
		Date:            date,
		CongestionLevel: req.CongestionLevel,
		TotonoiLevel:    req.TotonoiLevel,
		CostPerformance: req.CostPerformance,
		Comment:         req.Comment,
		CreatedAt:       date,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.Usecase.CreateExperienceLog(c.Request().Context(), log); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "体験ログの保存に失敗しました")
	}

	return c.JSON(http.StatusCreated, ExperienceLogResponse{
		ID:      string(log.ID),
		Message: "体験ログの投稿に成功しました",
	})
}
