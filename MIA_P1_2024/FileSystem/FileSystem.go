package FileSystem

import (
	"encoding/binary"
	"fmt"
	"go/Structs"
	"go/Utilities"
	"os"
	"strings"
)

func Mkfs(id string, type_ string, fs_ string) {
	fmt.Println("MKFS with:")
	fmt.Println("Id: ", id)
	fmt.Println("Type: ", type_)
	fmt.Println("Fs: ", fs_)
	creationDate := Utilities.GettingDate()

	driveletter := string(id[0])

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var tempMBR Structs.MRB
	// Read bin file
	if err := Utilities.ReadObject(file, &tempMBR, 0); err != nil {
		return
	}

	//printing object
	Structs.PrintMBR(tempMBR)

	//Strating the index in -1 'cause if we add, starts in 0
	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if tempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(tempMBR.Partitions[i].Id[:]), id) {
				fmt.Println("Partition found")
				if strings.Contains(string(tempMBR.Partitions[i].Status[:]), "1") {
					fmt.Println("Partiton mounted")
					index = i
				} else {
					fmt.Println("Partition isn't Mounted")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		Structs.PrintPartition(tempMBR.Partitions[index])
	} else {
		fmt.Println("Partition not found")
		return
	}

	// Numerator = (partition_mounted.size - sizeof(Structs::Superblock))
	// Denominator base = (4 + sizeof(Structs::Inodes) + 3 * sizeof(Structs::Fileblock))
	// Temp = "2" ? 0 : sizeof(Structs::Journaling)
	// denominator = base + temp
	// n = floor( numerator / denominator )

	numerator := int32(tempMBR.Partitions[index].Size - int32(binary.Size(Structs.Superblock{})))
	denominator_base := int32(4 + int32(binary.Size((Structs.Inode{}))) + 3*int32(binary.Size(Structs.Fileblock{})))
	var temp int32 = 0
	if fs_ == "2fs" {
		temp = 0
	} else {
		temp = int32(binary.Size(Structs.Journaling{}))
	}
	denominator := denominator_base + temp
	n := int32(numerator / denominator)

	fmt.Println("N:", n)

	// Var newBR Structs MBR
	var newSuperblock Structs.Superblock
	newSuperblock.S_inodes_count = 0
	newSuperblock.S_blocks_count = 0

	newSuperblock.S_free_blocks_count = 3 * n
	newSuperblock.S_free_inodes_count = n

	copy(newSuperblock.S_mtime[:], creationDate)
	copy(newSuperblock.S_umtime[:], creationDate)
	newSuperblock.S_mnt_count = 0

	if fs_ == "2fs" {
		create_ext2(n, tempMBR.Partitions[index], newSuperblock, creationDate, file)
	} else {
		fmt.Println("EXT3")
	}

	//Close bin file
	defer file.Close()
	fmt.Println("Bye MKFS!")
}

func create_ext2(n int32, partition Structs.Partition, newSuperBlock Structs.Superblock, date string, file *os.File) {
	fmt.Println("Creating Ext2 with data:")
	fmt.Println("N:", n)
	fmt.Println("Superblock:", newSuperBlock)
	fmt.Println("Date:", date)

	newSuperBlock.S_filesystem_type = 2
	newSuperBlock.S_bm_inode_start = partition.Start + int32(binary.Size(Structs.Superblock{}))
	newSuperBlock.S_bm_block_start = newSuperBlock.S_bm_inode_start + n
	newSuperBlock.S_inode_start = newSuperBlock.S_bm_block_start + 3*n
	newSuperBlock.S_block_start = newSuperBlock.S_inode_start + n*int32(binary.Size(Structs.Inode{}))

	newSuperBlock.S_free_inodes_count -= 1
	newSuperBlock.S_free_blocks_count -= 1
	newSuperBlock.S_free_inodes_count -= 1
	newSuperBlock.S_free_blocks_count -= 1

	for i := int32(0); i < n; i++ {
		err := Utilities.WriteObject(file, byte(0), int64(newSuperBlock.S_bm_inode_start+i))
		if err != nil {
			fmt.Println("Error witing object EXT2: ", err)
		}
	}

	//Starting all blocks in -1
	var newInode Structs.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err := Utilities.WriteObject(file, newInode, int64(newSuperBlock.S_inode_start+i*int32(binary.Size(Structs.Inode{}))))
		if err != nil {
			fmt.Println("Error writing object 2 EXT2: ", err)
		}
	}

	//New File Block
	var newFileBlock Structs.Fileblock
	for i := int32(0); i < 3*n; i++ {
		err := Utilities.WriteObject(file, newFileBlock, int64(newSuperBlock.S_block_start+i*int32(binary.Size(Structs.Fileblock{}))))
		if err != nil {
			fmt.Println("Error writing object 3 EXT2: ", err)
		}
	}

	//Creating first Inode
	var Inode0 Structs.Inode
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_size = 0
	//Copy all data inside
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	copy(Inode0.I_perm[:], "0")
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0

	//Creating folder of blocks
	var Folderblock0 Structs.Folderblock //0 Block -> Folder
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[1].B_inodo = 1
	copy(Folderblock0.B_content[1].B_name[:], "users.txt")

	var Inode1 Structs.Inode //Inode 1
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_size = int32(binary.Size(Structs.Folderblock{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	copy(Inode1.I_perm[:], "0")
	copy(Inode1.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}

	Inode1.I_block[0] = 1

	data := "1,g,root\n1,U,root,root,123\n"

	//Starting the fileBlock
	var Fileblock1 Structs.Fileblock // 1th Block -> File
	copy(Fileblock1.B_content[:], data)

	// 0 Inode  -> 0 Block
	// 1 Inode  -> 1 Block
	// Create the root folder
	// Create the user.txt with the data path

	//Write superblock
	err := Utilities.WriteObject(file, newSuperBlock, int64(partition.Start))

	//write bitmap inodes
	err = Utilities.WriteObject(file, byte(1), int64(newSuperBlock.S_bm_inode_start))
	err = Utilities.WriteObject(file, byte(1), int64(newSuperBlock.S_bm_inode_start+1))

	//write bitmap blocks
	err = Utilities.WriteObject(file, byte(1), int64(newSuperBlock.S_bm_block_start))
	err = Utilities.WriteObject(file, byte(1), int64(newSuperBlock.S_bm_block_start+1))

	fmt.Println("Inode 0:", int64(newSuperBlock.S_inode_start))
	fmt.Println("Inode 1:", int64(newSuperBlock.S_inode_start+int32(binary.Size(Structs.Inode{}))))

	//write inodes
	err = Utilities.WriteObject(file, Inode0, int64(newSuperBlock.S_bm_inode_start))                                     //Inode 0
	err = Utilities.WriteObject(file, Inode1, int64(newSuperBlock.S_bm_inode_start+int32(binary.Size(Structs.Inode{})))) //Inode 1

	//write blocks
	err = Utilities.WriteObject(file, Folderblock0, int64(newSuperBlock.S_block_start))                                       //Block 0
	err = Utilities.WriteObject(file, Fileblock1, int64(newSuperBlock.S_block_start+int32(binary.Size(Structs.Fileblock{})))) //Block 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	//mkfs -type=full -id=A119

	fmt.Println("Ending creating ext2")
}
