package sitemap

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
)

func entryParser(decoder *xml.Decoder, se *xml.StartElement, consume EntryConsumer) error {
	if se.Name.Local == "url" {
		entry := newSitemapEntry()

		decodeError := decoder.DecodeElement(entry, se)
		if decodeError != nil {
			return decodeError
		}

		consumerError := consume(entry)
		if consumerError != nil {
			return consumerError
		}
	}

	return nil
}

func indexEntryParser(decoder *xml.Decoder, se *xml.StartElement, consume IndexEntryConsumer) error {
	if se.Name.Local == "sitemap" {
		entry := new(sitemapIndexEntry)

		decodeError := decoder.DecodeElement(entry, se)
		if decodeError != nil {
			return decodeError
		}

		consumerError := consume(entry)
		if consumerError != nil {
			return consumerError
		}
	}

	return nil
}

type elementParser func(*xml.Decoder, *xml.StartElement) error

func parseLoop(ctx context.Context, reader io.Reader, parser elementParser) error {
	decoder := xml.NewDecoder(reader)

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			t, tokenError := decoder.Token()

			if tokenError == io.EOF {
				return nil
			} else if tokenError != nil {
				fmt.Println("tokenError:", tokenError)
				return tokenError
			}

			se, isStartElement := t.(xml.StartElement)
			if !isStartElement {
				continue
			}

			parserError := parser(decoder, &se)
			if parserError != nil {
				fmt.Println("parserError:", parserError)
				return parserError
			}
		}
	}
}
