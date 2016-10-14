package thread

import (
	"testing"
)

func TestParseRepoName(t *testing.T) {
	repoOwner := "Samangan" // TODO: Generate this so that the tests dont fail on a fork
	repoName := "wtf"

	actualOwner, actualName := parseRepoName()

	if actualName != repoName {
		t.Errorf("Want %v to equal %v", actualName, repoName)
	}

	if actualOwner != repoOwner {
		t.Errorf("Want %v to equal %v", actualOwner, repoOwner)
	}
}
