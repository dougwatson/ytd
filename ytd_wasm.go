package main

func writeFile(destFile string, b []byte) error {
	return fs.AddFile(destFile, string(b))
}
