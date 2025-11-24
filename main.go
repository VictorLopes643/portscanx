package main

import (
    "bufio"
    "bytes"
    "fmt"
    "net"
    "os"
    "os/exec"
    "strings"
    "sync"
    "time"
)

type Job struct {
    IP string
}

func hasNmap() bool {
    _, err := exec.LookPath("nmap")
    return err == nil
}

// Carrega IPs de um arquivo
func carregarIPs(caminho string) ([]string, error) {
    arquivo, err := os.Open(caminho)
    if err != nil {
        return nil, err
    }
    defer arquivo.Close()

    var ips []string
    scanner := bufio.NewScanner(arquivo)
    for scanner.Scan() {
        linha := strings.TrimSpace(scanner.Text())
        if linha != "" {
            ips = append(ips, linha)
        }
    }
    return ips, scanner.Err()
}

func portaAberta(ip string, porta int) bool {
    destino := fmt.Sprintf("%s:%d", ip, porta)
    timeout := 700 * time.Millisecond
    conn, err := net.DialTimeout("tcp", destino, timeout)
    if err != nil {
        return false
    }
    conn.Close()
    return true
}

func portasToString(portas map[int]string) string {
	lista := make([]string, 0, len(portas))
	for porta := range portas {
		lista = append(lista, fmt.Sprintf("%d", porta))
	}
	return strings.Join(lista, ",")
}


func portasNmap(ip string) ([]int, string, error) {

    portas := map[int]string{
        21: "FTP",
        2121: "FTP2",
        22: "SSH",
        80: "HTTP",
        8080: "HTTP2",
        8443: "HTTP3",
        139: "SMB",
        443: "HTTPS",
        445: "SMB",
        3306: "MySQL",
        3389: "RDP",
        25: "SMTP",
        23: "TELNET",
        2049: "NFS",
    }

    portaString := portasToString(portas)

    // PASTA DO HOST
    pastaHost := fmt.Sprintf("%s/hosts/%s", globalPastaPrincipal, ip)
    os.MkdirAll(pastaHost, 0755)

    // XML dentro da pasta do host
    arquivoXML := fmt.Sprintf("%s/Nmap_%s.xml", pastaHost, ip)

    cmd := exec.Command("nmap", "-sS", "-Pn", "-v", "-p"+portaString, "-oX", arquivoXML, ip)

    var out bytes.Buffer
    var stderr bytes.Buffer

    cmd.Stdout = &out
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        return nil, pastaHost, err
    }

    // Parse das portas abertas
    linhas := strings.Split(out.String(), "\n")
    var portasAbertas []int

    for _, linha := range linhas {
        if strings.Contains(linha, "/tcp") && strings.Contains(linha, "open") {
            campos := strings.Fields(linha)
            portaStr := strings.Split(campos[0], "/")[0]
            var porta int
            fmt.Sscanf(portaStr, "%d", &porta)
            portasAbertas = append(portasAbertas, porta)
        }
    }

    return portasAbertas, pastaHost, nil
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup, lock *sync.Mutex) {
    defer wg.Done()

    for job := range jobs {

        portas, pastaHost, err := portasNmap(job.IP)
        if err != nil {
            fmt.Printf("[Worker %d] Erro no NMAP %s: %v\n", id, job.IP, err)
            continue
        }

        for _, p := range portas {
            salvarHost(p, job.IP, pastaHost, lock)
        }
    }
}

func salvarHost(porta int, ip string, pastaHost string, lock *sync.Mutex) {
    if porta < 1 {
        return
    }
    lock.Lock()
    defer lock.Unlock()

    // Arquivo dentro do HOST
    arquivoHost := fmt.Sprintf("%s/porta_%d_open.txt", pastaHost, porta)

    f1, _ := os.OpenFile(arquivoHost, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    f1.WriteString(ip + "\n")
    f1.Close()

    // Arquivo global da porta
    arquivoGlobal := fmt.Sprintf("%s/portas/porta_%d_open_all.txt", globalPastaPrincipal, porta)

    f2, _ := os.OpenFile(arquivoGlobal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    f2.WriteString(ip + "\n")
    f2.Close()
}


var globalPastaPrincipal string

func main() {
	Banner()
	if len(os.Args) < 2 {
		red := "\033[31m"
		reset := "\033[0m"
		yellow := "\033[33m"

		fmt.Printf("%s[-] Uso incorreto!%s\n", red, reset)
		fmt.Printf("%sUso correto:%s\n", yellow, reset)
		fmt.Println("  PortScanX <IP> ou <arquivo.txt>\n")
		fmt.Println("Exemplos:")
		fmt.Println("  PortScanX 192.168.0.1")
		fmt.Println("  PortScanX hosts.txt")
		return
	}

    if hasNmap() {
        fmt.Println("[+] Nmap detected on the system.")
    } else {
        fmt.Println("[-] Nmap is NOT installed.")
    }

    entrada := os.Args[1]
    var ips []string
	
    if strings.HasSuffix(entrada, ".txt") {
        lista, err := carregarIPs(entrada)
        if err != nil {
            fmt.Println("Erro ao carregar arquivo:", err)
            return
        }
        ips = lista
    } else {
        ips = append(ips, entrada)
    }

    // -----------------------------
    // Criar PASTA PRINCIPAL da execução
    // -----------------------------
    dataExecucao := time.Now().Format("2006-01-02_15-04-05")
    pastaPrincipal := fmt.Sprintf("scan_%s", dataExecucao)

    os.MkdirAll(pastaPrincipal+"/hosts", 0755)
    os.MkdirAll(pastaPrincipal+"/portas", 0755)

    // Variável global acessível no programa inteiro
    globalPastaPrincipal = pastaPrincipal

    fmt.Println("Pasta principal criada em:", globalPastaPrincipal)
    fmt.Println("Iniciando scans NMAP por host...\n")

    // -----------------------------
    // CONFIG WORKERS
    // -----------------------------
    numWorkers := 10
    jobs := make(chan Job, 100)
    var wg sync.WaitGroup
    var lock sync.Mutex

    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, jobs, &wg, &lock)
    }

    // Enviar IPs para os workers
    for _, ip := range ips {
        jobs <- Job{IP: ip}
    }

    close(jobs)
    wg.Wait()

    fmt.Printf("\nConcluído! Resultados salvos em: %s\n", globalPastaPrincipal)
}

func Banner() {
    fmt.Println(`
██████╗  ██████╗ ██████╗ ████████╗███████╗ ██████╗ █████╗ ███╗   ██╗██╗  ██╗
██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝██╔════╝██╔══██╗████╗  ██║╚██╗██╔╝
██████╔╝██║   ██║██████╔╝   ██║   ███████╗██║     ███████║██╔██╗ ██║ ╚███╔╝ 
██╔═══╝ ██║   ██║██╔══██╗   ██║   ╚════██║██║     ██╔══██║██║╚██╗██║ ██╔██╗ 
██║     ╚██████╔╝██║  ██║   ██║   ███████║╚██████╗██║  ██║██║ ╚████║██╔╝ ██╗
╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝                         
                                                                                                     
        [ Programa ]  PortScanX v1.0
        [ Autor    ]  Victor Lopes
        [ LinkedIn ]  linkedin.com/in/victorlopes643
        [ GitHub   ]  github.com/victorlopes643

        Ferramenta focada em análise e segurança ofensiva.
        Desenvolvida para fins educacionais e pesquisas.
    `)
}