package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddJobLevelPermissions(t *testing.T) {
	const inputDirectory = "./testfiles/joblevelpermskb/input"
	const outputDirectory = "./testfiles/joblevelpermskb/output"
	files, err := ioutil.ReadDir(inputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		input, err := ioutil.ReadFile(path.Join(inputDirectory, f.Name()))

		if err != nil {
			log.Fatal(err)
		}

		fixWorkflowPermsResponse, err := AddJobLevelPermissions(string(input))
		output := fixWorkflowPermsResponse.FinalOutput
		jobErrors := fixWorkflowPermsResponse.JobErrors

		// some test cases return a job error for known issues
		if len(jobErrors) > 0 {
			for _, je := range jobErrors {
				if strings.Compare(je.JobName, "job-with-error") == 0 {
					if strings.Contains(je.Errors[0], "KnownIssue") {
						output = je.Errors[0]
					} else {
						t.Errorf("test failed. unexpected job error %s, error: %v", f.Name(), jobErrors)
					}
				}
			}

		}

		if fixWorkflowPermsResponse.AlreadyHasPermissions {
			output = errorAlreadyHasPermissions
		}

		if fixWorkflowPermsResponse.IncorrectYaml {
			output = errorIncorrectYaml
		}

		if err != nil {
			t.Errorf("test failed %s, error: %v", f.Name(), err)
		}

		expectedOutput, err := ioutil.ReadFile(path.Join(outputDirectory, f.Name()))

		if err != nil {
			log.Fatal(err)
		}

		if output != string(expectedOutput) {
			t.Errorf("test failed %s did not match expected output\n%s", f.Name(), output)
		}
	}
}

func TestStarterWorflowPermissions(t *testing.T) {
	const inputDirectory = "/Users/varunsharma/go/src/github.com/varunsh-coder/starter-workflows/"
	const solvableWorkflowsFile = "/Users/varunsharma/go/src/github.com/varunsh-coder/starter-workflows/solvableWorkflows.txt"
	const nonSolvableWorkflowsFile = "/Users/varunsharma/go/src/github.com/varunsh-coder/starter-workflows/nonSolvableWorkflows.txt"
	const missingActionsFile = "/Users/varunsharma/go/src/github.com/varunsh-coder/starter-workflows/actions.txt"
	missingActions := make(map[string]bool)
	solvableWorkflows := make(map[string]bool)
	nonSolvableWorkflows := make(map[string]bool)
	err := filepath.Walk(inputDirectory, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		ext := filepath.Ext(f.Name())
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		input, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		fixWorkflowPermsResponse, err := AddJobLevelPermissions(string(input))

		if err != nil {
			return err
		}

		if len(fixWorkflowPermsResponse.MissingActions) > 0 {
			nonSolvableWorkflows[path] = true
			for _, v := range fixWorkflowPermsResponse.MissingActions {
				actionkey := strings.Split(v, "@")
				if len(actionkey) > 1 {
					missingActions[actionkey[0]] = true
				}
			}
		}

		if fixWorkflowPermsResponse.HasErrors == false && fixWorkflowPermsResponse.IsChanged {
			solvableWorkflows[path] = true
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	actions := ""
	for k := range missingActions {
		actions += k + "\n"
	}
	ioutil.WriteFile(missingActionsFile, []byte(actions), 0644)

	nonSolvableWorkflowsList := ""
	for k := range nonSolvableWorkflows {
		nonSolvableWorkflowsList += k + "\n"
	}
	ioutil.WriteFile(nonSolvableWorkflowsFile, []byte(nonSolvableWorkflowsList), 0644)

	solvableWorkflowsList := ""
	for k := range solvableWorkflows {
		solvableWorkflowsList += k + "\n"
	}
	ioutil.WriteFile(solvableWorkflowsFile, []byte(solvableWorkflowsList), 0644)

}

func Test_addPermissions(t *testing.T) {
	type args struct {
		inputYaml   string
		jobName     string
		permissions []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "bad yaml",
			args: args{
				inputYaml: "123",
			}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addPermissions(tt.args.inputYaml, tt.args.jobName, tt.args.permissions)
			if (err != nil) != tt.wantErr {
				t.Errorf("addPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("addPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddWorkflowLevelPermissions(t *testing.T) {
	const inputDirectory = "./testfiles/toplevelperms/input"
	const outputDirectory = "./testfiles/toplevelperms/output"
	files, err := ioutil.ReadDir(inputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		input, err := ioutil.ReadFile(path.Join(inputDirectory, f.Name()))

		if err != nil {
			log.Fatal(err)
		}

		output, err := AddWorkflowLevelPermissions(string(input))

		if err != nil {
			t.Errorf("Error not expected")
		}

		expectedOutput, err := ioutil.ReadFile(path.Join(outputDirectory, f.Name()))

		if err != nil {
			log.Fatal(err)
		}

		if output != string(expectedOutput) {
			t.Errorf("test failed %s did not match expected output\n%s", f.Name(), output)
		}
	}

}
