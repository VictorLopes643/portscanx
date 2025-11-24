# 🔍 PortScanX

markdown

PortScanX is a **fast and simple port scanner** written in Go, designed with a focus on:

- Offensive security  
- High performance  
- TCP port scanning  
- Educational and lab environments  

Perfect for pentesters, security analysts, and students.

---

## Features

✔ Scan single ports or port ranges  
✔ Supports direct input via **IP** or **file list**  
✔ Clean and organized output  
✔ Stylish ASCII banner  
✔ Simple and easy-to-extend code  
✔ High speed using goroutines  

---

## Installation

### Clone the repository

```sh
git clone https://github.com/VictorLopes643/portscanx.git
cd portscanx
Build the binary
sh
Copy code
go build -o portscanx
Usage
Simple mode (single target)
sh
Copy code
./portscanx 192.168.0.1
Using a target list
sh
Copy code
./portscanx targets.txt
Example of targets.txt
Copy code
192.168.0.1
10.0.0.5