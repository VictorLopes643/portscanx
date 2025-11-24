📘 README.md — PortScanX
# 🔍 PortScanX
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

PortScanX é um **scanner de portas rápido e simples**, desenvolvido em Go, com foco em:

- 🛡️ Segurança ofensiva  
- ⚡ Alta performance  
- 📡 Escaneamento de portas TCP  
- 🧪 Fins educacionais e laboratoriais  

Ideal para pentesters, analistas de segurança e estudantes.

---

## 📦 **Características**

✔ Escaneia portas individuais ou ranges  
✔ Suporte a input direto via **IP** ou via **arquivo**  
✔ Saída clara e organizada  
✔ Banner estiloso  
✔ Código simples e fácil de estender  
✔ Alta velocidade graças à goroutines  

---

## 🛠️ **Instalação**

### **Clonando o repositório**

```sh
git clone https://github.com/VictorLopes643/portscanx.git
cd portscanx

Compilando
go build -o portscanx

🚀 Uso
Modo simples (um alvo):
./portscanx 192.168.0.1

Alvo via lista:
./portscanx alvos.txt

Exemplo de arquivo alvos.txt:
192.168.0.1
10.0.0.5
scanme.nmap.org

