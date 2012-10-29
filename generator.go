package main

import (
	"encoding/csv"
	"fmt"
	"github.com/Claymore/go-config/config"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	configFile, err := os.Open("prefectures.cfg")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer configFile.Close()
	reader := config.NewReader(configFile)
	sections, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	codesFile, err := os.Open("codes.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer codesFile.Close()
	codesReader := csv.NewReader(codesFile)
	codesReader.Comma = '\t'
	codes := make(map[string]string)
	prefectures := make(map[string]string)
	for {
		record, err := codesReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		code := fmt.Sprintf("JP-%s:%s", record[0][0:2], record[2])
		codes[code] = record[0]
		prefectures[record[0][0:2]] = record[1]
	}

	for prefecture, section := range sections {
		filename := fmt.Sprintf("data/%s.csv", prefecture)
		templateFilename := fmt.Sprintf("templates/%s.txt", prefecture)
		templateFile, err := os.OpenFile(templateFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 644)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer templateFile.Close()

		fmt.Fprintf(templateFile, "{{#switch:{{{1}}}\n  |источник = <ref>{{cite web|date=%s|url=%s|title=%s|publisher=%s|accessdate=%s|lang=jp|description=%s}}</ref>\n", section["publication_date"], section["url"], section["title"], section["publisher"], section["access_date"], section["description"])
		fmt.Fprintf(templateFile, "  |дата    = %s\n", section["data_date"])

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()
		reader := csv.NewReader(file)
		reader.TrailingComma = true
		reader.TrimLeadingSpace = true
		municipalities := 0
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
				return
			}
			nameId, err := strconv.Atoi(section["name_column_id"])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			populationId, err := strconv.Atoi(section["population_column_id"])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			name := strings.Replace(strings.TrimSpace(record[nameId]), " ", "", -1)
			population := record[populationId]
			_, err = strconv.Atoi(population)
			if err != nil {
				continue
			}
			code := fmt.Sprintf("%s:%s", prefecture, name)

			if name == section["total_name"] {
				fmt.Fprintf(templateFile, "  |%s   = %s\n", prefecture, population)
			}

			if strings.HasSuffix(name, "市") || strings.HasSuffix(name, "町") || strings.HasSuffix(name, "村") {
				fmt.Fprintf(templateFile, "  |%s = %7s <!-- %s -->\n", codes[code], population, name)
				municipalities++
			}
		}
		fmt.Fprintf(templateFile, "  |муниципалитетов = %d\n", municipalities)
		fmt.Fprintf(templateFile, "  |#default = <span style=\"color:red;\">Неверный параметр</span><includeonly>[[Категория:Википедия:Населённые пункты Японии с ошибочным параметром шаблона]]</includeonly>\n")
		fmt.Fprintf(templateFile, "}}<noinclude>\n")
		fmt.Fprintf(templateFile, "{{doc|Население административной единицы Японии/doc}}\n\n")
		fmt.Fprintf(templateFile, "[[Категория:Шаблоны:Население административных единиц Японии|%s]]\n\n", prefecture)
		fmt.Fprintf(templateFile, "[[ja:Template:自治体人口/%s]]\n</noinclude>", prefectures[prefecture[3:]])
	}
}
