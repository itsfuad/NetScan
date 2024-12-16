# NetScan

NetScan is a simple network scanning tool written in Go. It allows you to scan a subnet for open ports on IP addresses within that subnet.

## Features

- Detects the local IP address and subnet.
- Scans a specified range of ports on all IP addresses within the subnet.
- Outputs the IP addresses with open ports.

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/itsfuad/netscan.git
    ```
2. Navigate to the project directory:
    ```sh
    cd netscan
    ```
3. Build the project:
    ```sh
    go build -o netscan main.go
    ```

## Usage

Run the `netscan` executable with the desired options:

```sh
./netscan [options]
```

### Options

- `-ip`: IP address to scan (default: local IP address)
- `-start-port`: Start port for scanning (default: 1)
- `-end-port`: End port for scanning (default: 1024)
- `-timeout`: Timeout for port scanning (default: 2s)

### Example

Scan the local subnet for open ports in the range 1-1024:

```sh
./netscan -start-port 1 -end-port 1024
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgements

- [Go Programming Language](https://golang.org/)
- [net package](https://pkg.go.dev/net)
