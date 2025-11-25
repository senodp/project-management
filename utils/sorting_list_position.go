package utils

import (
	"github.com/google/uuid"
	"github.com/senodp/project-management/models"
)

func SortListByPosition(lists []models.List, order []uuid.UUID) []models.List {
	if len(order) == 0{
		return lists
	}
	
	ordered := make([]models.List, 0, len(order))

	listMap := make(map[uuid.UUID]models.List)
	for _, l := range lists {
		listMap[l.PublicID] = l
	}

	// urutkan sesuai order nya
	for _, id := range order {
		if list, ok := listMap[id]; ok {
			ordered = append(ordered, list)
		}
	}
	return ordered
}