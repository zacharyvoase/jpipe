package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	keySep        = kingpin.Flag("key", "Key separator character").Short('k').Default("/").String()
	columnSep     = kingpin.Flag("column", "Column separator character").Short('c').Default("\t").String()
	nullDelimited = kingpin.Flag("null", "Separate key-value pairs with \\0 instead of newlines").Short('z').Bool()
	ascii         = kingpin.Flag("ascii", "Ensure JSON strings are output only in ASCII format").Short('a').Bool()
)

type jsonNode struct {
	keyPath      []string
	encodedValue string
}

func extendKeyPath(keyPath []string, key string) []string {
	newKeyPath := make([]string, len(keyPath), len(keyPath)+1)
	copy(newKeyPath, keyPath)
	return append(newKeyPath, key)
}

func unwrap(keyPath []string, obj interface{}, output chan jsonNode) {
	switch obj := obj.(type) {
	default:
		// This should never happen
		log.Println("Can't decode object %s", obj)
	case string:
		if *ascii {
			output <- jsonNode{keyPath, encodeStringAsASCII(obj)}
		} else {
			marshalledBytes, _ := json.Marshal(obj)
			output <- jsonNode{keyPath, string(marshalledBytes)}
		}
	case bool, nil:
		marshalledBytes, _ := json.Marshal(obj)
		output <- jsonNode{keyPath, string(marshalledBytes)}
	case json.Number:
		output <- jsonNode{keyPath, obj.String()}
	case map[string]interface{}:
		output <- jsonNode{keyPath, "{}"}
		for k, v := range obj {
			if strings.Contains(k, *keySep) {
				log.Println(fmt.Sprintf(
					"key \"%s\" contains key separator \"%s\"", k, *keySep))
				os.Exit(1)
			}
			unwrap(extendKeyPath(keyPath, k), v, output)
		}
	case []interface{}:
		output <- jsonNode{keyPath, "[]"}
		for i, v := range obj {
			unwrap(extendKeyPath(keyPath, strconv.Itoa(i)), v, output)
		}
	}
}

var hex = "0123456789abcdef"

func encodeStringAsASCII(str string) string {
	var output bytes.Buffer
	output.WriteByte('"')
	for _, b := range bytes.Runes([]byte(str)) {
		if b < utf8.RuneSelf {
			switch b {
			case '\\', '"':
				output.WriteByte('\\')
				output.WriteByte(byte(b))
			case '\n':
				output.WriteByte('\\')
				output.WriteByte('n')
			case '\r':
				output.WriteByte('\\')
				output.WriteByte('r')
			case '\t':
				output.WriteByte('\\')
				output.WriteByte('t')
			default:
				if b < 0x20 || b == '<' || b == '>' || b == '&' {
					output.WriteString(`\u00`)
					output.WriteByte(hex[b>>4])
					output.WriteByte(hex[b&0xF])
				} else {
					output.WriteByte(byte(b))
				}
			}
		} else {
			output.WriteString(fmt.Sprintf("\\u%04x", b))
		}
	}
	output.WriteByte('"')
	return output.String()
}

func main() {
	kingpin.Parse()

	dec := json.NewDecoder(os.Stdin)
	dec.UseNumber()
	output := make(chan jsonNode)
	done := make(chan int32)
	go func() {
		for node := range output {
			os.Stdout.WriteString(*keySep)
			os.Stdout.WriteString(strings.Join(node.keyPath, *keySep))
			os.Stdout.WriteString(*columnSep)
			os.Stdout.WriteString(node.encodedValue)
			if *nullDelimited {
				os.Stdout.WriteString("\x00")
			} else {
				os.Stdout.WriteString("\n")
			}
			os.Stdout.Sync()
		}
		done <- 1
	}()
	for {
		var v interface{}
		if err := dec.Decode(&v); err != nil {
			if err != io.EOF {
				log.Println("ERROR: %s", err)
			}
			close(output)
			break
		}
		unwrap(make([]string, 0), v, output)
	}
	<-done
}
