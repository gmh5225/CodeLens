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
	case ".go.mod":
		return "gomod"
	case ".go.sum":
		return "gosum"
	case ".go.work":
		return "gowork"
	case ".py":
		return "python"
	case ".pyi":
		return "pythoni"
	case ".pyx", ".pxd":
		return "cython"
	case ".js", ".jsx", ".ts", ".tsx":
		return "javascript"
	case ".mjs", ".cjs":
		return "javascript"
	case ".dts", ".d.ts":
		return "typescript"
	case ".java":
		return "java"
	case ".cpp", ".cc", ".cxx", ".c++", ".hpp":
		return "cpp"
	case ".tpp", ".tcc", ".inl":
		return "cpp"
	case ".c", ".h":
		return "c"
	case ".cs":
		return "csharp"
	case ".csx":
		return "csharp"
	case ".rb":
		return "ruby"
	case ".rbs", ".gemspec":
		return "ruby"
	case ".rake", ".ru":
		return "ruby"
	case ".php":
		return "php"
	case ".phpt", ".phtml", ".phar":
		return "php"
	case ".blade.php":
		return "blade"
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
	case ".ink":
		return "ink"
	case ".scilla":
		return "scilla"
	case ".aes":
		return "sophia"
	case ".ride":
		return "ride"
	case ".leo":
		return "leo"
	case ".clar":
		return "clarity"

	// Markup languages and configuration files
	case ".md", ".markdown":
		return "markdown"
	case ".mdx":
		return "mdx"
	case ".rst":
		return "restructuredtext"
	case ".asciidoc", ".adoc":
		return "asciidoc"
	case ".tex", ".sty", ".cls":
		return "latex"
	case ".json":
		return "json"
	case ".jsonc":
		return "jsonc"
	case ".json5":
		return "json5"
	case ".jsonnet", ".libsonnet":
		return "jsonnet"
	case ".jsonld":
		return "jsonld"
	case ".xml":
		return "xml"
	case ".xaml":
		return "xaml"
	case ".plist":
		return "plist"
	case ".svg":
		return "svg"
	case ".yaml", ".yml":
		return "yaml"
	case ".eyaml", ".eyml":
		return "eyaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"
	case ".conf", ".cfg":
		return "config"
	case ".properties":
		return "properties"
	case ".env":
		return "dotenv"

	// Configuration and build files
	case "Dockerfile":
		return "dockerfile"
	case ".dockerignore":
		return "dockerignore"
	case "Makefile", "makefile", "GNUmakefile":
		return "makefile"
	case "CMakeLists.txt":
		return "cmake"
	case ".cmake":
		return "cmake"
	case "Vagrantfile":
		return "ruby"
	case "Jenkinsfile":
		return "groovy"
	case ".gitlab-ci.yml":
		return "yaml"
	case ".travis.yml":
		return "yaml"
	case "azure-pipelines.yml":
		return "yaml"
	case ".circleci/config.yml":
		return "yaml"
	case ".github/workflows/*.yml":
		return "yaml"
	case "package.json":
		return "json"
	case "composer.json":
		return "json"
	case "cargo.toml":
		return "toml"
	case "poetry.lock":
		return "toml"

	// Database and query languages
	case ".sql":
		return "sql"
	case ".mysql":
		return "mysql"
	case ".pgsql":
		return "postgresql"
	case ".plsql":
		return "plsql"
	case ".tsql":
		return "tsql"
	case ".hql":
		return "hive"
	case ".cypher":
		return "cypher"
	case ".sparql":
		return "sparql"

	// Web template languages
	case ".html", ".htm":
		return "html"
	case ".xhtml":
		return "xhtml"
	case ".shtml":
		return "shtml"
	case ".cshtml":
		return "cshtml"
	case ".jinja", ".j2":
		return "jinja"
	case ".njk":
		return "nunjucks"
	case ".webc":
		return "webc"
	case ".kit":
		return "kit"

	// Style languages
	case ".css", ".scss", ".sass", ".less":
		return "css"
	case ".styl":
		return "stylus"
	case ".pcss", ".postcss":
		return "postcss"
	case ".sss":
		return "sugarss"

	// Web frameworks
	case ".vue":
		return "vue"
	case ".svelte":
		return "svelte"
	case ".astro":
		return "astro"
	case ".solid":
		return "solid"
	case ".qwik":
		return "qwik"
	case ".marko":
		return "marko"
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
