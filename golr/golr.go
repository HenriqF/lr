package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var sections map[string][]string = make(map[string][]string)
var sections_life map[string]int = make(map[string]int)

var program_output bytes.Buffer
var program_input []string

func arquivo_existe(nome string) bool {
	_, err := os.Stat(nome)
	if err != nil {
		return false
	}
	return true
}

func ler_arquivo(nome string) (string, []byte) {
	if !arquivo_existe(nome) {
		log.Fatalf("Arquivo não existe: [%v]", nome)
	}
	cont, err := os.ReadFile(nome)
	if err != nil {
		return "", nil
	}
	return string(cont), cont
}

func out_contains(arg string) {
	last_out := program_output.String()
	if !strings.Contains(last_out, arg) {
		log.Fatalf("Ultimo resultado não tinha [%v]", arg)
	}
}

func run_control(arg string, silent bool) {
	program_output.Reset()
	cmd := exec.Command("cmd", "/C", arg)

	cmd.Stdout = &program_output
	cmd.Stderr = &program_output
	stdin, _ := cmd.StdinPipe()

	if !silent {
		fmt.Printf("\nrun_control==============\n")
	}
	cmd.Start()
	// if err != nil {
	// 	log.Fatalf("deu merda %v", err)
	// }
	for i := 0; i < len(program_input); i++ {
		stdin.Write([]byte(program_input[i] + "\n"))
	}
	program_input = program_input[:0]

	cmd.Wait()

	fmt.Printf("%v", program_output.String())
	if !silent {
		fmt.Printf("\n=========================\n")
	}
}

func run(arg string, silent bool) {
	cmd := exec.Command("cmd", "/C", arg)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if !silent {
		fmt.Printf("\nrun======================\n")
	}
	cmd.Run()
	// if err != nil {
	// 	log.Fatalf("deu merda %v", err)
	// }
	if !silent {
		fmt.Printf("\n=========================\n")
	}
}

func read_section(section_name string) {
	fmt.Printf("Executando: [%v]\n", section_name)
	section := sections[section_name]
	if sections_life[section_name] <= 0 {
		log.Fatalf("[%v] não pode ser usado novamente/ não existe", section_name)
	}
	sections_life[section_name] -= 1

	for i := 0; i < len(section); i++ {
		cmd_end := 0

		if strings.TrimSpace(section[i]) == "" {
			continue
		}

		if !strings.Contains(section[i], ":") {
			log.Fatalf("comando sem ':' %v", section[i])
		}

		for ; cmd_end < len(section[i]); cmd_end++ {
			if section[i][cmd_end] == ':' {
				break
			}
		}
		command := section[i][0:cmd_end]
		arg := strings.TrimSpace(section[i][cmd_end+1 : len(section[i])])

		switch command {
		case "goto":
			read_section(arg)

		case "run":
			run(arg, false)

		case "run_control":
			run_control(arg, false)

		case "srun":
			run(arg, true)

		case "srun_control":
			run_control(arg, true)

		case "out_contains":
			out_contains(arg)

		case "prog_input":
			program_input = append(program_input, arg)

		case "wait":
			tts, err := strconv.Atoi(arg)
			if err != nil {
				log.Fatalf("Quantia de tempo imprópria: [%v]", arg)
			}
			time.Sleep(time.Duration(tts) * time.Millisecond)

		case "sleep":
			fmt.Println("Enter para continuar...")
			s := bufio.NewScanner(os.Stdin)
			s.Scan()

		case "log":
			fmt.Println(arg)

		default:
			fmt.Println(command)
		}

	}
}

func main() {
	log.SetFlags(0)
	arq_lr_nome := "lr"
	args := os.Args

	entry_section := "main"

	if len(args) > 1 {
		if args[1] == "init" && !arquivo_existe(arq_lr_nome) {
			fmt.Printf("Criando lr...")
			os.WriteFile(arq_lr_nome, []byte("[main]\n    log: init\n    end"), 0666)
			os.Exit(0)
		}
		entry_section = args[1]
	}

	conteudo, _ := ler_arquivo(arq_lr_nome)
	conteudo = strings.TrimSpace(conteudo)
	i := 0
	section_start := 0

	for ; i < len(conteudo); i++ {
		if conteudo[i] == '[' {
			section_start = i
			for ; i < len(conteudo); i++ {
				if conteudo[i] == ']' {
					break
				}
				if i == len(conteudo)-1 && conteudo[i-1] != ']' {
					log.Fatalf("fix par aberto")
				}
			}
			section_name := conteudo[section_start+1 : i]
			section_content := []string{}

			for _, v := range strings.Split(conteudo[i+1:], "\n") {
				line := strings.TrimSpace(v)
				if line == "end" {
					break
				}
				section_content = append(section_content, line)
			}

			sections[section_name] = section_content
			sections_life[section_name] = 1
		}
	}

	if len(sections[entry_section]) == 0 {
		log.Fatalf("[%v] não existe.\n", entry_section)
	}

	read_section(entry_section)

}
