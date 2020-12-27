package etymology

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Word struct {
	Language string
	Word     string
}

type Etymology struct {
	derivedFromEdges map[Word][]Word
	etymologyEdges   map[Word][]Word
}

func New(data io.Reader) (Etymology, error) {
	etymologyEdges := make(map[Word][]Word)
	derivedFromEdges := make(map[Word][]Word)

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		components := strings.Split(scanner.Text(), "\t")

		var m map[Word][]Word

		switch components[1] {
		case "rel:etymology":
			m = etymologyEdges
		case "rel:is_derived_from":
			m = derivedFromEdges
		default:
			continue
		}

		from := strings.Split(components[0], ": ")
		to := strings.Split(components[2], ": ")

		fromWord := Word{
			Language: from[0],
			Word:     from[1],
		}

		if components[1] == "rel:is_derived_from" && m[fromWord] != nil {
			fmt.Println("dupe: ", from)
		}

		m[fromWord] = append(m[fromWord], Word{
			Language: to[0],
			Word:     to[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return Etymology{}, fmt.Errorf("failed to scan file: %w", err)
	}

	return Etymology{etymologyEdges: etymologyEdges, derivedFromEdges: derivedFromEdges}, nil
}

type Relation int

const (
	RelationEtymology = iota
	RelationDerivedFrom
)

type Node struct {
	Word        Word
	DerivedFrom []Node
	Etymology   []Node
}

func (e Etymology) Lookup(word Word) Node {
	node := Node{
		Word: word,
	}

	for _, derivedFrom := range e.derivedFromEdges[word] {
		n := e.Lookup(derivedFrom)
		node.DerivedFrom = append(node.DerivedFrom, n)
	}

	for _, etym := range e.etymologyEdges[word] {
		n := e.Lookup(etym)
		node.Etymology = append(node.Etymology, n)
	}

	return node
}