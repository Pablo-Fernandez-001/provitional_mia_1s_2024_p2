package Utilities

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Creating bin-file
func CreateFile(name string) error {
	//if the directory exist
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error creating directory: ", err)
		return err
	}
	//Creating file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Error creating file: ", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Open file to read or write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil, err
	}
	return file, nil
}

// Writing file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Error Writing object: ", err)
		return err
	}
	return nil
}

// Reading files
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Error Reading Object: ", err)
		return err
	}
	return nil
}

// Round a number
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// Current date
func GettingDate() string {
	currentTime := time.Now()
	date := currentTime.Format("01-02-2006")
	return date

}

// Finding if the file exists
func ReturnFileName(path string) string {
	namesSlice := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	var namesGetted []string
	f, err := os.Open(path)
	if err != nil {
		//fmt.Println("Error finding directory: ", err)
		//fmt.Println("Using default letter")
		return "A"
	}
	files, err := f.ReadDir(0)
	if err != nil {
		//fmt.Println("Error reading directory: ", err)
		//fmt.Println("Using default letter")
		return "A"
	}

	for _, v := range files {
		name := strings.Split(v.Name(), ".")
		namesGetted = append(namesGetted, name[0])
	}

	for i := 0; i < len(namesSlice); i++ {
		letter := namesSlice[i]
		//fmt.Println("Letter to find:", letter)
		//fmt.Println("Array into to search:", namesGetted)
		data := findArr(namesGetted, letter)
		//fmt.Println("data getted: ", data)
		if !data {
			return strings.ToUpper(letter)
		}
		if len(namesGetted) == len(namesSlice) {
			return "LastOne"
		}
	}
	return "A"
}

func WriteTxtFile(path string, content []byte) error {
	content2, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("error while reading the file. %v", err)
	}

	//  new content

	// append the content
	content2 = append(content2, content...)

	// overwrite the content of hello.txt
	err = ioutil.WriteFile(path, content2, 0777)

	if err != nil {
		log.Fatalf("error while writing the file. %v", err)
	}
	return nil
}

// finding into array
func findArr(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Deleting file
func DeleteFile(driveletter string) (bool, string) {

	path := "./test/" + strings.ToUpper(driveletter) + ".dsk"
	if _, err := os.Stat(path); err != nil {
		e := os.Remove(path)
		if e != nil {
			return true, "Remove Error"
		}
		return true, "Error finding path, path doesn't exists"
	}
	return false, ""
}

// Splitting to get the first data from the paramethers
func SplittingOneParam(params string) string {
	Splitting := strings.Split(params, "=")
	data := Splitting[1]
	return data
}

// File exists
func FileExists(fileName string) {
	_, error := os.Stat(fileName)

	if os.IsNotExist(error) {
		CreateFile(fileName)
	}
}

// Find data
func FindData(path string, data string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error trying reading data")
	}

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)

	for Scanner.Scan() {
		word := Scanner.Text()
		if word == data {
			return true
		}
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return false

}

func InsertUsers(path string, dataToFind string) {
	fmt.Println("Data to Find: \"", dataToFind, "\"")
	chooser := FindData(path, dataToFind)
	fmt.Println("Chooser: ", chooser)
	if !chooser {
		data := dataToFind + "\n"
		if err := WriteTxtFile(path, []byte(data)); err != nil {
			log.Fatal(err)
		}

	}
}

func InsertGroups(path string, dataToFind string) {

	chooser := FindData(path, dataToFind)
	if dataToFind != "1,G,root" {
		numb := returningGroupandUsersQty(path, "g")

		if !chooser {
			data := strconv.Itoa(numb) + "," + dataToFind + "\n"
			if err := WriteTxtFile(path, []byte(data)); err != nil {
				log.Fatal(err)
			}
		}
	}
	if !chooser {
		data := dataToFind + "\n"
		if err := WriteTxtFile(path, []byte(data)); err != nil {
			log.Fatal(err)
		}
	}
}

func returningGroupandUsersQty(path string, dataToCount string) int {
	counter := 0
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error trying reading data")
	}

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)
	if strings.ToLower(dataToCount) == "g" {

		for Scanner.Scan() {
			word := Scanner.Text()
			Splitting := strings.Split(word, ",G,")
			if len(Splitting) != 1 {
				counter++
				fmt.Println("Group Splitted: ", Splitting)
				fmt.Println("QTY: ", counter)
			}
		}
		if err := Scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	if strings.ToLower(dataToCount) == "u" {
		for Scanner.Scan() {
			word := Scanner.Text()
			Splitting := strings.Split(word, ",U,")
			if len(Splitting) != 1 {
				counter++
				fmt.Println("User Splitted: ", Splitting)
				fmt.Println("QTY: ", counter)
			}
		}
		if err := Scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return counter
}
