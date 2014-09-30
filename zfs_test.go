package zfs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/theairkit/runcmd"
)

var (
	fs = []string{
		"zroot",
		"zroot/src",
	}

	snap = []string{
		"first",
		"second",
	}
)

func TestCreateSnap(t *testing.T) {
	lRunner, err := NewZfs(runcmd.NewLocalRunner())
	if err != nil {
		t.Error(err)
	}

	// Create valid snapshot:
	if err := lRunner.CreateSnap(fs[0] + "@" + snap[0]); err != nil {
		fmt.Println("test")
		t.Error(err)
	}

	// Create invalid snapshot, error is normal:
	if err := lRunner.CreateSnap(fs[0] + "blah" + "@" + snap[0]); err != nil {
		fmt.Println(err.Error())
	}
}

func TestDestroySnap(t *testing.T) {
	lRunner, err := NewZfs(runcmd.NewLocalRunner())
	if err != nil {
		t.Error(err)
	}

	// Delete valid snapshot:
	if err := lRunner.DestroySnap(fs[0] + "@" + snap[0]); err != nil {
		t.Error(err)
	}

	// Delete invalid snapshot, error is normal:
	if err := lRunner.DestroySnap(fs[0] + "blah" + "@" + snap[0]); err != nil {
		fmt.Println(err.Error())
	}
}

func TestExistSnap(t *testing.T) {
	lRunner, err := NewZfs(runcmd.NewLocalRunner())
	if err != nil {
		t.Error(err)
	}

	// Check exists valid snapshot:
	exists, err := lRunner.ExistSnap(fs[0] + "@" + snap[0])
	if err != nil {
		t.Error(err)
	}
	if exists {
		fmt.Println(fs[0] + "@" + snap[0] + " exists")
	} else {
		fmt.Println(fs[0] + "@" + snap[0] + " does not exists")
	}

	// Checks exists invalid snapshot:
	if err := lRunner.DestroySnap(fs[0] + "blah" + "@" + snap[0]); err != nil {
		fmt.Println(err.Error())
	}
}

func TestProperty(t *testing.T) {
	lRunner, err := NewZfs(runcmd.NewLocalRunner())
	if err != nil {
		t.Error(err)
	}

	// Get valid property:
	val, err := lRunner.Property(fs[0], "readonly")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(val)

	// Get invalid property:
	val, err = lRunner.Property(fs[0], "readonly-blah")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(val)
		t.Error("must no property")
	}
}

func TestListFs(t *testing.T) {
	lRunner, err := NewZfs(runcmd.NewLocalRunner())
	if err != nil {
		t.Error(err)
	}

	// List fs: non-recursive
	list, err := lRunner.ListFs(fs[0], FS, false)
	if err != nil {
		t.Error(err)
	}
	if len(list) > 1 {
		fmt.Println(list)
		t.Error(errors.New("error list fs non-recursive: more than one fs: "))
	}
	fmt.Println(list[0])

	// List fs: recursive
	list, err = lRunner.ListFs(fs[0], FS, true)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 1 {
		fmt.Println(list)
		t.Error(errors.New("error list fs recursive: only one fs: "))
	}
	for _, fs := range list {
		fmt.Println(fs)
	}

	// List snap: non-recursive
	list, err = lRunner.ListFs(fs[1]+"@"+snap[1], SNAP, false)
	if err != nil {
		t.Error(err)
	}
	if len(list) > 1 {
		fmt.Println(list)
		t.Error(errors.New("error list snap non-recursive: more than one fs: "))
	}
	fmt.Println(list[0])

	// List snap: recursive
	list, err = lRunner.ListFs(fs[0], SNAP, true)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 1 {
		fmt.Println(list)
		t.Error(errors.New("error list fs recursive: only one fs: "))
	}
	for _, fs := range list {
		fmt.Println(fs)
	}
}
