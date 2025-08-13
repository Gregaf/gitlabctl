package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func (ch *CommandHandler) CommandVerify(args []string) {
	flagSetVerify := flag.NewFlagSet("verify", flag.ExitOnError)
	flagProjects := flagSetVerify.String("projects", "", "Gitlab Projects to verify, comma separated")
	flagBranchName := flagSetVerify.String("branch", "", "Gitlab Branch name")
	flagSetVerify.Parse(args[1:])

	ch.logger.Debug("Verify command executed", "args", flagSetVerify.Args())

	projects := []string{}
	if *flagProjects != "" {
		for _, p := range strings.Split(*flagProjects, ",") {
			projects = append(projects, strings.TrimSpace(p))
		}
	}

	gitlabProjects, _ := ch.gitlabService.FindGitlabProjects(projects)

	missingBranchProjects, err := ch.gitlabService.BulkHasBranch(gitlabProjects, *flagBranchName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	ch.logger.Debug("Found projects missing branch", "targetBranch", *flagBranchName, "missing", missingBranchProjects, "projects", gitlabProjects)

	projectNames := []string{}
	for _, e := range missingBranchProjects {
		projectNames = append(projectNames, e.NamespacePath)
	}

	out := strings.Join(projectNames, " ")

	if len(projectNames) <= 0 {
		out = "All projects pass validation"
	}

	fmt.Println(out)
}
