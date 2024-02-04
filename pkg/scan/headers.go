/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"errors"
)

var (
	ErrWrongHeaderFormat = errors.New("header with wrong format")
)

/*
// GetHeaders returns the headers map.
func GetHeaders(r *Runner) (map[string]interface{}, error) {
	var (
		headers []string
		result  map[string]interface{}
	)

	if r.Options.HeadersFile != "" {
		h, err := readHeadersFile(r.Options.HeadersFile)
		if err != nil {
			return result, err
		}

		headers = h
	}

	if len(r.Options.Headers) != 0 {
		headers = r.Options.Headers
	}

	return readHeaders(headers)
}

func readHeaders(input []string) (map[string]interface{}, error) {
	headers := map[string]interface{}{}

	if len(input) == 0 {
		return map[string]interface{}{}, nil
	}

	for _, h := range input {
		hName, hValue, err := splitHeader(h)
		if err != nil {
			return nil, err
		}

		headers[hName] = hValue
	}

	return headers, nil
}

func splitHeader(header string) (string, string, error) {
	splitted := strings.Split(header, ":")
	if len(splitted) == 1 {
		return "", "", fmt.Errorf("%s: %w", splitted[0], ErrWrongHeaderFormat)
	}

	headerName := splitted[0]
	headerValue := splitted[1]

	return headerName, headerValue, nil
}

func readHeadersFile(inputFile string) ([]string, error) {
	var text []string

	file, err := os.Open(inputFile)
	if err != nil {
		return text, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		dir := scanner.Text()
		if len(dir) > 0 {
			text = append(text, dir)
		}
	}

	file.Close()

	return text, nil
}

*/
