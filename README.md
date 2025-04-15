

# goqr 
## Simple Go QR Code Generator

A simple Go project that generates QR codes from a given URI. Optionally, it should (but does not) [TODO](https://github.com/Cheezmeister/goqr/issues/1) overlay a logo in the center of the QR code.

## Installation

To install the project, clone the repository and run the following command to download the necessary dependencies:

## Usage

To generate a QR code, run the following command:

```bash
go run qr-code-generator.go -uri "https://example.com" -logo "logo.svg" -output "output.png" -size 256 
```

- `-uri`: The URI to encode in the QR code (required).
- `-logo`: Path to a logo image (SVG or PNG) to place in the center (optional).
- `-output`: Output filename (default is `qrcode.png`).
- `-size`: Size of the QR code in pixels (default is 256).

## Contributing

Please feel free to submit a pull request or open an issue if you have suggestions or improvements.

## License

This project is licensed under the BML License.

See the LICENSE file for details.

## Contact

[Here.](/Cheezmeister)
