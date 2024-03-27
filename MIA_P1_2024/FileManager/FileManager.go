package FileManager

import (
	"encoding/binary"
	"fmt"
	"go/Global"
	"go/Structs"
	"go/Utilities"
	"go/UtilitiesInodes"
	"os"
	"strings"
)

func Mkusr(user string, pass string, grp string) {
	fmt.Println("Starting Making User")
	fmt.Println("User: ", user)
	fmt.Println("Password: ", pass)
	fmt.Println("Groups: ", grp)

	if !Global.User.Status {
		fmt.Println("User already logged")
		return
	}

	driveletter := string(Global.User.ID[0])

	//Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	fmt.Println("Filepath:", filepath)
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMRB Structs.MRB
	//Reading objects bin file
	if err := Utilities.ReadObject(file, &TempMRB, 0); err != nil {
		return
	}

	ID := string(Global.User.ID[:])
	index := int(Global.User.ID[1])

	fmt.Println("ID:", ID)
	fmt.Println("Index:", index)

	var tempSuperblock Structs.Superblock
	//Reading objects from bin file
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMRB.Partitions[index].Start)); err != nil {
		return
	}

	// Initialize the search  initSearch -> /users.txt -> doesn't returns de I-node
	// initSearch -> 1
	indexInode := UtilitiesInodes.InitSearch("/users.txt", file, tempSuperblock)

	var crrInode Structs.Inode
	// Read object from bin file
	if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
		return
	}

	// read file data
	data := UtilitiesInodes.GetInodeFileData(crrInode, file, tempSuperblock)

	fmt.Println("data:", data)
	// UID , Tipo , Grupo , Nombre , Contraseña

	// read number of users

	// UID -> read number of users + 1

	// write new user -> validate if data > 64 -> create new block

	fmt.Println("indexInode:", indexInode)

	fmt.Println("Ending Make User")
}

func MkGrp(path string, groupName string) {
	//making in partition
	fmt.Println("Data: ", Global.User.ID)
	fmt.Println("Data: ", Global.User.Status)
	//Easy part
	Utilities.FileExists(path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer file.Close()
	dataToFind := "1,G,root"
	//Utilities.InsertGroups(path, dataToFind)
	dataToFind = "U,root\t,root\t,123"
	Utilities.InsertUsers(path, dataToFind)
	//dataToFind = "G," + groupName
	//Utilities.InsertGroups(path, dataToFind)
}
