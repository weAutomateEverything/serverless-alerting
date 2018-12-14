package common

func GetCorsHeaders() map[string]string{
	return map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Access-Control-Allow-Credentials": "true",
	}
}
