// This file is part of Netsak
// Copyright (C) 2018 matterbeam
// Use of this source is governed by a GPLv3

package nsfmt


// State represents the current state passed to the formatters.  It provides
// access to the io.Writer interface.
type State interface {

    // Write is the function to call to insert text in the place of the tag.
    // It writes at most len(p).  Returns the number of bytes written and
    // an error, if any.
    Write(p []byte) (n int, err error)

    // WriteString is like Write but allows to write a string directly.
    WriteString(s string) (n int, err error)

}

