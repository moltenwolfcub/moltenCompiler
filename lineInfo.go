package main

import "fmt"

type LineInfo struct {
	File string
	Line int
	Col  int
}

func NewLineInfo(fileName string) LineInfo {
	return LineInfo{
		File: fileName,
		Line: 1,
		Col:  1,
	}
}

func (l *LineInfo) NextLine() {
	l.Line++
	l.Col = 1
}

func (l *LineInfo) IncColumn() {
	l.Col++
}

func (l *LineInfo) IncWord(word []rune) {
	l.Col += len(word)
}

func (l LineInfo) PositionedError(message string) error {
	return fmt.Errorf("%s:%d:%d: %s", l.File, l.Line, l.Col, message)
}
