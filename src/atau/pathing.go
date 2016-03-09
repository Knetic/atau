package atau

import (
	"strings"
	"fmt"
)

/*
	Interpolates querystring parameters for the given [fullPath].
*/
func interpolatePath(api *API, method Method, fullPath string) string {

	var placeholder, replacement string

	// replace parameters as referenced in paths
	for _, parameter := range method.PathParameters {

		placeholder = fmt.Sprintf("{%s}", parameter)
		replacement = fmt.Sprintf("\"+%s+\"", parameter)
		fullPath = strings.Replace(fullPath, placeholder, replacement, -1)
	}

	return fullPath
}

/*
	Paths can have parameter placeholders baked into them.
	This provides a full path to a specific resource given an API's base path, the method's path,
	and interpolated with the correct variable names for all parameters.
*/
func resolvePath(api *API, method Method) string {

	if(strings.HasSuffix(api.BaseURL, "/")) {
		return api.BaseURL + method.Path
	}
	return api.BaseURL + "/" + method.Path
}

/*
	Appends the appropriate querystring to the given [fullPath]
*/
func appendQuerystringPath(api *API, method Method, fullPath string) string {

	var querystrings []string

	// set querystring
	for _, key := range method.QueryParameters {
		querystrings = append(querystrings, fmt.Sprintf("%s=\"+%s+\"", key, key))
	}

	if(len(querystrings) > 0) {
		fullPath = fullPath + "?" + strings.Join(querystrings, "&")
	}

	return fullPath
}
