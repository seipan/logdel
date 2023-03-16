package main

// func main() {
// 	if err := run(); err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
// }

// func run() error {
// 	logdel.Analyzer.Flags = flag.NewFlagSet(logdel.Analyzer.Name, flag.ExitOnError)
// 	logdel.Analyzer.Flags.Parse(os.Args[1:])

// 	if logdel.Analyzer.Flags.NArg() < 1 {
// 		return errors.New("patterns of packages must be specified")
// 	}

// 	pkgs, err := packages.Load(logdel.Analyzer.Config, logdel.Analyzer.Flags.Args()...)
// 	if err != nil {
// 		return err
// 	}

// 	for _, pkg := range pkgs {
// 		prog, srcFuncs, err := internal.BuildSSA(pkg, logdel.Analyzer.SSABuilderMode)
// 		if err != nil {
// 			return err
// 		}

// 		pass := &internal.Pass{
// 			Package:  pkg,
// 			SSA:      prog,
// 			SrcFuncs: srcFuncs,
// 			Stdin:    os.Stdin,
// 			Stdout:   os.Stdout,
// 			Stderr:   os.Stderr,
// 		}

// 		if err := logdel.Analyzer.Run(pass); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
