package run

import (
	"reflect"
	"testing"

	oct "github.com/kyma-incubator/octopus/pkg/apis/testing/v1alpha1"
	"github.com/kyma-project/cli/pkg/kyma/cmd/test"
	"github.com/kyma-project/cli/pkg/kyma/cmd/test/client"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_matchTestDefinitionNames(t *testing.T) {
	testData := []struct {
		testName        string
		shouldFail      bool
		testNames       []string
		testDefinitions []oct.TestDefinition
		result          []oct.TestDefinition
	}{
		{
			testName:   "match all tests",
			shouldFail: false,
			testNames:  []string{"test1", "test2"},
			testDefinitions: []oct.TestDefinition{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
					TypeMeta: metav1.TypeMeta{
						APIVersion: "",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test2",
					},
					TypeMeta: metav1.TypeMeta{
						APIVersion: "",
					},
				},
			},
			result: []oct.TestDefinition{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
					TypeMeta: metav1.TypeMeta{
						APIVersion: "",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test2",
					},
					TypeMeta: metav1.TypeMeta{
						APIVersion: "",
					},
				},
			},
		},
	}

	for _, tt := range testData {
		result, err := matchTestDefinitionNames(tt.testNames, tt.testDefinitions)
		if tt.shouldFail {
			require.NotNil(t, err, tt.testName)
		} else {
			require.Nil(t, err, tt.testName)
			require.True(t, reflect.DeepEqual(result, tt.result), tt.testName)
		}
	}
}

func Test_generateTestsResource(t *testing.T) {
	testData := []struct {
		testName             string
		shouldFail           bool
		inputTestName        string
		inputTestDefinitions []oct.TestDefinition
		expectedResult       *oct.ClusterTestSuite
	}{
		{
			testName:      "create test with existing test definition",
			shouldFail:    false,
			inputTestName: "TestOneProper",
			inputTestDefinitions: []oct.TestDefinition{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test1",
						Namespace: "kyma-test",
					},
					TypeMeta: metav1.TypeMeta{
						APIVersion: "",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test2",
						Namespace: "kyma-system",
					},
					TypeMeta: metav1.TypeMeta{
						APIVersion: "",
					},
				},
			},
			expectedResult: &oct.ClusterTestSuite{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "testing.kyma-project.io/v1alpha1",
					Kind:       "ClusterTestSuite",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestOneProper",
					Namespace: test.NamespaceForTests,
				},
				Spec: oct.TestSuiteSpec{
					MaxRetries:  1,
					Concurrency: 1,
					Selectors: oct.TestsSelector{
						MatchNames: []oct.TestDefReference{
							{
								Name:      "test1",
								Namespace: "kyma-test",
							},
							{
								Name:      "test2",
								Namespace: "kyma-system",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range testData {
		result := generateTestsResource(
			tt.inputTestName,
			tt.inputTestDefinitions,
		)
		if tt.shouldFail {
			require.False(
				t,
				reflect.DeepEqual(result, tt.expectedResult),
				tt.testName)
		} else {
			require.True(
				t,
				reflect.DeepEqual(result, tt.expectedResult),
				tt.testName)
		}
	}
}

func Test_verifyIfTestNotExists(t *testing.T) {
	testData := []struct {
		testName       string
		inputSuiteName string
		inputSuites    []oct.ClusterTestSuite
		expectedExists bool
	}{
		{
			testName:       "verify existing test",
			inputSuiteName: "test1",
			inputSuites: []oct.ClusterTestSuite{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test2",
					},
				},
			},
			expectedExists: false,
		},
		{
			testName:       "verify non-existing test",
			inputSuiteName: "test1",
			inputSuites: []oct.ClusterTestSuite{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test3",
					},
				},
			},
			expectedExists: true,
		},
	}

	for _, tt := range testData {
		mCli := client.NewMockedTestRestClient(nil, &oct.ClusterTestSuiteList{
			Items: tt.inputSuites,
		})
		tExists, _ := verifyIfTestNotExists(tt.inputSuiteName, mCli)
		require.Equal(t, tExists, tt.expectedExists)
	}
}