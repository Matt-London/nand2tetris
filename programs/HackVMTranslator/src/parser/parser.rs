use crate::parser::Command;

use std::fs::{File, metadata, read_dir};
use std::path::Path;
use std::io::{prelude::*, BufReader};

use substring::Substring;

#[derive(PartialEq, Eq)]
enum FileType {
    None,
    File,
    Directory
}

pub struct Parser {
    vm_path: String,
    file_type: FileType,
    cmd_cnt: u32,
    command_list: Vec<Command>
}

impl Parser {
    pub fn new(fname: &str) -> Self {
        let mut parser = Parser {
            vm_path: fname.to_owned(),
            file_type: FileType::None,
            cmd_cnt: 0,
            command_list: Vec::new()
        };

        parser.process_project();

        return parser;
    }

    /// Writes the initializer code into the output list
    fn write_init(&mut self) {
        let mut command = Command::new("", self.cmd_cnt, "");
        self.cmd_cnt += 1;
        command.write_init();
        self.command_list.push(command);
        let sysinit = Command::new("call Sys.init 0", self.cmd_cnt, "");
        self.cmd_cnt += 1;
        self.command_list.push(sysinit);
    }

    /// Processes a vm file and appends its commands to the command_list
    fn process_file(&mut self, file_path: &str) {
        // Get base file name
        let path = Path::new(file_path);
        let full_name = path.file_name().unwrap().to_os_string().into_string().unwrap();

        let file_name = full_name.substring(0, full_name.find(".").unwrap());

        let vm_code = File::open(file_path).unwrap();
        let bf = BufReader::new(vm_code);

        for line  in bf.lines() {
            let command = Command::new(&line.unwrap(), self.cmd_cnt, file_name);

            if command.has_command() {
                self.command_list.push(command);
                self.cmd_cnt += 1;
            }
        }
    }

    /// Takes in a directory and loads in and processes each vm file within
    fn process_directory(&mut self, dir_path: &str) {
        let filenames = read_dir(dir_path).unwrap();
        
        // Loop through each file in the directory and list vm files
        for filename in filenames {
            let file = filename.unwrap().path().as_os_str().to_str().unwrap().to_owned();

            if file.contains(".vm") {
                self.process_file(file.as_str());
            }

        }

    }

    fn process_project(&mut self) {
        // Check if this is a file or directory
        let md = metadata(self.vm_path.as_str()).unwrap();
        if md.is_dir() {
            self.file_type = FileType::Directory;
        }
        else {
            self.file_type = FileType::File;
        }

        // Process for file
        match self.file_type {
            FileType::File => {
                self.process_file(self.vm_path.clone().as_str())
            },
            FileType::Directory => {
                self.write_init();
                self.process_directory(self.vm_path.clone().as_str())
            },
            _ => panic!("Input path ({}) is not a valid path", self.vm_path)
        }

    }

    pub fn output(&mut self, output_path: &str) {
        // Open file for outputting
        let mut asm_code = File::create(output_path).expect("Failed to open file");

        // Loop through processed commands
        for cmd in self.command_list.iter() {
            // Loop through each string in cmd
            for asm_cmd in cmd.get_processed().unwrap().iter() {
                asm_code.write_all(format!("{}\n", asm_cmd).as_bytes()).expect("Failed to write");
            }
        }
        // match bw.flush() {
        //     Ok(_)   => (),
        //     Err(_)  => println!("Failed to flush to file")
        // }

    }
}