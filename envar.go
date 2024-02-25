package envar

import (
	"bytes"
	"io"
	"os"
	"strings"
)

func Overload(filenames ...string) (err error) {
	filenames = defaultOrFilenames(filenames)

	for _, filename := range filenames {
		buf, err := getFileBuffer(filename)
		if err != nil {
			return err
		}
		envMap, err := parse(&buf)
		if err != nil {
			return err
		}

		err = loadEnv(envMap, true)
		if err != nil {
			return err
		}
	}

	return
}

func Load(filenames ...string) (err error) {
	filenames = defaultOrFilenames(filenames)

	for _, filename := range filenames {
		buf, err := getFileBuffer(filename)
		if err != nil {
			return err
		}
		envMap, err := parse(&buf)
		if err != nil {
			return err
		}

		err = loadEnv(envMap, false)
		if err != nil {
			return err
		}
	}

	return
}

func loadEnv(envMap map[string]string, overload bool) (err error) {
	rawEnv := os.Environ()
	envStateMap := make(map[string]bool)

	for _, envEntry := range rawEnv {
		key := strings.Split(envEntry, "=")[0]
		envStateMap[key] = true
	}

	for key, value := range envMap {
		if !envStateMap[key] || overload {
			err := os.Setenv(key, value)
			if err != nil {
				return err
			}
		}
	}

	return
}

func defaultOrFilenames(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}

	return filenames
}

func getFileBuffer(filename string) (buf bytes.Buffer, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(&buf, file)
	return
}
