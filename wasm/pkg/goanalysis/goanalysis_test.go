package goanalysis

import (
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAppendImports(t *testing.T) {
	existingContent := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}`

	type args struct {
		fileContent      string
		importStatements []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Add single import statement",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"strings"},
			},
			want: `package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "Add multiple import statements",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"st strings", "strconv", "os"},
			},
			want: `package main

import (
	"fmt"
	"os"
	"strconv"
	st "strings"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "Add multiple import statements with an existing one",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"st strings", "strconv", "os", "fmt"},
			},
			want: `package main

import (
	"fmt"
	"os"
	"strconv"
	st "strings"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "Add duplicate import statement",
			args: args{
				fileContent:      existingContent,
				importStatements: []string{"fmt"},
			},
			want: existingContent + "\n",
			err:  nil,
		},
		{
			name: "No import statement",
			args: args{
				fileContent: `package main

func main() {
	fmt.Println("Hello, world!")
}`,
				importStatements: []string{"fmt"},
			},
			want: `package main

import  "fmt"

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "No import statement and add two imports",
			args: args{
				fileContent: `package main

func main() {
	fmt.Println("Hello, world!")
}`,
				importStatements: []string{"fmt", "os"},
			},
			want: `package main

import (
	 "fmt"
	 "os"
)

func main() {
	fmt.Println("Hello, world!")
}
`,
			err: nil,
		},
		{
			name: "Add empty file content",
			args: args{
				fileContent:      "",
				importStatements: []string{"fmt"},
			},
			err: errors.New("1:1: expected 'package', found 'EOF'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendImports(tt.args.fileContent, tt.args.importStatements...)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAppendCode(t *testing.T) {
	existingContent := `package main

import (
    "fmt"
)

func main() {
    fmt.Println("Hello, world!")
}

func anotherFunction() bool {
    // Some code here
    fmt.Println("Another function")
    return true
}`

	type args struct {
		fileContent  string
		functionName string
		codeToInsert string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Append code to the end of the function",
			args: args{
				fileContent:  existingContent,
				functionName: "main",
				codeToInsert: "fmt.Println(\"Inserted code here\")",
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
	fmt.Println("Inserted code here")

}

func anotherFunction() bool {
	// Some code here
	fmt.Println("Another function")
	return true
}
`,
		},
		{
			name: "Append code with return statement",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				codeToInsert: "fmt.Println(\"Inserted code here\")",
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}

func anotherFunction() bool {
	// Some code here
	fmt.Println("Another function")
	fmt.Println("Inserted code here")

	return true
}
`,
		},
		{
			name: "Function not found",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				codeToInsert: "fmt.Println(\"Inserted code here\")",
			},
			err: errors.New("function nonexistentFunction not found"),
		},
		{
			name: "Invalid code",
			args: args{
				fileContent:  existingContent,
				functionName: "anotherFunction",
				codeToInsert: "%#)(u309f/..\"",
			},
			err: errors.New("1:1: expected operand, found '%' (and 2 more errors)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendCode(tt.args.fileContent, tt.args.functionName, tt.args.codeToInsert)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestReplaceReturn(t *testing.T) {
	existingContent := `package main

import (
    "fmt"
)

func main() {
    x := calculate()
    fmt.Println("Result:", x)
}

func calculate() int {
    return 42
}`

	type args struct {
		fileContent  string
		functionName string
		returnVars   []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Replace return statement with a single variable",
			args: args{
				fileContent:  existingContent,
				functionName: "calculate",
				returnVars:   []string{"result"},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	x := calculate()
	fmt.Println("Result:", x)
}

func calculate() int {
	return result

}
`,
		},
		{
			name: "Replace return statement with multiple variables",
			args: args{
				fileContent:  existingContent,
				functionName: "calculate",
				returnVars:   []string{"result", "err"},
			},
			want: `package main

import (
	"fmt"
)

func main() {
	x := calculate()
	fmt.Println("Result:", x)
}

func calculate() int {
	return result, err

}
`,
		},
		{
			name: "Function not found",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				returnVars:   []string{"result"},
			},
			err: errors.New("function nonexistentFunction not found"),
		},
		{
			name: "Invalid result",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				returnVars:   []string{"ae@@of..!\""},
			},
			err: errors.New("1:3: illegal character U+0040 '@' (and 1 more errors)"),
		},
		{
			name: "Reserved word",
			args: args{
				fileContent:  existingContent,
				functionName: "nonexistentFunction",
				returnVars:   []string{"range"},
			},
			err: errors.New("1:1: expected operand, found 'range'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceReturn(tt.args.fileContent, tt.args.functionName, tt.args.returnVars...)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestReplaceCode(t *testing.T) {
	var (
		newFunction     = `fmt.Println("This is the new function.")`
		existingContent = `package main

import (
    "fmt"
)

func main() {
    fmt.Println("Hello, world!")
}

func oldFunction() {
    fmt.Println("This is the old function.")
}`
	)

	type args struct {
		fileContent     string
		oldFunctionName string
		newFunction     string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Replace function implementation",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "oldFunction",
				newFunction:     newFunction,
			},
			want: `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}

func oldFunction() { fmt.Println("This is the new function.") }
`,
		},
		{
			name: "Replace main function implementation",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "main",
				newFunction:     newFunction,
			},
			want: `package main

import (
	"fmt"
)

func main() { fmt.Println("This is the new function.") }

func oldFunction() {
	fmt.Println("This is the old function.")
}
`,
		},
		{
			name: "Function not found",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "nonexistentFunction",
				newFunction:     newFunction,
			},
			err: errors.New("function nonexistentFunction not found in file content"),
		},
		{
			name: "Invalid new function",
			args: args{
				fileContent:     existingContent,
				oldFunctionName: "nonexistentFunction",
				newFunction:     "ae@@of..!\"",
			},
			err: errors.New("1:25: illegal character U+0040 '@' (and 2 more errors)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceCode(tt.args.fileContent, tt.args.oldFunctionName, tt.args.newFunction)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
