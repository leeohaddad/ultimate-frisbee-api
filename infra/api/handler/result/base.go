package result

type HTTP struct {
	StatusCode     int
	ResponseType   ResponseBodyType
	JSONResponse   interface{}
	StringResponse string
}

// ResponseBodyType is the type of the response body.
type ResponseBodyType string

// responseBodyTypeList represents a list of possible types for response bodies.
type responseBodyTypeList struct {
	JSON   ResponseBodyType
	String ResponseBodyType
}

// ResponseBodyTypes represents the names of possible types for response bodies.
var ResponseBodyTypes = &responseBodyTypeList{
	JSON:   "JSON",
	String: "String",
}
