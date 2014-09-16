package zfs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/theairkit/runcmd"
)

var (
	fs = []string{
		"zroot/src",       //fs must be exists
		"zroot/src/host1", //fs must be exists
		"zroot/blah",      //fs must be not exists
		"zroot/src*",      //valid mask
		"zroot/src/*",     //valid mask
		"zroot/blah*",     //invalid mask
	}
)

func TestListFs(t *testing.T) {
	lRunner, err := NewZfs(runcmd.NewLocalRunner())
	if err != nil {
		t.Error(err)
	}

	// Test ListFs: fs by name:
	out, err := lRunner.ListFs(fs[0], FS, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[0] + ", FS, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[2], FS, false)
	if err != nil {
		fmt.Println("lRunner.ListFs(" + fs[2] + ", FS, false)")
		fmt.Println(err.Error())
		fmt.Println("")
	} else {
		t.Error(errors.New("something wrong"))
	}

	out, err = lRunner.ListFs(fs[0], FS, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[0] + ", FS, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[2], FS, true)
	if err != nil {
		fmt.Println("lRunner.ListFs(" + fs[2] + ", FS, true)")
		fmt.Println(err.Error())
		fmt.Println("")
	} else {
		t.Error(errors.New("something wrong"))
	}

	// Test listFs: fs by mask:
	out, err = lRunner.ListFs(fs[3], FS, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[3] + ", FS, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[4], FS, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[4] + ", FS, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[5], FS, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[5] + ", FS, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[3], FS, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[3] + ", FS, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[4], FS, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[4] + ", FS, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[5], FS, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[5] + ", FS, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	// Test listFs: snapshot by name:
	out, err = lRunner.ListFs(fs[0], SNAP, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[0] + ", SNAP, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[1], SNAP, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[1] + ", SNAP, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[2], SNAP, false)
	if err != nil {
		fmt.Println("lRunner.ListFs(" + fs[2] + ", SNAP, false)")
		fmt.Println(err.Error())
		fmt.Println("")
	} else {
		t.Error(errors.New("something wrong"))
	}

	out, err = lRunner.ListFs(fs[0], SNAP, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[0] + ", SNAP, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[1], SNAP, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[1] + ", SNAP, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[2], SNAP, true)
	if err != nil {
		fmt.Println("lRunner.ListFs(" + fs[2] + ", SNAP, true)")
		fmt.Println(err.Error())
		fmt.Println("")
	} else {
		t.Error(errors.New("something wrong"))
	}

	// Test ListFs: snapshot by mask:
	out, err = lRunner.ListFs(fs[3], SNAP, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[3] + ", SNAP, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[4], SNAP, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[4] + ", SNAP, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[5], SNAP, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[5] + ", SNAP, false)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[3], SNAP, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[3] + ", SNAP, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[4], SNAP, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[4] + ", SNAP, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")

	out, err = lRunner.ListFs(fs[5], SNAP, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("lRunner.ListFs(" + fs[5] + ", SNAP, true)")
	for _, f := range out {
		fmt.Println(f)
	}
	fmt.Println("")
}
