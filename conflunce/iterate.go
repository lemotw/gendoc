package conflunce

import "github.com/lemotw/gendoc/model"

func WalkAllChildPage(url, pageid string, fn func(c *model.ChildPage) error, auth *model.ConflunceAccount) error {
	pageStack := []string{pageid}

	for len(pageStack) > 0 {
		pid := pageStack[len(pageStack)-1]
		pageStack = pageStack[:len(pageStack)-1]

		childPage, err := FetchConflunceChildPage(url, pid, auth)
		if err != nil {
			return err
		}

		for i := 0; i < len(childPage.Results); i++ {
			if err := fn(&childPage.Results[i]); err != nil {
				return err
			}
			pageStack = append(pageStack, childPage.Results[i].ID)
		}
	}

	return nil
}
