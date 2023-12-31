package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Student struct {
	StudentID   int
	FirstName   string
	LastName    string
	Certificate string
	Notes       string
}

func readCSV(filePath string) (map[int][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	resultMap := make(map[int][]string)
	//skip header
	for _, record := range records[1:] {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		resultMap[id] = record[1:]
	}
	return resultMap, nil
}

func joinCSVs(file1, file2 string) (map[int]Student, error) {
	data1, err := readCSV(file1)
	if err != nil {
		return nil, err
	}

	data2, err := readCSV(file2)
	if err != nil {
		return nil, err
	}

	result := make(map[int]Student)
	for id, values := range data1 {
		if notes, ok := data2[id]; ok {
			result[id] = Student{
				StudentID:   id,
				FirstName:   values[0],
				LastName:    values[1],
				Certificate: values[2],
				Notes:       strings.Join(notes, ", "),
			}
		}
	}
	return result, nil
}

func saveToCSV(students map[int]Student, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Headers, not necessary
	header := []string{"StudentID", "FirstName", "LastName", "Certificate", "Notes"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Sort ID, also not necessary
	var ids []int
	for id := range students {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for _, id := range ids {
		student := students[id]
		record := []string{
			strconv.Itoa(student.StudentID),
			student.FirstName,
			student.LastName,
			student.Certificate,
			student.Notes,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// update field with safety check
func updateStudentField(students map[int]Student, studentID int, field, newValue string) error {
	// Check if the studentID exists in the map
	student, exists := students[studentID]
	if !exists {
		return fmt.Errorf("student does not exist %d", studentID)
	}

	switch field {
	case "FirstName":
		student.FirstName = newValue
	case "LastName":
		student.LastName = newValue
	case "Certificate":
		student.Certificate = newValue
	case "Notes":
		student.Notes = newValue
	default:
		return fmt.Errorf("invalid field: %s", field)
	}

	students[studentID] = student
	return nil
}

func main() {
	startTime := time.Now()
	var mStart runtime.MemStats
	runtime.ReadMemStats(&mStart)

	students, err := joinCSVs("student_list_big.csv", "student_notes_big.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	elapsed := time.Since(startTime)
	var mEnd runtime.MemStats
	runtime.ReadMemStats(&mEnd)

	// Calculate memory usage
	memoryUsed := mEnd.Alloc - mStart.Alloc // bytes

	//print all, no order
	// fmt.Println("print all")
	// for id, student := range students {
	// 	fmt.Printf("StudentID: %d, Name: %s, Surname: %s, Certificate: %s, Notes: %s\n", id, student.FirstName, student.LastName, student.Certificate, student.Notes)
	// }

	//print all, with order
	fmt.Println("\nprint all, with order")
	for id := 1; id <= len(students); id++ {
		fmt.Printf("StudentID: %d, Name: %s, Surname: %s, Certificate: %s, Notes: %s\n", id, students[id].FirstName, students[id].LastName, students[id].Certificate, students[id].Notes)
	}

	//update field
	updateStudentField(students, 1, "FirstName", "John")
	updateStudentField(students, 1, "Notes", "THIS IS JOHNY BEEEEEE")

	//print fews
	fmt.Println("\nprint fews")
	fmt.Println(students[1], students[100], students[1000])

	//print specific field
	fmt.Println("print specific field")
	fmt.Println(students[1].FirstName, students[1].LastName)

	//mem used
	fmt.Printf("\nResource usages: \n")
	fmt.Printf("Time taken: %s\n", elapsed)
	fmt.Printf("Memory used: %.2f MB\n", float64(memoryUsed)/1024/1024)

	//can also write to existing CSV, it will replace all data
	saveToCSV(students, "result.csv")
}
