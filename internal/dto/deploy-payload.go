package dto

type DeployPayloadDTO struct {
	Owner string
	Name  string
	Url   string
}

func MapNativeToDeployPayloadDTO(native interface{}) DeployPayloadDTO {
	nativeMap, _ := native.(map[string]interface{})

	return DeployPayloadDTO{
		Owner: nativeMap["owner"].(string),
		Name:  nativeMap["name"].(string),
		Url:   nativeMap["url"].(string),
	}
}
