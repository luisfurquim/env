package env

import (
	"os"
	"fmt"
	"time"
	"regexp"
	"errors"
	"reflect"
	"strings"
	"github.com/luisfurquim/goose"
)


type GooseG struct {
	Env goose.Alert
}

var timeFormat reflect.Type = reflect.TypeOf(time.Time{})
var durationFormat reflect.Type = reflect.TypeOf(time.Duration(0))

var reFloatFormat *regexp.Regexp = regexp.MustCompile(`^\%[beEfFgGxX]$`)

var Goose GooseG = GooseG{
	Env: goose.Alert(3),
}

var ErrInvalidType error = errors.New("Invalid type")
var ErrInvalidFormat error = errors.New("Invalid format")
var ErrMissingRequiredVariable  error = errors.New("Missing required variable")
var ErrInvalidRequiredPolicy error = errors.New("Invalid required policy")

func Read(d interface{}) error {
	var data reflect.Value
	var fld reflect.StructField
	var typ, fldTyp reflect.Type
	var i int
	var ok bool
	var tag string
	var sval string
	var ival int64
	var uival uint64
	var fval float64
	var tmval time.Time
	var dval time.Duration
	var format string
	var fldRef reflect.Value
	var err error

	data = reflect.ValueOf(d)
	if data.Kind() != reflect.Pointer {
		Goose.Env.Logf(1, "Input parameter must be called by reference: %s", ErrInvalidType)
		return ErrInvalidType
	}

	data = data.Elem()

	if data.Kind() != reflect.Struct {
		Goose.Env.Logf(1, "Input parameter must point to a struct type: %s", ErrInvalidType)
		return ErrInvalidType
	}

	typ = data.Type()

	for i=0; i<typ.NumField(); i++ {
		fld = typ.Field(i)
		if tag, ok = fld.Tag.Lookup("env"); !ok {
			continue
		}

		if tag == "" {
			continue
		}

		if sval, ok = os.LookupEnv(tag); !ok {
			if sval, ok = fld.Tag.Lookup("default"); !ok {
				sval, ok = fld.Tag.Lookup("required")
				if !ok {
					continue
				}

				sval = strings.ToLower(sval)
				if sval == "yes" || sval == "true" {
					return ErrMissingRequiredVariable
				}

				if sval == "no" || sval == "false" {
					continue
				}

				return ErrInvalidRequiredPolicy
			}
		}

		fldRef = data.Field(i)
		fldTyp = fld.Type
		if fldTyp.Kind() == reflect.Pointer {
			fldRef = fldRef.Elem()
			fldTyp = fldTyp.Elem()
		}
		
		if fldTyp == timeFormat {
			if format, ok = fld.Tag.Lookup("format"); !ok {
				format = "2006-01-02 15:04:05"
			}

			tmval, err = time.Parse(format, sval)
			if err != nil {
				Goose.Env.Logf(1, "Error scanning %s time value for %s environment variable: %s", sval, tag, err)
				return err
			}

			fldRef.Set(reflect.ValueOf(tmval))

			continue
		}

		if fldTyp == durationFormat {
			dval, err = time.ParseDuration(sval)
			if err != nil {
				Goose.Env.Logf(1, "Error scanning %s duration value for %s environment variable: %s", sval, tag, err)
				return err
			}

			fldRef.SetInt(int64(dval))

			continue
		}

		switch fldTyp.Kind() {
		case reflect.String:
			fldRef.SetString(sval)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			_, err = fmt.Sscanf(sval, "%d", &ival)
			if err != nil {
				Goose.Env.Logf(1, "Error scanning %s int value for %s environment variable: %s", sval, tag, err)
				return err
			}
			fldRef.SetInt(ival)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			_, err = fmt.Sscanf(sval, "%u", &uival)
			if err != nil {
				Goose.Env.Logf(1, "Error scanning %s uint value for %s environment variable: %s", sval, tag, err)
				return err
			}
			fldRef.SetUint(uival)

		case reflect.Float32, reflect.Float64:
			format, ok = fld.Tag.Lookup("format")
			if ok {
				if !reFloatFormat.MatchString(format) {
					Goose.Env.Logf(1,"Error scanning %s float format for %s environment variable: %s", format, tag, ErrInvalidFormat)
					return ErrInvalidFormat
				}
			} else {
				format = "%f"
			}
			_, err = fmt.Sscanf(sval, format, &fval)
			if err != nil {
				Goose.Env.Logf(1, "Error scanning %s float value for %s environment variable: %s", sval, tag, err)
				return err
			}
			fldRef.SetFloat(fval)

		case reflect.Bool:
			sval = strings.ToLower(sval)
			if sval == "yes" || sval == "true" {
				fldRef.SetBool(true)
			}

			if sval == "no" || sval == "false" {
				fldRef.SetBool(false)
			}

		default:
			Goose.Env.Logf(1, "Error scanning %s value for %s environment variable: %s", sval, tag, ErrInvalidType)
			return ErrInvalidType
		}
	}

	return nil
}
