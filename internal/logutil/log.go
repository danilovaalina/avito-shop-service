package logutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = errorStackMarshaller
	log.Logger = log.Output(beautifulWriter(io.Writer(os.Stderr)))
}

func beautifulWriter(out io.Writer) zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{
		Out:           out,
		TimeFormat:    time.RFC3339Nano,
		FieldsExclude: []string{zerolog.ErrorFieldName, zerolog.ErrorStackFieldName},
		FormatExtra: func(fields map[string]any, buf *bytes.Buffer) error {
			for _, f := range []string{zerolog.ErrorFieldName, zerolog.ErrorStackFieldName} {
				v, ok := fields[f]
				if !ok {
					continue
				}

				err := formatField(buf, f, v)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func formatField(buf *bytes.Buffer, field string, value any) error {
	s := fmt.Sprintf("\n\x1b[%dm%v=\x1b[0m\"\n%v\"", 36, field, value)

	_, err := buf.WriteString(s)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// errorStackMarshaller for cockroachdb/errors.
func errorStackMarshaller(err error) any {
	return fmt.Sprintf("%+v", err)
}
