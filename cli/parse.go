package cli

import (
	"path"
	"path/filepath"
	"strings"
)

func Parse(a *App, args []string) (appName string, argsResult []string, flagResult map[string]string, err error) {
	_, appName = path.Split(args[0])
	flagResult = make(map[string]string)
	argsResult = make([]string, 0)
	// Remove the path and extension of the executable
	appName = filepath.Base(appName)
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))

	args = args[1:]
	lenArgs := len(args)

	i := -1
	for {
		i++
		if i >= lenArgs {
			break
		}
		arg, ok := getArg(args, i)
		if !ok {
			break
		}

		isFlag, isHaveVal := checkFlag(arg)
		if isFlag {
			name, val, ok := getValueFlag(arg)
			if isHaveVal && ok {
				flagResult[name] = val
				continue
			}
			arg2, ok := getArg(args, i+1)
			if !ok {
				flagResult[name] = ""
				continue
			}
			isFlag, _ := checkFlag(arg2)
			if isFlag {
				flagResult[name] = ""
				continue
			}
			flagResult[name] = arg2
			i++
			continue
		}

		argsResult = append(argsResult, arg)
	}

	return appName, argsResult, flagResult, nil
}

func getArg(args []string, i int) (string, bool) {
	if len(args) <= i {
		return "", false
	}
	return strings.TrimSpace(args[i]), true
}

func checkFlag(name string) (bool, bool) {
	if strings.HasPrefix(name, "-") {
		if strings.HasPrefix(name, "=") {
			return true, false
		}
		return true, true
	}
	return false, false
}

func getValueFlag(name string) (string, string, bool) {
	arr1 := strings.Split(name, "-")
	arr := strings.Split(arr1[len(arr1)-1], "=")
	if len(arr) > 1 {
		return arr[0], arr[1], true
	}
	return arr[0], "", false
}
