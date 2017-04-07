package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// файл с конфигурацией

type Configuration struct {
	RootFolder     string   `json:"rootFolder"`
	ExcludeFolders []string `json:"excludeFolders"`
	Debug          bool     `json:"debug"`
}

func main() {
	// попытка прочитать конфиг
	configuration, err := readConfig()
	// пробуем открыть папку из конфига
	dir, err := os.Open(configuration.RootFolder)
	defer dir.Close()
	// если нет - то фаталь
	logFatal(err)

	// читаем содержимое
	fileInfos, err := dir.Readdir(-1)
	// фаталь
	logFatal(err)

	for _, fi := range fileInfos {
		// печатаем имя файла
		if configuration.Debug {
			log.Printf("%s\n", configuration.RootFolder+"/"+fi.Name())
		}
		// если это папка
		if fi.IsDir() {
			if configuration.Debug {
				log.Printf("======== \\%s\\%s ======\n", configuration.RootFolder, fi.Name())
			}
			// для вывода имени папки нужны преобразования ....
			var buffer bytes.Buffer
			buffer.WriteString(configuration.RootFolder)
			buffer.WriteString("/")
			buffer.WriteString(fi.Name())
			// пошли в рекурсию
			readDirRec(buffer.String(), configuration.Debug, configuration.ExcludeFolders)
			// удялем последнюю папку в каталоге если она не в исалючениях!
			if !checkExclude(buffer.String(), configuration.ExcludeFolders) {
				err := os.Remove(buffer.String())
				logError(err)
			}

		} else {
			// если не папка - сразу удаляем
			err := os.Remove(configuration.RootFolder + "/" + fi.Name())
			logError(err)
		}
	}
}

// есть дублирование кода, в  main - но пока хз как лучше
func readDirRec(mydir string, debug bool, exclude []string) {

	dir, err := os.Open(mydir)

	logError(err)

	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)

	logError(err)

	for _, fi := range fileInfos {

		if debug {
			log.Printf("%s\n", mydir+"/"+fi.Name())
		}
		if fi.IsDir() {
			if debug {
				log.Printf("======== \\%s\\%s ======\n", mydir, fi.Name())
			}
			var buffer bytes.Buffer
			buffer.WriteString(mydir)
			buffer.WriteString("/")
			buffer.WriteString(fi.Name())
			readDirRec(buffer.String(), debug, exclude)
			// при выходе нужно удалить и папку если она не в исалючениях!
			if !checkExclude(mydir, exclude) {
				err := os.Remove(buffer.String())
				logError(err)
			}

		} else {
			err := os.Remove(mydir + "/" + fi.Name())
			logError(err)
		}
	}

}

// логгер ошибок
func logError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

// обработчик фатальных ошибок
func logFatal(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// читаем файл с конфигом, и возвращаем структуру
func readConfig() (Configuration, error) {
	// пробуем открыть файл конфига
	file, err := os.Open("conf.json")
        logFatal(err)

	var config Configuration

	jsonParser := json.NewDecoder(file)

	if err = jsonParser.Decode(&config); err != nil {
		fmt.Println("parsing config file", err.Error())
	}

	return config, err

}

// проверяем папки к исключению
func checkExclude(dir string, exclude []string) bool {

	for _, value := range exclude {
		// fmt.Println(dir)
		// fmt.Println(value)
		if value == dir {
			return true
		}

	}

	return false

}

// DISABLE: если не можем - запрашиваем с терминала
//  var config_file string
//  // запросим путь к нему в интерактиве
//  reader := bufio.NewReader(os.Stdin)
//  fmt.Print("Enter full path to config file [default conf.json]: ")
//  config_file, _ := reader.ReadString('\n')
//  // режем финальный перевод строки
//  config_file = strings.TrimSpace(config_file)
//  // пробуем открыть его
//  readConfig(config_file)
// }
