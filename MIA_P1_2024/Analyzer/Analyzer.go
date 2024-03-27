package analyzer

import (
	"flag"
	"fmt"
	"go/DiskManagement"
	"go/FileManager"
	"go/FileSystem"
	"go/User"
	"go/Utilities"
	"os"
	"regexp"
	"strings"

	"github.com/schollz/progressbar/v3"
)

//Example of commands

//mkdisk -size=3000 -unit=K -fit=BF
//fdisk -size=300 -driveletter=A -name=Particion1
//mount -driveletter=A -name=Particion1
//unmount -id=A111
//mkfs -type=full -id=A111
//print -driveletter=A
//clsdsk -driveletter=A
//need to do ----
//rmdisk ✓
//rmdisk -driveletter=A
//login ✓
//login -user -pass -id
//login -user=root -pass=123 -id=A111
//logout ✓
//logout <--- code
//mkgrp
//mkgrp -name=usuarios
//rmgrp
//mkusr ✓
//mkusr -user=user1 -pass=usuario -grp=usuarios2
//rmusr
//mkfile
//cat
//remove
//edit
//rename
//mkdir
//copy
//move
//find
//chown
//chgrp
//chmod
//pause
//recovery
//loss
//execute
//rep
//read -> to read a file

// regular expretions
var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

// Main function
func Analyze(input string) {
	fmt.Println("Data to analyze: ", input)
	command, params := getCommandAndParams(input)
	fmt.Println("Commands: ->", command, " and params: ->", params)
	choosingCommand(command, params)
}

// Separate params
func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

// To choose the commad to use
func choosingCommand(command string, params string) {
	comm := strings.ToLower(command)
	switch comm {
	case "mkdisk", "mkdsk":
		fmt.Println("make disk")
		fn_mkdisk(params)
	case "fdisk", "fdsk":
		fmt.Println("Format disk")
		fn_fdisk(params)
	case "mount":
		fmt.Println("Mount disk")
		fn_mount(params)
	case "unmount":
		fmt.Println("Unmount disk")
		fn_unmount(params)
	case "mkfs":
		fmt.Println("Making file system")
		fn_mkfs(params)
	case "login":
		fmt.Println("Login")
		fn_login(params)
	case "logout":
		fmt.Println("Logout")
		fn_logout()
	case "print":
		fmt.Println("Print")
		fn_printing(params)
	case "clsdsk", "clsdisk":
		fmt.Println("Clear Disk")
		fn_clsdsk(params)
	case "mkusr":
		fmt.Println("Making User")
		fn_mkusr(params)
	case "rmdsk", "rmdisk":
		fmt.Println("RM Disk")
		fn_rmdsk(params)
	case "mkgrp":
		fmt.Println("making group")
		fn_mkgrp(params)
	default:
		fmt.Println("Command doesn't found, try again")
	}
}

func fn_mkgrp(params string) {
	groupName := strings.Split(params, "=")
	group := groupName[1]
	fmt.Println("Group: ", group)
	path := "./test/users.text"
	FileManager.MkGrp(path, group)
}

func fn_rmdsk(params string) {

	driveletter := Utilities.SplittingOneParam(params)

	bar := progressbar.Default(100, "Deletting file...")
	state, types := Utilities.DeleteFile(driveletter)
	if !state {
		bar.Add(100)
		fmt.Println("File deleted successfully", types)
	} else {
		fmt.Println("Error")
	}
}

func fn_mkusr(input string) {
	// Define flags
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "User")
	pass := fs.String("pass", "", "Password")
	grp := fs.String("grp", "", "Group")

	//Parsing flags
	fs.Parse(os.Args[1:])

	//Finding the flags into the input
	matches := re.FindAllStringSubmatch(input, -1)

	fmt.Println(user, pass, matches, grp)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "user", "pass", "grp":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag doesn't foundes")
		}
	}

	//Calling function
	FileManager.Mkusr(*user, *pass, *grp)
}

func fn_logout() {
	User.Logout()
}

func fn_printing(input string) {
	driveLetter := Utilities.SplittingOneParam(input)
	DiskManagement.Print(driveLetter)
}

func fn_clsdsk(input string) {
	driveletter := Utilities.SplittingOneParam(input)
	DiskManagement.Clsdsk(driveletter)
}

// Function Login
func fn_login(input string) {
	//Defining flags
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "User")
	pass := fs.String("pass", "", "Password")
	id := fs.String("id", "", "ID")

	// Parsing the flags
	fs.Parse(os.Args[1:])

	//finding the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	//Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag wasn't founded")
		}
	}

	//Calling function
	User.Login(*user, *pass, *id)
}

// Function Making file system
func fn_mkfs(input string) {
	//Define flags
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "", "Type")
	fs_ := fs.String("fs", "2fs", "Fs")

	//Parse the flags
	fs.Parse(os.Args[1:])

	// find flags into the input
	matches := re.FindAllStringSubmatch(input, -1)

	// find the flags in the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error, Flag doesn't founded")
		}
	}

	// Call the function
	FileSystem.Mkfs(*id, *type_, *fs_)
}

// Function mount disk
func fn_mount(input string) {
	//Define flags
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	driveLetter := fs.String("driveletter", "", "Letter")
	name := fs.String("name", "", "Name")

	// Parsing all flags
	fs.Parse(os.Args[1:])

	//find the flags in input
	matches := re.FindAllStringSubmatch(input, -1)

	//input process
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error, Flag doesn't founded")
		}
	}

	//Calling function
	DiskManagement.Mount(*driveLetter, *name)
}

func fn_unmount(input string) {
	fmt.Println("Unmounting...")
	fmt.Println("Input: ", input)
	splitting := strings.Split(input, "=")
	nemonic := splitting[0]
	id := splitting[1]

	fmt.Println("Nemonic: \"", nemonic, "\"")
	fmt.Println("ID: \"", id, "\"")

	splitted := strings.Split(id, "")

	fmt.Println("Splitted: ", splitted)

	//disk  partition   identifier
	// A    1/X         11
	driveLetter := splitted[0]
	partition := splitted[1]
	identifier := splitted[2] + splitted[3]
	DiskManagement.Unmount(driveLetter, partition, identifier)
}

// Function format disk
func fn_fdisk(input string) {
	//Define flags

	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Size")
	driveletter := fs.String("driveletter", "", "Letter")
	name := fs.String("name", "", "Name")
	unit := fs.String("unit", "m", "Unit")
	type_ := fs.String("type", "p", "Type")
	fit := fs.String("fit", "f", "Fit")

	//Parsing flags
	fs.Parse(os.Args[1:])

	//finding flags to imput
	matches := re.FindAllStringSubmatch(input, -1)

	//Process to input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "driveletter", "name", "type":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error, flag doesn't founded")
		}
	}
	//Calling function
	DiskManagement.Fdisk(*size, *driveletter, *name, *unit, *type_, *fit)

}

// Function making disk
func fn_mkdisk(params string) {
	//Define Flags
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Size")
	fit := fs.String("fit", "f", "Fit")
	unit := fs.String("unit", "m", "Unit")

	//Parsing Flags
	fs.Parse(os.Args[1:])

	//Finding flags into the input
	matches := re.FindAllStringSubmatch(params, -1)

	//Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error, flag doesn't founded")
		}
	}
	//Calling function
	DiskManagement.Mkdisk(*size, *fit, *unit)
}
