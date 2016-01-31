package atau

import (
	"github.com/Knetic/presilo"
)

/*
	Represents a "parameters" section, as found in every 'method' and sometimes the top-level document.
*/
type ParameterList struct {

	Parameters map[string]presilo.TypeSchema
}

/*
	Generates an overarching object schema which can be used to collect all values in this list as an object.
*/
func (this ParameterList) GenerateWrapperSchema() presilo.TypeSchema {

	var ret *presilo.ObjectSchema

	return ret
}

/*
	Gets a list of all parameters, with all required parameters first, followed by the rest.
	Both required and not-required parameters are returned in the order that most closely matches the "orderedParameters" list.
*/
func (this ParameterList) GetOrderedParameters() []string {

	var orderedParameters presilo.SortableStringArray

	for key, _ := range this.Parameters {
		orderedParameters = append(orderedParameters, key)
	}

	orderedParameters.Sort()
	return orderedParameters
}

/*
	Gets a list of all parameters, in the order that they should appear as given by the "orderedParameters" list.
*/
func (this ParameterList) GetOrderedParametersVerbatim() []string {
	return this.GetOrderedParameters()
}

/*
	Returns a list of only the parameters listed as "required".
	If these parameters can be found in the "orderedParameters" list, this will return the requested order. Or, if not all are in order,
	this returns the ordered ones first, followed by the rest.
*/
func (this ParameterList) GetRequiredParameters() []string {
	return []string{}
}
