// Copyright 2023 IronCore authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package matchers_test

import (
	"fmt"

	. "github.com/ironcore-dev/controller-utils/testutils/matchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/semantic"
)

var _ = Describe("Matchers", func() {
	Context("EqualitiesEqualMatcher", func() {
		Describe("Match", func() {
			It("should match using the supplied equalities", func() {
				matcher := EqualitiesEqualMatcher{
					Equalities: semantic.EqualitiesOrDie(func(s1 string, s2 string) bool {
						if s1 == "*" || s2 == "*" {
							return true
						}
						return s1 == s2
					}),
					Expected: "foo",
				}

				Expect(matcher.Match("*")).To(BeTrue(), "* should match")
				Expect(matcher.Match("foo")).To(BeTrue(), "foo should match")
				Expect(matcher.Match("bar")).To(BeFalse(), "bar should not match")
				Expect(matcher.Match(fmt.Errorf("foo"))).To(BeFalse(), "an error should not match")
			})

			It("should error if the equalities are not set", func() {
				matcher := EqualitiesEqualMatcher{
					Expected: "foo",
				}
				_, err := matcher.Match("foo")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("EqualitiesDerivativeMatcher", func() {
		type Struct struct {
			A string
			B string
		}

		Describe("Match", func() {
			It("should match using the supplied equalities", func() {
				matcher := EqualitiesDerivativeMatcher{
					Equalities: semantic.EqualitiesOrDie(func(s1 string, s2 string) bool {
						if s1 == "" || s1 == "*" || s2 == "*" {
							return true
						}
						return s1 == s2
					}),
					Expected: Struct{
						A: "foo",
						B: "bar",
					},
				}

				Expect(matcher.Match(Struct{A: "*"})).To(BeTrue(), "A:* should match")
				Expect(matcher.Match(Struct{A: "foo"})).To(BeTrue(), "A:foo should match")
				Expect(matcher.Match(Struct{A: "foo", B: "bar"})).To(BeTrue(), "A:foo,B:bar should match")
				Expect(matcher.Match(Struct{A: "bar"})).To(BeFalse(), "A:bar should not match")
				Expect(matcher.Match(fmt.Errorf("foo"))).To(BeFalse(), "an error should not match")
			})

			It("should error if the equalities are not set", func() {
				matcher := EqualitiesDerivativeMatcher{
					Expected: "foo",
				}
				_, err := matcher.Match("foo")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("ErrorFuncMatcher", func() {
		Describe("Match", func() {
			It("should match using the given function", func() {
				matcher := ErrorFuncMatcher{
					Func: apierrors.IsNotFound,
				}

				Expect(matcher.Match(fmt.Errorf("custom"))).To(BeFalse(), "custom error should not match")
				Expect(matcher.Match(apierrors.NewNotFound(schema.GroupResource{}, ""))).To(BeTrue(), "not found should match")
				_, err := matcher.Match(1)
				Expect(err).To(HaveOccurred())
			})

			It("should error if the error function is not set", func() {
				matcher := ErrorFuncMatcher{}
				_, err := matcher.Match(fmt.Errorf("custom"))
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("ErrorMessage", func() {
			It("should report a correct error message", func() {
				matcher := ErrorFuncMatcher{
					Func: apierrors.IsNotFound,
				}

				Expect(matcher.FailureMessage(fmt.Errorf("custom"))).
					To(HavePrefix("expected an error matching k8s.io/apimachinery/pkg/api/errors.IsNotFound to have occurred but got"))
			})
		})

		Describe("NegatedErrorMessage", func() {
			It("should report a correct negated error message", func() {
				matcher := ErrorFuncMatcher{
					Func: apierrors.IsNotFound,
				}

				Expect(matcher.NegatedFailureMessage(fmt.Errorf("custom"))).
					To(HavePrefix("expected an error not matching k8s.io/apimachinery/pkg/api/errors.IsNotFound to have occurred but got"))
			})
		})
	})
})
