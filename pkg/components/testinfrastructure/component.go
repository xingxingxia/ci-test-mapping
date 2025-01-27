package testinfrastructure

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var TestInfrastructureComponent = Component{
	Component: &config.Component{
		Name:                 "Test Infrastructure",
		Operators:            []string{},
		DefaultJiraComponent: "Test Infrastructure",
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"Clone the correct source code into an image and tag it as src",
					"core files found",
					"Infrastructure - AWS simulate policy rate-limit",
					"Infrastructure - GCP quota exceeded (route to forum-gcp)",
					"initialize",
					"Inject an RPM repository that will point at the RPM server",
					"[sig-arch] cloud API quota should not be exceeded",
					"[sig-arch] should not see excessive pull back-off on registry.redhat.io",
					"Start a service that hosts the RPMs generated by this build",
					"Store build results into a layer on top of bin and save as rpms",
					"Store build results into a layer on top of src and save as bin",
				},
				Priority: 2,
			},
		},
		TestRenames: map[string]string{
			"[sig-arch] pathological event should not see excessive pull back-off on registry.redhat.io": "[sig-arch] should not see excessive pull back-off on registry.redhat.io",
		},
	},
}

func (c *Component) IdentifyTest(test *v1.TestInfo) (*v1.TestOwnership, error) {
	if matcher := c.FindMatch(test); matcher != nil {
		jira := matcher.JiraComponent
		if jira == "" {
			jira = c.DefaultJiraComponent
		}
		return &v1.TestOwnership{
			Name:          test.Name,
			Component:     c.Name,
			JIRAComponent: jira,
			Priority:      matcher.Priority,
			Capabilities:  append(matcher.Capabilities, identifyCapabilities(test)...),
		}, nil
	}

	return nil, nil
}

func (c *Component) StableID(test *v1.TestInfo) string {
	// Look up the stable name for our test in our renamed tests map.
	if stableName, ok := c.TestRenames[test.Name]; ok {
		return stableName
	}
	return test.Name
}

func (c *Component) JiraComponents() (components []string) {
	components = []string{c.DefaultJiraComponent}
	for _, m := range c.Matchers {
		components = append(components, m.JiraComponent)
	}

	return components
}
