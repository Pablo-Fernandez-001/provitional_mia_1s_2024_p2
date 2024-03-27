package main

import (
	"bufio"
	"fmt"
	Analyzer "go/Analyzer"
	"os"
	"strings"
)

func main() {
	var loop int = 0
	fmt.Println("Enter command")
	fmt.Println("Or type -h to have some help")
	for loop != 1 {
		var input string
		fmt.Print("</>:")

		//Scanning input
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()
		choose := strings.ToLower(input)
		switch choose {
		case "-h", "h", "-help", "help":
			helpList()
		case "-bye", "bye", "-close", "close", "exit", "-exit":
			fmt.Println("Bye!")
			loop = 1
		default:
			Analyzer.Analyze(input)
		}

	}
}

// List of commands to help
func helpList() {
	fmt.Println("------------------------------------------------------------------------------HELP DATA------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println("Use the command without the \"<\", \" | \" and \">\" characters, and doesn't matter how do you write the command")
	fmt.Println()
	fmt.Println(" the \"*\" means obligatory command, except the nemonic")
	fmt.Println()
	fmt.Println("<-----------------------------------------------------------------------------All command list----------------------------------------------------------------------------->")
	fmt.Println()
	fmt.Println("------------------------------------------------------------------------------Exit Program------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println("< -bye | bye | -close | close | -exit | exit > -> To exit of the program")
	fmt.Println(" \t └ example: -close")
	fmt.Println(" \t └ example: -exit")
	fmt.Println(" \t └ example: -bye")
	fmt.Println(" \t └ example: close")
	fmt.Println(" \t └ example: exit")
	fmt.Println(" \t └ example: bye")
	fmt.Println()
	fmt.Println("------------------------------------------------------------------------------INIT------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println("< -i | i | -init | init > -> To initialize all the folders, and programs")
	fmt.Println(" \t └ example: i")
	fmt.Println(" \t └ example: -i")
	fmt.Println(" \t └ example: init")
	fmt.Println(" \t └ example: -init")
	fmt.Println()
	fmt.Println("------------------------------------------------------------------------------HELP------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println("< -h | h | -help | help > -> To get help into the program")
	fmt.Println(" \t └ example: h")
	fmt.Println(" \t └ example: -h")
	fmt.Println(" \t └ example: help")
	fmt.Println(" \t └ example: -help")
	fmt.Println()
	fmt.Println("------------------------------------------------------------------------------MKDIR------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println(" mkdir < -path >* < -r > -> To make a directory with path")
	fmt.Println(" \t └ -path: there you can put the path of the folder <Obligatory command>")
	fmt.Println(" \t └ -r: if the path doesn't exists, you can created")
	fmt.Println(" \t └ example: mkdir -r -path=/home/user/docs/usac")
	fmt.Println()
	//* it's not yet at all
	fmt.Println("------------------------------------------------------------------------------MKDISK------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println(" mkdisk < -size >* < -unit > < -fit > -> To make a new disk")
	fmt.Println(" \t └ -size: to put the defined size <Obligatory command>")
	fmt.Println(" \t └ -unit: to put the name of the unit")
	fmt.Println(" \t └ -fit: to use the type of fit, bf (best fit), ff (first fit), or wf (worst fit)")
	fmt.Println(" \t └ example: mkdisk -size=3000 -unit=K -fit=BF")
	fmt.Println()
	// it's not yet at all
	fmt.Println("------------------------------------------------------------------------------FDISK------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println(" fdisk < -size >* < -driveletter >* < -name >* < -unit > < -type > < -fit >  < -delete > < -add >-> To format disk disk")
	fmt.Println(" \t └ -size: to put the defined size <Obligatory command>")
	fmt.Println(" \t └ -driveletter: to find disk by the letter, if exists <Obligatory command>")
	fmt.Println(" \t └ -name: to find the name of partition <Obligatory command>")
	fmt.Println(" \t └ -unit: to recive the letter to indicate the units B (bytes), K (kilobytes), M (megabytes)")
	fmt.Println(" \t └ -type: to use the type of partition, p (primary), e (extended), or l (logic)")
	fmt.Println(" \t └ -fit: to use the type of fit, bf (best fit), ff (first fit), or wf (worst fit)")
	fmt.Println(" \t └ -delete: to delete the partition with \"-name\" & \"-path\" (delete the partition with those paramethers), with word full")
	fmt.Println(" \t └ -add: to delete the spaces into the partitions")
	fmt.Println(" \t └ example: fdisk -size=300 -driveletter=A -name=Particion1")
	fmt.Println()
	// it's not yet at all
	fmt.Println("------------------------------------------------------------------------------MOUNT------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println(" mount < -driveletter >* < -name >*-> To mount a disk disk")
	fmt.Println(" \t └ -driveletter: to find disk by the letter, if exists <Obligatory command>")
	fmt.Println(" \t └ -name: to find the name of partition <Obligatory command>")
	fmt.Println(" \t └ example: mount -driveletter=A -name=Particion1")
	fmt.Println()
	// it's not yet at all
	fmt.Println("------------------------------------------------------------------------------MKFS------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println(" mkfs < -id >* < -type >* < -fs > -> To create the format systen")
	fmt.Println(" \t └ -id: Indicates the generated id with the command mount <Obligatory command>")
	fmt.Println(" \t └ -type: Indicate the type of format, with word full <Obligatory command>")
	fmt.Println(" \t └ -fs: Indicates the kind of files to create, 2fs (to create EXT2), 3fs (to create EXT3)")
	fmt.Println(" \t └ example: mkfs -type=full -id=A111")
	fmt.Println()
	// it's not yet at all
	fmt.Println("------------------------------------------------------------------------------LOGIN------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println(" login < -user >* < -pass >* < -id >* -> To login a user")
	fmt.Println(" \t └ -user: To specify the name of user <Obligatory command>")
	fmt.Println(" \t └ -pass: To specify the password of the user <Obligatory command>")
	fmt.Println(" \t └ -id: To indicate the mounted partition who is initializated into the sesion. <Obligatory command>")
	fmt.Println(" \t └ example: login -user=root -pass=123 -id=A111")
	fmt.Println()
}
