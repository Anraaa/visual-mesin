package services

import (
	"github.com/anraaa/visual-mesin/internal/repository"
)

type GroupStat struct {
	GroupName     string `json:"group_name"`
	GroupColor    string `json:"group_color"`
	GroupIcon     string `json:"group_icon"`
	ResourceCount int    `json:"resource_count"`
	TotalRecords  int64  `json:"total_records"`
}

type DashboardSummary struct {
	TotalResources int              `json:"total_resources"`
	TotalRecords   int64            `json:"total_records"`
	GroupStats     []GroupStat      `json:"group_stats"`
}

type DashboardService struct {
	groupRepo *repository.ResourceGroupRepo
	querySvc  *ResourceQueryService
}

func NewDashboardService(groupRepo *repository.ResourceGroupRepo, querySvc *ResourceQueryService) *DashboardService {
	return &DashboardService{groupRepo: groupRepo, querySvc: querySvc}
}

func (s *DashboardService) GetSummary() (*DashboardSummary, error) {
	groups, err := s.groupRepo.List()
	if err != nil {
		return nil, err
	}

	summary := &DashboardSummary{
		TotalResources: 0,
		TotalRecords:   0,
		GroupStats:     make([]GroupStat, 0),
	}

	for _, group := range groups {
		gs := GroupStat{
			GroupName:  group.Name,
			GroupColor: group.Color,
			GroupIcon:  group.Icon,
		}

		for _, item := range group.Items {
			if !item.IsActive {
				continue
			}
			gs.ResourceCount++

			count, err := s.querySvc.GetResourceCount(item.ResourceName)
			if err == nil {
				gs.TotalRecords += count
			}
		}

		summary.TotalResources += gs.ResourceCount
		summary.TotalRecords += gs.TotalRecords
		summary.GroupStats = append(summary.GroupStats, gs)
	}

	return summary, nil
}
