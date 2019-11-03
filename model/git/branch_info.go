package git

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

type CustomDateTime struct {
	time.Time
}

func (cdt *CustomDateTime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02T15:04:05-07:00", strInput)
	if err != nil {
		return err
	}

	cdt.Time = newTime
	return nil
}

type BranchInformation struct {
	BranchName     string         `json:"branch_name"`
	CommiterName   string         `json:"commiter_name"`
	LastCommitDate CustomDateTime `json:"last_commit_date"`
}

func NewBranchOwnerMap(gitLogs []string) map[string][]BranchInformation {
	branchInformations := unmarshalLogs(gitLogs)
	branchOwnerMap := tieOldBranchToAuthor(branchInformations)
	sortBranchesByDate(&branchOwnerMap)
	return branchOwnerMap
}

func unmarshalLogs(gitLogs []string) []BranchInformation {
	var branchInformations []BranchInformation
	for _, gitLog := range gitLogs {
		branchInformation := BranchInformation{}
		json.Unmarshal([]byte(gitLog), &branchInformation)
		branchInformations = append(branchInformations, branchInformation)
	}
	return branchInformations
}

func tieOldBranchToAuthor(branchInformations []BranchInformation) map[string][]BranchInformation {
	now := time.Now()
	branchOwnerMap := make(map[string][]BranchInformation)
	for _, branchInformation := range branchInformations {
		days := int(now.Sub(branchInformation.LastCommitDate.Time).Hours()) / 24
		if days >= 14 {
			branchOwnerMap[branchInformation.CommiterName] = append(branchOwnerMap[branchInformation.CommiterName], branchInformation)
		}
	}
	return branchOwnerMap
}

func sortBranchesByDate(branchOwnerMap *map[string][]BranchInformation) {
	for _, branchInformations := range *branchOwnerMap {
		sort.SliceStable(branchInformations, func(i, j int) bool {
			return branchInformations[i].LastCommitDate.Before(branchInformations[j].LastCommitDate.Time)
		})
	}
}
