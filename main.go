// tcalc
// bitwise calculator that is run from the terminal.
// It supports basic variable assignments, and most arithmetic and bitwise operations.
package main

import (
	"flag"
	"os"
	"runtime/pprof"
	"slices"

	"fortio.org/cli"
	"fortio.org/log"
	"fortio.org/terminal/ansipixels"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	truecolorDefault := ansipixels.DetectColorMode().TrueColor
	fTrueColor := flag.Bool("truecolor", truecolorDefault,
		"Use true color (24-bit RGB) instead of 8-bit ANSI colors (default is true if COLORTERM is set)")
	fCpuprofile := flag.String("profile-cpu", "", "write cpu profile to `file`")
	fMemprofile := flag.String("profile-mem", "", "write memory profile to `file`")
	cli.Main()
	if *fCpuprofile != "" {
		f, err := os.Create(*fCpuprofile)
		if err != nil {
			return log.FErrf("can't open file for cpu profile: %v", err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			return log.FErrf("can't start cpu profile: %v", err)
		}
		log.Infof("Writing cpu profile to %s", *fCpuprofile)
		defer pprof.StopCPUProfile()
	}
	ap := ansipixels.NewAnsiPixels(0)
	c := configure(ap)

	ap.TrueColor = *fTrueColor
	if err := ap.Open(); err != nil {
		return 1 // error already logged
	}
	defer func() {
		c.AP.ShowCursor()
		c.AP.MouseClickOff()
		c.AP.Restore()
		c.AP.ClearScreen()
	}()
	c.AP.MouseClickOn()
	ap.SyncBackgroundColor()
	ap.OnResize = func() error {
		ap.StartSyncMode()
		c.Update()
		ap.EndSyncMode()
		return nil
	}
	_ = ap.OnResize() // initial draw.
	c.Update()
	for {
		errReading := c.AP.ReadOrResizeOrSignal()
		if errReading != nil {
			log.Errf("error getting read/resize/signal: %v", errReading)
			break
		}
		c.AP.StartSyncMode()
		if !c.Tick() {
			c.AP.EndSyncMode()
			break
		}
		c.Update()
		c.AP.EndSyncMode()
	}
	if *fMemprofile != "" {
		f, errMP := os.Create(*fMemprofile)
		if errMP != nil {
			return log.FErrf("can't open file for mem profile: %v", errMP)
		}
		errMP = pprof.WriteHeapProfile(f)
		if errMP != nil {
			return log.FErrf("can't write mem profile: %v", errMP)
		}
		log.Infof("Wrote memory profile to %s", *fMemprofile)
		_ = f.Close()
	}
	return 0
}

func (c *config) Update() {
	c.AP.ClearScreen()
	if c.AP.H > 19 {
		for i, str := range instructions {
			c.AP.WriteAtStr(0, i, str)
		}
	}
	strings := displayString(c.state.Ans, c.state.Prev, c.state.Err)
	y := c.AP.H - 13
	for i, str := range strings {
		c.AP.WriteAtStr(0, y+i, str)
	}
	for i := range 27 {
		c.AP.WriteAtStr(i, c.AP.H, ansipixels.Horizontal)
	}
	c.AP.WriteAtStr(0, c.AP.H-2, c.input)
	c.DrawHistory()
	c.AP.MoveCursor(c.index, c.AP.H-2)
}

func (c *config) Tick() bool {
	c.AP.MoveCursor(c.index+1, c.AP.H-2)
	if c.AP.LeftClick() && c.AP.MouseRelease() {
		x, y := c.AP.Mx, c.AP.My
		if slices.Contains(validClickXs, x) && y < c.AP.H-2 && y >= c.AP.H-6 {
			bit := c.determineBitFromXY(x, c.AP.H-2-y)
			c.clicked = true
			c.state.Ans = (c.state.Ans) ^ (1 << bit)
		}
	}
	diff := len(c.history) - (c.AP.H / 2) + 1
	if diff > 0 {
		c.history = c.history[diff:]
	}
	return c.handleInput()
}
