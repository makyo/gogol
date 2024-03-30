package rle

import (
	"fmt"
	"strconv"
	"strings"
)

type RLEField struct {
	Width, Height, Top, Left  int
	Field                     [][]bool
	Name, Origin              string
	Comments, ExtendedRLEData []string
	Survive, Born             []int
}

// Marshal generates the contents of an RLE file from a given field.
func (f *RLEField) Marshal() string {
	var out strings.Builder

	// Write the # lines
	if f.Name != "" {
		fmt.Fprintf(&out, "#N %s\n", f.Name)
	}
	if f.Origin != "" {
		fmt.Fprintf(&out, "#O %s\n", f.Origin)
	}
	if len(f.Comments) > 0 {
		for _, comment := range f.Comments {
			fmt.Fprintf(&out, "#C %s\n", comment)
		}
	}
	fmt.Fprintf(&out, "#R %d  %d\n", f.Left, f.Top)

	// Write the header
	fmt.Fprintf(&out, "x = %d, y = %d, rule = B", f.Width, f.Height)
	for _, b := range f.Born {
		fmt.Fprintf(&out, "%d", b)
	}
	fmt.Fprint(&out, "/S")
	for _, s := range f.Survive {
		fmt.Fprintf(&out, "%d", s)
	}
	fmt.Fprint(&out, "\n")

	// Write the content
	var count int
	var chunks []string
	for _, row := range f.Field {
		count = 1
		for x, col := range row {

			// No need to notate dead cells up until the end of the line (or empty lines).
			toContinue := false
			for i := x; i < len(row); i++ {
				toContinue = toContinue || row[i]
			}
			if !toContinue {
				break
			}

			if x < len(row)-1 && row[x+1] == col {
				count++
				continue
			}

			// 'o' for living cells, 'b' for dead cells. For some reason ?.?
			var curr string
			if col {
				curr = "o"
			} else {
				curr = "b"
			}

			// Print the count and cell state (or just cell state if count is 1).
			if count == 1 {
				chunks = append(chunks, curr)
			} else {
				chunks = append(chunks, fmt.Sprintf("%d%s", count, curr))
				count = 1
			}
		}
		chunks = append(chunks, "$")
	}

	// Write the rule chunks to the output, collapsing repeated EOLs and adding linebreaks.
	count = 0
	lineLen := 0
	for _, chunk := range chunks {
		if chunk == "$" {
			count++
			continue
		} else {
			if count > 1 {
				chunk = fmt.Sprintf("%d$", count) + chunk
			} else if count == 1 {
				chunk = "$" + chunk
			}
			count = 0
			lineLen += len(chunk)
			if lineLen > 70 {
				chunk = "\n" + chunk
				lineLen = len(chunk)
			}
			fmt.Fprint(&out, chunk)
		}
	}

	// End with a bang.
	fmt.Fprint(&out, "!\n")
	return out.String()
}

// Unmarshal builds a field from the contents of an RLE file.
func Unmarshal(contents string) (*RLEField, error) {
	f := &RLEField{
		Width:   -1,
		Height:  -1,
		Born:    []int{3},
		Survive: []int{2, 3},
	}

	// Keep track of whether the header has been seen; don't process rules otherwise.
	headerSeen := false

	// Keep track of current coordinates for building the field
	x := 0
	y := 0

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for # lines
		if line[0] == byte('#') {
			header, content, found := strings.Cut(line, " ")
			if !found {
				return nil, fmt.Errorf("Malformed # line - contains no space: %q", line)
			}

			switch header {
			case "#C", "#c":
				// Comments.
				f.Comments = append(f.Comments, content)

			case "#CXRLE":
				// Extended RLE information (which we don't use, but okay).
				f.ExtendedRLEData = append(f.ExtendedRLEData, content)

			case "#N":
				// The name of the pattern
				f.Name = content

			case "#O":
				// When and by whom the file was created.
				f.Origin = content

			case "#R":
				// The coordinates of the top-left corner of the pattern.
				coords := strings.Split(content, " ")
				if len(coords) == 2 {
					return nil, fmt.Errorf("Malformed # line - #R line should contain integer X and Y values separated by a space: %q", content)
				}
				cx, err := strconv.Atoi(coords[0])
				if err != nil {
					return nil, fmt.Errorf("Malformed # line - #R line should contain integer X and Y values separated by a space: %q", content)
				}
				cy, err := strconv.Atoi(coords[1])
				if err != nil {
					return nil, fmt.Errorf("Malformed # line - #R line should contain integer X and Y values separated by a space: %q", content)
				}
				f.Left = cx
				f.Top = cy
			case "#r":
				// Additional rule stuff from XLife that we'll just discard for now.

			default:
				return nil, fmt.Errorf("Malformed # line - unknown header: %s", header)
			}
			continue
		}

		// Process the header rule
		if line[0] == byte('x') {
			headerSeen = true
			pairs := strings.Split(line, ",")
			for _, pair := range pairs {

				// Process key/value pairs
				k, v, found := strings.Cut(strings.ReplaceAll(strings.TrimSpace(pair), " ", ""), "=")
				if !found {
					return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S##': %q", line)
				}

				switch k {
				case "x":
					// Set width.
					width, err := strconv.Atoi(v)
					if err != nil || width < 1 {
						return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S#': %q", line)
					}
					f.Width = width

				case "y":
					// Set height.
					height, err := strconv.Atoi(v)
					if err != nil || height < 1 {
						return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S#': %q", line)
					}
					f.Height = height

				case "rule":
					// Parse the rule in birth/survival notation (see: https://conwaylife.com/wiki/Rulestring )
					born, survive, found := strings.Cut(v, "/")
					if !found {
						return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S#': %q", line)
					}

					// Get the parts
					_, born, found = strings.Cut(born, "B")
					if !found {
						return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S#': %q", line)
					}
					_, survive, found = strings.Cut(survive, "S")
					if !found {
						return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S#': %q", line)
					}

					// Build the list of values
					f.Born = []int{}
					f.Survive = []int{}
					if len(born) > 0 {
						_, err := strconv.Atoi(born)
						if err != nil {
							return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S#': %q", line)
						}
						for _, b := range born {
							s, _ := strconv.Atoi(string(b))
							f.Born = append(f.Born, s)
						}
					}
					if len(survive) > 0 {
						_, err := strconv.Atoi(survive)
						if err != nil {
							return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S#': %q", line)
						}
						for _, s := range survive {
							s, _ := strconv.Atoi(string(s))
							f.Survive = append(f.Survive, s)
						}
					}

				default:
					return nil, fmt.Errorf("Malformed header line - must take the form 'x = m, y = n' with an optional ',  rule = B#/S#': %q", line)
				}
			}

			// No x/y provided is an error.
			if f.Width < 1 || f.Height < 1 {
				return nil, fmt.Errorf("Malformed header line - width and height must be positive: %q", line)
			}

			// Make the field.
			f.Field = make([][]bool, f.Height)
			for i, _ := range f.Field {
				f.Field[i] = make([]bool, f.Width)
			}
			continue
		}

		// Process the rule itself, but only if we already have a header.
		if headerSeen {
			count := 0
			for _, char := range line {
				switch string(char) {
				case "b":
					// Dead.
					if count == 0 {
						count = 1
					}
					for i := 0; i < count; i++ {
						f.Field[y][x] = false
						x++
					}
					count = 0

				case "o":
					// Alive.
					if count == 0 {
						count = 1
					}
					for i := 0; i < count; i++ {
						f.Field[y][x] = true
						x++
					}
					count = 0

				case "$":
					// End of line.
					if count == 0 {
						count = 1
					}
					x = 0
					y += count
					count = 0

				case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
					// Multiplier.
					multiplier, _ := strconv.Atoi(string(char))
					count = count*10 + multiplier

				case "!":
					break

				default:
					return nil, fmt.Errorf("Malformed rule - unexpected character '%s' in rule definition", string(char))
				}
			}
		}
	}

	if !headerSeen {
		return nil, fmt.Errorf("Malformed rule - no header")
	}

	return f, nil
}
