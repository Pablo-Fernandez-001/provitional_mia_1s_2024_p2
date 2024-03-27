package DiskManagement

import (
	"encoding/binary"
	"fmt"
	"go/Structs"
	"go/Utilities"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func Unmount(driveLetter string, partition string, id string) {
	fmt.Println("Drive letter: ", driveLetter)
	fmt.Println("Partition: ", partition)
	fmt.Println("ID: ", id)
	//Who can I need to Do
	//Open disk/binary file
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}
	//Read disk/binary file
	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	fmt.Println("Printing the mbr")
	Structs.PrintMBR(TempMBR)
	//Finding the partition (data) inside disk
	//What can I have
	parId := driveLetter + partition + id
	fmt.Println("Print partiton id: ", parId)
	start := 0
	end := 0
	//Finding into partitions
	//var newMRB Structs.MRB
	for i := 0; i < 4; i++ {
		data := TempMBR.Partitions[i]

		temID := string(data.Id[:])
		if temID == parId {
			//OVERWRITE DATA
			//Moving from the begining of the file -> to -> the begining of the partition
			//file.Seek(int64(temStart), 0)
			var newMRB Structs.MRB
			partition := newMRB.Partitions[i]
			partition.Status = [1]byte{}
			partition.Type = [1]byte{}
			partition.Fit = [1]byte{}
			partition.Start = 0
			partition.Size = 0
			partition.Name = [16]byte{}
			partition.Correlative = 0
			partition.Id = [4]byte{}
			start = int(data.Start)
			if i != 3 {
				end = int(TempMBR.Partitions[i+1].Start) - 1
			} else {
				end = int(TempMBR.Partitions[i].Size) + int(TempMBR.Partitions[i].Start) - 1
			}
			//Writing data

			copy(TempMBR.Partitions[i].Status[:], partition.Status[:])
			copy(TempMBR.Partitions[i].Type[:], partition.Type[:])
			copy(TempMBR.Partitions[i].Fit[:], partition.Fit[:])
			TempMBR.Partitions[i].Start = partition.Start
			TempMBR.Partitions[i].Size = partition.Size
			copy(TempMBR.Partitions[i].Name[:], partition.Name[:])
			TempMBR.Partitions[i].Correlative = partition.Correlative
			copy(TempMBR.Partitions[i].Id[:], partition.Id[:])

			// Overwrite the MBR

			i = 4

		}

	}

	fmt.Println("Before")
	Structs.PrintMBR(TempMBR)
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		return
	}
	fmt.Println("After")
	Structs.PrintMBR(TempMBR)
	fmt.Println("Starting data: ", start)
	fmt.Println("Ending data: ", end)
	defer file.Close()
	fmt.Println("Unmounting succesfuly...")
}

func Print(driveletter string) {
	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)
	defer file.Close()

}

func Clsdsk(driveletter string) {
	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	for i := 0; i < 1; i++ {
		Utilities.WriteObject(file, TempMBR, int64(i))
	}

	// Print object
	Structs.PrintMBR(TempMBR)
	defer file.Close()

}

func Mount(driveletter string, name string) {
	fmt.Println("Mounting...")
	fmt.Println("Driveletter:", driveletter)
	fmt.Println("Name:", name)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)

	var index int = -1
	var count = 0
	// Iterate over the partitions
	bar := progressbar.Default(100, "Creating partitions...")
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			count++
			if strings.Contains(string(TempMBR.Partitions[i].Name[:]), name) {
				index = i
				break
			}
		}
	}
	bar.Add(100)

	if index != -1 {
		fmt.Println("Partition found")
		Structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition not found")
		return
	}

	// id = DriveLetter + Correlative + 11

	id := strings.ToUpper(driveletter) + strconv.Itoa(count) + "11"

	copy(TempMBR.Partitions[index].Status[:], "1")
	copy(TempMBR.Partitions[index].Id[:], id)

	// Overwrite the MBR
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR2)

	// Close bin file
	defer file.Close()

	fmt.Println("End mounting")
}

func Fdisk(size int, driveletter string, name string, unit string, type_ string, fit string) {
	fmt.Println("Start formatting disk: ")
	fmt.Println("Size:", size)
	fmt.Println("Driveletter:", driveletter)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)

	// validate fit equals to b/w/f
	if fit != "b" && fit != "w" && fit != "f" {
		fmt.Println("Error: Fit must be b, w or f")
		return
	}

	// validate size > 0
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// validate unit equals to b/k/m
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be b, k or m")
		return
	}

	// validate type equals to p/e/l
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Type must be p, e or l")
		return
	}

	// Set the size in bytes
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)

	var count = 0
	var gap = int32(0)
	// Iterate over the partitions
	bar := progressbar.Default(100, "Creating partitions...")
	for i := 0; i < 4; i++ {
		bar.Add(25)
		if TempMBR.Partitions[i].Size != 0 {
			count++
			gap = TempMBR.Partitions[i].Start + TempMBR.Partitions[i].Size
		}
	}

	//Formatting
	bar = progressbar.Default(100, "Formatting")
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size == 0 {
			TempMBR.Partitions[i].Size = int32(size)

			if count == 0 {
				TempMBR.Partitions[i].Start = int32(binary.Size(TempMBR))
			} else {
				TempMBR.Partitions[i].Start = gap
			}

			copy(TempMBR.Partitions[i].Name[:], name)
			copy(TempMBR.Partitions[i].Fit[:], fit)
			copy(TempMBR.Partitions[i].Status[:], "0")
			copy(TempMBR.Partitions[i].Type[:], type_)
			TempMBR.Partitions[i].Correlative = int32(count + 1)
			break
		}
	}
	bar.Add(100)

	// Overwrite the MBR
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR2)

	// Close bin file
	defer file.Close()

	fmt.Println("End formatting disk")
}

func Mkdisk(size int, fit string, unit string) {
	fmt.Println("Starting making disk")
	fmt.Println("Size: ", size)
	fmt.Println("Fit: ", fit)
	fmt.Println("Unit: ", unit)

	//Validate fit equials to best / worst / fit format
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit must to be b, w, or f")
		fmt.Println("best / worst / f")
		return
	}

	// validate size > 0
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	//Validate if the unit it's k/m (kilo/mega)
	if unit != "k" && unit != "m" {
		fmt.Println("Error unit doesn't format correctly (k/m)")
		fmt.Println("kilo/mega")
		return
	}

	driveletter := Utilities.ReturnFileName("./test/")
	//Creating file
	err := Utilities.CreateFile("./test/" + driveletter + ".dsk")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Set the size in bytes
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// Open bin file
	file, err := Utilities.OpenFile("./test/" + driveletter + ".dsk")
	if err != nil {
		return
	}
	// Write 0 binary data to the file

	// create array of byte(0)
	//Progress bar
	bar := progressbar.Default(100, "Creating Disk...")
	var size_ float64 = float64(size)
	var j float64 = 0
	var toSave int = 0
	var k float64 = 0
	for i := 0; i < size; i++ {
		err := Utilities.WriteObject(file, byte(0), int64(i))
		j = float64(i) + 1
		k = (j / size_) * 100.0
		rounded := Utilities.RoundFloat(k, 0)
		if int(rounded) != toSave {
			toSave = int(rounded)
			//increasing the progress bar
			bar.Add64(1)
		}
		if err != nil {
			fmt.Println("Error: ", err)
		}

	}

	// Create a new instance of MRB
	var newMRB Structs.MRB
	newMRB.MbrSize = int32(size)
	newMRB.Signature = 10 // Randomize
	copy(newMRB.Fit[:], fit)
	//Getting the current creation date
	creationDate := Utilities.GettingDate()
	copy(newMRB.CreationDate[:], creationDate)

	// Write object in bin file
	if err := Utilities.WriteObject(file, newMRB, 0); err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)

	// Close bin file
	defer file.Close()

	fmt.Println("Disk created succesfully")
}
