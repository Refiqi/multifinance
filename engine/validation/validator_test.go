package validation

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidator(t *testing.T) {
	Convey("Test Validator", t, func() {
		validator := NewValidation()
		Convey("Test validator Validate", func() {
			testStruct := struct {
				TestField string `validate:"required"`
			}{}

			Convey("error", func() {
				errMap := validator.Validate(testStruct)
				So(errMap, ShouldHaveLength, 1)
				So(errMap, ShouldContainKey, "TestField")
			})

			Convey("success", func() {
				testStruct.TestField = "testing"
				errMap := validator.Validate(testStruct)
				So(errMap, ShouldBeNil)
			})
		})
	})
}
