package form3

import "testing"

func Test_ToPointer(t *testing.T) {
	// Test with an integer value
	var intValue int = 42
	var intPointer *int = ToPointer(intValue)
	if *intPointer != intValue {
		t.Errorf("ToPointer did not return a pointer to the correct integer value. Expected %v but got %v", intValue, *intPointer)
	}

	// Test with a string value
	var strValue string = "Hello, world!"
	var strPointer *string = ToPointer(strValue)
	if *strPointer != strValue {
		t.Errorf("ToPointer did not return a pointer to the correct string value. Expected %v but got %v", strValue, *strPointer)
	}

	// Test with a boolean value
	var boolValue bool = true
	var boolPointer *bool = ToPointer(boolValue)
	if *boolPointer != boolValue {
		t.Errorf("ToPointer did not return a pointer to the correct boolean value. Expected %v but got %v", boolValue, *boolPointer)
	}

	// Test with a nil value
	var errValue error
	var errPointer *error = ToPointer(errValue)
	if *errPointer != nil {
		t.Errorf("ToPointer did not return a nil pointer for a nil value. Expected %v but got %v", nil, errPointer)
	}
}
