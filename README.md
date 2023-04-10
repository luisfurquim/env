# env
Simplify reading environment variables in golang

## Example

To read environment variables, first you need to create a struct type with, at least, the 'env' tag, where it defines the name of the environment variable. The optional 'required' and 'default' tags may also be defined. The 'required' tag accepts the true/yes/false/no values. The 'default' tag accepts any value, but its value must be a compatible representation of values acceptable by the type of the related struct field. The reflect.Kind of the struct field must be one of String, Int*, Uint*, Float*, Bool or time.Time or time.Duration (indeed it is of reflect.Int64 kind, but it has special treatment). Pointers to theses types are also accepted.

```Go


type Options struct {
	Threshold      *float64       `env:"THRESHOLD"`
	RunningContext *string        `env:"CONTEXT" required:"yes"`
	Verbose			 goose.Alert   `env:"VERBOSE" default:"4"`
	Id              int           `env:"ID" required:"yes"`
	StartDate       time.Time     `env:"BEGIN" required:"true"`
	Period          time.Duration `env:"PERIOD" required:"yes"`
}


```

Then you just need to call Read() function.

```Go

   .
   .
   .

   var options Options

	options.RunningContext   = &otherPackageOrStruct.RunningContext

	// You must pass options by reference, otherwise we couldn't set fields ther than pointer fields.
	err = env.Read(&options)
	if err != nil {
		...
	}

   .
   .
   .

```
