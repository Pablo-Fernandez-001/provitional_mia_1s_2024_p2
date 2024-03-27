package User

import (
	"encoding/binary"
	"fmt"
	"go/Global"
	"go/Structs"
	"go/Utilities"
	"strings"
)

func Login(user string, pass string, id string) {
	fmt.Println("Starting loging...")
	fmt.Println("User: ", user)
	fmt.Println("Password: ", pass)
	fmt.Println("ID:: ", id)

	driveletter := string(id[0])

	//Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}
	//Creating the MRB
	var TempMBR Structs.MRB
	//Reading object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print Object
	Structs.PrintMBR(TempMBR)

	//fmt.Println("--------")

	//Starting index
	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				fmt.Println("A Partition was founded")
				if strings.Contains(string(TempMBR.Partitions[i].Status[:]), "i") {
					fmt.Println("Partition it's mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return
				}
				break
			}
		}
	}
	if index != -1 {
		Structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition wasn't founded")
		return
	}

	//Starting the superblock var
	var tempSuperblock Structs.Superblock
	//Read objet from the bin file
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		return
	}

	//Starting Index Inode
	indexInode := int32(1)

	//Starting the Inode var
	var crrInode Structs.Inode
	//Read object from bin file
	if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
		return
	}

	//Starting Fileblock var
	var Fileblock Structs.Fileblock
	// Read object from file
	if err := Utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(Structs.Fileblock{})))); err != nil {
		return
	}

	fmt.Println("Fileblock...")
	data := string(Fileblock.B_content[:])
	// Divide the strings into lines
	lines := strings.Split(data, "\n")

	//Iterate throught the lines
	for _, line := range lines {
		// Printing each line
		fmt.Println(line)
	}

	//Printing objects
	fmt.Println("Inode", crrInode.I_block)

	//Close bin file
	defer file.Close()

	fmt.Println("End login")
}

func Logout() {
	fmt.Println("Starting Logout...")
	if Global.User.Status {
		Global.User.ID = ""
		Global.User.Status = false
		fmt.Println("User logged out")
	} else {
		fmt.Println("User doesn't logged")
	}
	fmt.Println("Ending Logout...")
}
