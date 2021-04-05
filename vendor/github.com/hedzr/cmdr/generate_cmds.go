/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

var (
	generatorCommands = &Command{
		BaseOpt: BaseOpt{
			// Name:        "generators",
			Group:       SysMgmtGroup,
			Short:       "g",
			Full:        "generate",
			Aliases:     []string{"gen"},
			Description: "generators for this app.",
			LongDescription: `
[cmdr] includes multiple generators like:

- linux man page generator
- shell completion script generator
- markdown generator
- more...

			`,
			Examples: `
$ {{.AppName}} gen sh --bash
			generate bash completion script
$ {{.AppName}} gen shell --auto
			generate shell completion script with detecting on current shell environment.
$ {{.AppName}} gen sh
			generate shell completion script with detecting on current shell environment.
$ {{.AppName}} gen man
			generate linux manual (man page)
$ {{.AppName}} gen doc
			generate document, default markdown.
$ {{.AppName}} gen doc --markdown
			generate markdown.
$ {{.AppName}} gen doc --pdf
			generate pdf.
$ {{.AppName}} gen markdown
			generate markdown.
$ {{.AppName}} gen pdf
			generate pdf.
			`,
		},
		SubCommands: []*Command{{
			BaseOpt: BaseOpt{
				Short:       "s",
				Full:        "shell",
				Aliases:     []string{"sh"},
				Description: "generate the bash/zsh auto-completion script or install it.",
				Action:      genShell,
			},
			Flags: []*Flag{
				{
					BaseOpt: BaseOpt{
						Short:       "b",
						Full:        "bash",
						Group:       "shell",
						Description: "generate auto completion script for Bash",
					},
					DefaultValue: false,
				},
				{
					BaseOpt: BaseOpt{
						Short:       "z",
						Full:        "zsh",
						Group:       "shell",
						Description: "generate auto completion script for Zsh",
					},
					DefaultValue: false,
				},
				{
					BaseOpt: BaseOpt{
						Short:       "a",
						Full:        "auto",
						Group:       "shell",
						Description: "generate auto completion script to fit for your current env.",
					},
					DefaultValue: true,
				},
				{
					BaseOpt: BaseOpt{
						Full:        "force-bash",
						Group:       "shell",
						Description: "just for --auto",
						Hidden:      true,
					},
					DefaultValue: true,
				},
			},
		}, {
			BaseOpt: BaseOpt{
				Short:       "m",
				Full:        "manual",
				Aliases:     []string{"man"},
				Description: "generate linux man page.",
				Action:      genManual,
			},
			Flags: []*Flag{
				{
					BaseOpt: BaseOpt{
						Short:       "d",
						Full:        "dir",
						Description: "the output directory",
						// Aliases:     []string{"mkd", "m"},
						// Group:       "output",
					},
					DefaultValue:            "./man1",
					DefaultValuePlaceholder: "DIR",
				},
			},
		}, {
			BaseOpt: BaseOpt{
				Short:       "d",
				Full:        "doc",
				Aliases:     []string{"markdown", "pdf", "docx", "tex"},
				Description: "generate a markdown document, or: pdf/TeX/...",
				Action:      genDoc,
			},
			Flags: []*Flag{
				{
					BaseOpt: BaseOpt{
						Short:       "d",
						Full:        "dir",
						Description: "the output directory",
						Group:       "output",
					},
					DefaultValue:            "./docs",
					DefaultValuePlaceholder: "DIR",
				},
				{
					BaseOpt: BaseOpt{
						Short:       "md",
						Full:        "markdown",
						Aliases:     []string{"mkd", "m"},
						Group:       "doc",
						Description: "generate mardown",
					},
					DefaultValue: false,
				},
				{
					BaseOpt: BaseOpt{
						Short:       "p",
						Full:        "pdf",
						Group:       "doc",
						Description: "generate pdf",
					},
					DefaultValue: false,
				},
				{
					BaseOpt: BaseOpt{
						Full:        "doc",
						Group:       "doc",
						Description: "generate word doc",
					},
					DefaultValue: false,
				},
				{
					BaseOpt: BaseOpt{
						Full:        "docx",
						Group:       "doc",
						Description: "generate word docx",
					},
					DefaultValue: false,
				},
				{
					BaseOpt: BaseOpt{
						Short:       "t",
						Full:        "tex",
						Group:       "doc",
						Description: "generate tex",
					},
					DefaultValue: true,
				},
			},
			// SubCommands: []*Command{
			// 	{
			// 		BaseOpt: BaseOpt{
			// 			Short:       "rt",
			// 			Full:        "runtime",
			// 			Description: "runtime",
			// 			Flags: []*Flag{
			// 				{
			// 					BaseOpt: BaseOpt{
			// 						Short:       "hi",
			// 						Full:        "hello",
			// 						Description: "world",
			// 					},
			// 				},
			// 				{
			// 					BaseOpt: BaseOpt{
			// 						Short:       "fi",
			// 						Full:        "fing",
			// 						Description: "finger",
			// 					},
			// 				},
			// 			},
			// 		},
			// 		SubCommands: []*Command{
			// 			{
			// 				BaseOpt: BaseOpt{
			// 					Short:       "ok",
			// 					Full:        "ready",
			// 					Description: "ok ready",
			// 					Flags: []*Flag{
			// 						{
			// 							BaseOpt: BaseOpt{
			// 								Short:       "a",
			// 								Full:        "hello",
			// 								Description: "hello world",
			// 							},
			// 						},
			// 						{
			// 							BaseOpt: BaseOpt{
			// 								Short:       "b",
			// 								Full:        "fing",
			// 								Description: "ready finger",
			// 							},
			// 						},
			// 					},
			// 				},
			// 				SubCommands: []*Command{
			// 					{
			// 						BaseOpt: BaseOpt{
			// 							Short:       "o1",
			// 							Full:        "ready1",
			// 							Description: "ok ready 1",
			// 							Flags: []*Flag{
			// 								{
			// 									BaseOpt: BaseOpt{
			// 										Short:       "a1",
			// 										Full:        "hello1",
			// 										Description: "hello world",
			// 									},
			// 								},
			// 								{
			// 									BaseOpt: BaseOpt{
			// 										Short:       "b1",
			// 										Full:        "fing1",
			// 										Description: "ready finger",
			// 									},
			// 								},
			// 							},
			// 						},
			// 					},
			// 					{
			// 						BaseOpt: BaseOpt{
			// 							Short:       "b1",
			// 							Full:        "bad1",
			// 							Description: "bad not ready 1",
			// 						},
			// 					},
			// 				},
			// 			},
			// 			{
			// 				BaseOpt: BaseOpt{
			// 					Short:       "b",
			// 					Full:        "bad",
			// 					Description: "bad not ready",
			// 				},
			// 			},
			// 		},
			// 	},
			// 	{
			// 		BaseOpt: BaseOpt{
			// 			Short:       "st",
			// 			Full:        "static",
			// 			Description: "static",
			// 			Flags: []*Flag{
			// 				{
			// 					BaseOpt: BaseOpt{
			// 						Short:       "hi",
			// 						Full:        "hello",
			// 						Description: "world",
			// 					},
			// 				},
			// 				{
			// 					BaseOpt: BaseOpt{
			// 						Short:       "fi",
			// 						Full:        "fing",
			// 						Description: "finger",
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },
		}},
	}
)
