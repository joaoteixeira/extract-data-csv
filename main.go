package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mohae/struct2csv"
)

type Data struct {
	DT_NOTIFIC string `csv:"NOTIFICACAO"`
	SG_UF_NOT  string `csv:"ESTADO"`
	CS_RACA    string `csv:"RACA"`
}

func main() {

	file, err := os.Open("database/data.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	uf := ""

	for {
		printMessage("Informe a UF que deseja extrair ou deixe em branco todos: ")
		fmt.Scanf("%s", &uf)

		if uf == "" {
			uf = "ALL"
			break
		}

		if len(uf) == 2 {
			uf = strings.ToUpper(uf)
			break
		}
		printMessageLn("[WARNING] Informe corretamente a UF. Tente novamente!")
	}

	printMessageLn("Extraindo dados do arquivo.")
	//lines, err := csv.NewReader(file).ReadAll()
	reader := csv.NewReader(file)

	//reader, err := csv.NewReader(bufio.NewReader(file)) //lê arquivo
	reader.Comma = ';'
	lines, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	var dataCSV []Data

	for i, line := range lines {

		if i > 0 {

			if line[4] == uf || uf == "ALL" {

				var raca string
				switch line[17] {
				case "1":
					raca = "PRETA"
				default:
					raca = "OUTRAS"
				}

				data := Data{
					DT_NOTIFIC: line[0],
					SG_UF_NOT:  line[4],
					CS_RACA:    raca,
				}

				dataCSV = append(dataCSV, data)
			}
		}

		//fmt.Println(data.Column_2)
	}

	if len(dataCSV) == 0 {
		printMessageLn("Não existem dados para a UF: " + uf)
		exit()
	}

	printMessageLn("Foram extraidas " + strconv.Itoa(len(dataCSV)) + " linhas.")
	outputData(uf, dataCSV)

}

func outputData(uf string, data []Data) {

	now := time.Now()

	if err := ensureDir("output"); err != nil {
		printMessageLn("[Error] Directory creation failed with error: " + err.Error())
		exit()
	}

	fileName := fmt.Sprintf("output/output_%s_%s.csv", uf, now.Format("20060102_150405"))

	csvFile, err := os.Create(fileName)

	if err != nil {
		printMessageLn("[Error] Failed creating file: " + err.Error())
		exit()
	}

	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = ';'

	enc := struct2csv.New()

	colhdrs, err := enc.GetColNames(data[0])

	if err != nil {
		fmt.Println(err)
	}

	_ = csvwriter.Write(colhdrs)

	printMessageLn("Gerando arquivo CSV...")
	printMessage("Adicionando linhas...")

	counter := 0
	counterErr := 0
	for _, v := range data {
		row, err := enc.GetRow(v)

		if err != nil {
			printMessageLn(err.Error())
			counterErr++
		} else {
			_ = csvwriter.Write(row)
			counter++
		}
	}

	fmt.Print("\n")

	csvwriter.Flush()
	csvFile.Close()
	printMessageLn("Registros incluidos ao CSV: " + strconv.Itoa(counter))
	printMessageLn("Registros Não incluidos ao CSV: " + strconv.Itoa(counterErr))
	printMessageLn("Arquivo CSV gerado com sucesso.")
	printMessageLn("Script finalizado.")

}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, 0755)

	if err == nil {
		return nil
	}

	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)

		if err != nil {
			return err
		}

		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}

		return nil
	}
	return err
}

func printMessage(m string) {
	now := time.Now()
	fmt.Printf("[%s] %s", now.Format("2006-01-02 15:04:05"), m)
}

func printMessageLn(m string) {
	now := time.Now()
	fmt.Printf("[%s] %s\n", now.Format("2006-01-02 15:04:05"), m)
}

func exit() {
	printMessageLn("Script finalizado.")
	os.Exit(0)
}
