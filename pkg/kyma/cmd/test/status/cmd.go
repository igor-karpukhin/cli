package status

import (
	"encoding/json"
	"fmt"
	"time"

	oct "github.com/kyma-incubator/octopus/pkg/apis/testing/v1alpha1"
	"github.com/kyma-project/cli/pkg/kyma/cmd/test"
	"github.com/kyma-project/cli/pkg/kyma/cmd/test/client"
	"github.com/kyma-project/cli/pkg/kyma/core"
	"github.com/spf13/cobra"
)

type command struct {
	opts *options
	core.Command
}

func NewCmd(o *options) *cobra.Command {
	cmd := command{
		Command: core.Command{Options: o.Options},
		opts:    o,
	}

	cobraCmd := &cobra.Command{
		Use:     "status",
		Short:   "Status of tests on a running Kyma cluster",
		Long:    `Status of tests on a running Kyma cluster`,
		RunE:    func(_ *cobra.Command, args []string) error { return cmd.Run(args) },
		Aliases: []string{"s"},
	}

	//cobraCmd.Flags().BoolVarP(&o.Wait, "wait", "w", false, "Wait for test execution to finish")
	//cobraCmd.Flags().StringVarP(&o.Logs, "logs", "l", "", "Test name to show logs")
	cobraCmd.Flags().BoolVarP(&o.Jsn, "raw", "r", false,
		"Print test status in raw json format")
	return cobraCmd
}

func (cmd *command) Run(args []string) error {
	cli, err := client.NewTestRESTClient(10 * time.Second)
	if err != nil {
		return fmt.Errorf("unable to create test REST client. E: %s", err)
	}

	switch len(args) {
	case 1:
		testSuite, err := cli.GetTestSuiteByName(args[0])
		if err != nil {
			return fmt.Errorf("unable to get test suite '%s'. E: %s",
				args[0], err.Error())
		}
		return cmd.printTestSuiteStatus(testSuite, cmd.opts.Jsn)
	case 0:
		testList, err := cli.ListTestSuites()
		if err != nil {
			return fmt.Errorf("unable to list test suites. E: %s", err.Error())
		}

		for _, t := range testList.Items {
			if err := cmd.printTestSuiteStatus(&t, cmd.opts.Jsn); err != nil {
				return err
			}
		}
	default:
		testsList, err := test.ListTestSuitesByName(cli, args)
		if err != nil {
			return fmt.Errorf("unable to list test suites. E: %s", err.Error())
		}

		for _, t := range testsList {
			if err := cmd.printTestSuiteStatus(&t, cmd.opts.Jsn); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cmd *command) printTestSuiteStatus(test *oct.ClusterTestSuite, raw bool) error {
	if test == nil {
		return fmt.Errorf("unable to print test suite. Nil pointer\r\n")
	}
	if raw {
		d, err := json.MarshalIndent(test, "", "\t")
		if err != nil {
			return fmt.Errorf("unable to marshal test suite '%s'. E: %s\r\n",
				test.GetName(), err.Error())
		}
		fmt.Println(string(d))
		return nil
	}
	fmt.Printf("Name:\t\t%s\r\n", test.GetName())
	fmt.Printf("Concurrency:\t%d\r\n", test.Spec.Concurrency)
	fmt.Printf("MaxRetries:\t%d\r\n", test.Spec.MaxRetries)
	if test.Status.StartTime != nil {
		fmt.Printf("StartTime:\t%s\r\n", test.Status.StartTime.String())
	} else {
		fmt.Printf("StartTime:\t%s\r\n", "not started yet")
	}
	if test.Status.CompletionTime != nil {
		fmt.Printf("EndTime:\t%s\r\n", test.Status.CompletionTime)
	} else {
		fmt.Printf("EndTime:\t%s\r\n", "not finished yet")
	}

	fmt.Printf("Condition:\t%s\r\n", test.Status.Conditions[len(test.Status.Conditions)-1].Type)
	fmt.Printf("Tests:\r\n")
	for _, t := range test.Status.Results {
		fmt.Printf("\t%s - %s\r\n", t.Name, t.Status)
	}
	return nil
}
