package data

import (
	"fmt"
	"io"
	"strconv"
)

// Title is a language-specific string used as a name
// or descriptor in other models.
type Title struct {
	String   string
	Language string
	Priority TitlePriority
}

// TitlePriority is an enum that describes the priority of
// a Title within a set of Titles.
type TitlePriority int

const (
	// TitlePriorityPrimary means the Title is a primary
	// one in a set.
	TitlePriorityPrimary       = 0
	titlePriorityPrimaryString = "Primary"
	// TitlePrioritySecondary means the Title is a
	// secondary one in a set.
	TitlePrioritySecondary       = 1
	titlePrioritySecondaryString = "Secondary"
	// TitlePriorityOther means the Title is a tertiary
	// or other one in a set.
	TitlePriorityOther       = 2
	titlePriorityOtherString = "Other"
)

// IsValid checks if the TitlePriority has a value that
// is a valid one.
func (p TitlePriority) IsValid() bool {
	switch p {
	case TitlePriorityPrimary, TitlePrioritySecondary, TitlePriorityOther:
		return true
	}
	return false
}

// String returns the written name of the TitlePriority.
func (p TitlePriority) String() string {
	switch p {
	case TitlePriorityPrimary:
		return titlePriorityPrimaryString
	case TitlePrioritySecondary:
		return titlePrioritySecondaryString
	case TitlePriorityOther:
		return titlePriorityOtherString
	}
	return fmt.Sprintf("%d", int(p))
}

// UnmarshalGQL casts the type of the given value to
// a TitlePriority.
func (p *TitlePriority) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("%v: %w", v, errInvalid)
	}

	switch str {
	case titlePriorityPrimaryString:
		*p = TitlePriorityPrimary
	case titlePrioritySecondaryString:
		*p = TitlePrioritySecondary
	case titlePriorityOtherString:
		*p = TitlePriorityOther
	default:
		return fmt.Errorf("%s: %w", str, errInvalid)
	}
	return nil
}

// MarshalGQL serializes the Priority into a GraphQL
// readable form.
func (p TitlePriority) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(p.String()))
}

// TitleSetMatchLanguage returns all the Titles with the given
// language in the set.
func TitleSetMatchLanguage(set []Title, lang string) []Title {
	var ss []Title
	for _, t := range set {
		if t.Language == lang {
			ss = append(ss, t)
		}
	}

	if len(ss) == 0 {
		return []Title{}
	}

	return ss
}

// TitleSetMatchPrimary returns all the Titles with primary priority
// in the set.
func TitleSetMatchPrimary(set []Title) []Title {
	return titleSetPriorityFilter(set, TitlePriorityPrimary)
}

// TitleSetMatchSecondary returns all the Titles with secondary
// priority in the set.
func TitleSetMatchSecondary(set []Title) []Title {
	return titleSetPriorityFilter(set, TitlePrioritySecondary)
}

// TitleSetMatchOther returns all the Titles with other priority
// in the set.
func TitleSetMatchOther(set []Title) []Title {
	return titleSetPriorityFilter(set, TitlePriorityOther)
}

func titleSetPriorityFilter(set []Title, p TitlePriority) []Title {
	var ss []Title
	for _, t := range set {
		if t.Priority == p {
			ss = append(ss, t)
		}
	}

	if len(ss) == 0 {
		return []Title{}
	}

	return ss
}
