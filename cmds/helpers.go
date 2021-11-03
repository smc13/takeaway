package cmds

import "github.com/AlecAivazis/survey/v2"

var SelectIcons = survey.WithIcons(func(is *survey.IconSet) {
	is.UnmarkedOption.Text = "○"
	is.MarkedOption.Text = "●"
	is.MarkedOption.Format = "cyan+b"
})
