package dto

type ProjectStatusDTO struct {
	Id     string
	Status string
}

// MapNativeToProjectStatusDTO should get only native data checked by avro codec!
func MapNativeToProjectStatusDTO(native interface{}) ProjectStatusDTO {
	nativeMap, _ := native.(map[string]interface{})

	return ProjectStatusDTO{
		Id:     nativeMap["id"].(string),
		Status: nativeMap["status"].(string),
	}
}
