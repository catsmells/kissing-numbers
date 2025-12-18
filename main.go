package main

import (
	"fmt"
	"math"
	"strconv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ip(v int) *int { return &v }

type KissingData struct {
	Exact    *int
	Lower    *int
	Upper    *int
	Root     string
	Diagram  string
}

var kissingTable = map[int]KissingData{
	1: {Exact: ip(2), Root: "A1", Diagram: "o\n"},
	2: {Exact: ip(6), Root: "A2", Diagram: "o---o\n"},
	3: {Exact: ip(12), Root: "A3", Diagram: "o---o---o\n"},
	4: {
		Exact: ip(24), Root: "D4",
		Diagram: `
    o
    |
o---o---o
`,
	},
	5: {Lower: ip(40), Upper: ip(44)},
	6: {Lower: ip(72), Upper: ip(78)},
	7: {Lower: ip(126), Upper: ip(134)},
	8: {
		Exact: ip(240), Root: "E8",
		Diagram: `
o---o---o---o---o---o---o
                |
                o
`,
	},
	24: {
		Exact: ip(196560),
		Root:  "Leech lattice",
		Diagram: "(No Coxeter diagram exists)\n",
	},
}

func asymptoticBounds(n int) (float64, float64) {
	// Kabatiansky–Levenshtein bounds:
	// 2^{0.2075n} ≤ K(n) ≤ 2^{0.401n}
	lower := math.Pow(2, 0.2075*float64(n))
	upper := math.Pow(2, 0.401*float64(n))
	return lower, upper
}

func renderExact(n int, out *tview.TextView) {
	data := kissingTable[n]
	if data.Exact != nil {
		out.SetText(fmt.Sprintf(
			"[green]Dimension:[white] %d\n"+
				"[green]Exact Kissing Number:[white] %d\n"+
				"[green]Root System:[white] %s\n\n"+
				"[green]Coxeter Diagram:\n[white]%s",
			n, *data.Exact, data.Root, data.Diagram))
	} else if data.Lower != nil {
		out.SetText(fmt.Sprintf(
			"[yellow]Dimension %d (Unproven)\n\n"+
				"[green]Lower Bound:[white] %d\n"+
				"[red]Upper Bound:[white] %d\n\n"+
				"[yellow]Exact kissing number not known.",
			n, *data.Lower, *data.Upper))
	} else {
		out.SetText(fmt.Sprintf(
			"[yellow]Dimension %d\nNo finite bounds stored.",
			n))
	}
}

func renderAsymptotic(n int, out *tview.TextView) {
	l, u := asymptoticBounds(n)
	out.SetText(fmt.Sprintf(
		"[cyan]Asymptotic Bounds (Kabatiansky–Levenshtein)\n\n"+
			"[green]Lower Bound:[white] %.3e\n"+
			"[red]Upper Bound:[white] %.3e\n\n"+
			"[yellow]These bounds hold for sufficiently large dimensions.",
		l, u))
}

func main() {
	app := tview.NewApplication()

	input := tview.NewInputField().
		SetLabel("Dimension n: ").
		SetFieldWidth(10)

	output := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetBorder(true).
		SetTitle("Kissing Numbers")

	mode := "exact"
	var currentN int

	update := func() {
		if mode == "exact" {
			renderExact(currentN, output)
		} else {
			renderAsymptotic(currentN, output)
		}
	}

	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}
		n, err := strconv.Atoi(input.GetText())
		if err != nil || n < 1 {
			output.SetText("[red]Invalid dimension")
			return
		}
		currentN = n
		update()
	})

	app.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Rune() {
		case 'q', 'Q':
			app.Stop()
		case 'e', 'E':
			mode = "exact"
			update()
		case 'a', 'A':
			mode = "asymptotic"
			update()
		}
		return ev
	})

	help := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[E] Exact / Bounds   [A] Asymptotic   [Q] Quit")

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(input, 3, 1, true).
		AddItem(output, 0, 6, false).
		AddItem(help, 1, 1, false)

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
