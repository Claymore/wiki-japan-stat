package main

import (
    "github.com/Claymore/go-config/config"
    "encoding/csv"
    "fmt"
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

    for _, section := range sections {
        filename := fmt.Sprintf("data/%s.csv", section.Name)
        templateFilename := fmt.Sprintf("templates/%s.txt", section.Name)
        templateFile, err := os.OpenFile(templateFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 644)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        defer templateFile.Close()

        fmt.Fprintf(templateFile, "{{#switch:{{{1}}}\n  |источник = <ref>{{cite web|date=%s|url=%s|title=%s|publisher=%s|accessdate=%s|lang=jp|description=%s}}</ref>\n", section.Options["publication_date"], section.Options["url"], section.Options["title"], section.Options["publisher"], section.Options["access_date"], section.Options["description"])
        fmt.Fprintf(templateFile, "  |дата    = %s\n", section.Options["data_date"])

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
            nameId, err := strconv.Atoi(section.Options["name_column_id"])
            if err != nil {
                fmt.Println("Error:", err)
                return
            }
            populationId, err := strconv.Atoi(section.Options["population_column_id"])
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
            code := fmt.Sprintf("%s:%s", section.Name, name)

            if name == section.Options["total_name"] {
                fmt.Fprintf(templateFile, "  |%s   = %s\n", section.Name, population)
            }

            if strings.HasSuffix(name, "市") || strings.HasSuffix(name, "町") || strings.HasSuffix(name, "村") {
                fmt.Fprintf(templateFile, "  |%s = %7s <!-- %s -->\n", codes[code], population, name)
                municipalities = municipalities + 1
            }
        }
        fmt.Fprintf(templateFile, "  |муниципалитетов = %d\n", municipalities)
        fmt.Fprintf(templateFile, "  |#default = <span style=\"color:red;\">Неверный параметр</span><includeonly>[[Категория:Википедия:Населённые пункты Японии с ошибочным параметром шаблона]]</includeonly>\n")
        fmt.Fprintf(templateFile, "}}<noinclude>\n")
        fmt.Fprintf(templateFile, "{{doc|Население административной единицы Японии/doc}}\n\n")
        fmt.Fprintf(templateFile, "[[Категория:Шаблоны:Население административных единиц Японии|%s]]\n\n", section.Name)
        fmt.Fprintf(templateFile, "[[ja:Template:自治体人口/%s]]\n</noinclude>", prefectures[section.Name[3:]])
    }
}
