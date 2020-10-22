package concurrency_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/amahdavian/async-event-bus/v2/internal/concurrency"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

var assertThat = then.AssertThat
var equals = is.EqualTo
var contains = is.ValueContaining

func TestShouldAddValueForGivenKey(t *testing.T) {
	mapOfSlice := concurrency.NewMapOfSlice()
	mapOfSlice.AppendAt("SomeKey", "SomeItem")

	actual, found := mapOfSlice.Get("SomeKey")

	assertThat(t, found, equals(true))
	assertThat(t, actual, equals([]interface{}{"SomeItem"}))
}

func TestShouldAllowConcurrentAddsForSameKey(t *testing.T) {
	mapOfSlice := concurrency.NewMapOfSlice()

	completionChannel := make(chan bool)
	var expectedValues []interface{}
	for i := 0; i < 100; i++ {
		value := fmt.Sprintf("Item-%d", i)
		expectedValues = append(expectedValues, value)
		go func() {
			mapOfSlice.AppendAt("SomeKey", value)
			completionChannel <- true
		}()
	}

	for i := 0; i < 100; i++ {
		select {
		case <-completionChannel:
		case <-time.After(1 * time.Second):
			t.Errorf("did not complete adding to the map within 1 second")
		}
	}

	actual, found := mapOfSlice.Get("SomeKey")

	assertThat(t, found, equals(true))
	assertThat(t, len(actual), equals(100))
	assertThat(t, actual, contains(expectedValues))
}

func TestShouldAllowConcurrentAddsForDifferentKeys(t *testing.T) {
	mapOfSlice := concurrency.NewMapOfSlice()

	completionChannel := make(chan bool)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		go func() {
			mapOfSlice.AppendAt(key, value)
			completionChannel <- true
		}()
	}

	for i := 0; i < 100; i++ {
		select {
		case <-completionChannel:
		case <-time.After(1 * time.Second):
			t.Errorf("did not complete adding to the map within 1 second")
		}
	}

	for i := 0; i < 100; i++ {
		actual, found := mapOfSlice.Get(fmt.Sprintf("key-%d", i))

		assertThat(t, found, equals(true))
		assertThat(t, len(actual), equals(1))
		assertThat(t, actual, contains([]interface{}{fmt.Sprintf("value-%d", i)}))
	}
}

func TestShouldRemoveValueAndDeleteKeyWhenOnlyOneValueIsAssociatedWithGivenKey(t *testing.T) {
	mapOfSlice := concurrency.NewMapOfSlice()
	mapOfSlice.AppendAt("SomeKey", "SomeItem")

	actual, found := mapOfSlice.Get("SomeKey")

	assertThat(t, found, equals(true))
	assertThat(t, actual, equals([]interface{}{"SomeItem"}))

	mapOfSlice.RemoveAt("SomeKey", "SomeItem")

	actual, found = mapOfSlice.Get("SomeKey")

	assertThat(t, found, is.False())
	assertThat(t, actual == nil, is.True())
}

func TestShouldOnlyRemoveGivenValueWhenMultipleValuesAreAssociatedWithGivenKey(t *testing.T) {
	mapOfSlice := concurrency.NewMapOfSlice()
	mapOfSlice.AppendAt("SomeKey", "SomeItem")
	mapOfSlice.AppendAt("SomeKey", "AnotherItem")

	actual, found := mapOfSlice.Get("SomeKey")

	assertThat(t, found, equals(true))
	assertThat(t, actual, equals([]interface{}{"SomeItem", "AnotherItem"}))

	mapOfSlice.RemoveAt("SomeKey", "SomeItem")

	actual, found = mapOfSlice.Get("SomeKey")

	assertThat(t, found, is.True())
	assertThat(t, actual, equals([]interface{}{"AnotherItem"}))
}

func TestShouldDoNothingWhenRemovingNonExistingValueForGivenKey(t *testing.T) {
	mapOfSlice := concurrency.NewMapOfSlice()
	mapOfSlice.AppendAt("SomeKey", "SomeItem")

	actual, found := mapOfSlice.Get("SomeKey")

	assertThat(t, found, equals(true))
	assertThat(t, actual, equals([]interface{}{"SomeItem"}))

	mapOfSlice.RemoveAt("SomeKey", "AnotherItem")

	actual, found = mapOfSlice.Get("SomeKey")

	assertThat(t, found, is.True())
	assertThat(t, actual, equals([]interface{}{"SomeItem"}))
}

func TestShouldAllowConcurrentRemovalsOfValuesForGivenKey(t *testing.T) {
	mapOfSlice := concurrency.NewMapOfSlice()

	var addedValues []interface{}
	for i := 0; i < 100; i++ {
		value := fmt.Sprintf("Item-%d", i)
		addedValues = append(addedValues, value)
		mapOfSlice.AppendAt("SomeKey", value)
	}

	completionChannel := make(chan bool)
	for _, value := range addedValues {
		valueToRemove := value
		go func() {
			mapOfSlice.RemoveAt("SomeKey", valueToRemove)
			completionChannel <- true
		}()
	}

	for i := 0; i < 100; i++ {
		select {
		case <-completionChannel:
		case <-time.After(1 * time.Second):
			t.Errorf("did not complete adding to the map within 1 second")
		}
	}

	actual, found := mapOfSlice.Get("SomeKey")

	assertThat(t, found, is.False())
	assertThat(t, actual == nil, is.True())
}
