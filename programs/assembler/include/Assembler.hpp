//
// Created by mattl on 8/30/2022.
//

#ifndef HACKASSEMBLER_ASSEMBLER_HPP
#define HACKASSEMBLER_ASSEMBLER_HPP

#include <map>
#include <string>
#include <cstdint>
#include <fstream>
#include <vector>


class Assembler {
private:
    /// label, address
    std::map<std::string, std::uint16_t> symbolTable;
    unsigned int currIns;

    /// Source file of asm
    std::ifstream source;
    std::string sourcePath;

    /// Destination file in binary
    std::ofstream dest;
    std::string destPath;

    /// Buffer to hold operating code
    std::vector<std::string> buffer;
    std::vector<std::string> tmpBuffer;

    /// Built in comp instructions
    std::map<std::string, std::string> binFromComp = {
            {"0", "0101010"},
            {"1", "0111111"},
            {"-1", "0111010"},
            {"D", "0001100"},
            {"A", "0110000"},
            {"!D", "0001101"},
            {"!A", "0110001"},
            {"-D", "0001111"},
            {"-A", "0110011"},
            {"D+1", "0011111"},
            {"1+D", "0011111"},
            {"A+1", "0110111"},
            {"1+A", "0110111"},
            {"D-1", "0001110"},
            {"A-1", "0110010"},
            {"D+A", "0000010"},
            {"A+D", "0000010"},
            {"D-A", "0010011"},
            {"A-D", "0000111"},
            {"D&A", "0000000"},
            {"A&D", "0000000"},
            {"D|A", "0010101"},
            {"A|D", "0010101"},
            {"M", "1110000"},
            {"!M", "1110001"},
            {"-M", "1110011"},
            {"M+1", "1110111"},
            {"1+M", "1110111"},
            {"M-1", "1110010"},
            {"D+M", "1000010"},
            {"M+D", "1000010"},
            {"D-M", "1010011"},
            {"M-D", "1000111"},
            {"D&M", "1000000"},
            {"M&D", "1000000"},
            {"D|M", "1010101"},
            {"M|D", "1010101"}
    };

    /// Built in dest instructions
    std::map<std::string, std::string> binFromDest = {
            {"null", "000"},
            {"M", "001"},
            {"D", "010"},
            {"DM", "011"},
            {"MD", "011"},
            {"A", "100"},
            {"AM", "101"},
            {"MA", "101"},
            {"AD", "110"},
            {"DA", "110"},
            {"ADM", "111"},
            {"AMD", "111"},
            {"DMA", "111"},
            {"DAM", "111"},
            {"MAD", "111"},
            {"MDA", "111"}
    };

    /// Built in dest instructions
    std::map<std::string, std::string> binFromJump = {
            {"null", "000"},
            {"JGT", "001"},
            {"JEQ", "010"},
            {"JGE", "011"},
            {"JLT", "100"},
            {"JNE", "101"},
            {"JLE", "110"},
            {"JMP", "111"}
    };

    /**
     * Read the contents of source.txt into the buffer
     */
    void readIntoBuffer();

    /**
     * Drop all whitespace and lines that are only comments
     */
    void dropExtras();

    /**
     * Look for labels and replace label calls with the given instruction
     */
    void processLabels();

    /**
     * Initialize all predefined symbols and prepare to make variables
     */
    void initPredefined();

    /**
     * Convert a c instruction to binary
     * @param cIns C instruction
     * @return String of binary
     */
    std::string cToBin(std::string cIns);

    /**
     * Convert an a instruction to binary
     * @param aIns A instruction
     * @return String of binary
     */
    std::string aToBin(std::string aIns);

    /**
     * Convert the buffer to binary
     */
    void convertToBin();

    /**
     * Write the buffer to the output file
     */
    void writeBuffer();

public:
    /**
     * Constructor that will open the source.txt file and process into the dest file
     *
     * @param source Source path
     * @param dest Destination path
     */
    Assembler(std::string source, std::string dest);


};

#endif //HACKASSEMBLER_ASSEMBLER_HPP
