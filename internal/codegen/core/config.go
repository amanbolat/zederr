package core

type Config struct {
	SpecPath string
	ExportGo ExportGo
}

type ExportGo struct {
	PackageName string
	OutputPath  string
}
