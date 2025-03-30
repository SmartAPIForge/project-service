package dto

type ProjectStatusDTO struct {
	Id     string
	Status string
}

func MapNativeToProjectStatusDTO(native interface{}) ProjectStatusDTO {
	nativeMap, _ := native.(map[string]interface{})

	return ProjectStatusDTO{
		Id:     nativeMap["id"].(string),
		Status: nativeMap["status"].(string),
	}
}
