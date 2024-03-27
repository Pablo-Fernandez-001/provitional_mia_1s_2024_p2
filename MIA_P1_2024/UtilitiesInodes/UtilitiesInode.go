package UtilitiesInodes

import (
	"encoding/binary"
	"fmt"
	"go/Structs"
	"go/Utilities"
	"os"
	"strings"
)

func InitSearch(path string, file *os.File, tempSuperBlock Structs.Superblock) int32 {
	fmt.Println("Start Inode Search:")
	fmt.Println("path: ", path)
	//path = "/path/new"

	//Splitting the path by "/"
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:]

	fmt.Println("Steps path: ", StepsPath, " length: ", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("Step: ", step)
	}

	var Inode0 Structs.Inode
	// Read object from bin file
	if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperBlock.S_block_start)); err != nil {
		return -1
	}
	fmt.Println("Ending Initialize Search")
	//Temp return
	return SearchInodeByPath(StepsPath, Inode0, file, tempSuperBlock)
}

func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func SearchInodeByPath(StepsPath []string, Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock) int32 {
	fmt.Println("Startint Search Inode by path:")
	index := int32(0)
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("Searched Name: ", SearchedName)

	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {
				//Direct form

				var crrFolderBlock Structs.Folderblock
				//Reading object from the bin file
				if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_content {
					//fmt.Println("Folder Founded")
					fmt.Println("Folder == Name: ", string(folder.B_name[:]), "B_inodo: ", folder.B_inodo)

					if strings.Contains(string(folder.B_name[:]), SearchedName) {

						fmt.Println("length Steps Path: ", len(StepsPath), "StepsPath:", StepsPath)
						if len(StepsPath) == 0 {
							fmt.Println("Folder founded")
							return folder.B_inodo
						} else {
							fmt.Println("Next Inode")
							var NextInode Structs.Inode
							if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Structs.Inode{})))); err != nil {
								return -1
							}

							return SearchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}
			} else {
				//Indirect Form
			}
		}
		index++
	}
	fmt.Println("Ending Search Inode by path")
	return 0
}

func GetInodeFileData(Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock) string {
	fmt.Println("Starting getting node file data")
	index := int32(0)

	//Defining content as a string
	var content string

	//Iterate over I block from I-node
	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {
				//Direct form
				var crrFileBlock Structs.Fileblock
				//Reading object from bin file
				if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
					return ""
				}
				content += string(crrFileBlock.B_content[:])
			} else {
				//Indirect case
			}
		}
		index++
	}
	fmt.Println("Ending getting file data")
	//temporal return
	return content
}
