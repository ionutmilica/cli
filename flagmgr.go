package cli

type FlagMgr struct {
	arguments []*Flag
	options   map[string]*Flag
}

func newFlagMgr(flags []*Flag) *FlagMgr {
	mgr := &FlagMgr{make([]*Flag, 0), make(map[string]*Flag, 0)}

	for _, flag := range flags {
		if flag.kind == argumentFlag {
			mgr.arguments = append(mgr.arguments, flag)
		} else {
			mgr.options[flag.name] = flag
		}
	}

	return mgr
}

func (mgr *FlagMgr) requiredArgs() []string {
	flags := []string{}

	for _, arg := range mgr.arguments {
		if arg.isRequiredArgument() {
			flags = append(flags, arg.name)
		}
	}

	return flags
}

func (mgr *FlagMgr) hasOption(opt string) bool {
	if _, ok := mgr.options[opt]; !ok {
		return false
	}
	return true
}

func (mgr *FlagMgr) hasArgument(i int) bool {
	return hasIndex(len(mgr.arguments), i)
}

func (mgr *FlagMgr) argument(i int) *Flag {
	return mgr.arguments[i]
}

func (mgr *FlagMgr) option(opt string) *Flag {
	return mgr.options[opt]
}

func hasIndex(size int, i int) bool {
	if size == 0 {
		return false
	}
	if i > -1 && i < size {
		return true
	}

	return false
}
