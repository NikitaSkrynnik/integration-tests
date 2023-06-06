package parallel

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }

// var matchMethod = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")
var matchMethod = ""

func recoverAndFailOnPanic(t *testing.T) {
	r := recover()
	failOnPanic(t, r)
}

func failOnPanic(t *testing.T, r interface{}) {
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

func Run(t *testing.T, s suite.TestingSuite) {
	defer recoverAndFailOnPanic(t)
	var suiteSetupDone bool

	s.SetT(t)

	var stats *SuiteInformation
	if _, ok := s.(suite.WithStats); ok {
		stats = newSuiteInformation()
	}

	tests := []testing.InternalTest{}
	methodFinder := reflect.TypeOf(s)
	suiteName := methodFinder.Elem().Name()

	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)

		ok, err := methodFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}

		if !ok {
			continue
		}

		if !suiteSetupDone {
			if stats != nil {
				stats.Start = time.Now()
			}

			if setupAllSuite, ok := s.(suite.SetupAllSuite); ok {
				setupAllSuite.SetupSuite()
			}

			suiteSetupDone = true
		}

		test := testing.InternalTest{
			Name: method.Name,
			F: func(testingT *testing.T) {
				defer recoverAndFailOnPanic(t)

				testingT.Parallel()
				defer func() {
					r := recover()

					if stats != nil {
						passed := !t.Failed()
						stats.end(method.Name, passed)
					}

					if afterTestSuite, ok := s.(suite.AfterTest); ok {
						afterTestSuite.AfterTest(suiteName, method.Name)
					}

					if tearDownTestSuite, ok := s.(suite.TearDownTestSuite); ok {
						tearDownTestSuite.TearDownTest()
					}

					failOnPanic(t, r)
				}()

				if setupTestSuite, ok := s.(suite.SetupTestSuite); ok {
					setupTestSuite.SetupTest()
				}
				if beforeTestSuite, ok := s.(suite.BeforeTest); ok {
					beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
				}

				if stats != nil {
					stats.start(method.Name)
				}

				subS := reflect.New(reflect.ValueOf(s).Elem().Type())
				subS.MethodByName("SetT").Call([]reflect.Value{reflect.ValueOf(testingT)})

				method.Func.Call([]reflect.Value{subS})
			},
		}
		tests = append(tests, test)
	}

	if suiteSetupDone {
		if tearDownAllSuite, ok := s.(suite.TearDownAllSuite); ok {
			tearDownAllSuite.TearDownSuite()
		}

		if suiteWithStats, measureStats := s.(WithStats); measureStats {
			stats.End = time.Now()
			suiteWithStats.HandleStats(suiteName, stats)
		}
	}

	if len(tests) == 0 {
		t.Log("warning: no tests to run")
		return
	}

	// run sub-tests in a group so tearDownSuite is called in the right order
	for _, test := range tests {
		t.Run(test.Name, test.F)
	}

}

// Filtering method according to set regular expression
// specified command-line argument -m
func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Test", name); !ok {
		return false, nil
	}
	return regexp.MatchString(matchMethod, name)
}
