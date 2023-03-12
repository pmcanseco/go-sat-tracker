package display

type Printer interface {
	Print(s string)
	PrintAt(line int, s string, clear bool)
}
