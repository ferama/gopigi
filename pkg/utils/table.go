package utils

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func GetTableWriter() table.Writer {
	t := table.NewWriter()

	colors := table.ColorOptions{
		IndexColumn:  text.Colors{text.FgHiWhite, text.BgBlack},
		Footer:       text.Colors{text.FgBlack, text.BgWhite},
		Header:       text.Colors{text.BgHiBlue, text.FgBlack},
		Row:          text.Colors{text.FgHiWhite, text.BgBlack},
		RowAlternate: text.Colors{text.FgWhite, text.BgBlack},
	}
	options := table.Options{
		DrawBorder:      false,
		SeparateColumns: true,
		SeparateFooter:  false,
		SeparateHeader:  false,
		SeparateRows:    false,
	}
	styleColor := table.Style{
		Name:    "Custom",
		Box:     table.StyleBoxLight,
		Color:   colors,
		Format:  table.FormatOptionsDefault,
		HTML:    table.DefaultHTMLOptions,
		Options: options,
		Title:   table.TitleOptionsBlackOnBlue,
	}

	t.SetStyle(styleColor)
	return t
}
