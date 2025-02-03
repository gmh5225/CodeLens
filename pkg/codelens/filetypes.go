package codelens

import (
	"path/filepath"
	"strings"
)

// getFileType determines the file type based on extension
func getFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	// Programming languages
	case ".go":
		return "golang"
	case ".py":
		return "python"
	case ".js", ".jsx", ".ts", ".tsx":
		return "javascript"
	case ".java":
		return "java"
	case ".cpp", ".cc", ".cxx", ".c++", ".hpp":
		return "cpp"
	case ".c", ".h":
		return "c"
	case ".cs":
		return "csharp"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".rs":
		return "rust"
	case ".swift":
		return "swift"
	case ".kt", ".kts":
		return "kotlin"
	case ".sol":
		return "solidity"
	case ".cairo":
		return "cairo"
	case ".move":
		return "move"
	case ".vy", ".viper":
		return "vyper"
	case ".scala":
		return "scala"
	case ".dart":
		return "dart"
	case ".lua":
		return "lua"
	case ".r":
		return "r"
	case ".m":
		return "matlab"
	case ".pl", ".pm":
		return "perl"
	case ".ex", ".exs":
		return "elixir"
	case ".erl", ".hrl":
		return "erlang"
	case ".clj", ".cljs":
		return "clojure"
	case ".hs":
		return "haskell"
	case ".f", ".f90", ".f95", ".f03":
		return "fortran"
	case ".jl":
		return "julia"
	case ".v", ".vh", ".sv":
		return "verilog"
	case ".vhd", ".vhdl":
		return "vhdl"
	case ".asm", ".s":
		return "assembly"
	case ".pas":
		return "pascal"
	case ".ml", ".mli":
		return "ocaml"
	case ".fs", ".fsx", ".fsi":
		return "fsharp"
	case ".d":
		return "d"
	case ".groovy":
		return "groovy"
	case ".tcl":
		return "tcl"
	case ".vb":
		return "vb"
	case ".coffee":
		return "coffeescript"
	case ".elm":
		return "elm"
	case ".hx":
		return "haxe"
	case ".nim":
		return "nim"
	case ".cr":
		return "crystal"
	case ".zig":
		return "zig"
	case ".mo":
		return "motoko"
	case ".wasm", ".wat":
		return "webassembly"
	case ".rkt":
		return "racket"
	case ".scm":
		return "scheme"
	case ".lisp", ".lsp":
		return "lisp"
	case ".ada", ".adb", ".ads":
		return "ada"
	case ".cob", ".cbl":
		return "cobol"

	// Markup languages and configuration files
	case ".md", ".markdown":
		return "markdown"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"
	case ".proto":
		return "protobuf"
	case ".graphql", ".gql":
		return "graphql"
	case ".thrift":
		return "thrift"
	case ".avsc":
		return "avro"
	case ".capnp":
		return "capnproto"
	case ".swagger", ".openapi":
		return "openapi"
	case ".raml":
		return "raml"
	case ".wsdl":
		return "wsdl"
	case ".xsd":
		return "xsd"
	case ".dtd":
		return "dtd"

	// Scripting languages
	case ".sh", ".bash":
		return "shell"
	case ".ps1":
		return "powershell"
	case ".bat", ".cmd":
		return "batch"
	case ".fish":
		return "fish"
	case ".zsh":
		return "zsh"
	case ".ksh":
		return "ksh"
	case ".tcsh":
		return "tcsh"
	case ".csh":
		return "csh"
	case ".nu":
		return "nushell"

	// Other common file types
	case ".sql":
		return "sql"
	case ".html", ".htm":
		return "html"
	case ".css", ".scss", ".sass", ".less":
		return "css"
	case ".vue":
		return "vue"
	case ".svelte":
		return "svelte"
	case ".astro":
		return "astro"
	case ".liquid":
		return "liquid"
	case ".haml":
		return "haml"
	case ".pug", ".jade":
		return "pug"
	case ".ejs":
		return "ejs"
	case ".hbs", ".handlebars":
		return "handlebars"
	case ".twig":
		return "twig"
	case ".mustache":
		return "mustache"
	case ".erb":
		return "erb"

	case "":
		return "unknown"
	default:
		// Return extension without dot
		return ext[1:]
	}
}
