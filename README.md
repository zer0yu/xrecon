# xrecon

[![Go Version](https://img.shields.io/github/go-mod/go-version/zer0yu/xrecon)](https://github.com/zer0yu/xrecon)
[![License](https://img.shields.io/github/license/zer0yu/xrecon)](https://github.com/zer0yu/xrecon/blob/main/LICENSE)

xrecon is a powerful web fingerprinting tool with CDN detection capabilities. It assists security researchers and penetration testers in quickly identifying the technology stack of target websites and determining if they use a CDN.

## Features

- Support for single URL and bulk URL file inputs
- Web fingerprinting using the [fingers](https://github.com/chainreactors/fingers) library
- Integrated CDN detection with [cdncheck](https://github.com/ExploitSuite/cdncheck)
- Multiple output formats: terminal, CSV, and TXT
- Automatic URL completion (if http:// or https:// is not provided)

## Installation

Ensure you have Go 1.18 or higher installed, then run:

```bash
go get -u github.com/zer0yu/xrecon
```

## Usage

### Basic Usage

Scan a single URL:

```bash
xrecon -url https://example.com
```

Scan a list of URLs from a file:

```bash
xrecon -file urls.txt
```

### Output Options

By default, results are displayed in the terminal. Use the `-output` and `-o` options to specify the output format and file:

```bash
xrecon -url https://example.com -output csv -o results
xrecon -file urls.txt -output txt -o results
```

## Output Example

```
URL: https://example.com
Fingerprint: apache http server:wappalyzer||easy-software-ranzhi-oa:goby
CDN Info: CDN: true, Provider: Cloudflare, Type: cdn

URL: https://google.com
Fingerprint: Google Web Server:wappalyzer
CDN Info: CDN: false, Provider: , Type: 
```

## Contributing

Contributions are welcome! Feel free to submit bug reports, feature requests, or pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Disclaimer

This tool is for educational and research purposes only. Unauthorized testing with this tool may violate laws. Users are responsible for all consequences of using this tool.

## Acknowledgements

- [fingers](https://github.com/chainreactors/fingers) - Web fingerprinting library
- [cdncheck](https://github.com/ExploitSuite/cdncheck) - CDN detection library [support china cdn list]

## Contact

If you have any questions or suggestions, please open an issue on GitHub.

---

Happy scanning!