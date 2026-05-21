package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime/pprof"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Program version
var Version = "unknown"

// Configuration options.
type ConfigurationOptions struct {
	cpuprofile          string
	debug               bool
	fieldSeparator      string
	files               []string
	bashHistory         bool
	forceOutput         bool
	ignoreCase          bool
	ignoreLeadingBlanks bool
	key                 KeyType
	memprofile          string
	multiline           int
	output              string
	printVersion        bool
	stableSort          bool
	uniq                UniqMode
}

func NewConfigurationOptions() ConfigurationOptions {
	return ConfigurationOptions{
		fieldSeparator: " ",
		files:          []string{},
		multiline:      1,
	}
}

var options = NewConfigurationOptions()

// initializeLogging initializes the logger.
func initializeLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func run(cmd *cobra.Command, args []string) {
	if options.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if options.printVersion {
		fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\n", Version)
		return
	}

	if options.cpuprofile != "" {
		f, err := os.Create(options.cpuprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if len(args) == 0 {
		options.files = append(options.files, "-")
	} else {
		options.files = append(options.files, args...)
	}

	if options.memprofile != "" {
		if options.bashHistory {
			RunBashHistoryMode(options.files)
		} else {
			var concatenatedContents ContentType = LoadInputFiles(options.files, options.key)
			var sortedContents ContentType = SortContents(concatenatedContents)
			UniqContents(sortedContents)
		}
		f, err := os.Create(options.memprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	var fd *os.File
	if options.output != "" {
		_, err := os.Stat(options.output)
		if err == nil && !options.forceOutput {
			log.Fatal().Msgf("output file %s already exists", options.output)
		}
		fd, err = os.Create(options.output)
		if err != nil {
			log.Fatal().Msgf("cannot open output file %s: %s", options.output, err.Error())
		}
		defer fd.Close()
	} else {
		fd = os.Stdout
	}

	// Use buffered writer for better performance
	writer := bufio.NewWriterSize(fd, 256*1024) // 256KB buffer

	if options.bashHistory {
		contents := RunBashHistoryMode(options.files)
		for _, record := range contents {
			for _, line := range record.Lines {
				writer.WriteString(line)
				writer.WriteByte('\n')
			}
		}
		writer.Flush()
		return
	}

	var concatenatedContents ContentType = LoadInputFiles(options.files, options.key)
	var sortedContents ContentType = SortContents(concatenatedContents)
	var uniqContents ContentType = UniqContents(sortedContents)

	for _, multiline := range uniqContents {
		for _, line := range multiline.Lines {
			writer.WriteString(line)
			writer.WriteByte('\n')
		}
	}
	writer.Flush()
}

func buildRootCmd() *cobra.Command {
	binaryName := path.Clean(os.Args[0])

	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [OPTION]... [FILE]...", binaryName),
		Short: "Sort concatenation of all FILE(s) to standard output",
		Long: `Write sorted concatenation of all FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

yet-another-sort interprets the input in 'multilines' which are line groupings
of one or multiple lines. This can be useful if the input contains header lines,
e.g. a timestamp. Examples are the bash history file, which contains lines such as

  #1692110031
  ls
  #1692110033
  yet-another-sort --multiline 2 --key 2,

Note that when the --key option is used, multiple adjacent field-separators are
treated as one separator, i.e. empty fields are ignored.

Key Specification:

  The keys are specified with F[,[F]] where F is a field (defined by the field
  separator) in the multiline used to sort.

  F      Use only field F for multiline comparisons
  F,     Use the remainder of the multiline starting with field F for comparison
  F1,F2  Use all fields between [F1,F2] for comparison`,
		Args:                  cobra.ArbitraryArgs,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			run(cmd, args)
			return nil
		},
		// Don't print usage on runtime errors
		SilenceUsage: true,
	}

	flags := rootCmd.Flags()
	flags.BoolVar(&options.bashHistory, "bash-history", false, "Process input as a bash history file: deduplicate command blocks keeping the latest timestamp, then sort by timestamp ascending")
	flags.BoolVar(&options.debug, "debug", false, "Print debugging output")
	flags.StringVar(&options.fieldSeparator, "field-separator", " ", "Use this field separator")
	flags.BoolVar(&options.forceOutput, "force", false, "Overwrite output file if it exists")
	flags.BoolVar(&options.ignoreCase, "ignore-case", false, "Ignore case for comparisons")
	flags.BoolVar(&options.ignoreLeadingBlanks, "ignore-leading-blanks", false, "Ignore leading whitespace")
	flags.BoolVar(&options.ignoreLeadingBlanks, "ignore-leading-whitespace", false, "Ignore leading whitespace, same as --ignore-leading-blanks")
	flags.VarP(&options.key, "key", "k", "Sort lines based on a particular field, see 'Key Specification' for details")
	flags.IntVar(&options.multiline, "multiline", 1, "Combine multiple lines before sorting")
	flags.StringVarP(&options.output, "output", "o", "", "Write output to file instead of standard out")
	flags.Var(&options.uniq, "uniq", `Uniq'ify the sorted multilines; keep ["first", "last"] of multiple identical lines`)
	flags.BoolVar(&options.printVersion, "version", false, "Print version and exit")
	flags.BoolVar(&options.stableSort, "stable-sort", false, "Use stable sort; preserves original input order among entries with equal keys, which is required for deterministic --uniq first/last behavior (slower)")

	flags.StringVar(&options.cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flags.StringVar(&options.memprofile, "memprofile", "", "Write memory profile to file")

	// Mark completion hints for --uniq values
	rootCmd.RegisterFlagCompletionFunc("uniq", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"first", "last"}, cobra.ShellCompDirectiveNoFileComp
	})
	// File-path completions for output/profile flags
	rootCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	})
	rootCmd.RegisterFlagCompletionFunc("cpuprofile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	})
	rootCmd.RegisterFlagCompletionFunc("memprofile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	})

	// Add the built-in completion subcommand
	rootCmd.AddCommand(buildCompletionCmd(binaryName))

	return rootCmd
}

func buildCompletionCmd(binaryName string) *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Long: fmt.Sprintf(`Generate a shell completion script for yet-another-sort.

Bash:
  source <(%s completion bash)

  # To load completions for each session, execute once:
  # Linux:
  %s completion bash > /etc/bash_completion.d/yet-another-sort
  # macOS:
  %s completion bash > $(brew --prefix)/etc/bash_completion.d/yet-another-sort

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  %s completion zsh > "${fpath[1]}/_yet-another-sort"

Fish:
  %s completion fish | source

  # To load completions for each session, execute once:
  %s completion fish > ~/.config/fish/completions/yet-another-sort.fish

PowerShell:
  %s completion powershell | Out-String | Invoke-Expression

  # To load completions for each session, add the output of the above command
  # to your PowerShell profile.
`, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
}

func main() {
	initializeLogging()
	rootCmd := buildRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
