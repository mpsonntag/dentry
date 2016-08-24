// Copyright (c) 2016, Michael Sonntag (sonntag@bio.lmu.de)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

const tagFileHeader = "!Tagnotes"

// TagEnt contains information stored in data and data associated with
// this information stored in tags
type TagEnt struct {
	Tags    []string
	Content string
}

// IsTagNote returns true if a byte array starts with a specific header sequence, false if not.
func IsTagNote(cont *[]byte) (bool, error) {
	r := bytes.NewReader(*cont)
	br := bufio.NewReader(r)
	l, _, err := br.ReadLine()
	if err != nil {
		return false, err
	}
	return strings.Index(string(l), tagFileHeader) == 0, nil
}

// TextToEnt scans a byte array and splits the content at '(#)' and removes the '(#)' occurrence.
// The resulting pieces are further split at '#)'. If '#)' exists, the first part is further
// split at ',' occurrences, the individual pieces are trimmed of whitespaces and
// stored in the tags field of a new tagEnt instance. The second part is stored in the
// body part of the tagEnt instance.
// All new tagEnt instances are stored in a tagEntList instance and returned if no error
// occurred.
func TextToEnt(cont *[]byte) (*[]TagEnt, error) {
	tmp := make([]TagEnt, 0, 32)

	r := bytes.NewReader(*cont)
	s := bufio.NewScanner(r)
	s.Split(splitOnHash)
	for s.Scan() {
		curr := strings.Replace(s.Text(), "(#)", "", -1)
		currParts := strings.Split(curr, "#)\n")
		if len(currParts) > 1 {
			currTags := strings.Split(currParts[0], ",")

			for i := range currTags {
				currTags[i] = strings.TrimSpace(currTags[i])
			}
			t := TagEnt{
				Tags:    currTags,
				Content: currParts[1],
			}
			tmp = append(tmp, t)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	// TODO for testing only, remove later
	for _, entry := range tmp {
		fmt.Printf("\tTags: '%s'\n\tcontent: '%s'\n", entry.Tags, entry.Content)
	}

	return &tmp, nil
}

// splitOnHash is a function satisfying bufio SplitFunc splitting a byte array at '\n(#)'.
func splitOnHash(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 1; i < len(data); i++ {
		if data[i] == '(' && data[i+1] == '#' && data[i+2] == ')' {
			// accept the split sign only at the beginning of a line
			tmp := string(data[i-1 : i+1])
			if tmp == "\n(" {
				return i + 3, data[:i+3], nil
			}
		}
	}
	return 0, data, bufio.ErrFinalToken
}
