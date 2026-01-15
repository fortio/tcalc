// tcalc
// bitwise calculator that is run from the terminal.
// It supports basic variable assignments, and most arithmetic and bitwise operations.
package main

import (
	"flag"
	"os"
	"runtime/pprof"
	"strings"

	"fortio.org/cli"
	"fortio.org/log"
	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	truecolorDefault := ansipixels.DetectColorMode().TrueColor
	fTrueColor := flag.Bool("truecolor", truecolorDefault,
		"Use true color (24-bit RGB) instead of 8-bit ANSI colors (default is true if COLORTERM is set)")
	fCpuprofile := flag.String("profile-cpu", "", "write cpu profile to `file`")
	fpsFlag := flag.Float64("fps", 60, "set fps for display refresh")
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
	ap := ansipixels.NewAnsiPixels(*fpsFlag)
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
	err := ap.FPSTicks(func() bool {
		if !c.Tick() {
			return false
		}
		c.Update()
		return true
	})
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
	if err != nil {
		log.Infof("Exiting on %v", err)
		return 1
	}
	return 0
}

func (c *config) Update() {
	diff := len(c.history) - (c.AP.H / 2) + 1
	if diff > 0 {
		c.history = c.history[diff:]
	}
	c.AP.ClearScreen()
	if c.AP.H < 14 {
		c.AP.WriteAtStr(0, 0, "Terminal too small")
		c.AP.ShowCursor()
		return
	}
	if c.AP.H > 20 {
		for i, str := range instructions {
			c.AP.WriteAtStr(0, i, str)
		}
	}
	c.strings = displayString(c.state.Ans, c.state.Prev, c.state.Err)
	y := c.AP.H - 13
	for i, str := range c.strings {
		c.AP.WriteAtStr(0, y+i, str)
	}
	for i := range 27 {
		c.AP.WriteAtStr(i, c.AP.H, ansipixels.Horizontal)
	}
	c.AP.WriteAtStr(0, c.AP.H-2, strings.Replace(c.input, "_ans_", italicPrefix+GREEN+"_ans_"+tcolor.Reset, -1))
	c.DrawHistory()
	if c.strings[0] == "" {
		c.AP.WriteAtStr(0, c.AP.H-13, c.notification)
	} else {
		c.AP.WriteAtStr(0, c.AP.H-14, c.notification)
	}
	c.AP.MoveCursor(c.index, c.AP.H-2)
}

func (c *config) Tick() bool {
	return c.handleInput()
}
