#include <iostream>

#include "Assembler.hpp"

int main(int argc, char* argv[]) {
    if (argc != 3) {
        std::cout << "Usage:" << std::endl;
        std::cout << "./assembler [input file] [output dest]" << std::endl;
        return 1;
    }

    Assembler as(argv[1], argv[2]);

    return 0;
}
