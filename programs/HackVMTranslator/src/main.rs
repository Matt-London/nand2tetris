mod constants;
mod operations;
mod parser;

use crate::operations::*;

use clap::Parser;

#[derive(Parser)]
struct Cli {
    /// Input file/directory path
    input_path: String,
    /// Output file path
    output_path: Option<String>,
}

fn main() {
    let args = Cli::parse();

    let mut parser = parser::Parser::new(&args.input_path);

    // Get the default name (input_path with .asm)

    let default_name;
    
    if args.input_path.contains(".vm") {
        default_name = args.input_path.replace(".vm", ".asm");
    }
    else {
        let unix_path = args.input_path.replace("\\", "/").to_owned();

        let name_stripped = unix_path.split("/").collect::<Vec<&str>>();

        let mut base_name = "";

        for word in name_stripped.iter().rev() {
            if word != &"" {
                base_name = word;
                break;
            }
        }

        default_name = format!("{}/{}.asm", args.input_path, base_name);
    }

    match args.output_path {
        Some(output) => parser.output(&output),
        None => parser.output(&default_name)
    }
    
}
