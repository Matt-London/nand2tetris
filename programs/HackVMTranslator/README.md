# HackVMTranslator
Rust implementation of the Nand2Tetris VM Translator

## Compiling
Use cargo to download dependencies and compile
```shell
$ cargo build
```

## Executing
Execute the compiled binary supplying the following parameters to translate VM to ASM
```shell
$ ./hack_vm_translator --help
Usage: hack_vm_translator <INPUT_PATH> [OUTPUT_PATH]

Arguments:
  <INPUT_PATH>   Input file path
  [OUTPUT_PATH]  Output file path

Options:
  -h, --help  Print help information
```