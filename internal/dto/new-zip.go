package dto

type NewZipDTO struct {
	Owner string
	Name  string
	Url   string
}

func MapNativeToNewZipDTO(native interface{}) NewZipDTO {
	nativeMap, _ := native.(map[string]interface{})

	return NewZipDTO{
		Owner: nativeMap["owner"].(string),
		Name:  nativeMap["name"].(string),
		Url:   nativeMap["url"].(string),
	}
}
