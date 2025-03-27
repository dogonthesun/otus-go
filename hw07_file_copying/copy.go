package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type AnsiProgressOutput struct {
	limit int64
	done  int64
}

func NewAnsiProgressOutput(limit int64) *AnsiProgressOutput {
	return &AnsiProgressOutput{
		limit: limit,
		done:  0,
	}
}

func (a *AnsiProgressOutput) Init() {
	fmt.Printf("writing...\n")
}

func (a *AnsiProgressOutput) Update(n int64) {
	a.done = min(a.limit, a.done+n)
	fmt.Printf("\033[1A\033[K%v%% written\n", int(float64(a.done)/float64(a.limit)*100))
}

type DummyProgressOutput struct {
	limit int64
	done  int64
	slots []int
}

func NewDummyProgressOutput(limit int64) *DummyProgressOutput {
	return &DummyProgressOutput{
		limit: limit,
		done:  0,
		slots: []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
	}
}

func (d *DummyProgressOutput) Init() {
	fmt.Printf("writing...")
}

func (d *DummyProgressOutput) Update(n int64) {
	if d.done < d.limit {
		d.done = min(d.limit, d.done+n)
		pct := int(float64(d.done) / float64(d.limit) * 100)
		for len(d.slots) > 0 && d.slots[0] <= pct {
			fmt.Printf(" %v%%", d.slots[0])
			d.slots = d.slots[1:]
		}
		if pct == 100 {
			fmt.Println(" done")
		}
	}
}

type ProgressOutput interface {
	Init()
	Update(int64)
}

type ProgressWriter struct {
	fd *os.File
	po ProgressOutput
}

func NewProgressWriter(fd *os.File, limit int64) *ProgressWriter {
	w := &ProgressWriter{
		fd: fd,
		po: nil,
	}
	if term.IsTerminal(int(os.Stdout.Fd())) {
		w.po = NewAnsiProgressOutput(limit)
	}
	if w.po == nil {
		w.po = NewDummyProgressOutput(limit)
	}
	w.po.Init()
	return w
}

func (w *ProgressWriter) Write(p []byte) (int, error) {
	n, err := w.fd.Write(p)
	if err == nil {
		w.po.Update(int64(n))
	}
	return n, err
}

var (
	ErrInvalidBoundaries = fmt.Errorf("invalid boundaries")
	ErrInvalidFileType   = fmt.Errorf("invalid file type")
)

type ErrInputOutput struct {
	msg string
	err error
}

func (e ErrInputOutput) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e ErrInputOutput) Unwrap() error {
	return e.err
}

func openPath(path string, flags int, perm os.FileMode) (*os.File, os.FileInfo, error) {
	fd, err := os.OpenFile(path, flags, perm)
	if err != nil {
		return nil, nil, err
	}
	fstat, err := fd.Stat()
	if err != nil {
		fd.Close()
		return nil, nil, err
	}
	return fd, fstat, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return fmt.Errorf("%w: empty file path input or output", ErrInvalidFileType)
	}

	ifd, ifstat, err := openPath(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		return ErrInputOutput{fmt.Sprintf("open %s", fromPath), err}
	}
	defer ifd.Close()

	if !ifstat.Mode().IsRegular() {
		return fmt.Errorf("%w: input is not a regular file", ErrInvalidFileType)
	}

	srcFileSize := ifstat.Size()
	if offset > srcFileSize {
		return fmt.Errorf("%w: offset %v exceeds file size %v", ErrInvalidBoundaries, offset, srcFileSize)
	}
	if offset < 0 {
		return fmt.Errorf("%w: negative offset", ErrInvalidBoundaries)
	}

	if _, err := ifd.Seek(offset, 0); err != nil {
		return ErrInputOutput{fmt.Sprintf("seek %v", offset), err}
	}

	ofd, ofstat, err := openPath(toPath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return ErrInputOutput{fmt.Sprintf("open %s", toPath), err}
	}
	defer ofd.Close()

	if os.SameFile(ifstat, ofstat) {
		return fmt.Errorf("%w: %s and %s are identical", ErrInvalidFileType, fromPath, toPath)
	}

	if err := ofd.Truncate(0); err != nil {
		return ErrInputOutput{fmt.Sprintf("truncate %s", toPath), err}
	}

	if limit < 0 {
		return fmt.Errorf("%w: negative limit", ErrInvalidBoundaries)
	}

	n := srcFileSize - offset
	if limit > 0 {
		n = min(limit, n)
	}

	reader := io.NewSectionReader(ifd, offset, n)
	writer := NewProgressWriter(ofd, n)

	if _, err = io.Copy(writer, reader); err != nil && !errors.Is(err, io.EOF) {
		return ErrInputOutput{fmt.Sprintf("coping from %s to %s", fromPath, toPath), err}
	}

	return nil
}
