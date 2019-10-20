package data

import "strings"

// Info represents some info for some other object
// that depends on language.
type Info struct {
	Data     string
	Language string
}

// infoClean cleens the given Info for storage.
func infoClean(e *Info) (err error) {
	e.Data = strings.Trim(e.Data, " ")
	e.Language = strings.Trim(e.Data, " ")
	return nil
}

// infoListClean cleans the given list of Info
// for storage.
func infoListClean(list []Info) (err error) {
	for _, i := range list {
		if err := infoClean(&i); err != nil {
			return err
		}
	}
	return nil
}
