package enum

import (
	"reflect"
	"testing"
)

// Custom types for registry testing to avoid conflicts with int
type (
	RegistryTestType1 int
	RegistryTestType2 int
	RegistryTestType3 int
)

func TestRegister(t *testing.T) {
	// Test registering labels for a custom type
	labels := []string{"option1", "option2", "option3"}
	Register[RegistryTestType1](labels...)

	// Verify registration worked
	retrieved := GetLabels[RegistryTestType1]()
	if !reflect.DeepEqual(retrieved, labels) {
		t.Errorf("expected labels %v, got %v", labels, retrieved)
	}
}

func TestGetLabels(t *testing.T) {
	// Test getting labels for unregistered type
	labels := GetLabels[RegistryTestType2]()
	if labels != nil {
		t.Errorf("expected nil for unregistered type, got %v", labels)
	}

	// Register and test retrieval
	expectedLabels := []string{"red", "green", "blue"}
	Register[RegistryTestType2](expectedLabels...)

	labels = GetLabels[RegistryTestType2]()
	if !reflect.DeepEqual(labels, expectedLabels) {
		t.Errorf("expected labels %v, got %v", expectedLabels, labels)
	}
}

func TestRegisterOverwrite(t *testing.T) {
	// Test that registering overwrites previous registration
	firstLabels := []string{"first", "set"}
	secondLabels := []string{"second", "set", "with", "more"}

	Register[RegistryTestType3](firstLabels...)
	retrieved := GetLabels[RegistryTestType3]()
	if !reflect.DeepEqual(retrieved, firstLabels) {
		t.Errorf("expected first labels %v, got %v", firstLabels, retrieved)
	}

	Register[RegistryTestType3](secondLabels...)
	retrieved = GetLabels[RegistryTestType3]()
	if !reflect.DeepEqual(retrieved, secondLabels) {
		t.Errorf("expected second labels %v, got %v", secondLabels, retrieved)
	}
}

func TestNewWrapperWithRegistry(t *testing.T) {
	// Test that NewWrapper registers the type
	labels := []string{"alpha", "beta", "gamma"}
	wrapper := NewWrapper[RegistryTestType1](labels...)

	// Verify the labels were registered
	registered := GetLabels[RegistryTestType1]()
	if !reflect.DeepEqual(registered, labels) {
		t.Errorf("expected registered labels %v, got %v", labels, registered)
	}

	// Verify the wrapper was created correctly
	if !reflect.DeepEqual(wrapper.Enum.labels, labels) {
		t.Errorf("expected wrapper labels %v, got %v", labels, wrapper.Enum.labels)
	}
}

func TestEnsureEnumWithRegistry(t *testing.T) {
	// Register labels for a type
	labels := []string{"morning", "afternoon", "evening"}
	Register[RegistryTestType2](labels...)

	// Create a wrapper with nil Enum but no local labels (simulating deserialization)
	wrapper := Wrapper[RegistryTestType2]{
		Enum:    nil,
		Current: 1,
		labels:  nil, // No local labels, should use registry
	}

	// Test that ensureEnum() uses the registry
	wrapper.ensureEnum()

	if wrapper.Enum == nil {
		t.Fatal("ensureEnum() did not initialize the Enum")
	}

	if !reflect.DeepEqual(wrapper.Enum.labels, labels) {
		t.Errorf("expected enum labels from registry %v, got %v", labels, wrapper.Enum.labels)
	}

	if !reflect.DeepEqual(wrapper.labels, labels) {
		t.Errorf("expected wrapper labels to be set from registry %v, got %v", labels, wrapper.labels)
	}
}

func TestEnsureEnumLocalVsRegistry(t *testing.T) {
	// Register labels for a type
	registryLabels := []string{"global1", "global2", "global3"}
	Register[RegistryTestType3](registryLabels...)

	// Create a wrapper with local labels (should prefer local over registry)
	localLabels := []string{"local1", "local2"}
	wrapper := Wrapper[RegistryTestType3]{
		Enum:    nil,
		Current: 0,
		labels:  localLabels,
	}

	// Test that ensureEnum() prefers local labels over registry
	wrapper.ensureEnum()

	if wrapper.Enum == nil {
		t.Fatal("ensureEnum() did not initialize the Enum")
	}

	if !reflect.DeepEqual(wrapper.Enum.labels, localLabels) {
		t.Errorf("expected enum labels from local %v, got %v", localLabels, wrapper.Enum.labels)
	}
}
